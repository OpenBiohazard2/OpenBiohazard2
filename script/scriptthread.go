package script

import (
	"log"

	"github.com/OpenBiohazard2/OpenBiohazard2/fileio"
)

type ScriptThread struct {
	RunStatus              bool
	WorkSetComponent       int
	WorkSetIndex           int
	ProgramCounter         int
	StackIndex             int
	SubLevel               int
	LevelState             []*LevelState
	OverrideProgramCounter bool
	FunctionIds            []int // Only used for debugging
}

type LevelState struct {
	IfElseCounter int
	LoopLevel     int
	ReturnAddress int
	Stack         []int
	LoopState     []*LoopState
}

type LoopState struct {
	Counter        int
	Break          int
	LevelIfCounter int
	StackValue     int
}

func NewLevelState() *LevelState {
	loopState := make([]*LoopState, 4)
	for i := 0; i < len(loopState); i++ {
		loopState[i] = NewLoopState()
	}

	return &LevelState{
		IfElseCounter: 0,
		LoopLevel:     0,
		ReturnAddress: 0,
		Stack:         make([]int, 8),
		LoopState:     loopState,
	}
}

func NewLoopState() *LoopState {
	return &LoopState{
		Counter:        0,
		Break:          0,
		LevelIfCounter: 0,
		StackValue:     0,
	}
}

func (loopState *LoopState) ResetLoopState() {
	loopState.Counter = 0
	loopState.Break = 0
	loopState.LevelIfCounter = 0
	loopState.StackValue = 0
}

func NewScriptThread() *ScriptThread {
	levelState := make([]*LevelState, 4)
	for i := 0; i < len(levelState); i++ {
		levelState[i] = NewLevelState()
	}
	levelState[0].IfElseCounter = -1
	levelState[0].LoopLevel = -1

	return &ScriptThread{
		RunStatus:              false,
		ProgramCounter:         0,
		StackIndex:             0,
		SubLevel:               0,
		LevelState:             levelState,
		OverrideProgramCounter: false,
		FunctionIds:            []int{-1},
	}
}

func (thread *ScriptThread) Reset() {
	thread.RunStatus = false
	thread.ProgramCounter = 0
	thread.StackIndex = 0
	thread.SubLevel = 0

	for i := 0; i < len(thread.LevelState); i++ {
		thread.LevelState[i].IfElseCounter = 0
		thread.LevelState[i].LoopLevel = 0
		thread.LevelState[i].ReturnAddress = 0
		for j := 0; j < len(thread.LevelState[i].Stack); j++ {
			thread.LevelState[i].Stack[j] = 0
		}

		for j := 0; j < len(thread.LevelState[i].LoopState); j++ {
			thread.LevelState[i].LoopState[j].ResetLoopState()
		}
	}
	thread.LevelState[0].IfElseCounter = -1
	thread.LevelState[0].LoopLevel = -1

	thread.OverrideProgramCounter = false
	thread.FunctionIds = []int{-1}
}

func (thread *ScriptThread) IncrementProgramCounter(opcode byte) {
	thread.ProgramCounter += fileio.InstructionSize[opcode]
}

func (scriptThread *ScriptThread) JumpToNextLocationOnStack() {
	scriptThread.ProgramCounter = scriptThread.PopStackTop()
	scriptThread.LevelState[scriptThread.SubLevel].IfElseCounter--
}

func (scriptThread *ScriptThread) PushStack(newPosition int) {
	scriptThread.LevelState[scriptThread.SubLevel].Stack[scriptThread.StackIndex] = newPosition
	scriptThread.StackIndex++
}

func (scriptThread *ScriptThread) PopStackTop() int {
	if scriptThread.StackIndex == 0 {
		log.Fatal("Script stack is empty")
	}

	scriptThread.StackIndex--
	return scriptThread.LevelState[scriptThread.SubLevel].Stack[scriptThread.StackIndex]
}

func (scriptThread *ScriptThread) ShouldTerminate(scriptReturnValue int) bool {
	return scriptReturnValue == INSTRUCTION_THREAD_END || scriptThread.LevelState[scriptThread.SubLevel].IfElseCounter < 0
}
