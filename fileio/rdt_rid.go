package fileio

// .rid - Camera position data

import (
	"encoding/binary"
	"io"
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

type RIDHeader struct {
	Flag             uint16
	DistanceToScreen uint16
	CameraFromX      int32
	CameraFromY      int32
	CameraFromZ      int32
	CameraToX        int32
	CameraToY        int32
	CameraToZ        int32
	MaskOffset       uint32
}

type CameraInfo struct {
	CameraFrom mgl32.Vec3
	CameraTo   mgl32.Vec3
	CameraFov  float32 // in degrees
}

type RIDOutput struct {
	CameraPositions []CameraInfo
	CameraMasks     [][]MaskRectangle
}

func LoadRDT_RID(r io.ReaderAt, fileLength int64, rdtHeader RDTHeader, offsets RDTOffsets) (*RIDOutput, error) {
	offset := int64(offsets.OffsetCameraPosition)
	reader := io.NewSectionReader(r, offset, fileLength-offset)

	// Read from file
	cameraPositions := make([]RIDHeader, int(rdtHeader.NumCameras))
	if err := binary.Read(reader, binary.LittleEndian, &cameraPositions); err != nil {
		return nil, err
	}

	// Convert camera positions to use floating point
	cameraInfos := make([]CameraInfo, len(cameraPositions))
	for i, cameraPosition := range cameraPositions {
		cameraFrom := mgl32.Vec3{float32(cameraPosition.CameraFromX), float32(cameraPosition.CameraFromY), float32(cameraPosition.CameraFromZ)}
		cameraTo := mgl32.Vec3{float32(cameraPosition.CameraToX), float32(cameraPosition.CameraToY), float32(cameraPosition.CameraToZ)}
		cameraInfos[i] = CameraInfo{
			CameraFrom: cameraFrom,
			CameraTo:   cameraTo,
			CameraFov:  CalculateFOVDegrees(int(cameraPosition.DistanceToScreen) >> 7),
		}
	}

	// Read background image masks
	cameraMasks := make([][]MaskRectangle, int(rdtHeader.NumCameras))
	for i := 0; i < int(rdtHeader.NumCameras); i++ {
		if cameraPositions[i].MaskOffset == 0xffffffff {
			cameraMasks[i] = make([]MaskRectangle, 0)
			continue
		}

		offset := int64(cameraPositions[i].MaskOffset)
		reader := io.NewSectionReader(r, offset, fileLength-offset)
		priOutput, err := LoadRDT_PRI(reader, fileLength)
		if err != nil {
			return nil, err
		}
		// Some cameras don't have image masks
		if priOutput != nil {
			cameraMasks[i] = priOutput.Masks
		} else {
			cameraMasks[i] = make([]MaskRectangle, 0)
		}
	}

	output := &RIDOutput{
		CameraPositions: cameraInfos,
		CameraMasks:     cameraMasks,
	}
	return output, nil
}

func CalculateFOVDegrees(distanceToScreen int) float32 {
	halfScreenHeight := 120.0 // assumes screen dimensions is 320x240
	// tan(fov / 2) = ((screen height) / 2) / distanceToScreen
	fovAngleRadians := 2.0 * math.Atan(halfScreenHeight/float64(distanceToScreen))
	return mgl32.RadToDeg(float32(fovAngleRadians))
}
