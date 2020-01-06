package fileio

// .lit - Camera light data

import (
	"encoding/binary"
	"io"
)

type LITPosition struct {
	X int16
	Y int16
	Z int16
}

type LITLightColor struct {
	R uint8
	G uint8
	B uint8
}

type LITCameraLight struct {
	LightType    [2]uint16
	Colors       [3]LITLightColor // Color for each light
	AmbientColor LITLightColor    // Ambient color for camera
	Positions    [3]LITPosition   // Position for each light
	Brightness   [3]uint16        // Brightness for each light
}

type LITOutput struct {
	Lights []LITCameraLight
}

func LoadRDT_LIT(r io.ReaderAt, fileLength int64, rdtHeader RDTHeader, offsets RDTOffsets) (*LITOutput, error) {
	offset := int64(offsets.OffsetLights)
	reader := io.NewSectionReader(r, offset, fileLength-offset)

	lights := make([]LITCameraLight, int(rdtHeader.NumCameras))
	if err := binary.Read(reader, binary.LittleEndian, &lights); err != nil {
		return nil, err
	}

	output := &LITOutput{
		Lights: lights,
	}
	return output, nil

}
