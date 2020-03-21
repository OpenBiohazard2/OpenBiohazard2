package game

import (
	"fmt"
	"github.com/go-gl/mathgl/mgl32"
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
	StageId          int
	RoomId           int
	CameraId         int
	MaxCamerasInRoom int
	StateStatus      int
	GameRoom         GameRoom
	AotManager       *AotManager
	Player           *Player
	ScriptBitArray   map[int]map[int]int
	ScriptVariable   map[int]int
}

func NewGame(stageId int, roomId int, cameraId int) *GameDef {
	return &GameDef{
		StageId:          stageId,
		RoomId:           roomId,
		CameraId:         cameraId,
		MaxCamerasInRoom: 0,
		StateStatus:      GAME_LOAD_ROOM,
		AotManager:       NewAotManager(),
		ScriptBitArray:   make(map[int]map[int]int),
		ScriptVariable:   make(map[int]int),
	}
}

func (gameDef *GameDef) ChangeCamera(newCamera int) {
	gameDef.StateStatus = GAME_LOAD_CAMERA
	gameDef.CameraId = newCamera
	if gameDef.CameraId >= gameDef.MaxCamerasInRoom {
		gameDef.CameraId = gameDef.MaxCamerasInRoom - 1
	}
	if gameDef.CameraId < 0 {
		gameDef.CameraId = 0
	}
}

func (gameDef *GameDef) HandleCameraSwitch(position mgl32.Vec3) {
	// Check is player entered a new region
	cameraSwitchHandler := gameDef.GameRoom.CameraSwitchHandler
	cameraSwitchNewRegion := cameraSwitchHandler.GetCameraSwitchNewRegion(gameDef.Player.Position, gameDef.CameraId)
	if cameraSwitchNewRegion != nil {
		// Switch to a new camera
		gameDef.ChangeCamera(int(cameraSwitchNewRegion.Cam1))
	}
}

func (gameDef *GameDef) HandleRoomSwitch(position mgl32.Vec3) {
	door := gameDef.AotManager.GetDoorNearPlayer(position)
	if door != nil {
		// Switch to a new room
		gameDef.StageId = 1 + int(door.Stage)
		gameDef.RoomId = int(door.Room)
		gameDef.CameraId = int(door.Camera)
		gameDef.Player.Position = mgl32.Vec3{float32(door.NextX), float32(door.NextY), float32(door.NextZ)}
		fmt.Println("New player position = ", gameDef.Player.Position)

		gameDef.StateStatus = GAME_LOAD_ROOM
		gameDef.AotManager = NewAotManager()
	}
}

func (gameDef *GameDef) GetBitArray(bitArrayIndex int, bitNumber int) int {
	bitArray, exists := gameDef.ScriptBitArray[bitArrayIndex]
	if !exists {
		gameDef.ScriptBitArray[bitArrayIndex] = make(map[int]int)
		bitArray = gameDef.ScriptBitArray[bitArrayIndex]
		fmt.Println("Initialize bit array index", bitArrayIndex)
	}
	value, exists := bitArray[bitNumber]
	if !exists {
		bitArray[bitNumber] = 0
		fmt.Println("Initialize bit array", bitArrayIndex, "with bit number ", bitNumber)
	}
	return value
}

func (gameDef *GameDef) SetBitArray(bitArrayIndex int, bitNumber int, value int) {
	_, exists := gameDef.ScriptBitArray[bitArrayIndex]
	if !exists {
		gameDef.ScriptBitArray[bitArrayIndex] = make(map[int]int)
	}
	gameDef.ScriptBitArray[bitArrayIndex][bitNumber] = value
}

func (gameDef *GameDef) GetScriptVariable(id int) int {
	return gameDef.ScriptVariable[id]
}

func (gameDef *GameDef) SetScriptVariable(id int, value int) {
	gameDef.ScriptVariable[id] = value
}
