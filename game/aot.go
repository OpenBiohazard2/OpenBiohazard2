package game

import (
	"fmt"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/samuelyuan/openbiohazard2/fileio"
	"github.com/samuelyuan/openbiohazard2/geometry"
)

// Handle script doors, items, events

const (
	AOT_DOOR  = 1
	AOT_ITEM  = 2
	AOT_EVENT = 5
)

type AotManager struct {
	Doors       []AotDoor
	Items       []AotItem
	Sprites     []fileio.ScriptInstrSceEsprOn
	AotTriggers []AotObject
}

type AotHeader struct {
	Aot   uint8
	Id    uint8
	Type  uint8
	Floor uint8
	Super uint8
}

type AotObject struct {
	Header AotHeader
	Bounds *geometry.Quad
	Data   [6]uint8
}

type AotDoor struct {
	Header                       AotHeader
	Bounds                       *geometry.Quad
	NextX, NextY, NextZ, NextDir int16
	Stage, Room, Camera          uint8
	NextFloor                    uint8
	TextureType                  uint8
	DoorType                     uint8
	KnockType                    uint8
	KeyId                        uint8
	KeyType                      uint8
	Free                         uint8
}

type AotItem struct {
	Header          AotHeader
	Bounds          *geometry.Quad
	ItemId          uint16
	Amount          uint16
	ItemPickedIndex uint16
	Md1ModelId      uint8
	Act             uint8
}

func NewAotManager() *AotManager {
	return &AotManager{
		Doors:       make([]AotDoor, 0),
		Items:       make([]AotItem, 0),
		Sprites:     make([]fileio.ScriptInstrSceEsprOn, 0),
		AotTriggers: make([]AotObject, 0),
	}
}

func (aotManager *AotManager) AddScriptSprite(sprite fileio.ScriptInstrSceEsprOn) {
	aotManager.Sprites = append(aotManager.Sprites, sprite)
}

func (aotManager *AotManager) GetDoorNearPlayer(position mgl32.Vec3) *AotDoor {
	for _, door := range aotManager.Doors {
		vertices := door.Bounds.Vertices
		if isPointInRectangle(position, vertices[0], vertices[1], vertices[2], vertices[3]) {
			return &door
		}
	}
	return nil
}

func (aotManager *AotManager) GetAotTriggerNearPlayer(position mgl32.Vec3) *AotObject {
	for _, aot := range aotManager.AotTriggers {
		vertices := aot.Bounds.Vertices
		if isPointInRectangle(position, vertices[0], vertices[1], vertices[2], vertices[3]) {
			return &aot
		}
	}
	return nil
}

func (aotManager *AotManager) AddDoorAot(aotInstruction fileio.ScriptInstrDoorAotSet) {
	aotHeader := AotHeader{
		Aot:   aotInstruction.Aot,
		Id:    aotInstruction.Id,
		Type:  aotInstruction.Type,
		Floor: aotInstruction.Floor,
		Super: aotInstruction.Super,
	}
	rect := geometry.NewRectangle(
		float32(aotInstruction.X), float32(aotInstruction.Z),
		float32(aotInstruction.Width), float32(aotInstruction.Depth))
	doorAot := AotDoor{
		Header:      aotHeader,
		Bounds:      rect,
		NextX:       aotInstruction.NextX,
		NextY:       aotInstruction.NextY,
		NextZ:       aotInstruction.NextZ,
		NextDir:     aotInstruction.NextDir,
		Stage:       aotInstruction.Stage,
		Room:        aotInstruction.Room,
		Camera:      aotInstruction.Camera,
		NextFloor:   aotInstruction.NextFloor,
		TextureType: aotInstruction.TextureType,
		DoorType:    aotInstruction.DoorType,
		KnockType:   aotInstruction.KnockType,
		KeyId:       aotInstruction.KeyId,
		KeyType:     aotInstruction.KeyType,
		Free:        aotInstruction.Free,
	}

	fmt.Println("Create new door aot", aotInstruction.Aot, "with aot type", aotInstruction.Id)
	aotManager.Doors = append(aotManager.Doors, doorAot)
}

func (aotManager *AotManager) AddDoorAot4p(aotInstruction fileio.ScriptInstrDoorAotSet4p) {
	aotHeader := AotHeader{
		Aot:   aotInstruction.Aot,
		Id:    aotInstruction.Id,
		Type:  aotInstruction.Type,
		Floor: aotInstruction.Floor,
		Super: aotInstruction.Super,
	}
	rect := geometry.NewQuadFourPoints([4][]float32{
		[]float32{float32(aotInstruction.X1), float32(aotInstruction.Z1)},
		[]float32{float32(aotInstruction.X2), float32(aotInstruction.Z2)},
		[]float32{float32(aotInstruction.X3), float32(aotInstruction.Z3)},
		[]float32{float32(aotInstruction.X4), float32(aotInstruction.Z4)},
	})
	doorAot := AotDoor{
		Header:      aotHeader,
		Bounds:      rect,
		NextX:       aotInstruction.NextX,
		NextY:       aotInstruction.NextY,
		NextZ:       aotInstruction.NextZ,
		NextDir:     aotInstruction.NextDir,
		Stage:       aotInstruction.Stage,
		Room:        aotInstruction.Room,
		Camera:      aotInstruction.Camera,
		NextFloor:   aotInstruction.NextFloor,
		TextureType: aotInstruction.TextureType,
		DoorType:    aotInstruction.DoorType,
		KnockType:   aotInstruction.KnockType,
		KeyId:       aotInstruction.KeyId,
		KeyType:     aotInstruction.KeyType,
		Free:        aotInstruction.Free,
	}

	fmt.Println("Create new door aot 4p", aotInstruction.Aot, "with aot type", aotInstruction.Id)
	aotManager.Doors = append(aotManager.Doors, doorAot)
}

func (aotManager *AotManager) AddItemAot(aotInstruction fileio.ScriptInstrItemAotSet) {
	aotHeader := AotHeader{
		Aot:   aotInstruction.Aot,
		Id:    aotInstruction.Id,
		Type:  aotInstruction.Type,
		Floor: aotInstruction.Floor,
		Super: aotInstruction.Super,
	}
	rect := geometry.NewRectangle(
		float32(aotInstruction.X), float32(aotInstruction.Z),
		float32(aotInstruction.Width), float32(aotInstruction.Depth))
	itemAot := AotItem{
		Header:          aotHeader,
		Bounds:          rect,
		ItemId:          aotInstruction.ItemId,
		Amount:          aotInstruction.Amount,
		ItemPickedIndex: aotInstruction.ItemPickedIndex,
		Md1ModelId:      aotInstruction.Md1ModelId,
		Act:             aotInstruction.Act,
	}

	fmt.Println("Create new item aot", aotInstruction.Aot, "with aot type", aotInstruction.Id)
	aotManager.Items = append(aotManager.Items, itemAot)
}

func (aotManager *AotManager) AddItemAot4p(aotInstruction fileio.ScriptInstrItemAotSet4p) {
	aotHeader := AotHeader{
		Aot:   aotInstruction.Aot,
		Id:    aotInstruction.Id,
		Type:  aotInstruction.Type,
		Floor: aotInstruction.Floor,
		Super: aotInstruction.Super,
	}
	rect := geometry.NewQuadFourPoints([4][]float32{
		[]float32{float32(aotInstruction.X1), float32(aotInstruction.Z1)},
		[]float32{float32(aotInstruction.X2), float32(aotInstruction.Z2)},
		[]float32{float32(aotInstruction.X3), float32(aotInstruction.Z3)},
		[]float32{float32(aotInstruction.X4), float32(aotInstruction.Z4)},
	})
	itemAot := AotItem{
		Header:          aotHeader,
		Bounds:          rect,
		ItemId:          aotInstruction.ItemId,
		Amount:          aotInstruction.Amount,
		ItemPickedIndex: aotInstruction.ItemPickedIndex,
		Md1ModelId:      aotInstruction.Md1ModelId,
		Act:             aotInstruction.Act,
	}

	fmt.Println("Create new item aot 4p", aotInstruction.Aot, "with aot type", aotInstruction.Id)
	aotManager.Items = append(aotManager.Items, itemAot)
}

func (aotManager *AotManager) AddAotTrigger(aotInstruction fileio.ScriptInstrAotSet) {
	aotHeader := AotHeader{
		Aot:   aotInstruction.Aot,
		Id:    aotInstruction.Id,
		Type:  aotInstruction.Type,
		Floor: aotInstruction.Floor,
		Super: aotInstruction.Super,
	}
	rect := geometry.NewRectangle(
		float32(aotInstruction.X), float32(aotInstruction.Z),
		float32(aotInstruction.Width), float32(aotInstruction.Depth))
	aotTrigger := AotObject{
		Header: aotHeader,
		Bounds: rect,
		Data:   aotInstruction.Data,
	}

	fmt.Println("Create new aot index", aotInstruction.Aot, "with aot type", aotInstruction.Id)
	aotManager.AotTriggers = append(aotManager.AotTriggers, aotTrigger)
}

func (aotManager *AotManager) AddAotTrigger4p(aotInstruction fileio.ScriptInstrAotSet4p) {
	aotHeader := AotHeader{
		Aot:   aotInstruction.Aot,
		Id:    aotInstruction.Id,
		Type:  aotInstruction.Type,
		Floor: aotInstruction.Floor,
		Super: aotInstruction.Super,
	}
	rect := geometry.NewQuadFourPoints([4][]float32{
		[]float32{float32(aotInstruction.X1), float32(aotInstruction.Z1)},
		[]float32{float32(aotInstruction.X2), float32(aotInstruction.Z2)},
		[]float32{float32(aotInstruction.X3), float32(aotInstruction.Z3)},
		[]float32{float32(aotInstruction.X4), float32(aotInstruction.Z4)},
	})
	aotTrigger := AotObject{
		Header: aotHeader,
		Bounds: rect,
		Data:   aotInstruction.Data,
	}

	fmt.Println("Create new aot 4p index", aotInstruction.Aot, "with aot type", aotInstruction.Id)
	aotManager.AotTriggers = append(aotManager.AotTriggers, aotTrigger)
}

func (aotManager *AotManager) ResetAotTrigger(aotInstruction fileio.ScriptInstrAotReset) {
	for i, aot := range aotManager.AotTriggers {
		if int(aot.Header.Aot) == int(aotInstruction.Aot) {
			fmt.Println("Reset aot index", aotInstruction.Aot, "with aot type", aotInstruction.Id)
			aot.Header.Aot = aotInstruction.Aot
			aot.Header.Id = aotInstruction.Id
			aot.Header.Type = aotInstruction.Type
			aot.Data = aotInstruction.Data
			aotManager.AotTriggers[i] = aot
			return
		}
	}

	fmt.Println("No existing aot found for", aotInstruction.Aot, ". Create new aot with aot type", aotInstruction.Id)
	aotTrigger := AotObject{
		Header: AotHeader{
			Aot:   aotInstruction.Aot,
			Id:    aotInstruction.Id,
			Type:  aotInstruction.Type,
			Floor: 0,
			Super: 0,
		},
		Bounds: geometry.NewRectangle(0, 0, 0, 0),
		Data:   aotInstruction.Data,
	}
	aotManager.AotTriggers = append(aotManager.AotTriggers, aotTrigger)
}
