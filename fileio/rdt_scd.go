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
	OP_RETURN          = 1
	OP_EVT_NEXT        = 2
	OP_EVT_EXEC        = 4
	OP_EVT_KILL        = 5
	OP_IF_START        = 6
	OP_ELSE_START      = 7
	OP_END_IF          = 8
	OP_SLEEP           = 9
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
	OP_KEEP_ITEM_CK    = 94
	OP_XA_VOL          = 95
	OP_KAGE_SET        = 96
	OP_SCE_ITEM_LOST   = 98
	OP_PLC_STOP        = 102
	OP_AOT_SET_4P      = 103
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
		OP_RETURN:          1,
		OP_EVT_NEXT:        1,
		OP_EVT_EXEC:        4,
		OP_EVT_KILL:        2,
		OP_IF_START:        4,
		OP_ELSE_START:      4,
		OP_END_IF:          1,
		OP_SLEEP:           4,
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
		OP_KEEP_ITEM_CK:    2,
		OP_XA_VOL:          2,
		OP_KAGE_SET:        14,
		OP_SCE_ITEM_LOST:   2,
		OP_PLC_STOP:        1,
		OP_AOT_SET_4P:      28,
		OP_LIGHT_KIDO_SET:  4,
		OP_SCE_SCR_MOVE:    4,
		OP_PARTS_SET:       6,
		OP_MOVIE_ON:        2,
		OP_SCE_PARTS_BOMB:  16,
		OP_SCE_PARTS_DOWN:  16,
	}
)

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
			programCounter += byteSize

			// return
			if opcode == OP_RETURN {
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
