package fileio

// .edd file - Animation data

import (
	"encoding/binary"
	"io"
)

type EDDHeaderObject struct {
	Count  uint16
	Offset uint16
}

type EDDTableElement struct {
	FrameId int
	Flag    int
}

type EDDOutput struct {
	AnimationIndexFrames [][]EDDTableElement
	NumFrames            int
}

func LoadEDDStream(r io.ReaderAt, fileLength int64) (*EDDOutput, error) {
	streamReader := io.NewSectionReader(r, int64(0), fileLength)
	// Read header count based on the offset of the first header
	// Everything before the first header offset is header data
	firstHeader := EDDHeaderObject{}
	if err := binary.Read(streamReader, binary.LittleEndian, &firstHeader); err != nil {
		return nil, err
	}
	headerCount := int(firstHeader.Offset) / 4

	// Read headers
	streamReader = io.NewSectionReader(r, int64(0), fileLength)
	eddHeaders := make([]EDDHeaderObject, int(headerCount))
	if err := binary.Read(streamReader, binary.LittleEndian, &eddHeaders); err != nil {
		return nil, err
	}

	animationIndexFrames := make([][]EDDTableElement, len(eddHeaders))
	bitReader := NewBitReader(streamReader)
	maxFrameNumber := 0
	for i := 0; i < len(eddHeaders); i++ {
		eddTable := make([]EDDTableElement, int(eddHeaders[i].Count))
		// Each element is 4 bytes (32 bits) in little endian
		for j := 0; j < int(eddHeaders[i].Count); j++ {
			frameId := int(bitReader.UnsafeReadNumBitsLittleEndian(12))
			flag := int(bitReader.UnsafeReadNumBitsLittleEndian(20))
			eddTable[j] = EDDTableElement{
				FrameId: frameId,
				Flag:    flag,
			}

			if frameId > maxFrameNumber {
				maxFrameNumber = frameId
			}
		}

		animationIndexFrames[i] = eddTable
	}

	output := &EDDOutput{
		AnimationIndexFrames: animationIndexFrames,
		NumFrames:            maxFrameNumber + 1,
	}
	return output, nil
}
