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
	Stack                  []int
	TotalTime              float64
	Check                  bool
	OverrideProgramCounter bool
}

type ScriptIfElseStart struct {
	Opcode      uint8 // 0x06
	Dummy       uint8
	BlockLength uint16
}

type ScriptCheckBitTest struct {
	Opcode   uint8 // 0x21
	BitArray uint8 // Index of array of bits to use
	Number   uint8 // Bit number to check
	Value    uint8 // Value to compare (0 or 1)
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

func NewScriptMemory() *ScriptMemory {
	return &ScriptMemory{
		ProgramCounter: 0,
		Stack:          []int{},
		TotalTime:      1000,
		Check:          false,
		OverrideProgramCounter: false,
	}
}

func (mem *ScriptMemory) IncrementProgramCounter(opcode byte) {
	mem.ProgramCounter += fileio.InstructionSize[opcode]
}

func (mem *ScriptMemory) PushStack(value int) {
	mem.Stack = append(mem.Stack, value)
}

func (mem *ScriptMemory) PopStack() int {
	if len(mem.Stack) == 0 {
		log.Fatal("Script stack is empty")
	}

	top := mem.Stack[len(mem.Stack)-1]
	// remove last element
	mem.Stack = mem.Stack[0 : len(mem.Stack)-1]
	return top
}

func (gameDef *GameDef) RunScript(scriptData []fileio.ScriptFunction, timeElapsedSeconds float64, init bool) {
	if init == false {
		gameDef.ScriptMemory.TotalTime += timeElapsedSeconds
		if gameDef.ScriptMemory.TotalTime < 1000 {
			return
		}
	}

	gameDef.ScriptMemory = NewScriptMemory()
	gameDef.ScriptMemory.ProgramCounter = 0
	for functionNum := 0; functionNum < len(scriptData); functionNum++ {
		programCounters := make([]int, 0)
		for counter, _ := range scriptData[functionNum].Instructions {
			programCounters = append(programCounters, counter)
		}
		sort.Ints(programCounters)

		lastProgramCounter := programCounters[len(programCounters)-1]
		for gameDef.ScriptMemory.ProgramCounter <= lastProgramCounter {
			lineData := scriptData[functionNum].Instructions[gameDef.ScriptMemory.ProgramCounter]
			opcode := lineData[0]

			switch opcode {
			case fileio.OP_IF_ELSE:
				gameDef.ScriptIfBlockStart(lineData)
			case fileio.OP_ELSE_CHECK:
				gameDef.ScriptElseCheck(lineData)
			case fileio.OP_CHECK:
				gameDef.ScriptCheckBit(lineData)
			case fileio.OP_DOOR_AOT_SET:
				gameDef.ScriptDoorAotSet(lineData)
			}

			if !gameDef.ScriptMemory.OverrideProgramCounter {
				gameDef.ScriptMemory.IncrementProgramCounter(opcode)
			}
			gameDef.ScriptMemory.OverrideProgramCounter = false
		}
	}
}

func (gameDef *GameDef) ScriptIfBlockStart(lineData []byte) {
	byteArr := bytes.NewBuffer(lineData)
	conditional := ScriptIfElseStart{}
	binary.Read(byteArr, binary.LittleEndian, &conditional)
	endIfBlock := gameDef.ScriptMemory.ProgramCounter + int(conditional.BlockLength)
	gameDef.ScriptMemory.PushStack(endIfBlock)
}

func (gameDef *GameDef) ScriptElseCheck(lineData []byte) {
	byteArr := bytes.NewBuffer(lineData)
	conditional := ScriptIfElseStart{}
	binary.Read(byteArr, binary.LittleEndian, &conditional)
	endElseBlock := gameDef.ScriptMemory.ProgramCounter + int(conditional.BlockLength)
	if gameDef.ScriptMemory.Check == true {
		// Skip else block
		gameDef.ScriptMemory.ProgramCounter = endElseBlock
	} else {
		// Execute else block
		opcode := lineData[0]
		gameDef.ScriptMemory.IncrementProgramCounter(opcode)
	}
	gameDef.ScriptMemory.OverrideProgramCounter = true
}

func (gameDef *GameDef) ScriptCheckBit(lineData []byte) {
	byteArr := bytes.NewBuffer(lineData)
	bitTest := ScriptCheckBitTest{}
	binary.Read(byteArr, binary.LittleEndian, &bitTest)
	if gameDef.GetBitArray(int(bitTest.BitArray), int(bitTest.Number)) != int(bitTest.Value) {
		gameDef.ScriptMemory.Check = false
		gameDef.ScriptMemory.ProgramCounter = gameDef.ScriptMemory.PopStack()
	} else {
		gameDef.ScriptMemory.Check = true
		opcode := lineData[0]
		gameDef.ScriptMemory.IncrementProgramCounter(opcode)
	}
	gameDef.ScriptMemory.OverrideProgramCounter = true
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
