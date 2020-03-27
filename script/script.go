package script

import (
	"bytes"
	"encoding/binary"
	"log"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/samuelyuan/openbiohazard2/fileio"
	"github.com/samuelyuan/openbiohazard2/game"
	// "sort"
)

const (
	SCRIPT_FRAMES_PER_SECOND = 30.0

	WORKSET_PLAYER = 1
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
	gameDef *game.GameDef) {
	for i := 0; i < len(scriptDef.ScriptThreads); i++ {
		scriptDef.RunScriptThread(scriptDef.ScriptThreads[i], scriptData, timeElapsedSeconds, gameDef)
	}
}

func (scriptDef *ScriptDef) RunScriptThread(
	curScriptThread *ScriptThread,
	scriptData fileio.ScriptFunction,
	timeElapsedSeconds float64,
	gameDef *game.GameDef) {
	/*programCounters := make([]int, 0)
	for counter, _ := range scriptData.Instructions {
		programCounters = append(programCounters, counter)
	}
	sort.Ints(programCounters)*/

	scriptDeltaTime += timeElapsedSeconds
	if scriptDeltaTime > 1.0/SCRIPT_FRAMES_PER_SECOND {
		scriptDeltaTime = 0.0
	} else {
		return
	}

	scriptThread = curScriptThread
	if scriptThread.RunStatus == false {
		return
	}

	for true {
		scriptReturnValue := 0
		for true {
			lineData := scriptData.Instructions[scriptThread.ProgramCounter]
			opcode := lineData[0]

			scriptThread.OverrideProgramCounter = false

			var returnValue int
			switch opcode {
			case fileio.OP_EVT_END:
				returnValue = scriptDef.ScriptEvtEnd(lineData)
			case fileio.OP_EVT_EXEC:
				returnValue = scriptDef.ScriptEvtExec(lineData, scriptData)
			case fileio.OP_IF_START:
				returnValue = scriptDef.ScriptIfBlockStart(lineData)
			case fileio.OP_ELSE_START:
				returnValue = scriptDef.ScriptElseCheck(lineData)
			case fileio.OP_END_IF:
				returnValue = scriptDef.ScriptEndIf()
			case fileio.OP_SLEEP:
				returnValue = scriptDef.ScriptSleep(lineData)
			case fileio.OP_SLEEPING:
				returnValue = scriptDef.ScriptSleeping(lineData)
			case fileio.OP_SWITCH:
				returnValue = scriptDef.ScriptSwitchBegin(lineData, scriptData.Instructions, gameDef)
			case fileio.OP_CASE:
				returnValue = 1
			case fileio.OP_DEFAULT:
				returnValue = 1
			case fileio.OP_END_SWITCH:
				returnValue = scriptDef.ScriptSwitchEnd()
			case fileio.OP_GOTO:
				returnValue = scriptDef.ScriptGoto(lineData)
			case fileio.OP_GOSUB:
				returnValue = scriptDef.ScriptGoSub(lineData, scriptData)
			case fileio.OP_BREAK:
				returnValue = scriptDef.ScriptBreak(lineData)
			case fileio.OP_CHECK:
				returnValue = scriptDef.ScriptCheckBit(lineData, gameDef)
			case fileio.OP_SET_BIT:
				returnValue = scriptDef.ScriptSetBit(lineData, gameDef)
			case fileio.OP_COMPARE:
				returnValue = scriptDef.ScriptCompare(lineData)
			case fileio.OP_CUT_CHG:
				returnValue = scriptDef.ScriptCameraChange(lineData, gameDef)
			case fileio.OP_AOT_SET:
				returnValue = scriptDef.ScriptAotSet(lineData, gameDef)
			case fileio.OP_OBJ_MODEL_SET:
				returnValue = scriptDef.ScriptObjectModelSet(lineData)
			case fileio.OP_WORK_SET:
				returnValue = scriptDef.ScriptWorkSet(lineData)
			case fileio.OP_POS_SET:
				returnValue = scriptDef.ScriptPositionSet(lineData, gameDef)
			case fileio.OP_SCA_ID_SET:
				returnValue = scriptDef.ScriptScaIdSet(lineData, gameDef)
			case fileio.OP_SCE_ESPR_ON:
				scriptDef.ScriptSpriteOn(lineData, gameDef)
			case fileio.OP_DOOR_AOT_SET:
				returnValue = scriptDef.ScriptDoorAotSet(lineData, gameDef)
			case fileio.OP_MEMBER_CMP:
				scriptDef.ScriptMemberCompare(lineData)
			case fileio.OP_PLC_MOTION: // 0x3f
				returnValue = scriptDef.ScriptPlcMotion(lineData)
			case fileio.OP_PLC_DEST: // 0x40
				returnValue = scriptDef.ScriptPlcDest(lineData)
			case fileio.OP_PLC_NECK: // 0x41
				returnValue = scriptDef.ScriptPlcNeck(lineData)
			case fileio.OP_SCE_EM_SET: // 0x44
				returnValue = scriptDef.ScriptSceEmSet(lineData)
			case fileio.OP_AOT_RESET:
				returnValue = scriptDef.ScriptAotReset(lineData, gameDef)
			case fileio.OP_ITEM_AOT_SET:
				returnValue = scriptDef.ScriptItemAotSet(lineData, gameDef)
			case fileio.OP_AOT_SET_4P:
				returnValue = scriptDef.ScriptAotSet4p(lineData)
			default:
				returnValue = 1
			}

			if !scriptThread.OverrideProgramCounter {
				scriptThread.IncrementProgramCounter(opcode)
			}
			scriptThread.OverrideProgramCounter = false

			// Control flow is broken
			if returnValue != 1 {
				scriptReturnValue = returnValue
				break
			}
		}

		// End thread
		if scriptReturnValue == 2 || scriptThread.LevelState[scriptThread.SubLevel].IfElseCounter < 0 {
			break
		}

		if scriptThread.StackIndex == 0 {
			log.Fatal("Script stack is empty")
		}

		// pop stack
		scriptThread.StackIndex--
		stackTop := scriptThread.Stack[scriptThread.StackIndex]
		scriptThread.ProgramCounter = stackTop
		scriptThread.LevelState[scriptThread.SubLevel].IfElseCounter--
	}
}

func (scriptDef *ScriptDef) ScriptEvtEnd(lineData []byte) int {
	// The program is returning from a subroutine
	if scriptThread.SubLevel != 0 {
		ifElseCounter := scriptThread.LevelState[scriptThread.SubLevel].IfElseCounter
		scriptThread.SubLevel--
		scriptThread.ProgramCounter = scriptThread.LevelState[scriptThread.SubLevel].ReturnAddress
		scriptThread.OverrideProgramCounter = true
		scriptThread.StackIndex = (8 * scriptThread.SubLevel) + ifElseCounter + 1
		return 1
	}

	// The program is in the top level
	scriptThread.RunStatus = false
	return 2
}

func (scriptDef *ScriptDef) ScriptEvtExec(lineData []byte, scriptData fileio.ScriptFunction) int {
	byteArr := bytes.NewBuffer(lineData)
	instruction := fileio.ScriptInstrEventExec{}
	binary.Read(byteArr, binary.LittleEndian, &instruction)

	nextThreadNum := 0

	if int(instruction.ThreadNum) >= 0 && int(instruction.ThreadNum) < len(scriptDef.ScriptThreads) {
		// thread num is defined
		nextThreadNum = int(instruction.ThreadNum)
	} else {
		// assign next available thread
		for i := 0; i < len(scriptDef.ScriptThreads); i++ {
			if scriptDef.ScriptThreads[i].RunStatus == false {
				nextThreadNum = i
				break
			}
		}
	}

	scriptDef.ScriptThreads[nextThreadNum].RunStatus = true
	scriptDef.ScriptThreads[nextThreadNum].ProgramCounter = scriptData.StartProgramCounter[instruction.Event]
	scriptDef.ScriptThreads[nextThreadNum].LevelState[0].IfElseCounter = -1
	scriptDef.ScriptThreads[nextThreadNum].LevelState[0].LoopCounter = -1
	return 1
}

func (scriptDef *ScriptDef) ScriptIfBlockStart(lineData []byte) int {
	byteArr := bytes.NewBuffer(lineData)
	conditional := fileio.ScriptInstrIfElseStart{}
	binary.Read(byteArr, binary.LittleEndian, &conditional)

	opcode := lineData[0]
	scriptThread.LevelState[scriptThread.SubLevel].IfElseCounter++
	newPosition := (scriptThread.ProgramCounter + fileio.InstructionSize[opcode]) + int(conditional.BlockLength)
	scriptThread.Stack[scriptThread.StackIndex] = newPosition
	scriptThread.StackIndex++

	return 1
}

func (scriptDef *ScriptDef) ScriptElseCheck(lineData []byte) int {
	byteArr := bytes.NewBuffer(lineData)
	conditional := fileio.ScriptInstrElseStart{}
	binary.Read(byteArr, binary.LittleEndian, &conditional)

	scriptThread.StackIndex--
	scriptThread.ProgramCounter = scriptThread.ProgramCounter + int(conditional.BlockLength)
	scriptThread.LevelState[scriptThread.SubLevel].IfElseCounter--
	scriptThread.OverrideProgramCounter = true
	return 1
}

func (scriptDef *ScriptDef) ScriptEndIf() int {
	scriptThread.StackIndex--
	scriptThread.LevelState[scriptThread.SubLevel].IfElseCounter--
	return 1
}

func (scriptDef *ScriptDef) ScriptSleep(lineData []byte) int {
	byteArr := bytes.NewBuffer(lineData)
	instruction := fileio.ScriptInstrSleep{}
	binary.Read(byteArr, binary.LittleEndian, &instruction)

	// goes to sleeping instruction (0xa)
	scriptThread.ProgramCounter = scriptThread.ProgramCounter + 1
	scriptThread.OverrideProgramCounter = true
	scriptThread.LevelState[scriptThread.SubLevel].LoopCounter++
	sleepCounterIndex := scriptThread.LevelState[scriptThread.SubLevel].LoopCounter
	scriptThread.LevelState[scriptThread.SubLevel].SleepCounter[sleepCounterIndex] = int(instruction.Count)

	return 1
}

func (scriptDef *ScriptDef) ScriptSleeping(lineData []byte) int {
	opcode := lineData[0]
	sleepCounterIndex := scriptThread.LevelState[scriptThread.SubLevel].LoopCounter
	scriptThread.LevelState[scriptThread.SubLevel].SleepCounter[sleepCounterIndex]--

	if scriptThread.LevelState[scriptThread.SubLevel].SleepCounter[sleepCounterIndex] == 0 {
		scriptThread.ProgramCounter += fileio.InstructionSize[opcode]
		scriptThread.LevelState[scriptThread.SubLevel].LoopCounter--
	}

	scriptThread.OverrideProgramCounter = true

	return 2
}

func (scriptDef *ScriptDef) ScriptSwitchBegin(
	lineData []byte,
	instructions map[int][]byte,
	gameDef *game.GameDef) int {
	byteArr := bytes.NewBuffer(lineData)
	switchConditional := fileio.ScriptInstrSwitch{}
	binary.Read(byteArr, binary.LittleEndian, &switchConditional)

	opcode := lineData[0]
	scriptThread.LevelState[scriptThread.SubLevel].LoopCounter++
	newLoopCounter := scriptThread.LevelState[scriptThread.SubLevel].LoopCounter
	newProgramCounter := scriptThread.ProgramCounter + fileio.InstructionSize[opcode]
	scriptThread.LevelState[scriptThread.SubLevel].LoopBreak[newLoopCounter] = newProgramCounter + int(switchConditional.BlockLength)
	scriptThread.LevelState[scriptThread.SubLevel].LoopIfCounter[newLoopCounter] = scriptThread.LevelState[scriptThread.SubLevel].IfElseCounter

	for true {
		newLineData := instructions[newProgramCounter]
		newOpcode := newLineData[0]

		if newOpcode == fileio.OP_CASE {
			byteArr = bytes.NewBuffer(newLineData)
			caseInstruction := fileio.ScriptInstrSwitchCase{}
			binary.Read(byteArr, binary.LittleEndian, &caseInstruction)

			switchValue := gameDef.GetScriptVariable(int(switchConditional.VarId))
			// Case matches
			if int(caseInstruction.Value) == switchValue {
				scriptThread.ProgramCounter = newProgramCounter + fileio.InstructionSize[newOpcode]
				scriptThread.OverrideProgramCounter = true
				return 1
			} else {
				// Move to the next case statement
				newProgramCounter += fileio.InstructionSize[newOpcode] + int(caseInstruction.BlockLength)
			}
		} else if newOpcode == fileio.OP_DEFAULT {
			scriptThread.ProgramCounter = newProgramCounter + fileio.InstructionSize[newOpcode]
			scriptThread.OverrideProgramCounter = true
			return 1
		} else if newOpcode == fileio.OP_END_SWITCH {
			scriptThread.LevelState[scriptThread.SubLevel].LoopCounter--
			scriptThread.ProgramCounter = newProgramCounter + fileio.InstructionSize[newOpcode]
			scriptThread.OverrideProgramCounter = true
			return 1
		} else {
			log.Fatal("Switch statement has unknown opcode: ", newOpcode)
		}
	}

	return 1
}

func (scriptDef *ScriptDef) ScriptSwitchEnd() int {
	scriptThread.LevelState[scriptThread.SubLevel].LoopCounter--
	return 1
}

func (scriptDef *ScriptDef) ScriptGoto(lineData []byte) int {
	byteArr := bytes.NewBuffer(lineData)
	instruction := fileio.ScriptInstrGoto{}
	binary.Read(byteArr, binary.LittleEndian, &instruction)

	// Disable due to infinite loop
	/*scriptThread.LevelState[scriptThread.SubLevel].IfElseCounter = int(instruction.IfElseCounter)
	scriptThread.StackIndex = (8 * scriptThread.SubLevel) + int(instruction.IfElseCounter) + 1
	scriptThread.LevelState[scriptThread.SubLevel].LoopCounter = int(instruction.LoopCounter)
	scriptThread.ProgramCounter += int(instruction.Offset)
	scriptThread.OverrideProgramCounter = true*/

	return 1
}

func (scriptDef *ScriptDef) ScriptGoSub(lineData []byte, scriptData fileio.ScriptFunction) int {
	byteArr := bytes.NewBuffer(lineData)
	instruction := fileio.ScriptInstrGoSub{}
	binary.Read(byteArr, binary.LittleEndian, &instruction)

	opcode := lineData[0]
	scriptThread.LevelState[scriptThread.SubLevel].ReturnAddress = scriptThread.ProgramCounter + fileio.InstructionSize[opcode]
	scriptThread.LevelState[scriptThread.SubLevel+1].IfElseCounter = -1
	scriptThread.LevelState[scriptThread.SubLevel+1].LoopCounter = -1
	scriptThread.StackIndex = 8 * (scriptThread.SubLevel + 1)
	scriptThread.SubLevel++

	scriptThread.ProgramCounter = scriptData.StartProgramCounter[instruction.Event]
	scriptThread.OverrideProgramCounter = true
	return 1
}

func (scriptDef *ScriptDef) ScriptBreak(lineData []byte) int {
	loopCounter := scriptThread.LevelState[scriptThread.SubLevel].LoopCounter
	scriptThread.OverrideProgramCounter = true
	scriptThread.ProgramCounter = scriptThread.LevelState[scriptThread.SubLevel].LoopBreak[loopCounter]
	scriptThread.LevelState[scriptThread.SubLevel].IfElseCounter = scriptThread.LevelState[scriptThread.SubLevel].LoopIfCounter[loopCounter]
	scriptThread.LevelState[scriptThread.SubLevel].LoopCounter--
	return 1
}

func (scriptDef *ScriptDef) ScriptCheckBit(lineData []byte, gameDef *game.GameDef) int {
	byteArr := bytes.NewBuffer(lineData)
	bitTest := fileio.ScriptInstrCheckBitTest{}
	binary.Read(byteArr, binary.LittleEndian, &bitTest)

	if gameDef.GetBitArray(int(bitTest.BitArray), int(bitTest.Number)) == int(bitTest.Value) {
		return 1
	}
	return 0
}

func (scriptDef *ScriptDef) ScriptSetBit(lineData []byte, gameDef *game.GameDef) int {
	byteArr := bytes.NewBuffer(lineData)
	instruction := fileio.ScriptInstrSetBit{}
	binary.Read(byteArr, binary.LittleEndian, &instruction)

	if instruction.Operation == 0 {
		// Clear bit
		gameDef.SetBitArray(int(instruction.BitArray), int(instruction.BitNumber), 0)
	} else if instruction.Operation == 1 {
		// Set bit
		gameDef.SetBitArray(int(instruction.BitArray), int(instruction.BitNumber), 1)
	} else if instruction.Operation == 7 {
		// Flip bit
		currentBit := gameDef.GetBitArray(int(instruction.BitArray), int(instruction.BitNumber))
		gameDef.SetBitArray(int(instruction.BitArray), int(instruction.BitNumber), currentBit^1)
	} else {
		log.Fatal("Set bit operation ", instruction.Operation, " is invalid.")
	}

	return 1
}

func (scriptDef *ScriptDef) ScriptCompare(lineData []byte) int {
	byteArr := bytes.NewBuffer(lineData)
	instruction := fileio.ScriptInstrMemberCompare{}
	binary.Read(byteArr, binary.LittleEndian, &instruction)

	// TODO: Evaluate to true or false
	// Assumes the statement is false

	return 1
}

func (scriptDef *ScriptDef) ScriptCameraChange(lineData []byte, gameDef *game.GameDef) int {
	byteArr := bytes.NewBuffer(lineData)
	instruction := fileio.ScriptInstrCutChg{}
	binary.Read(byteArr, binary.LittleEndian, &instruction)

	gameDef.ChangeCamera(int(instruction.CameraId))
	return 1
}

func (scriptDef *ScriptDef) ScriptAotSet(lineData []byte, gameDef *game.GameDef) int {
	byteArr := bytes.NewBuffer(lineData)
	instruction := fileio.ScriptInstrAotSet{}
	binary.Read(byteArr, binary.LittleEndian, &instruction)

	gameDef.AotManager.AddAotTrigger(instruction)
	return 1
}

func (scriptDef *ScriptDef) ScriptObjectModelSet(lineData []byte) int {
	byteArr := bytes.NewBuffer(lineData)
	instruction := fileio.ScriptInstrObjModelSet{}
	binary.Read(byteArr, binary.LittleEndian, &instruction)

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

func (scriptDef *ScriptDef) ScriptScaIdSet(lineData []byte, gameDef *game.GameDef) int {
	byteArr := bytes.NewBuffer(lineData)
	instruction := fileio.ScriptInstrScaIdSet{}
	binary.Read(byteArr, binary.LittleEndian, &instruction)

	if instruction.Flag == 0 {
		gameDef.RemoveCollisionEntity(gameDef.GameRoom.CollisionEntities, int(instruction.Id))
	}
	return 1
}

func (scriptDef *ScriptDef) ScriptSpriteOn(lineData []byte, gameDef *game.GameDef) {
	byteArr := bytes.NewBuffer(lineData)
	scriptSprite := fileio.ScriptSprite{}
	binary.Read(byteArr, binary.LittleEndian, &scriptSprite)

	gameDef.AotManager.AddScriptSprite(scriptSprite)
}

func (scriptDef *ScriptDef) ScriptDoorAotSet(lineData []byte, gameDef *game.GameDef) int {
	byteArr := bytes.NewBuffer(lineData)
	door := fileio.ScriptInstrDoorAotSet{}
	err := binary.Read(byteArr, binary.LittleEndian, &door)
	if err != nil {
		log.Fatal("Error loading door")
	}

	if door.Id != game.AOT_DOOR {
		log.Fatal("Door has incorrect aot type ", door.Id)
	}

	gameDef.AotManager.AddDoorAot(door)
	return 1
}

func (scriptDef *ScriptDef) ScriptMemberCompare(lineData []byte) {
	byteArr := bytes.NewBuffer(lineData)
	instruction := fileio.ScriptInstrMemberCompare{}
	binary.Read(byteArr, binary.LittleEndian, &instruction)

	// TODO: Evaluate to true or false
	// Assumes the statement is false
}

func (scriptDef *ScriptDef) ScriptPlcMotion(lineData []byte) int {
	byteArr := bytes.NewBuffer(lineData)
	instruction := fileio.ScriptInstrPlcMotion{}
	binary.Read(byteArr, binary.LittleEndian, &instruction)

	return 1
}

func (scriptDef *ScriptDef) ScriptPlcDest(lineData []byte) int {
	byteArr := bytes.NewBuffer(lineData)
	instruction := fileio.ScriptInstrPlcDest{}
	binary.Read(byteArr, binary.LittleEndian, &instruction)

	return 1
}

func (scriptDef *ScriptDef) ScriptPlcNeck(lineData []byte) int {
	byteArr := bytes.NewBuffer(lineData)
	instruction := fileio.ScriptInstrPlcNeck{}
	binary.Read(byteArr, binary.LittleEndian, &instruction)

	return 1
}

func (scriptDef *ScriptDef) ScriptSceEmSet(lineData []byte) int {
	return 1
}

func (scriptDef *ScriptDef) ScriptAotReset(lineData []byte, gameDef *game.GameDef) int {
	byteArr := bytes.NewBuffer(lineData)
	instruction := fileio.ScriptInstrAotReset{}
	binary.Read(byteArr, binary.LittleEndian, &instruction)

	gameDef.AotManager.ResetAotTrigger(instruction)
	return 1
}

func (scriptDef *ScriptDef) ScriptItemAotSet(lineData []byte, gameDef *game.GameDef) int {
	byteArr := bytes.NewBuffer(lineData)
	item := fileio.ScriptInstrItemAotSet{}
	binary.Read(byteArr, binary.LittleEndian, &item)

	if item.Id != game.AOT_ITEM {
		log.Fatal("Item has incorrect aot type ", item.Id)
	}

	gameDef.AotManager.AddItemAot(item)
	return 1
}

func (scriptDef *ScriptDef) ScriptAotSet4p(lineData []byte) int {
	return 1
}
