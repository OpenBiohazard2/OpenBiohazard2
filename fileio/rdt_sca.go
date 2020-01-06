package fileio

// .sca - Collision data

import (
	"encoding/binary"
	"io"
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
	X       int16
	Z       int16
	Width   uint16
	Density uint16
	Flag    uint16
	Type    uint16
	Floor   uint32
}

type CollisionEntity struct {
	X       int
	Z       int
	Width   int
	Density int
	Shape   int
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
		collisionEntities[i] = CollisionEntity{
			X:       int(scaElement.X),
			Z:       int(scaElement.Z),
			Width:   int(scaElement.Width),
			Density: int(scaElement.Density),
			Shape:   int(shape),
		}
	}
	output := &SCAOutput{
		CollisionEntities: collisionEntities,
	}
	return output, nil
}
