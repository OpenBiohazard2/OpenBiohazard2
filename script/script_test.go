package script

import "testing"

// SCRIPT_FRAMES_PER_SECOND is 30.0, so the throttle threshold is 1/30 ≈ 0.0333s.

func TestShouldRunScriptFrame_AccumulatesUnderThreshold(t *testing.T) {
	newDelta, shouldRun := shouldRunScriptFrame(0.0, 0.01)

	if shouldRun {
		t.Error("Expected shouldRun to be false while under the frame threshold")
	}
	if newDelta != 0.01 {
		t.Errorf("Expected delta to accumulate to 0.01, got %f", newDelta)
	}
}

func TestShouldRunScriptFrame_RunsPastThreshold(t *testing.T) {
	newDelta, shouldRun := shouldRunScriptFrame(0.0, 0.05)

	if !shouldRun {
		t.Error("Expected shouldRun to be true once accumulated time exceeds 1/30s")
	}
	if newDelta != 0.0 {
		t.Errorf("Expected delta to reset to 0 after running a frame, got %f", newDelta)
	}
}

func TestShouldRunScriptFrame_AccumulatesAcrossMultipleCalls(t *testing.T) {
	delta := 0.0
	var shouldRun bool

	// Three calls of 0.02s: 0.02, 0.04 (already past 1/30 ≈ 0.0333, so this call runs)
	delta, shouldRun = shouldRunScriptFrame(delta, 0.02)
	if shouldRun {
		t.Fatalf("Expected shouldRun to be false after first call (delta=%f)", delta)
	}

	delta, shouldRun = shouldRunScriptFrame(delta, 0.02)
	if !shouldRun {
		t.Fatalf("Expected shouldRun to be true once accumulated delta (0.04) exceeds 1/30")
	}
	if delta != 0.0 {
		t.Errorf("Expected delta to reset to 0 after running a frame, got %f", delta)
	}
}

func TestShouldRunScriptFrame_ExactlyAtThresholdDoesNotRun(t *testing.T) {
	// The check is strictly greater-than, so landing exactly on the threshold should not run.
	threshold := 1.0 / SCRIPT_FRAMES_PER_SECOND
	newDelta, shouldRun := shouldRunScriptFrame(0.0, threshold)

	if shouldRun {
		t.Error("Expected shouldRun to be false when delta lands exactly on the threshold")
	}
	if newDelta != threshold {
		t.Errorf("Expected delta to remain %f, got %f", threshold, newDelta)
	}
}
