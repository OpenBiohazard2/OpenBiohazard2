package game

import (
	"math"
	"testing"

	"github.com/OpenBiohazard2/OpenBiohazard2/fileio"
	"github.com/go-gl/mathgl/mgl32"
)

func floatsClose(a, b, epsilon float32) bool {
	diff := a - b
	if diff < 0 {
		diff = -diff
	}
	return diff <= epsilon
}

func TestNewPlayer(t *testing.T) {
	pos := mgl32.Vec3{10, 20, 30}
	player := NewPlayer(pos, 45)

	if player.Position != pos {
		t.Errorf("Expected Position %v, got %v", pos, player.Position)
	}
	if player.RotationAngle != 45 {
		t.Errorf("Expected RotationAngle 45, got %f", player.RotationAngle)
	}
	if player.PoseNumber != PLAYER_IDLE_POSE {
		t.Errorf("Expected PoseNumber %d, got %d", PLAYER_IDLE_POSE, player.PoseNumber)
	}
}

func TestPredictPositionForward_NoRotation(t *testing.T) {
	player := &Player{Position: mgl32.Vec3{0, 0, 0}, RotationAngle: 0}
	result := player.PredictPositionForward(1.0)

	expected := mgl32.Vec3{PLAYER_FORWARD_SPEED, 0, 0}
	if result != expected {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestPredictPositionForward_ScalesWithElapsedTime(t *testing.T) {
	player := &Player{Position: mgl32.Vec3{0, 0, 0}, RotationAngle: 0}
	result := player.PredictPositionForward(0.5)

	expected := mgl32.Vec3{PLAYER_FORWARD_SPEED * 0.5, 0, 0}
	if result != expected {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestPredictPositionForward_PreservesDistanceUnderRotation(t *testing.T) {
	player := &Player{Position: mgl32.Vec3{100, 0, -50}, RotationAngle: 37}
	result := player.PredictPositionForward(1.0)

	movement := result.Sub(player.Position)
	if !floatsClose(movement.Len(), PLAYER_FORWARD_SPEED, 0.01) {
		t.Errorf("Expected movement distance %f, got %f", float32(PLAYER_FORWARD_SPEED), movement.Len())
	}
	// A rotation around the Y axis should never move the player vertically.
	if movement.Y() != 0 {
		t.Errorf("Expected Y component of movement to be 0, got %f", movement.Y())
	}
}

func TestPredictPositionBackward_NoRotation(t *testing.T) {
	player := &Player{Position: mgl32.Vec3{0, 0, 0}, RotationAngle: 0}
	result := player.PredictPositionBackward(1.0)

	expected := mgl32.Vec3{-1 * PLAYER_BACKWARD_SPEED, 0, 0}
	if result != expected {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestRotatePlayerLeft_NoWrap(t *testing.T) {
	player := &Player{RotationAngle: 100}
	player.RotatePlayerLeft(0.5)

	expected := float32(50) // 100 - 100*0.5
	if player.RotationAngle != expected {
		t.Errorf("Expected RotationAngle %f, got %f", expected, player.RotationAngle)
	}
}

func TestRotatePlayerLeft_WrapsBelowZero(t *testing.T) {
	player := &Player{RotationAngle: 10}
	player.RotatePlayerLeft(1.0) // 10 - 100 = -90 -> wraps to 270

	expected := float32(270)
	if player.RotationAngle != expected {
		t.Errorf("Expected RotationAngle to wrap to %f, got %f", expected, player.RotationAngle)
	}
}

func TestRotatePlayerRight_NoWrap(t *testing.T) {
	player := &Player{RotationAngle: 100}
	player.RotatePlayerRight(0.5)

	expected := float32(150) // 100 + 100*0.5
	if player.RotationAngle != expected {
		t.Errorf("Expected RotationAngle %f, got %f", expected, player.RotationAngle)
	}
}

func TestRotatePlayerRight_WrapsAboveThreeSixty(t *testing.T) {
	player := &Player{RotationAngle: 350}
	player.RotatePlayerRight(1.0) // 350 + 100 = 450 -> wraps to 90

	expected := float32(90)
	if player.RotationAngle != expected {
		t.Errorf("Expected RotationAngle to wrap to %f, got %f", expected, player.RotationAngle)
	}
}

func TestPredictPositionSlope_XAxisRamp(t *testing.T) {
	player := &Player{}
	slopedEntity := &fileio.CollisionEntity{
		SlopeType:   0,
		RampBottom:  0,
		Width:       100,
		SlopeHeight: 200,
	}
	flatPosition := mgl32.Vec3{50, 0, 0}

	result := player.PredictPositionSlope(flatPosition, slopedEntity)

	expected := mgl32.Vec3{50, 100, 0} // distance = |50-0|/100 = 0.5; y = 200*0.5 = 100
	if result != expected {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestPredictPositionSlope_ZAxisRamp(t *testing.T) {
	player := &Player{}
	slopedEntity := &fileio.CollisionEntity{
		SlopeType:   2,
		RampBottom:  0,
		Density:     50,
		SlopeHeight: 100,
	}
	flatPosition := mgl32.Vec3{0, 0, 25}

	result := player.PredictPositionSlope(flatPosition, slopedEntity)

	expected := mgl32.Vec3{0, 50, 25} // distance = |25-0|/50 = 0.5; y = 100*0.5 = 50
	if result != expected {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestPredictPositionClimbBox_ClimbsUpFromGround(t *testing.T) {
	player := &Player{Position: mgl32.Vec3{10, 0, 20}}
	result := player.PredictPositionClimbBox()

	expected := mgl32.Vec3{10, fileio.FLOOR_HEIGHT_UNIT, 20}
	if result != expected {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestPredictPositionClimbBox_ClimbsDownFromBox(t *testing.T) {
	player := &Player{Position: mgl32.Vec3{10, fileio.FLOOR_HEIGHT_UNIT, 20}}
	result := player.PredictPositionClimbBox()

	expected := mgl32.Vec3{10, 0, 20}
	if result != expected {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestPredictPositionClimbBox_NoOpBeyondSecondFloor(t *testing.T) {
	pos := mgl32.Vec3{10, 2 * fileio.FLOOR_HEIGHT_UNIT, 20}
	player := &Player{Position: pos}
	result := player.PredictPositionClimbBox()

	if result != pos {
		t.Errorf("Expected position to be unchanged at %v, got %v", pos, result)
	}
}

func TestHandlePlayerInputForward_NoCollisionMovesPlayer(t *testing.T) {
	player := &Player{Position: mgl32.Vec3{0, 0, 0}, RotationAngle: 0}
	player.HandlePlayerInputForward([]fileio.CollisionEntity{}, 1.0)

	expected := mgl32.Vec3{PLAYER_FORWARD_SPEED, 0, 0}
	if player.Position != expected {
		t.Errorf("Expected Position %v, got %v", expected, player.Position)
	}
	if player.PoseNumber != PLAYER_WALKING_POSE {
		t.Errorf("Expected PoseNumber %d, got %d", PLAYER_WALKING_POSE, player.PoseNumber)
	}
}

func TestHandlePlayerInputBackward_NoCollisionMovesPlayer(t *testing.T) {
	player := &Player{Position: mgl32.Vec3{0, 0, 0}, RotationAngle: 0}
	player.HandlePlayerInputBackward([]fileio.CollisionEntity{}, 1.0)

	expected := mgl32.Vec3{-1 * PLAYER_BACKWARD_SPEED, 0, 0}
	if player.Position != expected {
		t.Errorf("Expected Position %v, got %v", expected, player.Position)
	}
}

func TestGetModelMatrix_TranslatesToPosition(t *testing.T) {
	player := &Player{Position: mgl32.Vec3{5, 6, 7}, RotationAngle: 0}
	m := player.GetModelMatrix()

	// With no rotation, the model matrix should just be a translation - the origin maps
	// to the player's position.
	transformed := m.Mul4x1(mgl32.Vec4{0, 0, 0, 1})
	expected := mgl32.Vec4{5, 6, 7, 1}
	if math.Abs(float64(transformed.X()-expected.X())) > 0.01 ||
		math.Abs(float64(transformed.Y()-expected.Y())) > 0.01 ||
		math.Abs(float64(transformed.Z()-expected.Z())) > 0.01 {
		t.Errorf("Expected origin to map to %v, got %v", expected, transformed)
	}
}
