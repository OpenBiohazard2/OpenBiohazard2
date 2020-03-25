package game

import (
	"fmt"

	"github.com/samuelyuan/openbiohazard2/fileio"
)

type GameRoom struct {
	CameraPositionData  []fileio.CameraInfo
	CameraSwitchHandler *CameraSwitchHandler
	CollisionEntities   []fileio.CollisionEntity
	InitScriptData      fileio.ScriptFunction
	RoomScriptData      fileio.ScriptFunction
}

func (gameDef *GameDef) NewGameRoom(rdtOutput *fileio.RDTOutput) GameRoom {
	cameraSwitches := rdtOutput.CameraSwitchData.CameraSwitches

	return GameRoom{
		CameraSwitchHandler: NewCameraSwitchHandler(cameraSwitches, gameDef.MaxCamerasInRoom),
		CameraPositionData:  rdtOutput.RIDOutput.CameraPositions,
		CollisionEntities:   rdtOutput.CollisionData.CollisionEntities,
		InitScriptData:      rdtOutput.InitScriptData.ScriptData,
		RoomScriptData:      rdtOutput.RoomScriptData.ScriptData,
	}
}

// stage starts from 1
// room number is a hex from 0
// player number is 0 or 1
func (g *GameDef) GetRoomFilename(playerNum int) string {
	stage := g.StageId
	roomNumber := g.RoomId
	return fmt.Sprintf(RDT_FILE, playerNum, stage, roomNumber, playerNum)
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
