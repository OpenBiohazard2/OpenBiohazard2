package script

import (
	"bytes"
	"encoding/binary"

	"github.com/OpenBiohazard2/OpenBiohazard2/fileio"
	"github.com/OpenBiohazard2/OpenBiohazard2/game"
	"github.com/OpenBiohazard2/OpenBiohazard2/render"
)

func (scriptDef *ScriptDef) ScriptSceEsprOn(lineData []byte, gameDef *game.GameDef, renderDef *render.RenderDef) int {
	byteArr := bytes.NewBuffer(lineData)
	scriptSprite := fileio.ScriptInstrSceEsprOn{}
	binary.Read(byteArr, binary.LittleEndian, &scriptSprite)

	gameDef.GameWorld.AotManager.AddScriptSprite(scriptSprite)
	renderDef.AddSprite(scriptSprite)
	return 1
}

func (scriptDef *ScriptDef) ScriptSceEsprKill(lineData []byte) int {
	byteArr := bytes.NewBuffer(lineData)
	instruction := fileio.ScriptInstrSceEsprKill{}
	binary.Read(byteArr, binary.LittleEndian, &instruction)

	// TODO: implement

	return 1
}
