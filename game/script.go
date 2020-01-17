package game

import (
	"../fileio"
	"bytes"
	"encoding/binary"
	"log"
	"sort"
)

type ScriptMemory struct {
	ProgramCounter         int
	FunctionNum            int
	Stack                  [][]int
	TotalTime              float64
	OverrideProgramCounter bool
	EndScript              bool
}

type ScriptInstrEventExec struct {
	Opcode   uint8 // 0x04
	Cond     uint8
	ExOpcode uint8
	Event    uint8
}

type ScriptInstrIfElseStart struct {
	Opcode      uint8 // 0x06
	Dummy       uint8
	BlockLength uint16
}

type ScriptInstrElseStart struct {
	Opcode      uint8 // 0x07
	Dummy       uint8
	BlockLength uint16
}

type ScriptInstrSwitch struct {
	Opcode      uint8 // 0x13
	VarId       uint8
	BlockLength uint16
}

type ScriptInstrSwitchCase struct {
	Opcode      uint8 // 0x14
	Dummy       uint8
	BlockLength uint16
	Value       uint16
}

type ScriptInstrGoto struct {
	Opcode        uint8 // 0x17
	IfElseCounter uint8
	LoopCounter   uint8
	Unknown0      uint8
	Offset        uint8
	Unknown1      uint8
}

type ScriptInstrGoSub struct {
	Opcode uint8 // 0x18
	Event  uint8
}

type ScriptInstrCheckBitTest struct {
	Opcode   uint8 // 0x21
	BitArray uint8 // Index of array of bits to use
	Number   uint8 // Bit number to check
	Value    uint8 // Value to compare (0 or 1)
}

type ScriptSprite struct {
	Opcode   uint8
	Dummy    uint8
	SpriteId uint8
	Unknown0 [3]uint8
	Dir      int16
	X, Y, Z  int16
	Unknown1 uint16
}

type ScriptDoor struct {
	Opcode                       uint8 // 0x3b
	Id                           uint8 // Index of item in array of room objects list
	Unknown0                     [2]uint16
	X, Y                         int16 // Location of door
	Width, Height                int16 // Size of door
	NextX, NextY, NextZ, NextDir int16 // Position and direction of player after door entered
	Stage, Room, Camera          uint8 // Stage, room, camera after door entered
	Unknown1                     uint8
	DoorType                     uint8
	DoorLock                     uint8
	Unknown2                     uint8
	DoorLocked                   uint8
	DoorKey                      uint8
	Unknown3                     uint8
}

type ScriptInstrMemberCompare struct {
	Opcode   uint8 // 0x3e
	Unknown0 uint8
	Unknown1 uint8
	Compare  uint8
	Value    int16
}

func NewScriptMemory() *ScriptMemory {
	return &ScriptMemory{
		ProgramCounter:         0,
		FunctionNum:            0,
		Stack:                  make([][]int, 0),
		TotalTime:              60,
		OverrideProgramCounter: false,
		EndScript:              false,
	}
}

func (mem *ScriptMemory) IncrementProgramCounter(opcode byte) {
	mem.ProgramCounter += fileio.InstructionSize[opcode]
}

func (mem *ScriptMemory) PushStack(values []int) {
	mem.Stack = append(mem.Stack, values)
}

func (mem *ScriptMemory) PopStack() []int {
	if len(mem.Stack) == 0 {
		log.Fatal("Script stack is empty")
	}

	top := mem.Stack[len(mem.Stack)-1]
	// remove last element
	mem.Stack = mem.Stack[0 : len(mem.Stack)-1]
	return top
}

func (mem *ScriptMemory) PeekStackTop() []int {
	if len(mem.Stack) == 0 {
		log.Fatal("Script stack is empty")
	}

	return mem.Stack[len(mem.Stack)-1]
}

func (gameDef *GameDef) RunScript(scriptData fileio.ScriptFunction, timeElapsedSeconds float64, init bool, startFunction int) {
	if init == false {
		gameDef.ScriptMemory.TotalTime += timeElapsedSeconds
		if gameDef.ScriptMemory.TotalTime < 60 {
			return
		}
	}

	gameDef.ScriptMemory = NewScriptMemory()
	gameDef.ScriptMemory.TotalTime = 0

	programCounters := make([]int, 0)
	for counter, _ := range scriptData.Instructions {
		programCounters = append(programCounters, counter)
	}
	sort.Ints(programCounters)

	gameDef.ScriptMemory.ProgramCounter = scriptData.StartProgramCounter[startFunction]
	gameDef.ScriptMemory.FunctionNum = startFunction

	lastProgramCounter := programCounters[len(programCounters)-1]
	for gameDef.ScriptMemory.ProgramCounter <= lastProgramCounter {
		lineData := scriptData.Instructions[gameDef.ScriptMemory.ProgramCounter]
		opcode := lineData[0]

		switch opcode {
		case fileio.OP_RETURN:
			gameDef.ScriptReturn(startFunction)
		case fileio.OP_EVT_EXEC:
			gameDef.ScriptEvtExec(lineData, scriptData)
		case fileio.OP_IF_START:
			gameDef.ScriptIfBlockStart(lineData, scriptData)
		case fileio.OP_ELSE_START:
			gameDef.ScriptElseCheck(lineData)
		case fileio.OP_END_IF:
			gameDef.ScriptEndIf()
		case fileio.OP_SWITCH:
			gameDef.ScriptSwitchBegin(lineData)
		case fileio.OP_CASE:
			gameDef.ScriptSwitchCase(lineData)
		case fileio.OP_END_SWITCH:
			gameDef.ScriptSwitchEnd()
		case fileio.OP_GOSUB:
			gameDef.ScriptGoSub(lineData, scriptData)
		case fileio.OP_CHECK:
			gameDef.ScriptCheckBit(lineData)
		case fileio.OP_SET_BIT:
			gameDef.ScriptSetBit(lineData)
		case fileio.OP_COMPARE:
			gameDef.ScriptCompare(lineData)
		case fileio.OP_SCE_ESPR_ON:
			gameDef.ScriptSpriteOn(lineData)
		case fileio.OP_DOOR_AOT_SET:
			gameDef.ScriptDoorAotSet(lineData)
		case fileio.OP_MEMBER_CMP:
			gameDef.ScriptMemberCompare(lineData)
		}

		if !gameDef.ScriptMemory.OverrideProgramCounter {
			gameDef.ScriptMemory.IncrementProgramCounter(opcode)
		}
		gameDef.ScriptMemory.OverrideProgramCounter = false

		if gameDef.ScriptMemory.EndScript == true {
			break
		}
	}

}

func (gameDef *GameDef) ScriptReturn(startFunction int) {
	// End script
	if gameDef.ScriptMemory.FunctionNum == startFunction {
		gameDef.ScriptMemory.EndScript = true
		return
	}

	if len(gameDef.ScriptMemory.Stack) > 0 {
		stackTop := gameDef.ScriptMemory.PopStack()
		gameDef.ScriptMemory.FunctionNum = stackTop[2]
		gameDef.ScriptMemory.ProgramCounter = stackTop[1]
		gameDef.ScriptMemory.OverrideProgramCounter = true
	}
}

func (gameDef *GameDef) ScriptEvtExec(lineData []byte, scriptData fileio.ScriptFunction) {
	byteArr := bytes.NewBuffer(lineData)
	instruction := ScriptInstrEventExec{}
	binary.Read(byteArr, binary.LittleEndian, &instruction)

	// Call function
	// Save program counter to be after the function
	opcode := lineData[0]
	gameDef.ScriptMemory.IncrementProgramCounter(opcode)
	values := []int{int(opcode), gameDef.ScriptMemory.ProgramCounter, gameDef.ScriptMemory.FunctionNum}
	gameDef.ScriptMemory.PushStack(values)
	gameDef.ScriptMemory.FunctionNum = int(instruction.Event)
	gameDef.ScriptMemory.ProgramCounter = scriptData.StartProgramCounter[gameDef.ScriptMemory.FunctionNum]
	gameDef.ScriptMemory.OverrideProgramCounter = true
}

// Ends with else statement or end if statement
func (gameDef *GameDef) ScriptIfBlockStart(lineData []byte, scriptData fileio.ScriptFunction) {
	byteArr := bytes.NewBuffer(lineData)
	conditional := ScriptInstrIfElseStart{}
	binary.Read(byteArr, binary.LittleEndian, &conditional)

	opcode := lineData[0]
	gameDef.ScriptMemory.IncrementProgramCounter(opcode)
	endIfBlock := gameDef.ScriptMemory.ProgramCounter + int(conditional.BlockLength)
	gameDef.ScriptMemory.OverrideProgramCounter = true

	// It's either if/endif or if/else
	// There is no if/else/endif

	// If/endif
	// Endif is followed by no-op
	if value, exists := scriptData.Instructions[endIfBlock-2]; exists {
		if int(value[0]) != fileio.OP_END_IF {
			log.Fatal("If statement is missing endif statement")
		}
		gameDef.ScriptMemory.PushStack([]int{int(opcode), endIfBlock - 2})
		return
	}

	// If/else
	if value, exists := scriptData.Instructions[endIfBlock-4]; exists {
		if int(value[0]) != fileio.OP_ELSE_START {
			log.Fatal("If statement is missing else statement")
		}
		gameDef.ScriptMemory.PushStack([]int{int(opcode), endIfBlock - 4})
		return
	}
	log.Fatal("If statement does not have else statement or endif statement")
}

func (gameDef *GameDef) ScriptElseCheck(lineData []byte) {
	byteArr := bytes.NewBuffer(lineData)
	conditional := ScriptInstrElseStart{}
	binary.Read(byteArr, binary.LittleEndian, &conditional)

	endElseBlock := gameDef.ScriptMemory.ProgramCounter + int(conditional.BlockLength)
	stackTop := gameDef.ScriptMemory.PopStack()
	if stackTop[0] != fileio.OP_CHECK {
		log.Fatal("Else statement is missing check statement")
	}
	check := stackTop[1]
	// If-statement evaluated to true
	if check == 1 {
		// Skip else block
		gameDef.ScriptMemory.ProgramCounter = endElseBlock
		gameDef.ScriptMemory.OverrideProgramCounter = true
	}
	// Pop value from if statement
	stackTop = gameDef.ScriptMemory.PopStack()
	if stackTop[0] != fileio.OP_IF_START {
		log.Fatal("Endif statement is missing if statement")
	}
}

func (gameDef *GameDef) ScriptEndIf() {
	// Pop value from check statement
	stackTop := gameDef.ScriptMemory.PopStack()
	if stackTop[0] != fileio.OP_CHECK &&
		stackTop[0] != fileio.OP_COMPARE &&
		stackTop[0] != fileio.OP_MEMBER_CMP {
		log.Fatal("Endif statement is missing boolean")
	}
	// Pop value from if statement
	stackTop = gameDef.ScriptMemory.PopStack()
	if stackTop[0] != fileio.OP_IF_START {
		log.Fatal("End if statement is missing if statement")
	}
}

func (gameDef *GameDef) ScriptSwitchBegin(lineData []byte) {
	byteArr := bytes.NewBuffer(lineData)
	conditional := ScriptInstrSwitch{}
	binary.Read(byteArr, binary.LittleEndian, &conditional)

	opcode := int(lineData[0])
	endSwitch := gameDef.ScriptMemory.ProgramCounter + int(conditional.BlockLength)
	switchValue := gameDef.GetScriptVariable(int(conditional.VarId))
	gameDef.ScriptMemory.PushStack([]int{opcode, endSwitch, switchValue})
}

func (gameDef *GameDef) ScriptSwitchCase(lineData []byte) {
	byteArr := bytes.NewBuffer(lineData)
	conditional := ScriptInstrSwitchCase{}
	binary.Read(byteArr, binary.LittleEndian, &conditional)

	switchFallthrough := false
	stackTop := gameDef.ScriptMemory.PeekStackTop()
	// Check if the previous cases are true
	if stackTop[0] == fileio.OP_CASE {
		gameDef.ScriptMemory.PopStack()
		switchFallthrough = true
	}

	stackTop = gameDef.ScriptMemory.PeekStackTop()
	switchValue := stackTop[2]

	opcode := lineData[0]
	if int(conditional.Value) != switchValue && !switchFallthrough {
		// Move to the next case
		gameDef.ScriptMemory.ProgramCounter += int(conditional.BlockLength)
		gameDef.ScriptMemory.IncrementProgramCounter(opcode)
		gameDef.ScriptMemory.OverrideProgramCounter = true
	} else {
		// Switch statement fallthrough
		if int(conditional.BlockLength) == 0 {
			gameDef.ScriptMemory.PushStack([]int{int(opcode), 1})
		}
	}
}

func (gameDef *GameDef) ScriptSwitchEnd() {
	// Pop beginning of switch
	stackTop := gameDef.ScriptMemory.PopStack()
	if stackTop[0] != fileio.OP_SWITCH {
		log.Fatal("End switch statement does not have a begin switch statement")
	}
}

func (gameDef *GameDef) ScriptGoSub(lineData []byte, scriptData fileio.ScriptFunction) {
	byteArr := bytes.NewBuffer(lineData)
	instruction := ScriptInstrGoSub{}
	binary.Read(byteArr, binary.LittleEndian, &instruction)

	// Call function
	// Save program counter to be after the function
	opcode := lineData[0]
	gameDef.ScriptMemory.IncrementProgramCounter(opcode)
	values := []int{int(opcode), gameDef.ScriptMemory.ProgramCounter, gameDef.ScriptMemory.FunctionNum}
	gameDef.ScriptMemory.PushStack(values)
	gameDef.ScriptMemory.FunctionNum = int(instruction.Event)
	gameDef.ScriptMemory.ProgramCounter = scriptData.StartProgramCounter[gameDef.ScriptMemory.FunctionNum]
	gameDef.ScriptMemory.OverrideProgramCounter = true
}

func (gameDef *GameDef) ScriptCheckBit(lineData []byte) {
	byteArr := bytes.NewBuffer(lineData)
	bitTest := ScriptInstrCheckBitTest{}
	binary.Read(byteArr, binary.LittleEndian, &bitTest)

	opcode := int(lineData[0])
	if gameDef.GetBitArray(int(bitTest.BitArray), int(bitTest.Number)) != int(bitTest.Value) {
		stackTop := gameDef.ScriptMemory.PeekStackTop()
		// Skip if the check statement is missing an if statement before it
		if stackTop[0] != fileio.OP_IF_START {
			return
		}
		// stackTop = gameDef.ScriptMemory.PopStack()
		counter := stackTop[1]
		gameDef.ScriptMemory.ProgramCounter = counter
		gameDef.ScriptMemory.OverrideProgramCounter = true
		gameDef.ScriptMemory.PushStack([]int{opcode, 0})
	} else {
		gameDef.ScriptMemory.PushStack([]int{opcode, 1})
	}
}

func (gameDef *GameDef) ScriptSetBit(lineData []byte) {
	byteArr := bytes.NewBuffer(lineData)
	bitTest := ScriptInstrCheckBitTest{}
	binary.Read(byteArr, binary.LittleEndian, &bitTest)

	gameDef.SetBitArray(int(bitTest.BitArray), int(bitTest.Number), int(bitTest.Value))
}

func (gameDef *GameDef) ScriptCompare(lineData []byte) {
	byteArr := bytes.NewBuffer(lineData)
	instruction := ScriptInstrMemberCompare{}
	binary.Read(byteArr, binary.LittleEndian, &instruction)

	// TODO: Evaluate to true or false
	// Assumes the statement is false
	opcode := int(lineData[0])
	stackTop := gameDef.ScriptMemory.PeekStackTop()
	// Skip if the compare statement is missing an if statement before it
	if stackTop[0] != fileio.OP_IF_START {
		return
	}
	counter := stackTop[1]
	gameDef.ScriptMemory.ProgramCounter = counter
	gameDef.ScriptMemory.OverrideProgramCounter = true
	gameDef.ScriptMemory.PushStack([]int{opcode, 0})
}

func (gameDef *GameDef) ScriptSpriteOn(lineData []byte) {
	byteArr := bytes.NewBuffer(lineData)
	scriptSprite := ScriptSprite{}
	binary.Read(byteArr, binary.LittleEndian, &scriptSprite)

	gameDef.Sprites = append(gameDef.Sprites, scriptSprite)
}

func (gameDef *GameDef) ScriptDoorAotSet(lineData []byte) {
	byteArr := bytes.NewBuffer(lineData)
	var door ScriptDoor
	err := binary.Read(byteArr, binary.LittleEndian, &door)
	if err != nil {
		log.Fatal("Error loading door")
	}

	gameDef.Doors = append(gameDef.Doors, door)
}

func (gameDef *GameDef) ScriptMemberCompare(lineData []byte) {
	byteArr := bytes.NewBuffer(lineData)
	instruction := ScriptInstrMemberCompare{}
	binary.Read(byteArr, binary.LittleEndian, &instruction)

	// TODO: Evaluate to true or false
	// Assumes the statement is false
	opcode := int(lineData[0])
	stackTop := gameDef.ScriptMemory.PeekStackTop()
	// Skip if the compare statement is missing an if statement before it
	if stackTop[0] != fileio.OP_IF_START {
		return
	}
	counter := stackTop[1]
	gameDef.ScriptMemory.ProgramCounter = counter
	gameDef.ScriptMemory.OverrideProgramCounter = true
	gameDef.ScriptMemory.PushStack([]int{opcode, 0})
}
