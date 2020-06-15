package script

import (
	"bytes"
	"encoding/binary"
	"log"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/samuelyuan/openbiohazard2/fileio"
	"github.com/samuelyuan/openbiohazard2/game"
	"github.com/samuelyuan/openbiohazard2/render"
)

const (
	SCRIPT_FRAMES_PER_SECOND = 30.0

	INSTRUCTION_BREAK_FLOW = 0
	INSTRUCTION_NORMAL     = 1
	INSTRUCTION_THREAD_END = 2

	WORKSET_PLAYER = 1
	WORKSET_ENEMY  = 3
	WORKSET_OBJECT = 4
)

var (
	scriptThread    *ScriptThread
	scriptDeltaTime = 0.0
)

type ScriptDef struct {
	ScriptThreads []*ScriptThread
}

func NewScriptDef() *ScriptDef {
	scriptThreads := make([]*ScriptThread, 20)
	for i := 0; i < len(scriptThreads); i++ {
		scriptThreads[i] = NewScriptThread()
	}

	return &ScriptDef{
		ScriptThreads: scriptThreads,
	}
}

func (scriptDef *ScriptDef) Reset() {
	for i := 0; i < len(scriptDef.ScriptThreads); i++ {
		scriptDef.ScriptThreads[i].Reset()
	}
}

func (scriptDef *ScriptDef) InitScript(
	scriptData fileio.ScriptFunction,
	threadNum int,
	startFunction int) {
	scriptDef.ScriptThreads[threadNum].RunStatus = true
	scriptDef.ScriptThreads[threadNum].ProgramCounter = scriptData.StartProgramCounter[startFunction]
}

func (scriptDef *ScriptDef) RunScript(
	scriptData fileio.ScriptFunction,
	timeElapsedSeconds float64,
	gameDef *game.GameDef,
	renderDef *render.RenderDef) {
	for i := 0; i < len(scriptDef.ScriptThreads); i++ {
		// Regulate frames per second
		scriptDeltaTime += timeElapsedSeconds
		if scriptDeltaTime > 1.0/SCRIPT_FRAMES_PER_SECOND {
			scriptDeltaTime = 0.0
		} else {
			continue
		}

		scriptDef.RunScriptThread(scriptDef.ScriptThreads[i], scriptData, gameDef, renderDef)
	}
}

func (scriptDef *ScriptDef) RunScriptThread(
	curScriptThread *ScriptThread,
	scriptData fileio.ScriptFunction,
	gameDef *game.GameDef,
	renderDef *render.RenderDef) {
	scriptThread = curScriptThread

	// Thread should not run
	if scriptThread.RunStatus == false {
		return
	}

	for true {
		sectionReturnValue := scriptDef.RunScriptUntilBreakControlFlow(scriptData, gameDef, renderDef)

		// End thread
		if scriptThread.ShouldTerminate(sectionReturnValue) {
			break
		}

		scriptThread.JumpToNextLocationOnStack()
	}
}

func (scriptDef *ScriptDef) RunScriptUntilBreakControlFlow(
	scriptData fileio.ScriptFunction,
	gameDef *game.GameDef,
	renderDef *render.RenderDef) int {
	scriptReturnValue := 0
	for true {
		lineData := scriptData.Instructions[scriptThread.ProgramCounter]
		opcode := lineData[0]

		// Override can be modified during execution
		scriptThread.OverrideProgramCounter = false

		instructionReturnValue := scriptDef.ExecuteSingleInstruction(lineData, scriptData, gameDef, renderDef)

		if !scriptThread.OverrideProgramCounter {
			scriptThread.IncrementProgramCounter(opcode)
		}
		scriptThread.OverrideProgramCounter = false

		// Control flow is broken
		if instructionReturnValue != INSTRUCTION_NORMAL {
			scriptReturnValue = instructionReturnValue
			break
		}
	}
	return scriptReturnValue
}

func (scriptDef *ScriptDef) ExecuteSingleInstruction(
	lineData []byte,
	scriptData fileio.ScriptFunction,
	gameDef *game.GameDef,
	renderDef *render.RenderDef) int {
	var returnValue int

	opcode := lineData[0]
	switch opcode {
	case fileio.OP_EVT_END:
		returnValue = scriptDef.ScriptEvtEnd(lineData)
	case fileio.OP_EVT_EXEC:
		returnValue = scriptDef.ScriptEvtExec(lineData, scriptData)
	case fileio.OP_IF_START:
		returnValue = scriptDef.ScriptIfBlockStart(scriptThread, lineData)
	case fileio.OP_ELSE_START:
		returnValue = scriptDef.ScriptElseCheck(scriptThread, lineData)
	case fileio.OP_END_IF:
		returnValue = scriptDef.ScriptEndIf()
	case fileio.OP_SLEEP:
		returnValue = scriptDef.ScriptSleep(lineData)
	case fileio.OP_SLEEPING:
		returnValue = scriptDef.ScriptSleeping(lineData)
	case fileio.OP_FOR:
		returnValue = scriptDef.ScriptForLoopBegin(lineData)
	case fileio.OP_FOR_END:
		returnValue = scriptDef.ScriptForLoopEnd(lineData)
	case fileio.OP_SWITCH:
		returnValue = scriptDef.ScriptSwitchBegin(lineData, scriptData.Instructions, gameDef)
	case fileio.OP_CASE:
		returnValue = 1 // already implemented in switch statement
	case fileio.OP_DEFAULT:
		returnValue = 1 // already implemented in switch statement
	case fileio.OP_END_SWITCH:
		returnValue = scriptDef.ScriptSwitchEnd()
	case fileio.OP_GOTO:
		returnValue = scriptDef.ScriptGoto(lineData)
	case fileio.OP_GOSUB:
		returnValue = scriptDef.ScriptGoSub(lineData, scriptData)
	case fileio.OP_BREAK:
		returnValue = scriptDef.ScriptBreak(lineData)
	case fileio.OP_CHECK: // 0x21
		returnValue = scriptDef.ScriptCheckBit(lineData, gameDef)
	case fileio.OP_SET_BIT: // 0x22
		returnValue = scriptDef.ScriptSetBit(lineData, gameDef)
	case fileio.OP_COMPARE: // 0x23
		returnValue = scriptDef.ScriptCompare(lineData, gameDef)
	case fileio.OP_SAVE: // 0x24
		returnValue = scriptDef.ScriptSave(lineData, gameDef)
	case fileio.OP_COPY: // 0x25
		returnValue = scriptDef.ScriptCopy(lineData, gameDef)
	case fileio.OP_CALC: // 0x26
		returnValue = scriptDef.ScriptCalc(lineData, gameDef)
	case fileio.OP_CALC2: // 0x27
		returnValue = scriptDef.ScriptCalc(lineData, gameDef)
	case fileio.OP_CUT_CHG:
		returnValue = scriptDef.ScriptCameraChange(lineData, gameDef)
	case fileio.OP_AOT_SET:
		returnValue = scriptDef.ScriptAotSet(lineData, gameDef)
	case fileio.OP_OBJ_MODEL_SET:
		returnValue = scriptDef.ScriptObjectModelSet(lineData, renderDef)
	case fileio.OP_WORK_SET:
		returnValue = scriptDef.ScriptWorkSet(lineData)
	case fileio.OP_POS_SET:
		returnValue = scriptDef.ScriptPositionSet(lineData, gameDef)
	case fileio.OP_MEMBER_SET:
		returnValue = scriptDef.ScriptMemberSet(lineData, gameDef, renderDef)
	case fileio.OP_SCA_ID_SET:
		returnValue = scriptDef.ScriptScaIdSet(lineData, gameDef)
	case fileio.OP_SCE_ESPR_ON:
		returnValue = scriptDef.ScriptSceEsprOn(lineData, gameDef, renderDef)
	case fileio.OP_DOOR_AOT_SET:
		returnValue = scriptDef.ScriptDoorAotSet(lineData, gameDef)
	case fileio.OP_MEMBER_CMP:
		returnValue = scriptDef.ScriptMemberCompare(lineData)
	case fileio.OP_PLC_MOTION: // 0x3f
		returnValue = scriptDef.ScriptPlcMotion(lineData)
	case fileio.OP_PLC_DEST: // 0x40
		returnValue = scriptDef.ScriptPlcDest(lineData)
	case fileio.OP_PLC_NECK: // 0x41
		returnValue = scriptDef.ScriptPlcNeck(lineData)
	case fileio.OP_SCE_EM_SET: // 0x44
		returnValue = scriptDef.ScriptSceEmSet(lineData)
	case fileio.OP_AOT_RESET: // 0x46
		returnValue = scriptDef.ScriptAotReset(lineData, gameDef)
	case fileio.OP_SCE_ESPR_KILL: // 0x4c
		returnValue = scriptDef.ScriptSceEsprKill(lineData)
	case fileio.OP_ITEM_AOT_SET: // 0x4e
		returnValue = scriptDef.ScriptItemAotSet(lineData, gameDef)
	case fileio.OP_SCE_BGM_CONTROL: // 0x51
		returnValue = scriptDef.ScriptSceBgmControl(lineData)
	case fileio.OP_AOT_SET_4P:
		returnValue = scriptDef.ScriptAotSet4p(lineData, gameDef)
	case fileio.OP_DOOR_AOT_SET_4P:
		returnValue = scriptDef.ScriptDoorAotSet4p(lineData, gameDef)
	case fileio.OP_ITEM_AOT_SET_4P:
		returnValue = scriptDef.ScriptItemAotSet4p(lineData, gameDef)
	default:
		returnValue = 1
	}

	return returnValue
}

func (scriptDef *ScriptDef) ScriptSleep(lineData []byte) int {
	byteArr := bytes.NewBuffer(lineData)
	instruction := fileio.ScriptInstrSleep{}
	binary.Read(byteArr, binary.LittleEndian, &instruction)

	// goes to sleeping instruction (0xa)
	curLevelState := scriptThread.LevelState[scriptThread.SubLevel]

	scriptThread.ProgramCounter = scriptThread.ProgramCounter + 1
	scriptThread.OverrideProgramCounter = true

	curLevelState.LoopLevel++
	newLoopLevel := curLevelState.LoopLevel
	curLevelState.LoopState[newLoopLevel].Counter = int(instruction.Count)
	return 1
}

func (scriptDef *ScriptDef) ScriptSleeping(lineData []byte) int {
	opcode := lineData[0]
	curLevelState := scriptThread.LevelState[scriptThread.SubLevel]
	curLoopState := curLevelState.LoopState[curLevelState.LoopLevel]

	curLoopState.Counter--
	if curLoopState.Counter == 0 {
		scriptThread.ProgramCounter += fileio.InstructionSize[opcode]
		curLevelState.LoopLevel--
	}

	scriptThread.OverrideProgramCounter = true

	return INSTRUCTION_THREAD_END
}

func (scriptDef *ScriptDef) ScriptCheckBit(lineData []byte, gameDef *game.GameDef) int {
	byteArr := bytes.NewBuffer(lineData)
	bitTest := fileio.ScriptInstrCheckBitTest{}
	binary.Read(byteArr, binary.LittleEndian, &bitTest)

	if gameDef.GetBitArray(int(bitTest.BitArray), int(bitTest.Number)) == int(bitTest.Value) {
		return 1
	}
	return INSTRUCTION_BREAK_FLOW
}

func (scriptDef *ScriptDef) ScriptSetBit(lineData []byte, gameDef *game.GameDef) int {
	byteArr := bytes.NewBuffer(lineData)
	instruction := fileio.ScriptInstrSetBit{}
	binary.Read(byteArr, binary.LittleEndian, &instruction)

	switch int(instruction.Operation) {
	case 0:
		// Clear bit
		gameDef.SetBitArray(int(instruction.BitArray), int(instruction.BitNumber), 0)
	case 1:
		// Set bit
		gameDef.SetBitArray(int(instruction.BitArray), int(instruction.BitNumber), 1)
	case 7:
		// Flip bit
		currentBit := gameDef.GetBitArray(int(instruction.BitArray), int(instruction.BitNumber))
		gameDef.SetBitArray(int(instruction.BitArray), int(instruction.BitNumber), currentBit^1)
	default:
		log.Fatal("Set bit operation ", instruction.Operation, " is invalid.")
	}

	return 1
}

func (scriptDef *ScriptDef) ScriptCompare(lineData []byte, gameDef *game.GameDef) int {
	byteArr := bytes.NewBuffer(lineData)
	instruction := fileio.ScriptInstrCompare{}
	binary.Read(byteArr, binary.LittleEndian, &instruction)

	variableValue := gameDef.GetScriptVariable(int(instruction.VarId))
	otherValue := int(instruction.Value)

	switch int(instruction.Operation) {
	case 0:
		if variableValue == otherValue {
			return 1
		} else {
			return INSTRUCTION_BREAK_FLOW
		}
	case 1:
		// greater than
		if variableValue > otherValue {
			return 1
		} else {
			return INSTRUCTION_BREAK_FLOW
		}
	case 2:
		// greater than or equals to
		if variableValue >= otherValue {
			return 1
		} else {
			return INSTRUCTION_BREAK_FLOW
		}
	case 3:
		// less than
		if variableValue < otherValue {
			return 1
		} else {
			return INSTRUCTION_BREAK_FLOW
		}
	case 4:
		// less than or equals to
		if variableValue <= otherValue {
			return 1
		} else {
			return INSTRUCTION_BREAK_FLOW
		}
	case 5:
		// not equals
		if variableValue != otherValue {
			return 1
		} else {
			return INSTRUCTION_BREAK_FLOW
		}
	case 6:
		if variableValue&otherValue != 0 {
			return 1
		} else {
			return INSTRUCTION_BREAK_FLOW
		}
	}

	return 1
}

func (scriptDef *ScriptDef) ScriptSave(lineData []byte, gameDef *game.GameDef) int {
	byteArr := bytes.NewBuffer(lineData)
	instruction := fileio.ScriptInstrSave{}
	binary.Read(byteArr, binary.LittleEndian, &instruction)

	gameDef.SetScriptVariable(int(instruction.VarId), int(instruction.Value))
	return 1
}

func (scriptDef *ScriptDef) ScriptCopy(lineData []byte, gameDef *game.GameDef) int {
	byteArr := bytes.NewBuffer(lineData)
	instruction := fileio.ScriptInstrCopy{}
	binary.Read(byteArr, binary.LittleEndian, &instruction)

	sourceValue := gameDef.GetScriptVariable(int(instruction.SourceVarId))
	gameDef.SetScriptVariable(int(instruction.DestVarId), sourceValue)
	return 1
}

func (scriptDef *ScriptDef) ScriptCalc(lineData []byte, gameDef *game.GameDef) int {
	byteArr := bytes.NewBuffer(lineData)
	instruction := fileio.ScriptInstrCalc{}
	binary.Read(byteArr, binary.LittleEndian, &instruction)

	leftValue := int(gameDef.GetScriptVariable(int(instruction.VarId)))
	rightValue := int(instruction.Value)
	result := scriptDef.ScriptVariableCalculator(int(instruction.Operation), leftValue, rightValue)
	gameDef.SetScriptVariable(int(instruction.VarId), result)
	return 1
}

func (scriptDef *ScriptDef) ScriptCalc2(lineData []byte, gameDef *game.GameDef) int {
	byteArr := bytes.NewBuffer(lineData)
	instruction := fileio.ScriptInstrCalc2{}
	binary.Read(byteArr, binary.LittleEndian, &instruction)

	leftValue := int(gameDef.GetScriptVariable(int(instruction.VarId)))
	rightValue := int(gameDef.GetScriptVariable(int(instruction.SourceVarId)))
	result := scriptDef.ScriptVariableCalculator(int(instruction.Operation), leftValue, rightValue)
	gameDef.SetScriptVariable(int(instruction.VarId), result)
	return 1
}

func (scriptDef *ScriptDef) ScriptVariableCalculator(operation int, leftValue int, rightValue int) int {
	switch operation {
	case 0:
		return leftValue + rightValue
	case 1:
		return leftValue - rightValue
	case 2:
		return leftValue * rightValue
	case 3:
		return leftValue / rightValue
	case 4:
		return leftValue % rightValue
	case 5:
		return leftValue | rightValue
	case 6:
		return leftValue & rightValue
	case 7:
		return leftValue ^ rightValue
	case 8:
		return ^leftValue
	case 9:
		return leftValue << (rightValue % 32)
	case 10:
		return leftValue >> (rightValue % 32)
	case 11:
		return leftValue >> (rightValue % 32)
	default:
		log.Fatal("Script variable calculator operation ", operation, " is invalid.")
	}

	return INSTRUCTION_BREAK_FLOW
}

func (scriptDef *ScriptDef) ScriptCameraChange(lineData []byte, gameDef *game.GameDef) int {
	byteArr := bytes.NewBuffer(lineData)
	instruction := fileio.ScriptInstrCutChg{}
	binary.Read(byteArr, binary.LittleEndian, &instruction)

	gameDef.ChangeCamera(int(instruction.CameraId))
	return 1
}

func (scriptDef *ScriptDef) ScriptObjectModelSet(lineData []byte,
	renderDef *render.RenderDef) int {

	byteArr := bytes.NewBuffer(lineData)
	instruction := fileio.ScriptInstrObjModelSet{}
	binary.Read(byteArr, binary.LittleEndian, &instruction)

	renderDef.SetItemEntity(instruction)
	return 1
}

func (scriptDef *ScriptDef) ScriptWorkSet(lineData []byte) int {
	byteArr := bytes.NewBuffer(lineData)
	instruction := fileio.ScriptInstrWorkSet{}
	binary.Read(byteArr, binary.LittleEndian, &instruction)

	scriptThread.WorkSetComponent = int(instruction.Component)
	scriptThread.WorkSetIndex = int(instruction.Index)
	return 1
}

func (scriptDef *ScriptDef) ScriptPositionSet(lineData []byte, gameDef *game.GameDef) int {
	byteArr := bytes.NewBuffer(lineData)
	instruction := fileio.ScriptInstrPosSet{}
	binary.Read(byteArr, binary.LittleEndian, &instruction)

	if scriptThread.WorkSetComponent == WORKSET_PLAYER {
		gameDef.Player.Position = mgl32.Vec3{float32(instruction.X), float32(instruction.Y), float32(instruction.Z)}
	} else {
		// TODO: set position of object
	}

	return 1
}

func (scriptDef *ScriptDef) ScriptMemberSet(lineData []byte, gameDef *game.GameDef, renderDef *render.RenderDef) int {
	byteArr := bytes.NewBuffer(lineData)
	instruction := fileio.ScriptInstrMemberSet{}
	binary.Read(byteArr, binary.LittleEndian, &instruction)

	if scriptThread.WorkSetComponent == WORKSET_PLAYER {
		switch int(instruction.MemberIndex) {
		case 15:
			// convert to angle in degrees
			gameDef.Player.RotationAngle = (float32(instruction.Value) / 4096.0) * 360.0
		}
	} else if scriptThread.WorkSetComponent == WORKSET_OBJECT {
		modelObject := renderDef.ItemGroupEntity.ModelObjectData[int(scriptThread.WorkSetIndex)]
		switch int(instruction.MemberIndex) {
		case 15:
			// convert to angle in degrees
			modelObject.RotationAngle = (float32(instruction.Value) / 4096.0) * 360.0
		}
	} else {
		// TODO: set attribute of object
	}
	return 1
}

func (scriptDef *ScriptDef) ScriptScaIdSet(lineData []byte, gameDef *game.GameDef) int {
	byteArr := bytes.NewBuffer(lineData)
	instruction := fileio.ScriptInstrScaIdSet{}
	binary.Read(byteArr, binary.LittleEndian, &instruction)

	if instruction.Flag == 0 {
		gameDef.RemoveCollisionEntity(gameDef.GameRoom.CollisionEntities, int(instruction.Id))
	}
	return 1
}

func (scriptDef *ScriptDef) ScriptMemberCompare(lineData []byte) int {
	byteArr := bytes.NewBuffer(lineData)
	instruction := fileio.ScriptInstrMemberCompare{}
	binary.Read(byteArr, binary.LittleEndian, &instruction)

	return 1
}

func (scriptDef *ScriptDef) ScriptSceEmSet(lineData []byte) int {
	return 1
}

func (scriptDef *ScriptDef) ScriptSceBgmControl(lineData []byte) int {
	byteArr := bytes.NewBuffer(lineData)
	instruction := fileio.ScriptInstrSceBgmControl{}
	binary.Read(byteArr, binary.LittleEndian, &instruction)

	return 1
}
