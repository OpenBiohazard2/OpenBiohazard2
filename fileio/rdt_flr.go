package fileio

// .flr - Floor sound data

import (
	"encoding/binary"
	"io"
)

type FLRSound struct {
	X           int16
	Y           int16
	Width       uint16
	Depth       uint16
	SoundEffect uint16
	Height      uint16
}

type FLROutput struct {
	FloorSounds []FLRSound
}

func LoadRDT_FLRStream(r io.ReaderAt, fileLength int64, offsets RDTOffsets) (*FLROutput, error) {
	offset := int64(offsets.OffsetFloorSound)
	flrHeaderReader := io.NewSectionReader(r, offset, fileLength-offset)
	floorSoundCount := uint16(0)
	if err := binary.Read(flrHeaderReader, binary.LittleEndian, &floorSoundCount); err != nil {
		return nil, err
	}

	floorSounds := make([]FLRSound, floorSoundCount)
	if err := binary.Read(flrHeaderReader, binary.LittleEndian, &floorSounds); err != nil {
		return nil, err
	}

	output := &FLROutput{
		FloorSounds: floorSounds,
	}
	return output, nil
}
