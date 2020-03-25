package game

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/samuelyuan/openbiohazard2/fileio"
)

const (
	PLAYER_FORWARD_SPEED  = 4000
	PLAYER_BACKWARD_SPEED = 1000
)

type Player struct {
	Position      mgl32.Vec3
	RotationAngle float32
	PoseNumber    int
}

// Position is in world space
// Rotation angle is in degrees
func NewPlayer(initialPosition mgl32.Vec3, initialRotationAngle float32) *Player {
	return &Player{
		Position:      initialPosition,
		RotationAngle: initialRotationAngle,
		PoseNumber:    -1,
	}
}

func (p *Player) GetModelMatrix() mgl32.Mat4 {
	modelMatrix := mgl32.Ident4()
	modelMatrix = modelMatrix.Mul4(mgl32.Translate3D(p.Position.X(), p.Position.Y(), p.Position.Z()))
	modelMatrix = modelMatrix.Mul4(mgl32.HomogRotate3DY(mgl32.DegToRad(float32(p.RotationAngle))))
	return modelMatrix
}

func (gameDef *GameDef) HandlePlayerInputForward(collisionEntities []fileio.CollisionEntity, timeElapsedSeconds float64) {
	predictPosition := gameDef.PredictPositionForward(gameDef.Player.Position, gameDef.Player.RotationAngle, timeElapsedSeconds)
	collidingEntity := gameDef.CheckCollision(predictPosition, collisionEntities)
	if collidingEntity == nil {
		gameDef.Player.Position = predictPosition
		gameDef.Player.PoseNumber = 0
	} else {
		if gameDef.CheckRamp(collidingEntity) {
			predictPosition := gameDef.PredictPositionForwardSlope(gameDef.Player.Position, gameDef.Player.RotationAngle, collidingEntity, timeElapsedSeconds)
			gameDef.Player.Position = predictPosition
			gameDef.Player.PoseNumber = 0
		} else {
			gameDef.Player.PoseNumber = -1
		}
	}
}

func (gameDef *GameDef) HandlePlayerInputBackward(collisionEntities []fileio.CollisionEntity, timeElapsedSeconds float64) {
	predictPosition := gameDef.PredictPositionBackward(gameDef.Player.Position, gameDef.Player.RotationAngle, timeElapsedSeconds)
	collidingEntity := gameDef.CheckCollision(predictPosition, collisionEntities)
	if collidingEntity == nil {
		gameDef.Player.Position = predictPosition
		gameDef.Player.PoseNumber = 1
	} else {
		if gameDef.CheckRamp(collidingEntity) {
			predictPosition := gameDef.PredictPositionBackwardSlope(gameDef.Player.Position, gameDef.Player.RotationAngle, collidingEntity, timeElapsedSeconds)
			gameDef.Player.Position = predictPosition
			gameDef.Player.PoseNumber = 1
		} else {
			gameDef.Player.PoseNumber = -1
		}
	}
}

func (g *GameDef) PredictPositionForward(position mgl32.Vec3, rotationAngle float32, timeElapsedSeconds float64) mgl32.Vec3 {
	modelMatrix := mgl32.Ident4()
	modelMatrix = modelMatrix.Mul4(mgl32.HomogRotate3DY(mgl32.DegToRad(rotationAngle)))
	movementDelta := modelMatrix.Mul4x1(mgl32.Vec4{PLAYER_FORWARD_SPEED * float32(timeElapsedSeconds), 0.0, 0.0, 0.0})
	return position.Add(mgl32.Vec3{movementDelta.X(), movementDelta.Y(), movementDelta.Z()})
}

func (g *GameDef) PredictPositionBackward(position mgl32.Vec3, rotationAngle float32, timeElapsedSeconds float64) mgl32.Vec3 {
	modelMatrix := mgl32.Ident4()
	modelMatrix = modelMatrix.Mul4(mgl32.HomogRotate3DY(mgl32.DegToRad(rotationAngle)))
	movementDelta := modelMatrix.Mul4x1(mgl32.Vec4{-1 * PLAYER_BACKWARD_SPEED * float32(timeElapsedSeconds), 0.0, 0.0, 0.0})
	return position.Add(mgl32.Vec3{movementDelta.X(), movementDelta.Y(), movementDelta.Z()})
}

func (gameDef *GameDef) RotatePlayerLeft(timeElapsedSeconds float64) {
	gameDef.Player.RotationAngle -= 100 * float32(timeElapsedSeconds)
	if gameDef.Player.RotationAngle < 0 {
		gameDef.Player.RotationAngle += 360
	}
}

func (gameDef *GameDef) RotatePlayerRight(timeElapsedSeconds float64) {
	gameDef.Player.RotationAngle += 100 * float32(timeElapsedSeconds)
	if gameDef.Player.RotationAngle > 360 {
		gameDef.Player.RotationAngle -= 360
	}
}

func (g *GameDef) PredictPositionForwardSlope(
	position mgl32.Vec3,
	rotationAngle float32,
	slopedEntity *fileio.CollisionEntity,
	timeElapsedSeconds float64) mgl32.Vec3 {
	modelMatrix := mgl32.Ident4()
	modelMatrix = modelMatrix.Mul4(mgl32.HomogRotate3DY(mgl32.DegToRad(rotationAngle)))
	movementDelta := modelMatrix.Mul4x1(mgl32.Vec4{PLAYER_FORWARD_SPEED * float32(timeElapsedSeconds), 0.0, 0.0, 0.0})
	predictPositionFlat := position.Add(mgl32.Vec3{movementDelta.X(), movementDelta.Y(), movementDelta.Z()})

	distanceFromRampBottom := 0.0
	if slopedEntity.SlopeType == 0 || slopedEntity.SlopeType == 1 {
		// ramp bottom is on the x-axis
		distanceFromRampBottom = math.Abs(float64(predictPositionFlat.X()-slopedEntity.RampBottom)) / float64(slopedEntity.Width)
	} else if slopedEntity.SlopeType == 2 || slopedEntity.SlopeType == 3 {
		// ramp bottom is on the z-axis
		distanceFromRampBottom = math.Abs(float64(predictPositionFlat.Z()-slopedEntity.RampBottom)) / float64(slopedEntity.Density)
	}
	predictPositionY := float64(slopedEntity.SlopeHeight) * distanceFromRampBottom
	return mgl32.Vec3{predictPositionFlat.X(), float32(predictPositionY), predictPositionFlat.Z()}
}

func (g *GameDef) PredictPositionBackwardSlope(
	position mgl32.Vec3,
	rotationAngle float32,
	slopedEntity *fileio.CollisionEntity,
	timeElapsedSeconds float64) mgl32.Vec3 {
	modelMatrix := mgl32.Ident4()
	modelMatrix = modelMatrix.Mul4(mgl32.HomogRotate3DY(mgl32.DegToRad(rotationAngle)))
	movementDelta := modelMatrix.Mul4x1(mgl32.Vec4{-1 * PLAYER_BACKWARD_SPEED * float32(timeElapsedSeconds), 0.0, 0.0, 0.0})
	predictPositionFlat := position.Add(mgl32.Vec3{movementDelta.X(), movementDelta.Y(), movementDelta.Z()})
	distanceFromRampBottom := 0.0
	if slopedEntity.SlopeType == 0 || slopedEntity.SlopeType == 1 {
		// ramp bottom is on the x-axis
		distanceFromRampBottom = math.Abs(float64(predictPositionFlat.X()-slopedEntity.RampBottom)) / float64(slopedEntity.Width)
	} else if slopedEntity.SlopeType == 2 || slopedEntity.SlopeType == 3 {
		// ramp bottom is on the z-axis
		distanceFromRampBottom = math.Abs(float64(predictPositionFlat.Z()-slopedEntity.RampBottom)) / float64(slopedEntity.Density)
	}
	predictPositionY := float64(slopedEntity.SlopeHeight) * distanceFromRampBottom
	return mgl32.Vec3{predictPositionFlat.X(), float32(predictPositionY), predictPositionFlat.Z()}
}
