package client

import (
	"testing"

	"github.com/go-gl/glfw/v3.2/glfw"
)

// press simulates a key press via the same callback GLFW would invoke.
func press(handler *InputHandler, key glfw.Key) {
	handler.keyCallback(nil, key, 0, glfw.Press, 0)
}

// release simulates a key release via the same callback GLFW would invoke.
func release(handler *InputHandler, key glfw.Key) {
	handler.keyCallback(nil, key, 0, glfw.Release, 0)
}

func TestIsActive_ReflectsCurrentKeyState(t *testing.T) {
	handler := NewInputHandler()

	if handler.IsActive(ACTION_BUTTON) {
		t.Error("expected ACTION_BUTTON to be inactive before any key event")
	}

	press(handler, glfw.KeyEnter)
	if !handler.IsActive(ACTION_BUTTON) {
		t.Error("expected ACTION_BUTTON to be active while KeyEnter is held")
	}

	release(handler, glfw.KeyEnter)
	if handler.IsActive(ACTION_BUTTON) {
		t.Error("expected ACTION_BUTTON to be inactive after KeyEnter is released")
	}
}

func TestIsActiveOnce_FiresOnceOnPress(t *testing.T) {
	handler := NewInputHandler()

	// Frame 1: key gets pressed after the frame's snapshot is taken.
	handler.snapshotKeyState()
	press(handler, glfw.KeyEnter)
	if !handler.IsActiveOnce(ACTION_BUTTON) {
		t.Error("expected IsActiveOnce to fire on the frame the key transitions to pressed")
	}

	// Frame 2: key is still held, no new transition - should not fire again.
	handler.snapshotKeyState()
	if handler.IsActiveOnce(ACTION_BUTTON) {
		t.Error("expected IsActiveOnce to not fire again while the key stays held")
	}

	// Frame 3: still held.
	handler.snapshotKeyState()
	if handler.IsActiveOnce(ACTION_BUTTON) {
		t.Error("expected IsActiveOnce to remain false across multiple held frames")
	}
}

func TestIsActiveOnce_FiresAgainAfterReleaseAndRepress(t *testing.T) {
	handler := NewInputHandler()

	handler.snapshotKeyState()
	press(handler, glfw.KeyEnter)
	if !handler.IsActiveOnce(ACTION_BUTTON) {
		t.Fatal("expected first press to fire IsActiveOnce")
	}

	// Release the key.
	handler.snapshotKeyState()
	release(handler, glfw.KeyEnter)
	if handler.IsActiveOnce(ACTION_BUTTON) {
		t.Error("expected IsActiveOnce to be false on release")
	}

	// Press it again - should fire once more.
	handler.snapshotKeyState()
	press(handler, glfw.KeyEnter)
	if !handler.IsActiveOnce(ACTION_BUTTON) {
		t.Error("expected IsActiveOnce to fire again on a fresh press")
	}
}

func TestIsActiveOnce_KeyAlreadyHeldBeforeFirstSnapshot(t *testing.T) {
	handler := NewInputHandler()

	// Key pressed before any snapshot has been taken (e.g. held down since before the
	// handler existed). keysPressedLast starts zero-valued (false), so this should still
	// be treated as a fresh transition on the very first snapshot/check cycle.
	press(handler, glfw.KeyEnter)
	handler.snapshotKeyState() // snapshots the pressed state as "last"

	if handler.IsActiveOnce(ACTION_BUTTON) {
		t.Error("expected IsActiveOnce to be false once the pressed state has been snapshotted as last frame")
	}
}
