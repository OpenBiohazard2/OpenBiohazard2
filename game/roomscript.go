package game

import (
	"fmt"

	"github.com/OpenBiohazard2/OpenBiohazard2/fileio"
	"github.com/OpenBiohazard2/OpenBiohazard2/resource"
)

type RoomScript struct {
	InitScriptData fileio.ScriptFunction
	RoomScriptData fileio.ScriptFunction
}

func (gameDef *GameDef) NewRoomScript(rdtOutput *fileio.RDTOutput) RoomScript {
	return RoomScript{
		InitScriptData: rdtOutput.InitScriptData.ScriptData,
		RoomScriptData: rdtOutput.RoomScriptData.ScriptData,
	}
}

// stage starts from 1
// room number is a hex from 0
// player number is 0 or 1
func (g *GameDef) GetRoomFilename(playerNum int) string {
	stage := g.StageId
	roomNumber := g.RoomId
	return fmt.Sprintf(resource.RDT_FILE, stage, roomNumber, playerNum)
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
