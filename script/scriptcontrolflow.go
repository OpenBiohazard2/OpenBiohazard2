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

func (scriptDef *ScriptDef) ScriptEndIf(thread *ScriptThread) int {
	thread.StackIndex--
	thread.LevelState[thread.SubLevel].IfElseCounter--
	return 1
}

func (scriptDef *ScriptDef) ScriptForLoopBegin(thread *ScriptThread, lineData []byte) int {
	byteArr := bytes.NewBuffer(lineData)
	instruction := fileio.ScriptInstrForStart{}
	binary.Read(byteArr, binary.LittleEndian, &instruction)

	opcode := lineData[0]
	if instruction.Count != 0 {
		// Set the program counter to after the instruction
		// so that this instruction is only run once to initialize for loop
		newProgramCounter := thread.ProgramCounter + fileio.InstructionSize[opcode]
		curLevelState := thread.LevelState[thread.SubLevel]

		curLevelState.LoopLevel++
		newLoopState := curLevelState.LoopState[curLevelState.LoopLevel]
		newLoopState.Counter = int(instruction.Count)
		newLoopState.Break = newProgramCounter + int(instruction.BlockLength)
		newLoopState.StackValue = newProgramCounter
		newLoopState.LevelIfCounter = curLevelState.IfElseCounter

		thread.ProgramCounter = newProgramCounter
		thread.OverrideProgramCounter = true
		return 1
	}

	// Jump to end of for loop
	newProgramCounter := (thread.ProgramCounter + fileio.InstructionSize[opcode]) + int(instruction.BlockLength)
	thread.ProgramCounter = newProgramCounter

	thread.OverrideProgramCounter = true
	return 1
}

func (scriptDef *ScriptDef) ScriptForLoopEnd(thread *ScriptThread, lineData []byte) int {
	opcode := lineData[0]
	curLevelState := thread.LevelState[thread.SubLevel]
	curLoopState := curLevelState.LoopState[curLevelState.LoopLevel]
	curLoopState.Counter--

	if curLoopState.Counter != 0 {
		// Go back to beginning of for loop
		thread.ProgramCounter = curLoopState.StackValue
		thread.OverrideProgramCounter = true
		return 1
	}

	// Exit for loop block
	curLevelState.LoopLevel--
	thread.ProgramCounter += fileio.InstructionSize[opcode]
	thread.OverrideProgramCounter = true
	return 1
}

func (scriptDef *ScriptDef) ScriptSwitchBegin(
	thread *ScriptThread,
	lineData []byte,
	instructions map[int][]byte,
) int {

	byteArr := bytes.NewBuffer(lineData)
	switchConditional := fileio.ScriptInstrSwitch{}
	binary.Read(byteArr, binary.LittleEndian, &switchConditional)

	opcode := lineData[0]
	curLevelState := thread.LevelState[thread.SubLevel]

	curLevelState.LoopLevel++
	newLoopLevel := curLevelState.LoopLevel
	newProgramCounter := thread.ProgramCounter + fileio.InstructionSize[opcode]
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
				thread.ProgramCounter = newProgramCounter + fileio.InstructionSize[newOpcode]
				thread.OverrideProgramCounter = true
				return 1
			} else {
				// Move to the next case statement
				newProgramCounter += fileio.InstructionSize[newOpcode] + int(caseInstruction.BlockLength)
			}
		} else if newOpcode == fileio.OP_DEFAULT {
			thread.ProgramCounter = newProgramCounter + fileio.InstructionSize[newOpcode]
			thread.OverrideProgramCounter = true
			return 1
		} else if newOpcode == fileio.OP_END_SWITCH {
			curLevelState.LoopLevel--
			thread.ProgramCounter = newProgramCounter + fileio.InstructionSize[newOpcode]
			thread.OverrideProgramCounter = true
			return 1
		} else {
			log.Fatal("Switch statement has unknown opcode: ", newOpcode)
		}
	}

	return 1
}

func (scriptDef *ScriptDef) ScriptSwitchEnd(thread *ScriptThread) int {
	thread.LevelState[thread.SubLevel].LoopLevel--
	return 1
}

func (scriptDef *ScriptDef) ScriptGoto(thread *ScriptThread, lineData []byte) int {
	byteArr := bytes.NewBuffer(lineData)
	instruction := fileio.ScriptInstrGoto{}
	binary.Read(byteArr, binary.LittleEndian, &instruction)

	// Disable due to infinite loop
	/*thread.LevelState[thread.SubLevel].IfElseCounter = int(instruction.IfElseCounter)
	thread.StackIndex = int(instruction.IfElseCounter) + 1
	thread.LevelState[thread.SubLevel].LoopLevel = int(instruction.LoopLevel)
	thread.ProgramCounter += int(instruction.Offset)
	thread.OverrideProgramCounter = true*/

	return 1
}

func (scriptDef *ScriptDef) ScriptGoSub(thread *ScriptThread, lineData []byte, scriptData fileio.ScriptFunction) int {
	byteArr := bytes.NewBuffer(lineData)
	instruction := fileio.ScriptInstrGoSub{}
	binary.Read(byteArr, binary.LittleEndian, &instruction)

	opcode := lineData[0]
	scriptDef.ScriptDebugLine(fmt.Sprintf("(Gosub) Go to sub function %v", instruction.Event))
	thread.LevelState[thread.SubLevel].ReturnAddress = thread.ProgramCounter + fileio.InstructionSize[opcode]
	thread.LevelState[thread.SubLevel+1].IfElseCounter = -1
	thread.LevelState[thread.SubLevel+1].LoopLevel = -1
	thread.StackIndex = 0
	thread.SubLevel++

	thread.ProgramCounter = scriptData.StartProgramCounter[instruction.Event]
	thread.OverrideProgramCounter = true
	thread.FunctionIds = append(thread.FunctionIds, int(instruction.Event))
	return 1
}

func (scriptDef *ScriptDef) ScriptBreak(thread *ScriptThread, lineData []byte) int {
	curLevelState := thread.LevelState[thread.SubLevel]
	curLoopState := curLevelState.LoopState[curLevelState.LoopLevel]

	thread.OverrideProgramCounter = true
	thread.ProgramCounter = curLoopState.Break
	curLevelState.IfElseCounter = curLoopState.LevelIfCounter
	curLevelState.LoopLevel--
	return 1
}
