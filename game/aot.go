package game

import (
	"../fileio"
	"fmt"
	"github.com/go-gl/mathgl/mgl32"
)

// Handle script doors, items, events

const (
	AOT_EVENT = 5
)

type AotManager struct {
	Doors       []fileio.ScriptInstrDoorAotSet
	Items       []fileio.ScriptInstrItemAotSet
	Sprites     []fileio.ScriptSprite
	AotTriggers []fileio.ScriptInstrAotSet
}

func NewAotManager() *AotManager {
	return &AotManager{
		Doors:       make([]fileio.ScriptInstrDoorAotSet, 0),
		Items:       make([]fileio.ScriptInstrItemAotSet, 0),
		Sprites:     make([]fileio.ScriptSprite, 0),
		AotTriggers: make([]fileio.ScriptInstrAotSet, 0),
	}
}

func (aotManager *AotManager) GetDoorNearPlayer(position mgl32.Vec3) *fileio.ScriptInstrDoorAotSet {
	for _, door := range aotManager.Doors {
		corner1 := mgl32.Vec3{float32(door.X), 0, float32(door.Y)}
		corner2 := mgl32.Vec3{float32(door.X), 0, float32(door.Y + door.Height)}
		corner3 := mgl32.Vec3{float32(door.X + door.Width), 0, float32(door.Y + door.Height)}
		corner4 := mgl32.Vec3{float32(door.X + door.Width), 0, float32(door.Y)}
		if isPointInRectangle(position, corner1, corner2, corner3, corner4) {
			return &door
		}
	}
	return nil
}

func (aotManager *AotManager) GetAotTriggerNearPlayer(position mgl32.Vec3) *fileio.ScriptInstrAotSet {
	for _, aot := range aotManager.AotTriggers {
		corner1 := mgl32.Vec3{float32(aot.X), 0, float32(aot.Z)}
		corner2 := mgl32.Vec3{float32(aot.X), 0, float32(aot.Z + aot.Depth)}
		corner3 := mgl32.Vec3{float32(aot.X + aot.Width), 0, float32(aot.Z + aot.Depth)}
		corner4 := mgl32.Vec3{float32(aot.X + aot.Width), 0, float32(aot.Z)}
		if isPointInRectangle(position, corner1, corner2, corner3, corner4) {
			return &aot
		}
	}
	return nil
}

func (aotManager *AotManager) AddDoorAot(door fileio.ScriptInstrDoorAotSet) {
	aotManager.Doors = append(aotManager.Doors, door)
}

func (aotManager *AotManager) AddItemAot(item fileio.ScriptInstrItemAotSet) {
	aotManager.Items = append(aotManager.Items, item)
}

func (aotManager *AotManager) AddScriptSprite(sprite fileio.ScriptSprite) {
	aotManager.Sprites = append(aotManager.Sprites, sprite)
}

func (aotManager *AotManager) AddAotTrigger(aotTrigger fileio.ScriptInstrAotSet) {
	aotManager.AotTriggers = append(aotManager.AotTriggers, aotTrigger)
}

func (aotManager *AotManager) RemoveAotTrigger(aotIndex int) {
	for i, aot := range aotManager.AotTriggers {
		if int(aot.Aot) == aotIndex {
			fmt.Println("Remove aot index", aotIndex, ", aot type", aot.Id)
			aotManager.AotTriggers = append(aotManager.AotTriggers[:i], aotManager.AotTriggers[i+1:]...)
			return
		}
	}
}
