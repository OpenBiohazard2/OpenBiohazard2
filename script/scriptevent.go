package script

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/OpenBiohazard2/OpenBiohazard2/fileio"
)

func (scriptDef *ScriptDef) ScriptEvtEnd(lineData []byte) int {
	// The program is returning from a subroutine
	if scriptThread.SubLevel != 0 {
		ifElseCounter := scriptThread.LevelState[scriptThread.SubLevel].IfElseCounter
		scriptThread.SubLevel--
		scriptThread.ProgramCounter = scriptThread.LevelState[scriptThread.SubLevel].ReturnAddress
		scriptThread.OverrideProgramCounter = true
		scriptThread.StackIndex = ifElseCounter + 1

		scriptDef.ScriptDebugLine("Exit current function")
		return INSTRUCTION_NORMAL
	}

	// The program is in the top level
	scriptThread.RunStatus = false
	scriptDef.ScriptDebugLine(fmt.Sprintf("End script thread"))
	return INSTRUCTION_THREAD_END
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

	scriptDef.ScriptDebugLine(fmt.Sprintf("Start new script thread %v", nextThreadNum))

	scriptDef.ScriptThreads[nextThreadNum].RunStatus = true
	scriptDef.ScriptThreads[nextThreadNum].ProgramCounter = scriptData.StartProgramCounter[instruction.Event]
	scriptDef.ScriptThreads[nextThreadNum].LevelState[0].IfElseCounter = -1
	scriptDef.ScriptThreads[nextThreadNum].LevelState[0].LoopLevel = -1
	return INSTRUCTION_NORMAL
}
