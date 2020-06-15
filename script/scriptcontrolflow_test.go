package script

import (
	"testing"

	"github.com/samuelyuan/openbiohazard2/fileio"
)

func TestIfStatementTrue(t *testing.T) {
	scriptDef := NewScriptDef()
	scriptThread := scriptDef.ScriptThreads[0]
	scriptThread.ProgramCounter = 100

	dummy := byte(0)
	// little endian order
	blockLengthLeftHalf := byte(0)
	blockLengthRightHalf := byte(150)
	lineData := []byte{fileio.OP_IF_START, dummy, blockLengthRightHalf, blockLengthLeftHalf}
	returnValue := scriptDef.ScriptIfBlockStart(scriptThread, lineData)

	if returnValue != INSTRUCTION_NORMAL {
		t.Errorf("Instruction return value is incorrect, got: %d, want: %d.", returnValue, INSTRUCTION_NORMAL)
	}

	// Script has to evaluate next statement before jumping
	if scriptThread.ProgramCounter != 100 {
		t.Errorf("Script program counter is incorrect, got: %d, want: %d.", scriptThread.ProgramCounter, 100)
	}
}

func TestIfStatementFalse(t *testing.T) {
	scriptDef := NewScriptDef()
	scriptThread := scriptDef.ScriptThreads[0]
	scriptThread.ProgramCounter = 100

	dummy := byte(0)
	// little endian order
	blockLengthLeftHalf := byte(0)
	blockLengthRightHalf := byte(150)
	lineData := []byte{fileio.OP_IF_START, dummy, blockLengthRightHalf, blockLengthLeftHalf}
	returnValue := scriptDef.ScriptIfBlockStart(scriptThread, lineData)

	// Assume next statement evaluates to false
	returnValue = INSTRUCTION_BREAK_FLOW

	if scriptThread.ShouldTerminate(returnValue) != false {
		t.Errorf("Instruction return value is incorrect, got: %v, want: %v.", true, false)
	}

	if returnValue != INSTRUCTION_NORMAL {
		scriptThread.JumpToNextLocationOnStack()
	}

	if scriptThread.ProgramCounter != 254 {
		t.Errorf("Script program counter is incorrect, got: %d, want: %d.", scriptThread.ProgramCounter, 254)
	}
}

// If a statement is evaluated to true, it should execute all instructions in the if block
// and run this else statement to skip the else block
func TestElseStatementJump(t *testing.T) {
	scriptDef := NewScriptDef()
	scriptThread := scriptDef.ScriptThreads[0]
	scriptThread.ProgramCounter = 100

	dummy := byte(0)
	// little endian order
	blockLengthLeftHalf := byte(0)
	blockLengthRightHalf := byte(150)
	lineData := []byte{fileio.OP_ELSE_START, dummy, blockLengthRightHalf, blockLengthLeftHalf}
	returnValue := scriptDef.ScriptElseCheck(scriptThread, lineData)

	if returnValue != INSTRUCTION_NORMAL {
		t.Errorf("Instruction return value is incorrect, got: %d, want: %d.", returnValue, INSTRUCTION_NORMAL)
	}

	// Script should skip the else block
	if scriptThread.ProgramCounter != 250 {
		t.Errorf("Script program counter is incorrect, got: %d, want: %d.", scriptThread.ProgramCounter, 250)
	}
}
