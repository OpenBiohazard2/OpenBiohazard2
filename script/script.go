package script

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"

	"github.com/OpenBiohazard2/OpenBiohazard2/fileio"
	"github.com/OpenBiohazard2/OpenBiohazard2/game"
	"github.com/OpenBiohazard2/OpenBiohazard2/render"
	"github.com/OpenBiohazard2/OpenBiohazard2/world"
	"github.com/go-gl/mathgl/mgl32"
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
	scriptDeltaTime = 0.0
)

type ScriptDef struct {
	ScriptThreads  []*ScriptThread
	ScriptBitArray map[int]map[int]int
	ScriptVariable map[int]int
	DebugEnabled   bool
}

func NewScriptDef() *ScriptDef {
	scriptThreads := make([]*ScriptThread, 20)
	for i := 0; i < len(scriptThreads); i++ {
		scriptThreads[i] = NewScriptThread()
	}

	return &ScriptDef{
		ScriptThreads:  scriptThreads,
		ScriptBitArray: make(map[int]map[int]int),
		ScriptVariable: make(map[int]int),
		DebugEnabled:   false,
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
	scriptDef.ScriptDebugLine(fmt.Sprintf("Initialize script thread %v, start function %v", threadNum, startFunction))
	scriptDef.ScriptThreads[threadNum].RunStatus = true
	scriptDef.ScriptThreads[threadNum].ProgramCounter = scriptData.StartProgramCounter[startFunction]
	scriptDef.ScriptThreads[threadNum].FunctionIds = []int{startFunction}
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

		scriptDef.RunScriptThread(i, scriptDef.ScriptThreads[i], scriptData, gameDef, renderDef)
	}
}

func (scriptDef *ScriptDef) RunScriptThread(
	threadNum int,
	curScriptThread *ScriptThread,
	scriptData fileio.ScriptFunction,
	gameDef *game.GameDef,
	renderDef *render.RenderDef) {

	// Thread should not run
	if curScriptThread.RunStatus == false {
		return
	}

	for true {
		sectionReturnValue := scriptDef.RunScriptUntilBreakControlFlow(threadNum, curScriptThread, scriptData, gameDef, renderDef)

		// End thread
		if curScriptThread.ShouldTerminate(sectionReturnValue) {
			break
		}

		curScriptThread.JumpToNextLocationOnStack()
	}
}

func (scriptDef *ScriptDef) RunScriptUntilBreakControlFlow(
	threadNum int,
	curScriptThread *ScriptThread,
	scriptData fileio.ScriptFunction,
	gameDef *game.GameDef,
	renderDef *render.RenderDef) int {
	scriptReturnValue := 0
	for true {
		lineData := scriptData.Instructions[curScriptThread.ProgramCounter]
		if len(lineData) == 0 {
			curScriptThread.RunStatus = false
			log.Print("Warning: terminate script thread at program counter ", curScriptThread.ProgramCounter)
			break
		}

		opcode := lineData[0]

		// Override can be modified during execution
		curScriptThread.OverrideProgramCounter = false

		instructionReturnValue := scriptDef.ExecuteSingleInstruction(threadNum, curScriptThread, lineData, scriptData, gameDef, renderDef)

		if !curScriptThread.OverrideProgramCounter {
			curScriptThread.IncrementProgramCounter(opcode)
		}
		curScriptThread.OverrideProgramCounter = false

		// Control flow is broken
		if instructionReturnValue != INSTRUCTION_NORMAL {
			scriptReturnValue = instructionReturnValue
			break
		}
	}
	return scriptReturnValue
}

func (scriptDef *ScriptDef) ExecuteSingleInstruction(
	threadNum int,
	curScriptThread *ScriptThread,
	lineData []byte,
	scriptData fileio.ScriptFunction,
	gameDef *game.GameDef,
	renderDef *render.RenderDef) int {
	var returnValue int

	opcode := lineData[0]

	scriptDef.ScriptDebugFunction(threadNum, curScriptThread.FunctionIds, lineData)

	switch opcode {
	case fileio.OP_EVT_END:
		returnValue = scriptDef.ScriptEvtEnd(curScriptThread, lineData, threadNum)
	case fileio.OP_EVT_EXEC:
		returnValue = scriptDef.ScriptEvtExec(lineData, scriptData)
	case fileio.OP_IF_START:
		returnValue = scriptDef.ScriptIfBlockStart(curScriptThread, lineData)
	case fileio.OP_ELSE_START:
		returnValue = scriptDef.ScriptElseCheck(curScriptThread, lineData)
	case fileio.OP_END_IF:
		returnValue = scriptDef.ScriptEndIf(curScriptThread)
	case fileio.OP_SLEEP:
		returnValue = scriptDef.ScriptSleep(curScriptThread, lineData)
	case fileio.OP_SLEEPING:
		returnValue = scriptDef.ScriptSleeping(curScriptThread, lineData)
	case fileio.OP_FOR:
		returnValue = scriptDef.ScriptForLoopBegin(curScriptThread, lineData)
	case fileio.OP_FOR_END:
		returnValue = scriptDef.ScriptForLoopEnd(curScriptThread, lineData)
	case fileio.OP_SWITCH:
		returnValue = scriptDef.ScriptSwitchBegin(curScriptThread, lineData, scriptData.Instructions)
	case fileio.OP_CASE:
		returnValue = 1 // already implemented in switch statement
	case fileio.OP_DEFAULT:
		returnValue = 1 // already implemented in switch statement
	case fileio.OP_END_SWITCH:
		returnValue = scriptDef.ScriptSwitchEnd(curScriptThread)
	case fileio.OP_GOTO:
		returnValue = scriptDef.ScriptGoto(curScriptThread, lineData)
	case fileio.OP_GOSUB:
		returnValue = scriptDef.ScriptGoSub(curScriptThread, lineData, scriptData)
	case fileio.OP_BREAK:
		returnValue = scriptDef.ScriptBreak(curScriptThread, lineData)
	case fileio.OP_CHECK: // 0x21
		returnValue = scriptDef.ScriptCheckBit(lineData)
	case fileio.OP_SET_BIT: // 0x22
		returnValue = scriptDef.ScriptSetBit(lineData)
	case fileio.OP_COMPARE: // 0x23
		returnValue = scriptDef.ScriptCompare(lineData)
	case fileio.OP_SAVE: // 0x24
		returnValue = scriptDef.ScriptSave(lineData)
	case fileio.OP_COPY: // 0x25
		returnValue = scriptDef.ScriptCopy(lineData)
	case fileio.OP_CALC: // 0x26
		returnValue = scriptDef.ScriptCalc(lineData)
	case fileio.OP_CALC2: // 0x27
		returnValue = scriptDef.ScriptCalc(lineData)
	case fileio.OP_CUT_CHG:
		returnValue = scriptDef.ScriptCameraChange(lineData, gameDef)
	case fileio.OP_AOT_SET:
		returnValue = scriptDef.ScriptAotSet(lineData, gameDef)
	case fileio.OP_OBJ_MODEL_SET:
		returnValue = scriptDef.ScriptObjectModelSet(lineData, renderDef)
	case fileio.OP_WORK_SET:
		returnValue = scriptDef.ScriptWorkSet(curScriptThread, lineData)
	case fileio.OP_POS_SET:
		returnValue = scriptDef.ScriptPositionSet(curScriptThread, lineData, gameDef)
	case fileio.OP_MEMBER_SET:
		returnValue = scriptDef.ScriptMemberSet(curScriptThread, lineData, gameDef, renderDef)
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
		returnValue = scriptDef.ScriptSceEmSet(lineData, renderDef)
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

func (scriptDef *ScriptDef) ScriptSleep(thread *ScriptThread, lineData []byte) int {
	byteArr := bytes.NewBuffer(lineData)
	instruction := fileio.ScriptInstrSleep{}
	binary.Read(byteArr, binary.LittleEndian, &instruction)

	// goes to sleeping instruction (0xa)
	curLevelState := thread.LevelState[thread.SubLevel]

	thread.ProgramCounter = thread.ProgramCounter + 1
	thread.OverrideProgramCounter = true

	curLevelState.LoopLevel++
	newLoopLevel := curLevelState.LoopLevel
	curLevelState.LoopState[newLoopLevel].Counter = int(instruction.Count)
	return 1
}

func (scriptDef *ScriptDef) ScriptSleeping(thread *ScriptThread, lineData []byte) int {
	opcode := lineData[0]
	curLevelState := thread.LevelState[thread.SubLevel]
	curLoopState := curLevelState.LoopState[curLevelState.LoopLevel]

	curLoopState.Counter--
	if curLoopState.Counter == 0 {
		thread.ProgramCounter += fileio.InstructionSize[opcode]
		curLevelState.LoopLevel--
	}

	thread.OverrideProgramCounter = true

	return INSTRUCTION_THREAD_END
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

func (scriptDef *ScriptDef) ScriptWorkSet(thread *ScriptThread, lineData []byte) int {
	byteArr := bytes.NewBuffer(lineData)
	instruction := fileio.ScriptInstrWorkSet{}
	binary.Read(byteArr, binary.LittleEndian, &instruction)

	thread.WorkSetComponent = int(instruction.Component)
	thread.WorkSetIndex = int(instruction.Index)
	return 1
}

func (scriptDef *ScriptDef) ScriptPositionSet(thread *ScriptThread, lineData []byte, gameDef *game.GameDef) int {
	byteArr := bytes.NewBuffer(lineData)
	instruction := fileio.ScriptInstrPosSet{}
	binary.Read(byteArr, binary.LittleEndian, &instruction)

	if thread.WorkSetComponent == WORKSET_PLAYER {
		gameDef.Player.Position = mgl32.Vec3{float32(instruction.X), float32(instruction.Y), float32(instruction.Z)}
	} else {
		// TODO: set position of object
	}

	return 1
}

func (scriptDef *ScriptDef) ScriptMemberSet(thread *ScriptThread, lineData []byte, gameDef *game.GameDef, renderDef *render.RenderDef) int {
	byteArr := bytes.NewBuffer(lineData)
	instruction := fileio.ScriptInstrMemberSet{}
	binary.Read(byteArr, binary.LittleEndian, &instruction)

	if thread.WorkSetComponent == WORKSET_PLAYER {
		switch int(instruction.MemberIndex) {
		case 15:
			// convert to angle in degrees
			gameDef.Player.RotationAngle = (float32(instruction.Value) / 4096.0) * 360.0
		}
	} else if thread.WorkSetComponent == WORKSET_OBJECT {
		modelObject := renderDef.SceneSystem.ItemGroupEntity.ModelObjectData[int(thread.WorkSetIndex)]
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
		world.RemoveCollisionEntity(gameDef.GameWorld.GameRoom.CollisionEntities, int(instruction.Id))
	}
	return 1
}

func (scriptDef *ScriptDef) ScriptMemberCompare(lineData []byte) int {
	byteArr := bytes.NewBuffer(lineData)
	instruction := fileio.ScriptInstrMemberCompare{}
	binary.Read(byteArr, binary.LittleEndian, &instruction)

	return 1
}

func (scriptDef *ScriptDef) ScriptSceEmSet(lineData []byte, renderDef *render.RenderDef) int {
	byteArr := bytes.NewBuffer(lineData)
	instruction := fileio.ScriptInstrSceEmSet{}
	binary.Read(byteArr, binary.LittleEndian, &instruction)

	// Create enemy entity if we have valid data
	if instruction.Type != 0 && instruction.ModelType != 0 {
		// Load the EMD file based on the enemy type (3-digit hexadecimal)
		enemyEMDPath := fmt.Sprintf("data/PL0/EMD0/EM%03X.EMD", instruction.Type)
		
		// Load the enemy model data
		emdOutput := fileio.LoadEMDFile(enemyEMDPath)
		if emdOutput != nil {
			// Create enemy entity
			enemyEntity := render.NewEnemyEntity(emdOutput)
			enemyEntity.SetEnemyData(instruction)
			
			// Add to scene system
			renderDef.SceneSystem.EnemyGroupEntity.AddEnemy(enemyEntity)
			
			// Log enemy creation since there won't be too many enemies
			fmt.Printf("Created enemy type 0x%03X at position (%d, %d, %d)\n", 
				instruction.Type, instruction.X, instruction.Y, instruction.Z)
		} else {
			// Only log failures for debugging purposes
			fmt.Printf("Failed to load enemy model for type 0x%03X\n", instruction.Type)
		}
	}

	return 1
}

func (scriptDef *ScriptDef) ScriptSceBgmControl(lineData []byte) int {
	byteArr := bytes.NewBuffer(lineData)
	instruction := fileio.ScriptInstrSceBgmControl{}
	binary.Read(byteArr, binary.LittleEndian, &instruction)

	return 1
}
