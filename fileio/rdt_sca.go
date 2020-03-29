package fileio

// .sca - Collision data

import (
	"encoding/binary"
	"fmt"
	"io"
)

const (
	FLOOR_HEIGHT_UNIT = -1800

	SCA_TYPE_SLOPE = 11
	SCA_TYPE_STAIRS = 12
)

type SCAHeader struct {
	CeilingX       int16
	CeilingZ       int16
	Count          uint32
	CeilingY       int32
	CeilingWidth   uint16
	CeilingDensity uint16
}

type SCAElement struct {
	X            int16
	Z            int16
	Width        uint16
	Density      uint16
	Flag         uint16
	Type         uint16
	FloorNumFlag uint32
}

type CollisionEntity struct {
	ScaIndex    int
	X           int
	Z           int
	Width       int
	Density     int
	Shape       int
	SlopeHeight int
	SlopeType   int
	RampBottom  float32
	FloorCheck  []bool
}

type SCAOutput struct {
	CollisionEntities []CollisionEntity
}

func LoadRDT_SCA(r io.ReaderAt, fileLength int64, rdtHeader RDTHeader, offsets RDTOffsets) (*SCAOutput, error) {
	offset := offsets.OffsetCollisionData
	reader := io.NewSectionReader(r, int64(offset), fileLength-int64(offset))

	scaHeader := SCAHeader{}
	if err := binary.Read(reader, binary.LittleEndian, &scaHeader); err != nil {
		return nil, err
	}

	collisionEntities := make([]CollisionEntity, int(scaHeader.Count)-1)
	for i := 0; i < int(scaHeader.Count)-1; i++ {
		scaElement := SCAElement{}
		if err := binary.Read(reader, binary.LittleEndian, &scaElement); err != nil {
			return nil, err
		}

		shape := scaElement.Flag & 0x000F

		floorHeightMultiplier := int(scaElement.Type>>6) & 0x1F

		rampBottom := float32(0.0)
		slopeType := -1
		if shape == SCA_TYPE_SLOPE || shape == SCA_TYPE_STAIRS {
			elevationType := int(scaElement.Type>>4) & 3
			slopeType = elevationType
			switch slopeType{
			case 0:
				rampBottom = float32(scaElement.X)
			case 1:
				rampBottom = float32(scaElement.X) + float32(scaElement.Width)
			case 2:
				rampBottom = float32(scaElement.Z)
			case 3:
				rampBottom = float32(scaElement.Z) + float32(scaElement.Density)
			}
		}

		// Check if this floor has a collision entity
		// The boundaries can be different for each floor level
		floorCheck := make([]bool, 0)
		// Convert to binary string
		flags := fmt.Sprintf("%032b", int(scaElement.FloorNumFlag))
		for j := 31; j >= 0; j-- {
			if flags[j] == '1' {
				floorCheck = append(floorCheck, true)
			} else {
				floorCheck = append(floorCheck, false)
			}
		}

		collisionEntities[i] = CollisionEntity{
			ScaIndex:    i,
			X:           int(scaElement.X),
			Z:           int(scaElement.Z),
			Width:       int(scaElement.Width),
			Density:     int(scaElement.Density),
			Shape:       int(shape),
			SlopeHeight: floorHeightMultiplier * FLOOR_HEIGHT_UNIT,
			SlopeType:   slopeType,
			RampBottom:  rampBottom,
			FloorCheck:  floorCheck,
		}
	}
	output := &SCAOutput{
		CollisionEntities: collisionEntities,
	}
	return output, nil
}
