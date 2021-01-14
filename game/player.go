package game

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/samuelyuan/openbiohazard2/fileio"
	"github.com/samuelyuan/openbiohazard2/world"
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

func (player *Player) HandlePlayerInputForward(collisionEntities []fileio.CollisionEntity, timeElapsedSeconds float64) {
	predictPosition := player.PredictPositionForward(timeElapsedSeconds)
	collidingEntity := world.CheckCollision(predictPosition, collisionEntities)
	if collidingEntity == nil {
		player.Position = predictPosition
		player.PoseNumber = 0
	} else {
		if world.CheckRamp(collidingEntity) {
			player.Position = player.PredictPositionForwardSlope(collidingEntity, timeElapsedSeconds)
			player.PoseNumber = 0
		} else if collidingEntity.Shape == 9 || collidingEntity.Shape == 10 {
			player.Position = player.PredictPositionClimbBox()
		} else {
			player.PoseNumber = -1
		}
	}
}

func (player *Player) HandlePlayerInputBackward(collisionEntities []fileio.CollisionEntity, timeElapsedSeconds float64) {
	predictPosition := player.PredictPositionBackward(timeElapsedSeconds)
	collidingEntity := world.CheckCollision(predictPosition, collisionEntities)
	if collidingEntity == nil {
		player.Position = predictPosition
		player.PoseNumber = 1
	} else {
		if world.CheckRamp(collidingEntity) {
			player.Position = player.PredictPositionBackwardSlope(collidingEntity, timeElapsedSeconds)
			player.PoseNumber = 1
		} else {
			player.PoseNumber = -1
		}
	}
}

func (player *Player) PredictPositionForward(timeElapsedSeconds float64) mgl32.Vec3 {
	modelMatrix := mgl32.Ident4()
	modelMatrix = modelMatrix.Mul4(mgl32.HomogRotate3DY(mgl32.DegToRad(player.RotationAngle)))
	movementDelta := modelMatrix.Mul4x1(mgl32.Vec4{PLAYER_FORWARD_SPEED * float32(timeElapsedSeconds), 0.0, 0.0, 0.0})
	return player.Position.Add(mgl32.Vec3{movementDelta.X(), movementDelta.Y(), movementDelta.Z()})
}

func (player *Player) PredictPositionBackward(timeElapsedSeconds float64) mgl32.Vec3 {
	modelMatrix := mgl32.Ident4()
	modelMatrix = modelMatrix.Mul4(mgl32.HomogRotate3DY(mgl32.DegToRad(player.RotationAngle)))
	movementDelta := modelMatrix.Mul4x1(mgl32.Vec4{-1 * PLAYER_BACKWARD_SPEED * float32(timeElapsedSeconds), 0.0, 0.0, 0.0})
	return player.Position.Add(mgl32.Vec3{movementDelta.X(), movementDelta.Y(), movementDelta.Z()})
}

func (player *Player) RotatePlayerLeft(timeElapsedSeconds float64) {
	player.RotationAngle -= 100 * float32(timeElapsedSeconds)
	if player.RotationAngle < 0 {
		player.RotationAngle += 360
	}
}

func (player *Player) RotatePlayerRight(timeElapsedSeconds float64) {
	player.RotationAngle += 100 * float32(timeElapsedSeconds)
	if player.RotationAngle > 360 {
		player.RotationAngle -= 360
	}
}

func (player *Player) PredictPositionForwardSlope(
	slopedEntity *fileio.CollisionEntity,
	timeElapsedSeconds float64,
) mgl32.Vec3 {
	predictPositionFlat := player.PredictPositionForward(timeElapsedSeconds)
	return player.PredictPositionSlope(predictPositionFlat, slopedEntity)
}

func (player *Player) PredictPositionBackwardSlope(
	slopedEntity *fileio.CollisionEntity,
	timeElapsedSeconds float64,
) mgl32.Vec3 {
	predictPositionFlat := player.PredictPositionBackward(timeElapsedSeconds)
	return player.PredictPositionSlope(predictPositionFlat, slopedEntity)
}

// Player walks up or down the stairs or ramp
func (player *Player) PredictPositionSlope(predictPositionFlat mgl32.Vec3, slopedEntity *fileio.CollisionEntity) mgl32.Vec3 {
	distanceFromRampBottom := 0.0

	// Check slope type orientation
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

func (player *Player) PredictPositionClimbBox() mgl32.Vec3 {
	playerFloorNum := int(math.Round(float64(player.Position.Y()) / fileio.FLOOR_HEIGHT_UNIT))

	if playerFloorNum == 0 {
		// player is on the ground
		// climb up
		return mgl32.Vec3{player.Position.X(), fileio.FLOOR_HEIGHT_UNIT, player.Position.Z()}
	} else if playerFloorNum == 1 {
		// player is on the box
		// climb down
		return mgl32.Vec3{player.Position.X(), 0.0, player.Position.Z()}
	}

	return player.Position
}

func (gameDef *GameDef) HandlePlayerActionButton(collisionEntities []fileio.CollisionEntity) {
}
