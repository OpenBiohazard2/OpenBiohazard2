package world

import (
	"fmt"
	"github.com/samuelyuan/openbiohazard2/fileio"
)

type GameWorld struct {
	AotManager *AotManager
	GameRoom   *Room
}

type Room struct {
	CameraPositionData  []fileio.CameraInfo
	CameraSwitchHandler *CameraSwitchHandler
	CollisionEntities   []fileio.CollisionEntity
	MaxCamerasInRoom    int
}

func NewGameWorld() *GameWorld {
	return &GameWorld{
		AotManager: NewAotManager(),
	}
}

func (gameWorld *GameWorld) LoadNewRoom(rdtOutput *fileio.RDTOutput) {
	gameWorld.GameRoom = NewRoom(rdtOutput)
}

func NewRoom(rdtOutput *fileio.RDTOutput) *Room {
	maxCamerasInRoom := int(rdtOutput.Header.NumCameras)
	fmt.Println("Max cameras in room = ", maxCamerasInRoom)

	cameraSwitches := rdtOutput.CameraSwitchData.CameraSwitches

	return &Room{
		CameraSwitchHandler: NewCameraSwitchHandler(cameraSwitches, maxCamerasInRoom),
		CameraPositionData:  rdtOutput.RIDOutput.CameraPositions,
		CollisionEntities:   rdtOutput.CollisionData.CollisionEntities,
		MaxCamerasInRoom:    maxCamerasInRoom,
	}
}

func (room *Room) ClampNewCameraId(newCameraId int) int {
	cameraId := newCameraId
	if cameraId >= room.MaxCamerasInRoom {
		cameraId = room.MaxCamerasInRoom - 1
	}
	if cameraId < 0 {
		cameraId = 0
	}
	return cameraId
}
