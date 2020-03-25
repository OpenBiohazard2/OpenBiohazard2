package fileio

// .emr file - Skeleton data

import (
	"encoding/binary"
	"io"
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

type EMRHeader struct {
	OffsetArmatures uint16
	OffsetFrames    uint16
	Count           uint16
	ElementSize     uint16
}

type EMRArmature struct {
	Count  uint16 // Number of meshes linked to this one
	Offset uint16
}

type EMRRelativePosition struct {
	X int16
	Y int16
	Z int16
}

type EMRFrame struct {
	XOffset int16
	YOffset int16
	ZOffset int16
	XSpeed  int16
	YSpeed  int16
	ZSpeed  int16
}

type AnimationFrame struct {
	FrameHeader    EMRFrame
	RotationAngles []mgl32.Vec3
}

type EMROutput struct {
	RelativePositionData []EMRRelativePosition
	ArmatureData         []EMRArmature
	ArmatureChildren     [][]uint8
	MeshIdList           []uint8
	FrameData            []AnimationFrame
}

func LoadEMRStream(r io.ReaderAt, fileLength int64, animationData *EDDOutput) (*EMROutput, error) {
	streamReader := io.NewSectionReader(r, int64(0), fileLength)
	// Read header
	emrHeader := EMRHeader{}
	if err := binary.Read(streamReader, binary.LittleEndian, &emrHeader); err != nil {
		return nil, err
	}

	// Read relative positions
	relativePositions := make([]EMRRelativePosition, int(emrHeader.Count))
	if err := binary.Read(streamReader, binary.LittleEndian, &relativePositions); err != nil {
		return nil, err
	}

	// Read armature offsets
	streamReader = io.NewSectionReader(r, int64(emrHeader.OffsetArmatures), fileLength)
	armatures := make([]EMRArmature, int(emrHeader.Count))
	if err := binary.Read(streamReader, binary.LittleEndian, &armatures); err != nil {
		return nil, err
	}

	// Stores a hierarchy of the components
	// Used to calculated offset of each component based on its parent
	armatureChildren := make([][]uint8, int(emrHeader.Count))
	for i := 0; i < int(emrHeader.Count); i++ {
		streamReader = io.NewSectionReader(r, int64(emrHeader.OffsetArmatures)+int64(armatures[i].Offset), fileLength)

		armatureChildren[i] = make([]uint8, int(armatures[i].Count))
		if err := binary.Read(streamReader, binary.LittleEndian, &armatureChildren[i]); err != nil {
			return nil, err
		}
	}

	// List of all the component ids in this model
	streamReader = io.NewSectionReader(r, int64(emrHeader.OffsetArmatures)+int64(armatures[0].Offset), fileLength)
	meshList := make([]uint8, int(emrHeader.Count))
	if err := binary.Read(streamReader, binary.LittleEndian, &meshList); err != nil {
		return nil, err
	}

	// Read animation frames
	numFrames := animationData.NumFrames
	// EMR frame header is 12 bytes (6 uint16 values)
	// Remaining data contains an array of rotation angles
	remainSize := int(emrHeader.ElementSize) - 2*6

	frameData := make([]AnimationFrame, numFrames)
	for i := 0; i < numFrames; i++ {
		streamReader = io.NewSectionReader(r, int64(emrHeader.OffsetFrames)+int64(i*int(emrHeader.ElementSize)), fileLength)
		frameHeader := EMRFrame{}
		if err := binary.Read(streamReader, binary.LittleEndian, &frameHeader); err != nil {
			return nil, err
		}
		maxAngleCount := int(math.Floor(float64(remainSize) * 8.0 / 12.0))

		bitReader := NewBitReader(streamReader)
		angles := make([]mgl32.Vec3, int(maxAngleCount))
		for j := 0; j < int(maxAngleCount); j++ {
			// Each angle is stored as 12 bits
			angleX := convertFrameAngleToRadians(float32(bitReader.UnsafeReadNumBitsLittleEndian(12)))
			angleY := convertFrameAngleToRadians(float32(bitReader.UnsafeReadNumBitsLittleEndian(12)))
			angleZ := convertFrameAngleToRadians(float32(bitReader.UnsafeReadNumBitsLittleEndian(12)))
			angles[j] = mgl32.Vec3{angleX, angleY, angleZ}
		}

		frameData[i] = AnimationFrame{
			FrameHeader:    frameHeader,
			RotationAngles: angles,
		}
	}

	output := &EMROutput{
		RelativePositionData: relativePositions,
		ArmatureData:         armatures,
		ArmatureChildren:     armatureChildren,
		MeshIdList:           meshList,
		FrameData:            frameData,
	}
	return output, nil
}

// Convert 12 bit angle to radians
// This is different from converting degrees to radians
func convertFrameAngleToRadians(frameAngle float32) float32 {
	// maximum 12 bit number is 4095
	MAX_ANGLE := float32(4096.0)
	return (frameAngle / MAX_ANGLE) * (2.0 * math.Pi)
}
