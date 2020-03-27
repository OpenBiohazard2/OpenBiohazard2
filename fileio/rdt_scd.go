package fileio

// .scd - Script data

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
)

const (
	OP_NO_OP           = 0
	OP_EVT_END         = 1
	OP_EVT_NEXT        = 2
	OP_EVT_EXEC        = 4
	OP_EVT_KILL        = 5
	OP_IF_START        = 6
	OP_ELSE_START      = 7
	OP_END_IF          = 8
	OP_SLEEP           = 9
	OP_SLEEPING        = 10
	OP_WSLEEP          = 11
	OP_WSLEEPING       = 12
	OP_FOR             = 13
	OP_FOR_END         = 14
	OP_WHILE_START     = 15
	OP_WHILE_END       = 16
	OP_DO_START        = 17
	OP_DO_END          = 18
	OP_SWITCH          = 19
	OP_CASE            = 20
	OP_DEFAULT         = 21
	OP_END_SWITCH      = 22
	OP_GOTO            = 23
	OP_GOSUB           = 24
	OP_GOSUB_RETURN    = 25
	OP_BREAK           = 26
	OP_WORK_COPY       = 29
	OP_NO_OP2          = 32
	OP_CHECK           = 33
	OP_SET_BIT         = 34
	OP_COMPARE         = 35
	OP_SAVE            = 36
	OP_COPY            = 37
	OP_CALC            = 38
	OP_CALC2           = 39
	OP_SCE_RND         = 40
	OP_CUT_CHG         = 41
	OP_CUT_OLD         = 42
	OP_MESSAGE_ON      = 43
	OP_AOT_SET         = 44
	OP_OBJ_MODEL_SET   = 45
	OP_WORK_SET        = 46
	OP_SPEED_SET       = 47
	OP_ADD_SPEED       = 48
	OP_ADD_ASPEED      = 49
	OP_POS_SET         = 50
	OP_DIR_SET         = 51
	OP_MEMBER_SET      = 52
	OP_MEMBER_SET2     = 53
	OP_SE_ON           = 54
	OP_SCA_ID_SET      = 55
	OP_DIR_CK          = 57
	OP_SCE_ESPR_ON     = 58
	OP_DOOR_AOT_SET    = 59
	OP_CUT_AUTO        = 60
	OP_MEMBER_COPY     = 61
	OP_MEMBER_CMP      = 62
	OP_PLC_MOTION      = 63
	OP_PLC_DEST        = 64
	OP_PLC_NECK        = 65
	OP_PLC_RET         = 66
	OP_PLC_FLAG        = 67
	OP_SCE_EM_SET      = 68
	OP_AOT_RESET       = 70
	OP_AOT_ON          = 71
	OP_SUPER_SET       = 72
	OP_CUT_REPLACE     = 75
	OP_SCE_ESPR_KILL   = 76
	OP_ITEM_AOT_SET    = 78
	OP_SCE_TRG_CK      = 80
	OP_SCE_BGM_CONTROL = 81
	OP_SCE_FADE_SET    = 83
	OP_SCE_ESPR3D_ON   = 84
	OP_SCE_BGMTBL_SET  = 87
	OP_PLC_ROT         = 88
	OP_XA_ON           = 89
	OP_PLC_CNT         = 91
	OP_MIZU_DIV_SET    = 93
	OP_KEEP_ITEM_CK    = 94
	OP_XA_VOL          = 95
	OP_KAGE_SET        = 96
	OP_CUT_BE_SET      = 97
	OP_SCE_ITEM_LOST   = 98
	OP_PLC_STOP        = 102
	OP_AOT_SET_4P      = 103
	OP_DOOR_AOT_SET_4P = 104
	OP_ITEM_AOT_SET_4P = 105
	OP_LIGHT_KIDO_SET  = 107
	OP_SCE_SCR_MOVE    = 109
	OP_PARTS_SET       = 110
	OP_MOVIE_ON        = 111
	OP_SCE_PARTS_BOMB  = 122
	OP_SCE_PARTS_DOWN  = 123
)

var (
	InstructionSize = map[byte]int{
		OP_NO_OP:           1,
		OP_EVT_END:         1,
		OP_EVT_NEXT:        1,
		OP_EVT_EXEC:        4,
		OP_EVT_KILL:        2,
		OP_IF_START:        4,
		OP_ELSE_START:      4,
		OP_END_IF:          1,
		OP_SLEEP:           4,
		OP_SLEEPING:        3,
		OP_WSLEEP:          1,
		OP_WSLEEPING:       1,
		OP_FOR:             6,
		OP_FOR_END:         2,
		OP_WHILE_START:     4,
		OP_WHILE_END:       2,
		OP_DO_START:        4,
		OP_DO_END:          2,
		OP_SWITCH:          4,
		OP_CASE:            6,
		OP_DEFAULT:         2,
		OP_END_SWITCH:      2,
		OP_GOTO:            6,
		OP_GOSUB:           2,
		OP_GOSUB_RETURN:    2,
		OP_BREAK:           2,
		OP_WORK_COPY:       4,
		OP_NO_OP2:          1,
		OP_CHECK:           4,
		OP_SET_BIT:         4,
		OP_COMPARE:         6,
		OP_SAVE:            4,
		OP_COPY:            3,
		OP_CALC:            6,
		OP_CALC2:           4,
		OP_SCE_RND:         1,
		OP_CUT_CHG:         2,
		OP_CUT_OLD:         1,
		OP_MESSAGE_ON:      6,
		OP_AOT_SET:         20,
		OP_OBJ_MODEL_SET:   38,
		OP_WORK_SET:        3,
		OP_SPEED_SET:       4,
		OP_ADD_SPEED:       1,
		OP_ADD_ASPEED:      1,
		OP_POS_SET:         8,
		OP_DIR_SET:         8,
		OP_MEMBER_SET:      4,
		OP_MEMBER_SET2:     3,
		OP_SE_ON:           12,
		OP_SCA_ID_SET:      4,
		OP_DIR_CK:          8,
		OP_SCE_ESPR_ON:     16,
		OP_DOOR_AOT_SET:    32,
		OP_CUT_AUTO:        2,
		OP_MEMBER_COPY:     3,
		OP_MEMBER_CMP:      6,
		OP_PLC_MOTION:      4,
		OP_PLC_DEST:        8,
		OP_PLC_NECK:        10,
		OP_PLC_RET:         1,
		OP_PLC_FLAG:        4,
		OP_SCE_EM_SET:      22,
		OP_AOT_RESET:       10,
		OP_AOT_ON:          2,
		OP_SUPER_SET:       16,
		OP_CUT_REPLACE:     3,
		OP_SCE_ESPR_KILL:   5,
		OP_ITEM_AOT_SET:    22,
		OP_SCE_TRG_CK:      4,
		OP_SCE_BGM_CONTROL: 6,
		OP_SCE_FADE_SET:    6,
		OP_SCE_ESPR3D_ON:   22,
		OP_SCE_BGMTBL_SET:  8,
		OP_PLC_ROT:         4,
		OP_XA_ON:           4,
		OP_PLC_CNT:         2,
		OP_MIZU_DIV_SET:    2,
		OP_KEEP_ITEM_CK:    2,
		OP_XA_VOL:          2,
		OP_KAGE_SET:        14,
		OP_CUT_BE_SET:      4,
		OP_SCE_ITEM_LOST:   2,
		OP_PLC_STOP:        1,
		OP_AOT_SET_4P:      28,
		OP_DOOR_AOT_SET_4P: 40,
		OP_ITEM_AOT_SET_4P: 30,
		OP_LIGHT_KIDO_SET:  4,
		OP_SCE_SCR_MOVE:    4,
		OP_PARTS_SET:       6,
		OP_MOVIE_ON:        2,
		OP_SCE_PARTS_BOMB:  16,
		OP_SCE_PARTS_DOWN:  16,
	}
)

type ScriptInstrEventExec struct {
	Opcode    uint8 // 0x04
	ThreadNum uint8
	ExOpcode  uint8
	Event     uint8
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

type ScriptInstrSleep struct {
	Opcode uint8 // 0x09
	Dummy  uint8
	Count  uint16
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
	IfElseCounter int8
	LoopCounter   int8
	Unknown       uint8
	Offset        int16
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

type ScriptInstrSetBit struct {
	Opcode    uint8 // 0x22
	BitArray  uint8 // Index of array of bits to use
	BitNumber uint8 // Bit number to check
	Operation uint8 // 0x0: clear, 0x1: set, 0x2-0x6: invalid, 0x7: flip bit
}

type ScriptInstrObjModelSet struct {
	Opcode      uint8 // 0x2d
	ObjectIndex uint8
	ObjectId    uint8
	Counter     uint8
	Wait        uint8
	Num         uint8
	Floor       uint8
	Flag0       uint8
	Type        uint16
	Flag1       uint16
	Attribute   int16
	Position    [3]int16
	Direction   [3]int16
	Offset      [3]int16
	Dimensions  [3]uint16
}

type ScriptInstrWorkSet struct {
	Opcode    uint8 // 0x2e
	Component uint8
	Index     uint8
}

type ScriptInstrPosSet struct {
	Opcode uint8 // 0x32
	Dummy  uint8
	X      int16
	Y      int16
	Z      int16
}

type ScriptInstrScaIdSet struct {
	Opcode uint8 // 0x37
	Id     uint8
	Flag   uint16
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

type ScriptInstrDoorAotSet struct {
	Opcode                       uint8 // 0x3b
	Aot                          uint8 // Index of item in array of room objects list
	Id                           uint8
	Type                         uint8
	Floor                        uint8
	Super                        uint8
	X, Y                         int16 // Location of door
	Width, Height                int16 // Size of door
	NextX, NextY, NextZ, NextDir int16 // Position and direction of player after door entered
	Stage, Room, Camera          uint8 // Stage, room, camera after door entered
	NextFloor                    uint8
	TextureType                  uint8
	DoorType                     uint8
	KnockType                    uint8
	KeyId                        uint8
	KeyType                      uint8
	Free                         uint8
}

type ScriptInstrCutAuto struct {
	Opcode uint8 // 0x3c
	FlagOn uint8
}

type ScriptInstrMemberCompare struct {
	Opcode   uint8 // 0x3e
	Unknown0 uint8
	Unknown1 uint8
	Compare  uint8
	Value    int16
}

type ScriptInstrCutChg struct {
	Opcode   uint8 // 0x29
	CameraId uint8
}

type ScriptInstrAotSet struct {
	Opcode       uint8 // 0x2c
	Aot          uint8
	Id           uint8
	Type         uint8
	Floor        uint8
	Super        uint8
	X, Z         int16
	Width, Depth int16
	Data         [6]uint8
}

type ScriptInstrPlcMotion struct {
	Opcode     uint8 // 0x3f
	Action     uint8
	MoveNumber uint8
	SceneFlag  uint8
}

type ScriptInstrPlcDest struct {
	Opcode     uint8 // 0x40
	Dummy      uint8
	Action     uint8
	FlagNumber uint8
	DestX      int16
	DestZ      int16
}

type ScriptInstrPlcNeck struct {
	Opcode    uint8 // 0x41
	Operation uint8
	NeckX     int16
	NeckY     int16
	NeckZ     int16
	Unknown   [2]int8
}

type ScriptInstrPlcFlag struct {
	Opcode    uint8 // 0x43
	Operation uint8 // 0: OR, 1: Set, 2: XOR
	Flag      uint16
}

type ScriptInstrAotReset struct {
	Opcode uint8 // 0x46
	Aot    uint8
	Id     uint8
	Type   uint8
	Data   [6]uint8
}

type ScriptInstrItemAotSet struct {
	Opcode          uint8 // 0x4e
	Aot             uint8
	Id              uint8
	Type            uint8
	Floor           uint8
	Super           uint8
	X, Z            int16
	Width, Depth    int16
	ItemId          uint16
	Amount          uint16
	ItemPickedIndex uint16 // flag to check if item is picked up
	Md1ModelId      uint8
	Act             uint8
}

type ScriptInstrSceBgmControl struct {
	Opcode      uint8 // 0x51
	Id          uint8 // 0: Main, 1: sub0, 2: sub1
	Operation   uint8 // 0: nop, 1: start, 2: stop, 3: restart, 4: pause, 5: fadeout
	Type        uint8 // 0: MAIN_VOL, 1: PROG0_VOL, 2: PROG1_VOL, 3: PROG2_VOL
	LeftVolume  uint8
	RightVolume uint8
}

type ScriptInstrPlcRot struct {
	Opcode uint8 // 0x58
	Index  uint8 // 0 or 1
	Value  int16
}

type ScriptInstrXaOn struct {
	Opcode  uint8 // 0x59
	Channel uint8 // channel on which to play sound
	Id      int16 // ID of sound to play
}

type ScriptInstrMizuDivSet struct {
	Opcode     uint8 // 0x5d
	MizuDivMax uint8
}

type ScriptInstrKageSet struct {
	Opcode           uint8 // 0x60
	WorkSetComponent uint8
	WorkSetIndex     uint8
	Color            [3]uint8
	HalfX            int16
	HalfZ            int16
	OffsetX          int16
	OffsetZ          int16
}

type SCDOutput struct {
	ScriptData ScriptFunction
}

type ScriptFunction struct {
	Instructions        map[int][]byte // key is program counter, value is command
	StartProgramCounter []int          // set per function
}

func LoadRDT_SCDStream(fileReader io.ReaderAt, fileLength int64) (*SCDOutput, error) {
	streamReader := io.NewSectionReader(fileReader, int64(0), fileLength)
	firstOffset := uint16(0)
	if err := binary.Read(streamReader, binary.LittleEndian, &firstOffset); err != nil {
		return nil, err
	}

	functionOffsets := make([]uint16, 0)
	functionOffsets = append(functionOffsets, firstOffset)
	for i := 2; i < int(firstOffset); i += 2 {
		nextOffset := uint16(0)
		if err := binary.Read(streamReader, binary.LittleEndian, &nextOffset); err != nil {
			return nil, err
		}
		functionOffsets = append(functionOffsets, nextOffset)
	}

	programCounter := 0
	scriptData := ScriptFunction{}
	scriptData.Instructions = make(map[int][]byte)
	scriptData.StartProgramCounter = make([]int, 0)
	for functionNum := 0; functionNum < len(functionOffsets); functionNum++ {
		scriptData.StartProgramCounter = append(scriptData.StartProgramCounter, programCounter)

		var functionLength int64
		if functionNum != len(functionOffsets)-1 {
			functionLength = int64(functionOffsets[functionNum+1]) - int64(functionOffsets[functionNum])
		} else {
			functionLength = fileLength - int64(functionOffsets[functionNum])
		}

		streamReader = io.NewSectionReader(fileReader, int64(functionOffsets[functionNum]), functionLength)
		for lineNum := 0; lineNum < int(functionLength); lineNum++ {
			opcode := byte(0)
			if err := binary.Read(streamReader, binary.LittleEndian, &opcode); err != nil {
				return nil, err
			}

			byteSize, exists := InstructionSize[opcode]
			if !exists {
				fmt.Println("Unknown opcode:", opcode)
			}

			scriptData.Instructions[programCounter] = generateScriptLine(streamReader, byteSize, opcode)
			// Sleep contains sleep and sleeping commands
			if opcode == OP_SLEEP {
				scriptData.Instructions[programCounter+1] = scriptData.Instructions[programCounter][1:]
			}

			programCounter += byteSize

			// return
			if opcode == OP_EVT_END {
				break
			}
		}
	}

	output := &SCDOutput{
		ScriptData: scriptData,
	}
	return output, nil
}

func generateScriptLine(streamReader *io.SectionReader, totalByteSize int, opcode byte) []byte {
	scriptLine := make([]byte, 0)
	scriptLine = append(scriptLine, opcode)

	if totalByteSize == 1 {
		return scriptLine
	}

	parameters, err := readRemainingBytes(streamReader, totalByteSize-1)
	if err != nil {
		log.Fatal("Error reading script for opcode %v\n", opcode)
	}
	scriptLine = append(scriptLine, parameters...)
	return scriptLine
}

func readRemainingBytes(streamReader *io.SectionReader, byteSize int) ([]byte, error) {
	parameters := make([]byte, byteSize)
	if err := binary.Read(streamReader, binary.LittleEndian, &parameters); err != nil {
		return nil, err
	}
	return parameters, nil
}
