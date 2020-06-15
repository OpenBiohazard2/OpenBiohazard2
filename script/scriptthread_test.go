package script

import (
	"testing"
)

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
