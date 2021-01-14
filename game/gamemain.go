package game

import (
	"fmt"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/samuelyuan/openbiohazard2/world"
)

const (
	PLAYER_LEON   = 0
	PLAYER_CLAIRE = 1

	DIFFICULTY_EASY   = 0
	DIFFICULTY_NORMAL = 1

	GAME_LOAD_ROOM   = 0
	GAME_LOAD_CAMERA = 1
	GAME_LOOP        = 2
)

type GameDef struct {
	StageId     int
	RoomId      int
	CameraId    int
	StateStatus int
	RoomScript  RoomScript
	GameWorld   *world.GameWorld
	Player      *Player
}

func NewGame(stageId int, roomId int, cameraId int) *GameDef {
	return &GameDef{
		StageId:     stageId,
		RoomId:      roomId,
		CameraId:    cameraId,
		StateStatus: GAME_LOAD_ROOM,
		GameWorld:   world.NewGameWorld(),
	}
}

func (gameDef *GameDef) ChangeCamera(newCamera int) {
	gameDef.StateStatus = GAME_LOAD_CAMERA
	gameDef.CameraId = gameDef.GameWorld.GameRoom.ClampNewCameraId(newCamera)
}

func (gameDef *GameDef) HandleCameraSwitch(position mgl32.Vec3) {
	// Check is player entered a new region
	cameraSwitchHandler := gameDef.GameWorld.GameRoom.CameraSwitchHandler
	cameraSwitchNewRegion := cameraSwitchHandler.GetCameraSwitchNewRegion(gameDef.Player.Position, gameDef.CameraId)
	if cameraSwitchNewRegion != nil {
		// Switch to a new camera
		gameDef.ChangeCamera(int(cameraSwitchNewRegion.Cam1))
	}
}

func (gameDef *GameDef) HandleRoomSwitch(position mgl32.Vec3) {
	door := gameDef.GameWorld.AotManager.GetDoorNearPlayer(position)
	if door != nil {
		// Switch to a new room
		gameDef.StageId = 1 + int(door.Stage)
		gameDef.RoomId = int(door.Room)
		gameDef.CameraId = int(door.Camera)
		gameDef.Player.Position = mgl32.Vec3{float32(door.NextX), float32(door.NextY), float32(door.NextZ)}
		fmt.Println("New player position = ", gameDef.Player.Position)

		gameDef.StateStatus = GAME_LOAD_ROOM
		gameDef.GameWorld.AotManager = world.NewAotManager()
	}
}
