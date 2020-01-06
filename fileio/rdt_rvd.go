package fileio

// .rvd - Camera switch data

import (
	"encoding/binary"
	"io"
)

type RVDHeader struct {
	Flag  byte
	Floor byte
	Cam0  uint8
	Cam1  uint8
	X1    int16
	Z1    int16
	X2    int16
	Z2    int16
	X3    int16
	Z3    int16
	X4    int16
	Z4    int16
}

type RVDOutput struct {
	CameraSwitches []RVDHeader
}

func LoadRDT_RVD(r io.ReaderAt, fileLength int64, rdtHeader RDTHeader, offsets RDTOffsets) (*RVDOutput, error) {
	reader := io.NewSectionReader(r, int64(offsets.OffsetCameraSwitches), fileLength-int64(offsets.OffsetCameraSwitches))

	cameraSwitches := make([]RVDHeader, 0)
	for i := 0; i < 100; i++ {
		rvdHeader := RVDHeader{}
		if err := binary.Read(reader, binary.LittleEndian, &rvdHeader); err != nil {
			return nil, err
		}

		if rvdHeader.Flag == 255 && rvdHeader.Floor == 255 && rvdHeader.Cam0 == 255 && rvdHeader.Cam1 == 255 {
			break
		}

		cameraSwitches = append(cameraSwitches, rvdHeader)
	}
	output := &RVDOutput{
		CameraSwitches: cameraSwitches,
	}
	return output, nil
}
