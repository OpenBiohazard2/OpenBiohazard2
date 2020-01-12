package game

import (
	"../fileio"
	"fmt"
	"github.com/go-gl/mathgl/mgl32"
	"math"
)

const (
	PLAYER_LEON           = 0
	PLAYER_CLAIRE         = 1
	PLAYER_FORWARD_SPEED  = 25
	PLAYER_BACKWARD_SPEED = 15
	ROOMCUT_FILE          = "data/Common/bin/roomcut.bin"
	LEON_MODEL_FILE       = "data/Pl0/PLD/PL00.PLD"
)

type GameDef struct {
	StageId          int
	RoomId           int
	CameraId         int
	MaxCamerasInRoom int
	IsCameraLoaded   bool
	IsRoomLoaded     bool
	Doors            []ScriptDoor
	Player           *Player
}

func NewGame(stageId int, roomId int, cameraId int) *GameDef {
	return &GameDef{
		StageId:          stageId,
		RoomId:           roomId,
		CameraId:         cameraId,
		MaxCamerasInRoom: 0,
		IsCameraLoaded:   false,
		IsRoomLoaded:     false,
		Doors:            make([]ScriptDoor, 0),
	}
}

// stage starts from 1
// room number is a hex from 0
// player number is 0 or 1
func (g *GameDef) GetRoomFilename(playerNum int) string {
	stage := g.StageId
	roomNumber := g.RoomId
	return fmt.Sprintf("data/Pl%v/Rdu/ROOM%01d%02x%01d.RDT", playerNum, stage, roomNumber, playerNum)
}

func (g *GameDef) GetBackgroundImageNumber() int {
	stage := g.StageId
	roomNumber := g.RoomId
	cameraNum := g.CameraId
	return ((stage - 1) * 512) + (roomNumber * 16) + cameraNum
}

func (g *GameDef) PredictPositionForward(position mgl32.Vec3, rotationAngle float32) mgl32.Vec3 {
	modelMatrix := mgl32.Ident4()
	modelMatrix = modelMatrix.Mul4(mgl32.HomogRotate3DY(mgl32.DegToRad(rotationAngle)))
	movementDelta := modelMatrix.Mul4x1(mgl32.Vec4{PLAYER_FORWARD_SPEED, 0.0, 0.0, 0.0})
	return position.Add(mgl32.Vec3{movementDelta.X(), movementDelta.Y(), movementDelta.Z()})
}

func (g *GameDef) PredictPositionBackward(position mgl32.Vec3, rotationAngle float32) mgl32.Vec3 {
	modelMatrix := mgl32.Ident4()
	modelMatrix = modelMatrix.Mul4(mgl32.HomogRotate3DY(mgl32.DegToRad(rotationAngle)))
	movementDelta := modelMatrix.Mul4x1(mgl32.Vec4{-1 * PLAYER_BACKWARD_SPEED, 0.0, 0.0, 0.0})
	return position.Add(mgl32.Vec3{movementDelta.X(), movementDelta.Y(), movementDelta.Z()})
}

func (g *GameDef) PredictPositionForwardSlope(position mgl32.Vec3, rotationAngle float32, slopedEntity *fileio.CollisionEntity) mgl32.Vec3 {
	modelMatrix := mgl32.Ident4()
	modelMatrix = modelMatrix.Mul4(mgl32.HomogRotate3DY(mgl32.DegToRad(rotationAngle)))
	movementDelta := modelMatrix.Mul4x1(mgl32.Vec4{PLAYER_FORWARD_SPEED, 0.0, 0.0, 0.0})
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

func (g *GameDef) PredictPositionBackwardSlope(position mgl32.Vec3, rotationAngle float32, slopedEntity *fileio.CollisionEntity) mgl32.Vec3 {
	modelMatrix := mgl32.Ident4()
	modelMatrix = modelMatrix.Mul4(mgl32.HomogRotate3DY(mgl32.DegToRad(rotationAngle)))
	movementDelta := modelMatrix.Mul4x1(mgl32.Vec4{-1 * PLAYER_BACKWARD_SPEED, 0.0, 0.0, 0.0})
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

func (gameDef *GameDef) NextRoom() {
	gameDef.CameraId = 0
	gameDef.RoomId++
	if gameDef.RoomId >= 32 {
		gameDef.RoomId = 31
	}
}

func (gameDef *GameDef) PrevRoom() {
	gameDef.CameraId = 0
	gameDef.RoomId--
	if gameDef.RoomId < 0 {
		gameDef.RoomId = 0
	}
}

// Shows which regions are reachable from the current camera
// The key is the camera id
// The value is an array of switches that are reachable
func (gameDef *GameDef) GenerateCameraSwitchTransitions(cameraSwitches []fileio.RVDHeader) map[int][]int {
	cameraSwitchTransitions := make(map[int][]int, 0)
	for roomCameraId := 0; roomCameraId < gameDef.MaxCamerasInRoom; roomCameraId++ {
		cam1ZeroIndices := make([]int, 0)
		checkSwitchesIndices := make([]int, 0)
		for switchIndex, cameraSwitch := range cameraSwitches {
			// Cam0 is the current camera
			if int(cameraSwitch.Cam0) == roomCameraId {
				// The first cam1 = 0 is used for a different purpose
				// The second cam1 = 0 is the real camera switch
				if int(cameraSwitch.Cam1) == 0 {
					cam1ZeroIndices = append(cam1ZeroIndices, switchIndex)
				} else {
					checkSwitchesIndices = append(checkSwitchesIndices, switchIndex)
				}
			}
		}

		if len(cam1ZeroIndices) >= 2 {
			transitionRegion := cam1ZeroIndices[len(cam1ZeroIndices)-1]
			checkSwitchesIndices = append(checkSwitchesIndices, transitionRegion)
		}

		cameraSwitchTransitions[roomCameraId] = checkSwitchesIndices
	}
	return cameraSwitchTransitions
}

func (gameDef *GameDef) HandleCameraSwitch(position mgl32.Vec3, cameraSwitches []fileio.RVDHeader,
	cameraSwitchTransitions map[int][]int) {
	for _, regionIndex := range cameraSwitchTransitions[gameDef.CameraId] {
		region := cameraSwitches[regionIndex]
		corner1 := mgl32.Vec3{float32(region.X1), 0, float32(region.Z1)}
		corner2 := mgl32.Vec3{float32(region.X2), 0, float32(region.Z2)}
		corner3 := mgl32.Vec3{float32(region.X3), 0, float32(region.Z3)}
		corner4 := mgl32.Vec3{float32(region.X4), 0, float32(region.Z4)}

		if isPointInRectangle(position, corner1, corner2, corner3, corner4) {
			// Switch to a new camera
			gameDef.CameraId = int(region.Cam1)
			gameDef.IsCameraLoaded = false

			if gameDef.CameraId >= gameDef.MaxCamerasInRoom {
				gameDef.CameraId = gameDef.MaxCamerasInRoom - 1
			}
			if gameDef.CameraId < 0 {
				gameDef.CameraId = 0
			}
		}
	}
}

func (gameDef *GameDef) HandleRoomSwitch(position mgl32.Vec3) {
	for _, door := range gameDef.Doors {
		corner1 := mgl32.Vec3{float32(door.X), 0, float32(door.Y)}
		corner2 := mgl32.Vec3{float32(door.X), 0, float32(door.Y + door.Height)}
		corner3 := mgl32.Vec3{float32(door.X + door.Width), 0, float32(door.Y + door.Height)}
		corner4 := mgl32.Vec3{float32(door.X + door.Width), 0, float32(door.Y)}
		if isPointInRectangle(position, corner1, corner2, corner3, corner4) {
			// Switch to a new room
			gameDef.StageId = 1 + int(door.Stage)
			gameDef.RoomId = int(door.Room)
			gameDef.CameraId = int(door.Camera)
			gameDef.Player.Position = mgl32.Vec3{float32(door.NextX), float32(door.NextY), float32(door.NextZ)}

			gameDef.IsRoomLoaded = false
			gameDef.IsCameraLoaded = false
		}
	}
}

func (gameDef *GameDef) CheckCollision(newPosition mgl32.Vec3, collisionEntities []fileio.CollisionEntity) *fileio.CollisionEntity {
	for _, entity := range collisionEntities {
		switch entity.Shape {
		case 0:
			// Rectangle
			corner1 := mgl32.Vec3{float32(entity.X), 0, float32(entity.Z)}
			corner2 := mgl32.Vec3{float32(entity.X), 0, float32(entity.Z) + float32(entity.Density)}
			corner3 := mgl32.Vec3{float32(entity.X) + float32(entity.Width), 0, float32(entity.Z) + float32(entity.Density)}
			corner4 := mgl32.Vec3{float32(entity.X) + float32(entity.Width), 0, float32(entity.Z)}
			if isPointInRectangle(newPosition, corner1, corner2, corner3, corner4) {
				return &entity
			}
		case 1:
			// Triangle \\|
			vertex1 := mgl32.Vec3{float32(entity.X), 0, float32(entity.Z + entity.Density)}
			vertex2 := mgl32.Vec3{float32(entity.X + entity.Width), 0, float32(entity.Z + entity.Density)}
			vertex3 := mgl32.Vec3{float32(entity.X + entity.Width), 0, float32(entity.Z)}
			if isPointInTriangle(newPosition, vertex1, vertex2, vertex3) {
				return &entity
			}
		case 2:
			// Triangle |/
			vertex1 := mgl32.Vec3{float32(entity.X), 0, float32(entity.Z)}
			vertex2 := mgl32.Vec3{float32(entity.X), 0, float32(entity.Z + entity.Density)}
			vertex3 := mgl32.Vec3{float32(entity.X + entity.Width), 0, float32(entity.Z + entity.Density)}
			if isPointInTriangle(newPosition, vertex1, vertex2, vertex3) {
				return &entity
			}
		case 3:
			// Triangle /|
			vertex1 := mgl32.Vec3{float32(entity.X), 0, float32(entity.Z)}
			vertex2 := mgl32.Vec3{float32(entity.X + entity.Width), 0, float32(entity.Z + entity.Density)}
			vertex3 := mgl32.Vec3{float32(entity.X + entity.Width), 0, float32(entity.Z)}
			if isPointInTriangle(newPosition, vertex1, vertex2, vertex3) {
				return &entity
			}
		case 6:
			// Circle
			radius := float32(entity.Width) / 2.0
			center := mgl32.Vec3{float32(entity.X) + radius, 0, float32(entity.Z) + radius}
			if isPointInCircle(newPosition, center, radius) {
				return &entity
			}
		case 7:
			// Ellipse, rectangle with rounded corners on the x-axis
			majorAxis := float32(entity.Width) / 2.0
			minorAxis := float32(entity.Density) / 2.0
			center := mgl32.Vec3{float32(entity.X) + majorAxis, 0, float32(entity.Z) + minorAxis}
			if isPointInEllipseXAxisMajor(newPosition, center, majorAxis, minorAxis) {
				return &entity
			}
		case 8:
			// Ellipse, rectangle with rounded corners on the z-axis
			majorAxis := float32(entity.Density) / 2.0
			minorAxis := float32(entity.Width) / 2.0
			center := mgl32.Vec3{float32(entity.X) + minorAxis, 0, float32(entity.Z) + majorAxis}
			if isPointInEllipseZAxisMajor(newPosition, center, majorAxis, minorAxis) {
				return &entity
			}
		case 9:
			// Rectangle climb up
			corner1 := mgl32.Vec3{float32(entity.X), 0, float32(entity.Z)}
			corner2 := mgl32.Vec3{float32(entity.X), 0, float32(entity.Z) + float32(entity.Density)}
			corner3 := mgl32.Vec3{float32(entity.X) + float32(entity.Width), 0, float32(entity.Z) + float32(entity.Density)}
			corner4 := mgl32.Vec3{float32(entity.X) + float32(entity.Width), 0, float32(entity.Z)}
			if isPointInRectangle(newPosition, corner1, corner2, corner3, corner4) {
				return &entity
			}
		case 10:
			// Rectangle jump down
			corner1 := mgl32.Vec3{float32(entity.X), 0, float32(entity.Z)}
			corner2 := mgl32.Vec3{float32(entity.X), 0, float32(entity.Z) + float32(entity.Density)}
			corner3 := mgl32.Vec3{float32(entity.X) + float32(entity.Width), 0, float32(entity.Z) + float32(entity.Density)}
			corner4 := mgl32.Vec3{float32(entity.X) + float32(entity.Width), 0, float32(entity.Z)}
			if isPointInRectangle(newPosition, corner1, corner2, corner3, corner4) {
				return &entity
			}
		case 11:
			// Ramp
			corner1 := mgl32.Vec3{float32(entity.X), 0, float32(entity.Z)}
			corner2 := mgl32.Vec3{float32(entity.X), 0, float32(entity.Z) + float32(entity.Density)}
			corner3 := mgl32.Vec3{float32(entity.X) + float32(entity.Width), 0, float32(entity.Z) + float32(entity.Density)}
			corner4 := mgl32.Vec3{float32(entity.X) + float32(entity.Width), 0, float32(entity.Z)}
			if isPointInRectangle(newPosition, corner1, corner2, corner3, corner4) {
				return &entity
			}
		case 12:
			// Stairs
			corner1 := mgl32.Vec3{float32(entity.X), 0, float32(entity.Z)}
			corner2 := mgl32.Vec3{float32(entity.X), 0, float32(entity.Z) + float32(entity.Density)}
			corner3 := mgl32.Vec3{float32(entity.X) + float32(entity.Width), 0, float32(entity.Z) + float32(entity.Density)}
			corner4 := mgl32.Vec3{float32(entity.X) + float32(entity.Width), 0, float32(entity.Z)}
			if isPointInRectangle(newPosition, corner1, corner2, corner3, corner4) {
				return &entity
			}
		}
	}
	return nil
}

func (gameDef *GameDef) CheckRamp(entity *fileio.CollisionEntity) bool {
	return entity.Shape == 11 || entity.Shape == 12
}

func isPointInTriangle(point mgl32.Vec3, corner1 mgl32.Vec3, corner2 mgl32.Vec3, corner3 mgl32.Vec3) bool {
	// area of triangle ABC
	area := triangleArea(corner1, corner2, corner3)
	// area of PBC
	area1 := triangleArea(point, corner2, corner3)
	// area of APC
	area2 := triangleArea(corner1, point, corner3)
	// area of ABP
	area3 := triangleArea(corner1, corner2, point)

	// areas should be equal if point is in triangle
	areaDifference := area - (area1 + area2 + area3)
	return math.Abs(float64(areaDifference)) <= 0.01
}

// Find the area of triangle formed by p1, p2 and p3
func triangleArea(p1 mgl32.Vec3, p2 mgl32.Vec3, p3 mgl32.Vec3) float32 {
	return float32(math.Abs(float64((p1.X()*(p2.Z()-p3.Z()) + p2.X()*(p3.Z()-p1.Z()) + p3.X()*(p1.Z()-p2.Z())) / 2.0)))
}

func isPointInRectangle(point mgl32.Vec3, corner1 mgl32.Vec3, corner2 mgl32.Vec3, corner3 mgl32.Vec3, corner4 mgl32.Vec3) bool {
	x := point.X()
	z := point.Z()
	x1 := corner1.X()
	z1 := corner1.Z()

	x2 := corner2.X()
	z2 := corner2.Z()

	x3 := corner3.X()
	z3 := corner3.Z()

	x4 := corner4.X()
	z4 := corner4.Z()

	a := (x2-x1)*(z-z1) - (z2-z1)*(x-x1)
	b := (x3-x2)*(z-z2) - (z3-z2)*(x-x2)
	c := (x4-x3)*(z-z3) - (z4-z3)*(x-x3)
	d := (x1-x4)*(z-z4) - (z1-z4)*(x-x4)

	if (a > 0 && b > 0 && c > 0 && d > 0) ||
		(a < 0 && b < 0 && c < 0 && d < 0) {
		return true
	}
	return false
}

func isPointInCircle(point mgl32.Vec3, circleCenter mgl32.Vec3, radius float32) bool {
	distance := point.Sub(circleCenter).Len()
	return distance <= radius
}

func isPointInEllipseXAxisMajor(point mgl32.Vec3, ellipseCenter mgl32.Vec3, majorAxis float32, minorAxis float32) bool {
	xDistance := math.Pow(float64(point.X()-ellipseCenter.X()), 2) / float64(majorAxis*majorAxis)
	zDistance := math.Pow(float64(point.Z()-ellipseCenter.Z()), 2) / float64(minorAxis*minorAxis)
	return xDistance+zDistance <= 1.0
}

func isPointInEllipseZAxisMajor(point mgl32.Vec3, ellipseCenter mgl32.Vec3, majorAxis float32, minorAxis float32) bool {
	xDistance := math.Pow(float64(point.X()-ellipseCenter.X()), 2) / float64(minorAxis*minorAxis)
	zDistance := math.Pow(float64(point.Z()-ellipseCenter.Z()), 2) / float64(majorAxis*majorAxis)
	return xDistance+zDistance <= 1.0
}
