package game

import (
	"../fileio"
	"fmt"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	PLAYER_LEON           = 0
	PLAYER_CLAIRE         = 1
	PLAYER_FORWARD_SPEED  = 25
	PLAYER_BACKWARD_SPEED = 15
	ROOMCUT_FILE          = "data/Common/bin/roomcut.bin"
	LEON_MODEL_FILE       = "data/Pl0/PLD/PL00.PLD"
	DIFFICULTY_EASY       = 0
	DIFFICULTY_NORMAL     = 1
)

type GameDef struct {
	StageId          int
	RoomId           int
	CameraId         int
	MaxCamerasInRoom int
	IsCameraLoaded   bool
	IsRoomLoaded     bool
	GameRoom         GameRoom
	Doors            []ScriptDoor
	Items            []ScriptItemAotSet
	Sprites          []ScriptSprite
	Player           *Player
	ScriptMemory     *ScriptMemory
	ScriptBitArray   map[int]map[int]int
	ScriptVariable   map[int]int
}

type GameRoom struct {
	CameraPositionData      []fileio.CameraInfo
	CameraSwitches          []fileio.RVDHeader
	CameraSwitchTransitions map[int][]int
	CameraMaskData          [][]fileio.MaskRectangle
	CollisionEntities       []fileio.CollisionEntity
	LightData               []fileio.LITCameraLight
	InitScriptData          fileio.ScriptFunction
	RoomScriptData          fileio.ScriptFunction
	ItemTextureData         []*fileio.TIMOutput
	ItemModelData           []*fileio.MD1Output
	SpriteData              []fileio.SpriteData
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
		Items:            make([]ScriptItemAotSet, 0),
		Sprites:          make([]ScriptSprite, 0),
		ScriptMemory:     NewScriptMemory(),
		ScriptBitArray:   make(map[int]map[int]int),
		ScriptVariable:   make(map[int]int),
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
			gameDef.Doors = make([]ScriptDoor, 0)
			gameDef.Items = make([]ScriptItemAotSet, 0)
			gameDef.Sprites = make([]ScriptSprite, 0)
			gameDef.ScriptMemory = NewScriptMemory()
		}
	}
}

func (gameDef *GameDef) LoadNewRoom(rdtOutput *fileio.RDTOutput) {
	gameDef.MaxCamerasInRoom = int(rdtOutput.Header.NumCameras)
	fmt.Println("Max cameras in room = ", gameDef.MaxCamerasInRoom)

	gameDef.GameRoom = GameRoom{}
	gameDef.GameRoom.CameraSwitches = rdtOutput.CameraSwitchData.CameraSwitches
	gameDef.GameRoom.CameraSwitchTransitions = gameDef.GenerateCameraSwitchTransitions(gameDef.GameRoom.CameraSwitches)
	gameDef.GameRoom.CameraPositionData = rdtOutput.RIDOutput.CameraPositions
	gameDef.GameRoom.CameraMaskData = rdtOutput.RIDOutput.CameraMasks
	gameDef.GameRoom.CollisionEntities = rdtOutput.CollisionData.CollisionEntities
	gameDef.GameRoom.LightData = rdtOutput.LightData.Lights
	gameDef.GameRoom.InitScriptData = rdtOutput.InitScriptData.ScriptData
	gameDef.GameRoom.RoomScriptData = rdtOutput.RoomScriptData.ScriptData
	gameDef.GameRoom.ItemTextureData = rdtOutput.ItemTextureData
	gameDef.GameRoom.ItemModelData = rdtOutput.ItemModelData
	gameDef.GameRoom.SpriteData = rdtOutput.SpriteOutput.SpriteData
	gameDef.RunScript(gameDef.GameRoom.InitScriptData, -1, true, 0)
}

func (gameDef *GameDef) GetBitArray(bitArrayIndex int, bitNumber int) int {
	bitArray, exists := gameDef.ScriptBitArray[bitArrayIndex]
	if !exists {
		gameDef.ScriptBitArray[bitArrayIndex] = make(map[int]int)
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
