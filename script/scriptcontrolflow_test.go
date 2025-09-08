package script

import (
	"testing"

	"github.com/OpenBiohazard2/OpenBiohazard2/fileio"
)

// Helper function to create test line data for if statements
func createIfLineData(blockLength uint16) []byte {
	return []byte{
		fileio.OP_IF_START,
		0,                      // dummy
		byte(blockLength),      // little endian - low byte
		byte(blockLength >> 8), // little endian - high byte
	}
}

// Helper function to create test line data for else statements
func createElseLineData(blockLength uint16) []byte {
	return []byte{
		fileio.OP_ELSE_START,
		0,                      // dummy
		byte(blockLength),      // little endian - low byte
		byte(blockLength >> 8), // little endian - high byte
	}
}

func TestScriptIfBlockStart(t *testing.T) {
	scriptDef := NewScriptDef()
	scriptThread := scriptDef.ScriptThreads[0]
	scriptThread.ProgramCounter = 100

	// Test with block length 150
	lineData := createIfLineData(150)
	returnValue := scriptDef.ScriptIfBlockStart(scriptThread, lineData)

	if returnValue != INSTRUCTION_NORMAL {
		t.Errorf("Expected return value %d, got %d", INSTRUCTION_NORMAL, returnValue)
	}

	// Check that IfElseCounter was incremented (from -1 to 0)
	if scriptThread.LevelState[scriptThread.SubLevel].IfElseCounter != 0 {
		t.Errorf("Expected IfElseCounter to be 0, got %d", scriptThread.LevelState[scriptThread.SubLevel].IfElseCounter)
	}

	// Check that position was pushed to stack
	if scriptThread.StackIndex != 1 {
		t.Errorf("Expected StackIndex to be 1, got %d", scriptThread.StackIndex)
	}

	// Check that program counter wasn't changed
	if scriptThread.ProgramCounter != 100 {
		t.Errorf("Expected ProgramCounter to be 100, got %d", scriptThread.ProgramCounter)
	}
}

func TestScriptElseCheck(t *testing.T) {
	scriptDef := NewScriptDef()
	scriptThread := scriptDef.ScriptThreads[0]
	scriptThread.ProgramCounter = 100

	// Set up initial state
	scriptThread.LevelState[scriptThread.SubLevel].IfElseCounter = 1
	scriptThread.PushStack(250) // Push a return address

	// Test else check
	lineData := createElseLineData(150)
	returnValue := scriptDef.ScriptElseCheck(scriptThread, lineData)

	if returnValue != INSTRUCTION_NORMAL {
		t.Errorf("Expected return value %d, got %d", INSTRUCTION_NORMAL, returnValue)
	}

	// Check that program counter was updated to skip else block
	if scriptThread.ProgramCounter != 250 {
		t.Errorf("Expected ProgramCounter to be 250, got %d", scriptThread.ProgramCounter)
	}

	// Check that StackIndex was decremented
	if scriptThread.StackIndex != 0 {
		t.Errorf("Expected StackIndex to be 0, got %d", scriptThread.StackIndex)
	}
}

// Legacy tests for backward compatibility
func TestIfStatementTrue(t *testing.T) {
	scriptDef := NewScriptDef()
	scriptThread := scriptDef.ScriptThreads[0]
	scriptThread.ProgramCounter = 100

	lineData := createIfLineData(150)
	returnValue := scriptDef.ScriptIfBlockStart(scriptThread, lineData)

	if returnValue != INSTRUCTION_NORMAL {
		t.Errorf("Instruction return value is incorrect, got: %d, want: %d.", returnValue, INSTRUCTION_NORMAL)
	}

	if scriptThread.ProgramCounter != 100 {
		t.Errorf("Script program counter is incorrect, got: %d, want: %d.", scriptThread.ProgramCounter, 100)
	}
}

func TestIfStatementFalse(t *testing.T) {
	scriptDef := NewScriptDef()
	scriptThread := scriptDef.ScriptThreads[0]
	scriptThread.ProgramCounter = 100

	lineData := createIfLineData(150)
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

func TestElseStatementJump(t *testing.T) {
	scriptDef := NewScriptDef()
	scriptThread := scriptDef.ScriptThreads[0]
	scriptThread.ProgramCounter = 100

	// Set up initial state
	scriptThread.LevelState[scriptThread.SubLevel].IfElseCounter = 1
	scriptThread.PushStack(250)

	lineData := createElseLineData(150)
	returnValue := scriptDef.ScriptElseCheck(scriptThread, lineData)

	if returnValue != INSTRUCTION_NORMAL {
		t.Errorf("Instruction return value is incorrect, got: %d, want: %d.", returnValue, INSTRUCTION_NORMAL)
	}

	if scriptThread.ProgramCounter != 250 {
		t.Errorf("Script program counter is incorrect, got: %d, want: %d.", scriptThread.ProgramCounter, 250)
	}
}
