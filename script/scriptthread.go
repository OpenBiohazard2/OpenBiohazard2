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
	LoopCounter   int
	ReturnAddress int
	SleepCounter  []int
	LoopBreak     []int
	LoopIfCounter []int
}

func NewLevelState() *LevelState {
	return &LevelState{
		IfElseCounter: 0,
		LoopCounter:   0,
		ReturnAddress: 0,
		LoopBreak:     make([]int, 4),
		SleepCounter:  make([]int, 4),
		LoopIfCounter: make([]int, 4),
	}
}

func NewScriptThread() *ScriptThread {
	levelState := make([]*LevelState, 4)
	for i := 0; i < len(levelState); i++ {
		levelState[i] = NewLevelState()
	}
	levelState[0].IfElseCounter = -1
	levelState[0].LoopCounter = -1

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
		thread.LevelState[i].LoopCounter = 0
		thread.LevelState[i].ReturnAddress = 0

		for j := 0; j < len(thread.LevelState[i].LoopBreak); j++ {
			thread.LevelState[i].LoopBreak[j] = 0
		}
		for j := 0; j < len(thread.LevelState[i].SleepCounter); j++ {
			thread.LevelState[i].SleepCounter[j] = 0
		}
		for j := 0; j < len(thread.LevelState[i].LoopIfCounter); j++ {
			thread.LevelState[i].LoopIfCounter[j] = 0
		}
	}
	thread.LevelState[0].IfElseCounter = -1
	thread.LevelState[0].LoopCounter = -1

	thread.OverrideProgramCounter = false
}

func (thread *ScriptThread) IncrementProgramCounter(opcode byte) {
	thread.ProgramCounter += fileio.InstructionSize[opcode]
}
