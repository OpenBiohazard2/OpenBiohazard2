package render

import "testing"

func TestAdvanceSpriteFrame_AccumulatesUnderThreshold(t *testing.T) {
	newRuntime, newFrame := advanceSpriteFrame(0.0, 0, 4, 0.1)

	if newRuntime != 0.1 {
		t.Errorf("Expected runtime to accumulate to 0.1, got %f", newRuntime)
	}
	if newFrame != 0 {
		t.Errorf("Expected frame to stay at 0 while under SPRITE_FRAME_TIME, got %d", newFrame)
	}
}

func TestAdvanceSpriteFrame_AdvancesPastThreshold(t *testing.T) {
	// SPRITE_FRAME_TIME is 0.5 seconds
	newRuntime, newFrame := advanceSpriteFrame(0.0, 0, 4, 0.6)

	if newRuntime != 0.0 {
		t.Errorf("Expected runtime to reset to 0 after crossing the threshold, got %f", newRuntime)
	}
	if newFrame != 1 {
		t.Errorf("Expected frame to advance to 1, got %d", newFrame)
	}
}

func TestAdvanceSpriteFrame_WrapsAroundAtLastFrame(t *testing.T) {
	numFrames := 4
	// Already on the last frame; crossing the threshold should wrap back to 0.
	newRuntime, newFrame := advanceSpriteFrame(0.0, numFrames-1, numFrames, 0.6)

	if newRuntime != 0.0 {
		t.Errorf("Expected runtime to reset to 0, got %f", newRuntime)
	}
	if newFrame != 0 {
		t.Errorf("Expected frame to wrap around to 0, got %d", newFrame)
	}
}

func TestAdvanceSpriteFrame_SingleFrameAnimation(t *testing.T) {
	// A 1-frame animation should always wrap back to frame 0.
	newRuntime, newFrame := advanceSpriteFrame(0.0, 0, 1, 0.6)

	if newRuntime != 0.0 {
		t.Errorf("Expected runtime to reset to 0, got %f", newRuntime)
	}
	if newFrame != 0 {
		t.Errorf("Expected single-frame animation to stay at frame 0, got %d", newFrame)
	}
}

func TestAdvanceSpriteFrame_AccumulatesAcrossMultipleCalls(t *testing.T) {
	runtime, frame := 0.0, 0

	// Three calls of 0.2s each: 0.2, 0.4 (still under 0.5), then 0.6 crosses the threshold.
	runtime, frame = advanceSpriteFrame(runtime, frame, 4, 0.2)
	if frame != 0 {
		t.Fatalf("Expected frame to remain 0 after first call, got %d", frame)
	}

	runtime, frame = advanceSpriteFrame(runtime, frame, 4, 0.2)
	if frame != 0 {
		t.Fatalf("Expected frame to remain 0 after second call, got %d", frame)
	}

	runtime, frame = advanceSpriteFrame(runtime, frame, 4, 0.2)
	if runtime != 0.0 {
		t.Errorf("Expected runtime to reset after crossing threshold, got %f", runtime)
	}
	if frame != 1 {
		t.Errorf("Expected frame to advance to 1 after accumulated time crosses threshold, got %d", frame)
	}
}
