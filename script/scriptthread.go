package script

import (
	"github.com/samuelyuan/openbiohazard2/fileio"
)

type ScriptThread struct {
	RunStatus              bool
	WorkSetComponent       int
	WorkSetIndex           int
	ProgramCounter         int
	StackIndex             int
	SubLevel               int
	LevelState             []*LevelState
	Stack                  []int
	OverrideProgramCounter bool
}

type LevelState struct {
	IfElseCounter int
	LoopLevel     int
	ReturnAddress int
	LoopState     []*LoopState
}

type LoopState struct {
	Counter   int
	Break     int
	IfCounter int
	Stack     int
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
		LoopState:     loopState,
	}
}

func NewLoopState() *LoopState {
	return &LoopState{
		Counter:   0,
		Break:     0,
		IfCounter: 0,
		Stack:     0,
	}
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
		Stack:                  make([]int, 32),
		StackIndex:             0,
		SubLevel:               0,
		LevelState:             levelState,
		OverrideProgramCounter: false,
	}
}

func (thread *ScriptThread) Reset() {
	thread.RunStatus = false
	thread.ProgramCounter = 0
	for i := 0; i < len(thread.Stack); i++ {
		thread.Stack[i] = 0
	}
	thread.StackIndex = 0
	thread.SubLevel = 0

	for i := 0; i < len(thread.LevelState); i++ {
		thread.LevelState[i].IfElseCounter = 0
		thread.LevelState[i].LoopLevel = 0
		thread.LevelState[i].ReturnAddress = 0

		for j := 0; j < len(thread.LevelState[i].LoopState); j++ {
			thread.LevelState[i].LoopState[j].Counter = 0
			thread.LevelState[i].LoopState[j].Break = 0
			thread.LevelState[i].LoopState[j].IfCounter = 0
			thread.LevelState[i].LoopState[j].Stack = 0
		}
	}
	thread.LevelState[0].IfElseCounter = -1
	thread.LevelState[0].LoopLevel = -1

	thread.OverrideProgramCounter = false
}

func (thread *ScriptThread) IncrementProgramCounter(opcode byte) {
	thread.ProgramCounter += fileio.InstructionSize[opcode]
}
