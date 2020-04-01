package game

import (
	"github.com/go-gl/mathgl/mgl32"
)

type RoomMapKey struct {
	StageId int
	RoomId  int
}

var (
	// Jump to other rooms directly
	DebugLocations = map[RoomMapKey]mgl32.Vec3{
		RoomMapKey{1, 0x00}: mgl32.Vec3{18800, 0, -3160},
		RoomMapKey{1, 0x01}: mgl32.Vec3{-15047, 0, -11799},
		RoomMapKey{1, 0x02}: mgl32.Vec3{-25059, 0, 20944},
		RoomMapKey{1, 0x18}: mgl32.Vec3{-9212, 0, 2520},
		RoomMapKey{1, 0x19}: mgl32.Vec3{-6168, 0, -15445},
		RoomMapKey{1, 0x1a}: mgl32.Vec3{-7250, 0, -550},
		RoomMapKey{1, 0x1b}: mgl32.Vec3{-8353, 0, -20638},
	}
)
