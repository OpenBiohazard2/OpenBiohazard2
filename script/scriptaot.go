package script

import (
	"bytes"
	"encoding/binary"
	"log"

	"github.com/samuelyuan/openbiohazard2/fileio"
	"github.com/samuelyuan/openbiohazard2/game"
)

func (scriptDef *ScriptDef) ScriptAotSet(lineData []byte, gameDef *game.GameDef) int {
	byteArr := bytes.NewBuffer(lineData)
	instruction := fileio.ScriptInstrAotSet{}
	binary.Read(byteArr, binary.LittleEndian, &instruction)

	gameDef.AotManager.AddAotTrigger(instruction)
	return 1
}

func (scriptDef *ScriptDef) ScriptDoorAotSet(lineData []byte, gameDef *game.GameDef) int {
	byteArr := bytes.NewBuffer(lineData)
	door := fileio.ScriptInstrDoorAotSet{}
	err := binary.Read(byteArr, binary.LittleEndian, &door)
	if err != nil {
		log.Fatal("Error loading door")
	}

	if door.Id != game.AOT_DOOR {
		log.Fatal("Door has incorrect aot type ", door.Id)
	}

	gameDef.AotManager.AddDoorAot(door)
	return 1
}

func (scriptDef *ScriptDef) ScriptItemAotSet(lineData []byte, gameDef *game.GameDef) int {
	byteArr := bytes.NewBuffer(lineData)
	item := fileio.ScriptInstrItemAotSet{}
	binary.Read(byteArr, binary.LittleEndian, &item)

	if item.Id != game.AOT_ITEM {
		log.Fatal("Item has incorrect aot type ", item.Id)
	}

	gameDef.AotManager.AddItemAot(item)
	return 1
}

func (scriptDef *ScriptDef) ScriptAotReset(lineData []byte, gameDef *game.GameDef) int {
	byteArr := bytes.NewBuffer(lineData)
	instruction := fileio.ScriptInstrAotReset{}
	binary.Read(byteArr, binary.LittleEndian, &instruction)

	gameDef.AotManager.ResetAotTrigger(instruction)
	return 1
}

func (scriptDef *ScriptDef) ScriptAotSet4p(lineData []byte, gameDef *game.GameDef) int {
	byteArr := bytes.NewBuffer(lineData)
	instruction := fileio.ScriptInstrAotSet4p{}
	binary.Read(byteArr, binary.LittleEndian, &instruction)

	gameDef.AotManager.AddAotTrigger4p(instruction)
	return 1
}

func (scriptDef *ScriptDef) ScriptDoorAotSet4p(lineData []byte, gameDef *game.GameDef) int {
	byteArr := bytes.NewBuffer(lineData)
	door := fileio.ScriptInstrDoorAotSet4p{}
	err := binary.Read(byteArr, binary.LittleEndian, &door)
	if err != nil {
		log.Fatal("Error loading door")
	}

	if door.Id != game.AOT_DOOR {
		log.Fatal("Door has incorrect aot type ", door.Id)
	}

	gameDef.AotManager.AddDoorAot4p(door)
	return 1
}

func (scriptDef *ScriptDef) ScriptItemAotSet4p(lineData []byte, gameDef *game.GameDef) int {

	byteArr := bytes.NewBuffer(lineData)
	item := fileio.ScriptInstrItemAotSet4p{}
	binary.Read(byteArr, binary.LittleEndian, &item)

	if item.Id != game.AOT_ITEM {
		log.Fatal("Item has incorrect aot type ", item.Id)
	}

	gameDef.AotManager.AddItemAot4p(item)
	return 1
}
