package fileio

import (
	"encoding/binary"
	"fmt"
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

// Type-specific readers for common data types

// ReadUint8 reads a single uint8 from the stream
func (streamReader *StreamReader) ReadUint8() (uint8, error) {
	var value uint8
	if err := streamReader.ReadData(&value); err != nil {
		return 0, fmt.Errorf("failed to read uint8: %w", err)
	}
	return value, nil
}

// ReadUint16 reads a uint16 from the stream
func (streamReader *StreamReader) ReadUint16() (uint16, error) {
	var value uint16
	if err := streamReader.ReadData(&value); err != nil {
		return 0, fmt.Errorf("failed to read uint16: %w", err)
	}
	return value, nil
}

// ReadUint32 reads a uint32 from the stream
func (streamReader *StreamReader) ReadUint32() (uint32, error) {
	var value uint32
	if err := streamReader.ReadData(&value); err != nil {
		return 0, fmt.Errorf("failed to read uint32: %w", err)
	}
	return value, nil
}

// ReadInt8 reads a single int8 from the stream
func (streamReader *StreamReader) ReadInt8() (int8, error) {
	var value int8
	if err := streamReader.ReadData(&value); err != nil {
		return 0, fmt.Errorf("failed to read int8: %w", err)
	}
	return value, nil
}

// ReadInt16 reads an int16 from the stream
func (streamReader *StreamReader) ReadInt16() (int16, error) {
	var value int16
	if err := streamReader.ReadData(&value); err != nil {
		return 0, fmt.Errorf("failed to read int16: %w", err)
	}
	return value, nil
}

// ReadInt32 reads an int32 from the stream
func (streamReader *StreamReader) ReadInt32() (int32, error) {
	var value int32
	if err := streamReader.ReadData(&value); err != nil {
		return 0, fmt.Errorf("failed to read int32: %w", err)
	}
	return value, nil
}

// ReadFloat32 reads a float32 from the stream
func (streamReader *StreamReader) ReadFloat32() (float32, error) {
	var value float32
	if err := streamReader.ReadData(&value); err != nil {
		return 0, fmt.Errorf("failed to read float32: %w", err)
	}
	return value, nil
}

// ReadFloat64 reads a float64 from the stream
func (streamReader *StreamReader) ReadFloat64() (float64, error) {
	var value float64
	if err := streamReader.ReadData(&value); err != nil {
		return 0, fmt.Errorf("failed to read float64: %w", err)
	}
	return value, nil
}

// ReadBytes reads a slice of bytes from the stream
func (streamReader *StreamReader) ReadBytes(count int) ([]byte, error) {
	if count <= 0 {
		return nil, fmt.Errorf("invalid byte count: %d", count)
	}

	data := make([]byte, count)
	n, err := streamReader.reader.Read(data)
	if err != nil {
		return nil, fmt.Errorf("failed to read %d bytes: %w", count, err)
	}
	if n != count {
		return nil, fmt.Errorf("expected to read %d bytes, got %d", count, n)
	}
	return data, nil
}

// ReadString reads a null-terminated string from the stream
func (streamReader *StreamReader) ReadString(maxLength int) (string, error) {
	if maxLength <= 0 {
		return "", fmt.Errorf("invalid max length: %d", maxLength)
	}

	data := make([]byte, maxLength)
	n, err := streamReader.reader.Read(data)
	if err != nil {
		return "", fmt.Errorf("failed to read string: %w", err)
	}

	// Find null terminator
	for i := 0; i < n; i++ {
		if data[i] == 0 {
			return string(data[:i]), nil
		}
	}

	// No null terminator found, return what we have
	return string(data[:n]), nil
}

// ReadFixedString reads a fixed-length string from the stream
func (streamReader *StreamReader) ReadFixedString(length int) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("invalid string length: %d", length)
	}

	data, err := streamReader.ReadBytes(length)
	if err != nil {
		return "", fmt.Errorf("failed to read fixed string: %w", err)
	}

	// Remove null padding
	for i := len(data) - 1; i >= 0; i-- {
		if data[i] != 0 {
			return string(data[:i+1]), nil
		}
	}

	return "", nil
}
