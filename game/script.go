package game

import (
	"../fileio"
	"bytes"
	"encoding/binary"
	"log"
)

type ScriptDoor struct {
	Opcode                       uint8 // 0x3b
	Id                           uint8 // Index of item in array of room objects list
	Unknown0                     [2]uint16
	X, Y                         int16 // Location of door
	Width, Height                int16 // Size of door
	NextX, NextY, NextZ, NextDir int16 // Position and direction of player after door entered
	Stage, Room, Camera          uint8 // Stage, room, camera after door entered
	Unknown1                     uint8
	DoorType                     uint8
	DoorLock                     uint8
	Unknown2                     uint8
	DoorLocked                   uint8
	DoorKey                      uint8
	Unknown3                     uint8
}

func (gameDef *GameDef) RunScript(scriptData [][][]byte) {
	gameDef.Doors = make([]ScriptDoor, 0)

	for functionNum := 0; functionNum < len(scriptData); functionNum++ {
		for lineNum := 0; lineNum < len(scriptData[functionNum]); lineNum++ {
			lineData := scriptData[functionNum][lineNum]
			opcode := lineData[0]

			switch opcode {
			case fileio.OP_DOOR_AOT_SET:
				byteArr := bytes.NewBuffer(lineData)
				var door ScriptDoor
				err := binary.Read(byteArr, binary.LittleEndian, &door)
				if err != nil {
					log.Fatal("Error loading door")
				}
				gameDef.Doors = append(gameDef.Doors, door)
			}
		}
	}
}
