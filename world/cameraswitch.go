package world

import (
	"math"

	"github.com/OpenBiohazard2/OpenBiohazard2/fileio"
	"github.com/go-gl/mathgl/mgl32"
)

type CameraSwitchHandler struct {
	CameraSwitches          []fileio.RVDHeader
	CameraSwitchTransitions map[int][]int
}

func NewCameraSwitchHandler(cameraSwitches []fileio.RVDHeader, maxCamerasInRoom int) *CameraSwitchHandler {
	// cameraSwitchTransitions shows which regions are reachable from the current camera
	// The key is the camera id
	// The value is an array of switches that are reachable
	cameraSwitchTransitions := make(map[int][]int, 0)
	for roomCameraId := 0; roomCameraId < maxCamerasInRoom; roomCameraId++ {
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

	return &CameraSwitchHandler{
		CameraSwitches:          cameraSwitches,
		CameraSwitchTransitions: cameraSwitchTransitions,
	}
}

func (cameraSwitchHandler *CameraSwitchHandler) GetCameraSwitchNewRegion(position mgl32.Vec3, curCameraId int) *fileio.RVDHeader {
	playerFloorNum := int(math.Round(float64(position.Y()) / fileio.FLOOR_HEIGHT_UNIT))

	for _, regionIndex := range cameraSwitchHandler.CameraSwitchTransitions[curCameraId] {
		region := cameraSwitchHandler.CameraSwitches[regionIndex]
		corner1 := mgl32.Vec3{float32(region.X1), 0, float32(region.Z1)}
		corner2 := mgl32.Vec3{float32(region.X2), 0, float32(region.Z2)}
		corner3 := mgl32.Vec3{float32(region.X3), 0, float32(region.Z3)}
		corner4 := mgl32.Vec3{float32(region.X4), 0, float32(region.Z4)}

		// Check region floor for rooms with multiple floor heights
		if region.Floor != 255 && int(region.Floor) != playerFloorNum {
			continue
		}

		if isPointInRectangle(position, corner1, corner2, corner3, corner4) {
			return &region
		}
	}
	return nil
}
