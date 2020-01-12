package fileio

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
)

const (
	OP_NO_OP         = 0
	OP_RETURN        = 1
	OP_IF_ELSE       = 6
	OP_ELSE_CHECK    = 7
	OP_END_IF        = 8
	OP_CHECK         = 33
	OP_AOT_SET       = 44
	OP_OBJ_MODEL_SET = 45
	OP_DOOR_AOT_SET  = 59
	OP_SCE_EM_SET    = 68
	OP_AOT_RESET     = 70
	OP_ITEM_AOT_SET  = 78
	OP_AOT_SET_4P    = 103
)

var (
	instructionSize = map[byte]int{
		OP_NO_OP:         1,
		OP_RETURN:        1,
		OP_IF_ELSE:       4,
		OP_ELSE_CHECK:    4,
		OP_END_IF:        1,
		OP_CHECK:         4,
		OP_AOT_SET:       20,
		OP_OBJ_MODEL_SET: 38,
		OP_DOOR_AOT_SET:  32,
		OP_SCE_EM_SET:    22,
		OP_AOT_RESET:     10,
		OP_ITEM_AOT_SET:  22,
		OP_AOT_SET_4P:    28,
	}
)

type SCDOutput struct {
	ScriptData [][][]byte
}

func LoadRDT_SCDStream(fileReader io.ReaderAt, fileLength int64, rdtHeader RDTHeader, offsets RDTOffsets) (*SCDOutput, error) {
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

	scriptData := make([][][]byte, len(functionOffsets))

	for functionNum := 0; functionNum < len(functionOffsets); functionNum++ {
		functionData := make([][]byte, 0)

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

			byteSize, exists := instructionSize[opcode]
			if !exists {
				fmt.Println("Unknown opcode:", opcode)
			}

			functionData = append(functionData, generateScriptLine(streamReader, byteSize, opcode))

			// return
			if opcode == OP_RETURN {
				break
			}
		}

		scriptData[functionNum] = functionData
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
