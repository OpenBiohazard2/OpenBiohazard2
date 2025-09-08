package script

import (
	"testing"
)

func TestNewScriptThread(t *testing.T) {
	thread := NewScriptThread()

	// Test initial state
	if thread.RunStatus != false {
		t.Errorf("Expected RunStatus to be false, got %v", thread.RunStatus)
	}
	if thread.ProgramCounter != 0 {
		t.Errorf("Expected ProgramCounter to be 0, got %d", thread.ProgramCounter)
	}
	if thread.StackIndex != 0 {
		t.Errorf("Expected StackIndex to be 0, got %d", thread.StackIndex)
	}
	if thread.SubLevel != 0 {
		t.Errorf("Expected SubLevel to be 0, got %d", thread.SubLevel)
	}
	if thread.OverrideProgramCounter != false {
		t.Errorf("Expected OverrideProgramCounter to be false, got %v", thread.OverrideProgramCounter)
	}
	if len(thread.LevelState) != 4 {
		t.Errorf("Expected 4 LevelStates, got %d", len(thread.LevelState))
	}
	if len(thread.FunctionIds) != 1 || thread.FunctionIds[0] != -1 {
		t.Errorf("Expected FunctionIds to be [-1], got %v", thread.FunctionIds)
	}
}

func TestReset(t *testing.T) {
	thread := NewScriptThread()

	// Modify state
	thread.RunStatus = true
	thread.ProgramCounter = 100
	thread.StackIndex = 5
	thread.SubLevel = 2
	thread.OverrideProgramCounter = true
	thread.FunctionIds = []int{1, 2, 3}

	// Reset
	thread.Reset()

	// Verify reset state
	if thread.RunStatus != false {
		t.Errorf("Expected RunStatus to be false after reset, got %v", thread.RunStatus)
	}
	if thread.ProgramCounter != 0 {
		t.Errorf("Expected ProgramCounter to be 0 after reset, got %d", thread.ProgramCounter)
	}
	if thread.StackIndex != 0 {
		t.Errorf("Expected StackIndex to be 0 after reset, got %d", thread.StackIndex)
	}
	if thread.SubLevel != 0 {
		t.Errorf("Expected SubLevel to be 0 after reset, got %d", thread.SubLevel)
	}
	if thread.OverrideProgramCounter != false {
		t.Errorf("Expected OverrideProgramCounter to be false after reset, got %v", thread.OverrideProgramCounter)
	}
	if len(thread.FunctionIds) != 1 || thread.FunctionIds[0] != -1 {
		t.Errorf("Expected FunctionIds to be [-1] after reset, got %v", thread.FunctionIds)
	}
}

func TestPushAndPopStack(t *testing.T) {
	scriptThread := NewScriptThread()
	scriptThread.PushStack(100)
	scriptThread.PushStack(200)

	newPosition := scriptThread.PopStackTop()
	if newPosition != 200 {
		t.Errorf("Stack pop was incorrect, got: %d, want: %d.", newPosition, 200)
	}

	newPosition = scriptThread.PopStackTop()
	if newPosition != 100 {
		t.Errorf("Stack pop was incorrect, got: %d, want: %d.", newPosition, 100)
	}
}

func TestPushStack(t *testing.T) {
	thread := NewScriptThread()

	// Test pushing multiple values
	thread.PushStack(100)
	if thread.StackIndex != 1 {
		t.Errorf("Expected StackIndex to be 1, got %d", thread.StackIndex)
	}
	if thread.LevelState[0].Stack[0] != 100 {
		t.Errorf("Expected stack[0] to be 100, got %d", thread.LevelState[0].Stack[0])
	}

	thread.PushStack(200)
	if thread.StackIndex != 2 {
		t.Errorf("Expected StackIndex to be 2, got %d", thread.StackIndex)
	}
	if thread.LevelState[0].Stack[1] != 200 {
		t.Errorf("Expected stack[1] to be 200, got %d", thread.LevelState[0].Stack[1])
	}
}

func TestIncrementProgramCounter(t *testing.T) {
	thread := NewScriptThread()
	thread.ProgramCounter = 100

	// Test incrementing with opcode size 2 (0x10 appears to have size 2)
	thread.IncrementProgramCounter(0x10)
	if thread.ProgramCounter != 102 {
		t.Errorf("Expected ProgramCounter to be 102, got %d", thread.ProgramCounter)
	}

	// Test incrementing again
	thread.IncrementProgramCounter(0x10)
	if thread.ProgramCounter != 104 {
		t.Errorf("Expected ProgramCounter to be 104, got %d", thread.ProgramCounter)
	}
}

func TestJumpToNextLocationOnStack(t *testing.T) {
	thread := NewScriptThread()

	// Set up stack
	thread.PushStack(500)
	thread.PushStack(300)
	thread.LevelState[0].IfElseCounter = 2

	// Jump to next location
	thread.JumpToNextLocationOnStack()

	if thread.ProgramCounter != 300 {
		t.Errorf("Expected ProgramCounter to be 300, got %d", thread.ProgramCounter)
	}
	if thread.LevelState[0].IfElseCounter != 1 {
		t.Errorf("Expected IfElseCounter to be 1, got %d", thread.LevelState[0].IfElseCounter)
	}
}

func TestShouldTerminate(t *testing.T) {
	thread := NewScriptThread()

	// Test thread end - should terminate
	if !thread.ShouldTerminate(INSTRUCTION_THREAD_END) {
		t.Error("Expected to terminate with INSTRUCTION_THREAD_END")
	}

	// Test normal case with positive IfElseCounter - should not terminate
	thread.LevelState[0].IfElseCounter = 1
	if thread.ShouldTerminate(INSTRUCTION_NORMAL) {
		t.Error("Expected not to terminate with positive IfElseCounter")
	}

	// Test negative IfElseCounter - should terminate
	thread.LevelState[0].IfElseCounter = -1
	if !thread.ShouldTerminate(INSTRUCTION_NORMAL) {
		t.Error("Expected to terminate with negative IfElseCounter")
	}
}

func TestStackOperations(t *testing.T) {
	thread := NewScriptThread()

	// Test multiple push/pop operations
	values := []int{100, 200, 300, 400, 500}

	// Push all values
	for _, value := range values {
		thread.PushStack(value)
	}

	// Verify stack index
	if thread.StackIndex != len(values) {
		t.Errorf("Expected StackIndex to be %d, got %d", len(values), thread.StackIndex)
	}

	// Pop all values (in reverse order)
	for i := len(values) - 1; i >= 0; i-- {
		value := thread.PopStackTop()
		if value != values[i] {
			t.Errorf("Expected popped value to be %d, got %d", values[i], value)
		}
	}

	// Verify stack is empty
	if thread.StackIndex != 0 {
		t.Errorf("Expected StackIndex to be 0 after popping all values, got %d", thread.StackIndex)
	}
}

func TestSubLevelOperations(t *testing.T) {
	thread := NewScriptThread()

	// Test initial sublevel
	if thread.SubLevel != 0 {
		t.Errorf("Expected initial SubLevel to be 0, got %d", thread.SubLevel)
	}

	// Test changing sublevel
	thread.SubLevel = 2
	if thread.SubLevel != 2 {
		t.Errorf("Expected SubLevel to be 2, got %d", thread.SubLevel)
	}
}

func TestWorkSetOperations(t *testing.T) {
	thread := NewScriptThread()

	// Test setting WorkSetComponent
	thread.WorkSetComponent = WORKSET_PLAYER
	if thread.WorkSetComponent != WORKSET_PLAYER {
		t.Errorf("Expected WorkSetComponent to be %d, got %d", WORKSET_PLAYER, thread.WorkSetComponent)
	}

	// Test setting WorkSetIndex
	thread.WorkSetIndex = 5
	if thread.WorkSetIndex != 5 {
		t.Errorf("Expected WorkSetIndex to be 5, got %d", thread.WorkSetIndex)
	}
}

func TestFunctionIds(t *testing.T) {
	thread := NewScriptThread()

	// Test initial state (should have -1)
	if len(thread.FunctionIds) != 1 || thread.FunctionIds[0] != -1 {
		t.Errorf("Expected initial FunctionIds to be [-1], got %v", thread.FunctionIds)
	}

	// Test adding function IDs
	thread.FunctionIds = append(thread.FunctionIds, 1, 2, 3)
	if len(thread.FunctionIds) != 4 {
		t.Errorf("Expected 4 FunctionIds, got %d", len(thread.FunctionIds))
	}

	// Test removing function ID
	thread.FunctionIds = thread.FunctionIds[:len(thread.FunctionIds)-1]
	if len(thread.FunctionIds) != 3 {
		t.Errorf("Expected 3 FunctionIds after removal, got %d", len(thread.FunctionIds))
	}
}
