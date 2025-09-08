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

// Helper function to create test line data for for loops
func createForLineData(count uint16, blockLength uint16) []byte {
	return []byte{
		fileio.OP_FOR,          // Opcode
		0,                      // Dummy
		byte(blockLength),      // BlockLength - little endian - low byte
		byte(blockLength >> 8), // BlockLength - little endian - high byte
		byte(count),            // Count - little endian - low byte
		byte(count >> 8),       // Count - little endian - high byte
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

// New tests for parameter-based functions
func TestScriptEndIfWithThread(t *testing.T) {
	scriptDef := NewScriptDef()
	scriptThread := scriptDef.ScriptThreads[0]

	// Set up initial state
	scriptThread.LevelState[scriptThread.SubLevel].IfElseCounter = 2
	scriptThread.StackIndex = 3

	returnValue := scriptDef.ScriptEndIf(scriptThread)

	if returnValue != INSTRUCTION_NORMAL {
		t.Errorf("Expected return value %d, got %d", INSTRUCTION_NORMAL, returnValue)
	}

	// Check that IfElseCounter was decremented
	if scriptThread.LevelState[scriptThread.SubLevel].IfElseCounter != 1 {
		t.Errorf("Expected IfElseCounter to be 1, got %d", scriptThread.LevelState[scriptThread.SubLevel].IfElseCounter)
	}

	// Check that StackIndex was decremented
	if scriptThread.StackIndex != 2 {
		t.Errorf("Expected StackIndex to be 2, got %d", scriptThread.StackIndex)
	}
}

func TestScriptForLoopBeginWithThread(t *testing.T) {
	scriptDef := NewScriptDef()
	scriptThread := scriptDef.ScriptThreads[0]
	scriptThread.ProgramCounter = 100

	// Test with count > 0
	lineData := createForLineData(5, 50)
	returnValue := scriptDef.ScriptForLoopBegin(scriptThread, lineData)

	if returnValue != INSTRUCTION_NORMAL {
		t.Errorf("Expected return value %d, got %d", INSTRUCTION_NORMAL, returnValue)
	}

	// Check that loop level was incremented (from -1 to 0)
	curLevelState := scriptThread.LevelState[scriptThread.SubLevel]
	if curLevelState.LoopLevel != 0 {
		t.Errorf("Expected LoopLevel to be 0 (incremented from -1), got %d", curLevelState.LoopLevel)
	}

	// Check that loop state was set up correctly
	loopState := curLevelState.LoopState[curLevelState.LoopLevel]
	if loopState.Counter != 5 {
		t.Errorf("Expected Counter to be 5, got %d", loopState.Counter)
	}

	// Check that program counter was updated
	expectedPC := 100 + fileio.InstructionSize[fileio.OP_FOR]
	if scriptThread.ProgramCounter != expectedPC {
		t.Errorf("Expected ProgramCounter to be %d, got %d", expectedPC, scriptThread.ProgramCounter)
	}

	// Check that OverrideProgramCounter was set
	if !scriptThread.OverrideProgramCounter {
		t.Error("Expected OverrideProgramCounter to be true")
	}
}

func TestScriptForLoopEndWithThread(t *testing.T) {
	scriptDef := NewScriptDef()
	scriptThread := scriptDef.ScriptThreads[0]
	scriptThread.ProgramCounter = 200

	// Set up loop state
	curLevelState := scriptThread.LevelState[scriptThread.SubLevel]
	curLevelState.LoopLevel = 0 // Start from 0 (incremented from -1)
	loopState := curLevelState.LoopState[0]
	loopState.Counter = 3
	loopState.StackValue = 100

	lineData := []byte{fileio.OP_FOR_END}
	returnValue := scriptDef.ScriptForLoopEnd(scriptThread, lineData)

	if returnValue != INSTRUCTION_NORMAL {
		t.Errorf("Expected return value %d, got %d", INSTRUCTION_NORMAL, returnValue)
	}

	// Check that counter was decremented
	if loopState.Counter != 2 {
		t.Errorf("Expected Counter to be 2, got %d", loopState.Counter)
	}

	// Check that program counter was set to loop start
	if scriptThread.ProgramCounter != 100 {
		t.Errorf("Expected ProgramCounter to be 100, got %d", scriptThread.ProgramCounter)
	}

	// Check that OverrideProgramCounter was set
	if !scriptThread.OverrideProgramCounter {
		t.Error("Expected OverrideProgramCounter to be true")
	}
}

func TestScriptSwitchEndWithThread(t *testing.T) {
	scriptDef := NewScriptDef()
	scriptThread := scriptDef.ScriptThreads[0]

	// Set up initial state
	curLevelState := scriptThread.LevelState[scriptThread.SubLevel]
	curLevelState.LoopLevel = 2

	returnValue := scriptDef.ScriptSwitchEnd(scriptThread)

	if returnValue != INSTRUCTION_NORMAL {
		t.Errorf("Expected return value %d, got %d", INSTRUCTION_NORMAL, returnValue)
	}

	// Check that loop level was decremented
	if curLevelState.LoopLevel != 1 {
		t.Errorf("Expected LoopLevel to be 1, got %d", curLevelState.LoopLevel)
	}
}

func TestScriptBreakWithThread(t *testing.T) {
	scriptDef := NewScriptDef()
	scriptThread := scriptDef.ScriptThreads[0]
	scriptThread.ProgramCounter = 100

	// Set up loop state
	curLevelState := scriptThread.LevelState[scriptThread.SubLevel]
	curLevelState.LoopLevel = 0 // Start from 0 (incremented from -1)
	loopState := curLevelState.LoopState[0]
	loopState.Break = 250

	lineData := []byte{fileio.OP_BREAK}
	returnValue := scriptDef.ScriptBreak(scriptThread, lineData)

	if returnValue != INSTRUCTION_NORMAL {
		t.Errorf("Expected return value %d, got %d", INSTRUCTION_NORMAL, returnValue)
	}

	// Check that program counter was set to break address
	if scriptThread.ProgramCounter != 250 {
		t.Errorf("Expected ProgramCounter to be 250, got %d", scriptThread.ProgramCounter)
	}

	// Check that OverrideProgramCounter was set
	if !scriptThread.OverrideProgramCounter {
		t.Error("Expected OverrideProgramCounter to be true")
	}
}
