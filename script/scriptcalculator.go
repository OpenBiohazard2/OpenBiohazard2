package script

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"

	"github.com/OpenBiohazard2/OpenBiohazard2/fileio"
)

func (scriptDef *ScriptDef) GetBitArray(bitArrayIndex int, bitNumber int) int {
	bitArray, exists := scriptDef.ScriptBitArray[bitArrayIndex]
	if !exists {
		scriptDef.ScriptBitArray[bitArrayIndex] = make(map[int]int)
		bitArray = scriptDef.ScriptBitArray[bitArrayIndex]
		fmt.Println("Initialize bit array index", bitArrayIndex)
	}
	value, exists := bitArray[bitNumber]
	if !exists {
		bitArray[bitNumber] = 0
		fmt.Println("Initialize bit array", bitArrayIndex, "with bit number ", bitNumber)
	}
	return value
}

func (scriptDef *ScriptDef) SetBitArray(bitArrayIndex int, bitNumber int, value int) {
	_, exists := scriptDef.ScriptBitArray[bitArrayIndex]
	if !exists {
		scriptDef.ScriptBitArray[bitArrayIndex] = make(map[int]int)
	}
	scriptDef.ScriptBitArray[bitArrayIndex][bitNumber] = value
}

func (scriptDef *ScriptDef) GetScriptVariable(id int) int {
	return scriptDef.ScriptVariable[id]
}

func (scriptDef *ScriptDef) SetScriptVariable(id int, value int) {
	scriptDef.ScriptVariable[id] = value
}

func (scriptDef *ScriptDef) ScriptCheckBit(lineData []byte) int {
	byteArr := bytes.NewBuffer(lineData)
	bitTest := fileio.ScriptInstrCheckBitTest{}
	binary.Read(byteArr, binary.LittleEndian, &bitTest)

	if scriptDef.GetBitArray(int(bitTest.BitArray), int(bitTest.BitNumber)) == int(bitTest.Value) {
		return 1
	}
	return INSTRUCTION_BREAK_FLOW
}

func (scriptDef *ScriptDef) ScriptSetBit(lineData []byte) int {
	byteArr := bytes.NewBuffer(lineData)
	instruction := fileio.ScriptInstrSetBit{}
	binary.Read(byteArr, binary.LittleEndian, &instruction)

	switch int(instruction.Operation) {
	case 0:
		// Clear bit
		scriptDef.SetBitArray(int(instruction.BitArray), int(instruction.BitNumber), 0)
	case 1:
		// Set bit
		scriptDef.SetBitArray(int(instruction.BitArray), int(instruction.BitNumber), 1)
	case 7:
		// Flip bit
		currentBit := scriptDef.GetBitArray(int(instruction.BitArray), int(instruction.BitNumber))
		scriptDef.SetBitArray(int(instruction.BitArray), int(instruction.BitNumber), currentBit^1)
	default:
		log.Fatal("Set bit operation ", instruction.Operation, " is invalid.")
	}

	return 1
}

func (scriptDef *ScriptDef) ScriptCompare(lineData []byte) int {
	byteArr := bytes.NewBuffer(lineData)
	instruction := fileio.ScriptInstrCompare{}
	binary.Read(byteArr, binary.LittleEndian, &instruction)

	variableValue := scriptDef.GetScriptVariable(int(instruction.VarId))
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

func (scriptDef *ScriptDef) ScriptSave(lineData []byte) int {
	byteArr := bytes.NewBuffer(lineData)
	instruction := fileio.ScriptInstrSave{}
	binary.Read(byteArr, binary.LittleEndian, &instruction)

	scriptDef.SetScriptVariable(int(instruction.VarId), int(instruction.Value))
	return 1
}

func (scriptDef *ScriptDef) ScriptCopy(lineData []byte) int {
	byteArr := bytes.NewBuffer(lineData)
	instruction := fileio.ScriptInstrCopy{}
	binary.Read(byteArr, binary.LittleEndian, &instruction)

	sourceValue := scriptDef.GetScriptVariable(int(instruction.SourceVarId))
	scriptDef.SetScriptVariable(int(instruction.DestVarId), sourceValue)
	return 1
}

func (scriptDef *ScriptDef) ScriptCalc(lineData []byte) int {
	byteArr := bytes.NewBuffer(lineData)
	instruction := fileio.ScriptInstrCalc{}
	binary.Read(byteArr, binary.LittleEndian, &instruction)

	leftValue := int(scriptDef.GetScriptVariable(int(instruction.VarId)))
	rightValue := int(instruction.Value)
	result := scriptDef.ScriptVariableCalculator(int(instruction.Operation), leftValue, rightValue)
	scriptDef.SetScriptVariable(int(instruction.VarId), result)
	return 1
}

func (scriptDef *ScriptDef) ScriptCalc2(lineData []byte) int {
	byteArr := bytes.NewBuffer(lineData)
	instruction := fileio.ScriptInstrCalc2{}
	binary.Read(byteArr, binary.LittleEndian, &instruction)

	leftValue := int(scriptDef.GetScriptVariable(int(instruction.VarId)))
	rightValue := int(scriptDef.GetScriptVariable(int(instruction.SourceVarId)))
	result := scriptDef.ScriptVariableCalculator(int(instruction.Operation), leftValue, rightValue)
	scriptDef.SetScriptVariable(int(instruction.VarId), result)
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
