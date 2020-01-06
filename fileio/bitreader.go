package fileio

import (
	"encoding/binary"
	"io"
)

type BitReader struct {
	reader *io.SectionReader
	byte   byte
	offset byte
}

func NewBitReader(r *io.SectionReader) *BitReader {
	return &BitReader{r, 0, 0}
}

// Reads the next bit in the file
func (r *BitReader) ReadBit() (int, error) {
	if r.offset == 8 {
		r.offset = 0
	}
	if r.offset == 0 {
		r.byte = byte(0)
		if err := binary.Read(r.reader, binary.LittleEndian, &r.byte); err != nil {
			return 0, err
		}
	}
	hasBit := (r.byte & (0x80 >> r.offset)) != 0
	var bit int
	if hasBit {
		bit = 1
	} else {
		bit = 0
	}
	r.offset++
	return bit, nil
}

// Reads a sequence of bits in sequential order
// Do not use this function if the data is in little endian
func (r *BitReader) ReadNumBits(numBits int) (uint64, error) {
	var result uint64
	for i := numBits - 1; i >= 0; i-- {
		bit, err := r.ReadBit()
		if err != nil {
			return 0, err
		}
		if bit == 1 {
			result |= 1 << uint(i)
		}
	}

	return result, nil
}

func (r *BitReader) ReadBitLittleEndian() (int, error) {
	if r.offset == 8 {
		r.offset = 0
	}
	if r.offset == 0 {
		r.byte = byte(0)
		if err := binary.Read(r.reader, binary.LittleEndian, &r.byte); err != nil {
			return 0, err
		}
	}
	hasBit := (r.byte & (0x80 >> (7 - r.offset))) != 0
	var bit int
	if hasBit {
		bit = 1
	} else {
		bit = 0
	}
	r.offset++
	return bit, nil
}

func (r *BitReader) ReadNumBitsLittleEndian(numBits int) (uint64, error) {
	var result uint64
	for i := 0; i < numBits; i++ {
		bit, err := r.ReadBitLittleEndian()
		if err != nil {
			return 0, err
		}
		if bit == 1 {
			result |= 1 << uint(i)
		}
	}

	return result, nil
}

// No error handling
func (r *BitReader) UnsafeReadBit() int {
	bit, _ := r.ReadBit()
	return bit
}

func (r *BitReader) UnsafeReadNumBits(numBits int) uint64 {
	bits, _ := r.ReadNumBits(numBits)
	return bits
}

func (r *BitReader) UnsafeReadNumBitsLittleEndian(numBits int) uint64 {
	bits, _ := r.ReadNumBitsLittleEndian(numBits)
	return bits
}

func (r *BitReader) UnsafeReadByte() byte {
	return byte(r.UnsafeReadNumBits(8))
}
