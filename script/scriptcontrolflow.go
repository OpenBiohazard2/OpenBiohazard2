package script

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"

	"github.com/OpenBiohazard2/OpenBiohazard2/fileio"
)

func (scriptDef *ScriptDef) ScriptIfBlockStart(scriptThread *ScriptThread, lineData []byte) int {
	byteArr := bytes.NewBuffer(lineData)
	conditional := fileio.ScriptInstrIfElseStart{}
	binary.Read(byteArr, binary.LittleEndian, &conditional)

	opcode := lineData[0]
	scriptThread.LevelState[scriptThread.SubLevel].IfElseCounter++
	newPosition := (scriptThread.ProgramCounter + fileio.InstructionSize[opcode]) + int(conditional.BlockLength)
	scriptThread.PushStack(newPosition)

	return 1
}

func (scriptDef *ScriptDef) ScriptElseCheck(scriptThread *ScriptThread, lineData []byte) int {
	byteArr := bytes.NewBuffer(lineData)
	conditional := fileio.ScriptInstrElseStart{}
	binary.Read(byteArr, binary.LittleEndian, &conditional)

	scriptThread.StackIndex--
	scriptThread.LevelState[scriptThread.SubLevel].IfElseCounter--

	// Jump to position after the else block
	scriptThread.ProgramCounter = scriptThread.ProgramCounter + int(conditional.BlockLength)
	scriptThread.OverrideProgramCounter = true
	return 1
}

func (scriptDef *ScriptDef) ScriptEndIf() int {
	scriptThread.StackIndex--
	scriptThread.LevelState[scriptThread.SubLevel].IfElseCounter--
	return 1
}

func (scriptDef *ScriptDef) ScriptForLoopBegin(lineData []byte) int {
	byteArr := bytes.NewBuffer(lineData)
	instruction := fileio.ScriptInstrForStart{}
	binary.Read(byteArr, binary.LittleEndian, &instruction)

	opcode := lineData[0]
	if instruction.Count != 0 {
		// Set the program counter to after the instruction
		// so that this instruction is only run once to initialize for loop
		newProgramCounter := scriptThread.ProgramCounter + fileio.InstructionSize[opcode]
		curLevelState := scriptThread.LevelState[scriptThread.SubLevel]

		curLevelState.LoopLevel++
		newLoopState := curLevelState.LoopState[curLevelState.LoopLevel]
		newLoopState.Counter = int(instruction.Count)
		newLoopState.Break = newProgramCounter + int(instruction.BlockLength)
		newLoopState.StackValue = newProgramCounter
		newLoopState.LevelIfCounter = curLevelState.IfElseCounter

		scriptThread.ProgramCounter = newProgramCounter
		scriptThread.OverrideProgramCounter = true
		return 1
	}

	// Jump to end of for loop
	newProgramCounter := (scriptThread.ProgramCounter + fileio.InstructionSize[opcode]) + int(instruction.BlockLength)
	scriptThread.ProgramCounter = newProgramCounter

	scriptThread.OverrideProgramCounter = true
	return 1
}

func (scriptDef *ScriptDef) ScriptForLoopEnd(lineData []byte) int {
	opcode := lineData[0]
	curLevelState := scriptThread.LevelState[scriptThread.SubLevel]
	curLoopState := curLevelState.LoopState[curLevelState.LoopLevel]
	curLoopState.Counter--

	if curLoopState.Counter != 0 {
		// Go back to beginning of for loop
		scriptThread.ProgramCounter = curLoopState.StackValue
		scriptThread.OverrideProgramCounter = true
		return 1
	}

	// Exit for loop block
	curLevelState.LoopLevel--
	scriptThread.ProgramCounter += fileio.InstructionSize[opcode]
	scriptThread.OverrideProgramCounter = true
	return 1
}

func (scriptDef *ScriptDef) ScriptSwitchBegin(
	lineData []byte,
	instructions map[int][]byte,
) int {

	byteArr := bytes.NewBuffer(lineData)
	switchConditional := fileio.ScriptInstrSwitch{}
	binary.Read(byteArr, binary.LittleEndian, &switchConditional)

	opcode := lineData[0]
	curLevelState := scriptThread.LevelState[scriptThread.SubLevel]

	curLevelState.LoopLevel++
	newLoopLevel := curLevelState.LoopLevel
	newProgramCounter := scriptThread.ProgramCounter + fileio.InstructionSize[opcode]
	curLevelState.LoopState[newLoopLevel].Break = newProgramCounter + int(switchConditional.BlockLength)
	curLevelState.LoopState[newLoopLevel].LevelIfCounter = curLevelState.IfElseCounter

	for true {
		newLineData := instructions[newProgramCounter]
		newOpcode := newLineData[0]

		if newOpcode == fileio.OP_CASE {
			byteArr = bytes.NewBuffer(newLineData)
			caseInstruction := fileio.ScriptInstrSwitchCase{}
			binary.Read(byteArr, binary.LittleEndian, &caseInstruction)

			switchValue := scriptDef.GetScriptVariable(int(switchConditional.VarId))
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
			curLevelState.LoopLevel--
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
	scriptThread.LevelState[scriptThread.SubLevel].LoopLevel--
	return 1
}

func (scriptDef *ScriptDef) ScriptGoto(lineData []byte) int {
	byteArr := bytes.NewBuffer(lineData)
	instruction := fileio.ScriptInstrGoto{}
	binary.Read(byteArr, binary.LittleEndian, &instruction)

	// Disable due to infinite loop
	/*scriptThread.LevelState[scriptThread.SubLevel].IfElseCounter = int(instruction.IfElseCounter)
	scriptThread.StackIndex = int(instruction.IfElseCounter) + 1
	scriptThread.LevelState[scriptThread.SubLevel].LoopLevel = int(instruction.LoopLevel)
	scriptThread.ProgramCounter += int(instruction.Offset)
	scriptThread.OverrideProgramCounter = true*/

	return 1
}

func (scriptDef *ScriptDef) ScriptGoSub(lineData []byte, scriptData fileio.ScriptFunction) int {
	byteArr := bytes.NewBuffer(lineData)
	instruction := fileio.ScriptInstrGoSub{}
	binary.Read(byteArr, binary.LittleEndian, &instruction)

	opcode := lineData[0]
	scriptDef.ScriptDebugLine(fmt.Sprintf("(Gosub) Go to sub function %v", instruction.Event))
	scriptThread.LevelState[scriptThread.SubLevel].ReturnAddress = scriptThread.ProgramCounter + fileio.InstructionSize[opcode]
	scriptThread.LevelState[scriptThread.SubLevel+1].IfElseCounter = -1
	scriptThread.LevelState[scriptThread.SubLevel+1].LoopLevel = -1
	scriptThread.StackIndex = 0
	scriptThread.SubLevel++

	scriptThread.ProgramCounter = scriptData.StartProgramCounter[instruction.Event]
	scriptThread.OverrideProgramCounter = true
	scriptThread.FunctionIds = append(scriptThread.FunctionIds, int(instruction.Event))
	return 1
}

func (scriptDef *ScriptDef) ScriptBreak(lineData []byte) int {
	curLevelState := scriptThread.LevelState[scriptThread.SubLevel]
	curLoopState := curLevelState.LoopState[curLevelState.LoopLevel]

	scriptThread.OverrideProgramCounter = true
	scriptThread.ProgramCounter = curLoopState.Break
	curLevelState.IfElseCounter = curLoopState.LevelIfCounter
	curLevelState.LoopLevel--
	return 1
}
