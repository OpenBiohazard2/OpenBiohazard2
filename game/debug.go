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
		{1, 0x00}: {18800, 0, -3160},
		{1, 0x01}: {-15047, 0, -11799},
		{1, 0x02}: {-25059, 0, 20944},
		{1, 0x18}: {-9212, 0, 2520},
		{1, 0x19}: {-6168, 0, -15445},
		{1, 0x1a}: {-7250, 0, -550},
		{1, 0x1b}: {-8353, 0, -20638},
	}
)
