package fileio

import (
	"encoding/binary"
	"io"
)

type StreamReader struct {
	reader *io.SectionReader
}

func NewStreamReader(r *io.SectionReader) *StreamReader {
	return &StreamReader{r}
}

func (streamReader *StreamReader) SetPosition(newPosition int64) {
	streamReader.reader.Seek(newPosition, io.SeekStart)
}

func (streamReader *StreamReader) ReadData(data interface{}) error {
	// The data is little endian by default
	return binary.Read(streamReader.reader, binary.LittleEndian, data)
}