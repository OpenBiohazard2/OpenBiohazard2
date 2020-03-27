package fileio

import (
	"io"
)

func LoadRDT_VABStream(r io.ReaderAt, fileLength int64, offsets RDTOffsets) (*VABOutput, error) {
	offset := int64(offsets.OffsetRoomVABHeader)
	vabHeaderReader := io.NewSectionReader(r, offset, fileLength-offset)
	vabHeaderOutput, err := LoadVABHeaderStream(vabHeaderReader, fileLength)
	if err != nil {
		return nil, err
	}

	offset = int64(offsets.OffsetRoomVABData)
	vabDataReader := io.NewSectionReader(r, int64(offset), fileLength-int64(offset))
	_, err = LoadVABDataStream(vabDataReader, fileLength, vabHeaderOutput)
	if err != nil {
		return nil, err
	}

	return &VABOutput{}, nil
}
