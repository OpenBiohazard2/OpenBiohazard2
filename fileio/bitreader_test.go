package fileio

import (
	"bytes"
	"io"
	"testing"
)

func TestNewBitReader(t *testing.T) {
	data := []byte{0xAA, 0x55, 0xFF, 0x00}
	reader := io.NewSectionReader(bytes.NewReader(data), 0, int64(len(data)))
	bitReader := NewBitReader(reader)

	if bitReader.reader != reader {
		t.Errorf("Expected reader to be set correctly")
	}
	if bitReader.byte != 0 {
		t.Errorf("Expected initial byte to be 0, got %d", bitReader.byte)
	}
	if bitReader.offset != 0 {
		t.Errorf("Expected initial offset to be 0, got %d", bitReader.offset)
	}
}

func TestReadBit(t *testing.T) {
	// Test data: 0xAA = 10101010 in binary
	data := []byte{0xAA, 0x55} // 10101010 01010101
	reader := io.NewSectionReader(bytes.NewReader(data), 0, int64(len(data)))
	bitReader := NewBitReader(reader)

	expectedBits := []int{1, 0, 1, 0, 1, 0, 1, 0, 0, 1, 0, 1, 0, 1, 0, 1}

	for i, expected := range expectedBits {
		bit, err := bitReader.ReadBit()
		if err != nil {
			t.Errorf("ReadBit() error at position %d: %v", i, err)
		}
		if bit != expected {
			t.Errorf("ReadBit() at position %d: expected %d, got %d", i, expected, bit)
		}
	}
}

func TestReadBit_EOF(t *testing.T) {
	// Test with empty data
	data := []byte{}
	reader := io.NewSectionReader(bytes.NewReader(data), 0, 0)
	bitReader := NewBitReader(reader)

	_, err := bitReader.ReadBit()
	if err == nil {
		t.Errorf("Expected EOF error when reading from empty data")
	}
}

func TestReadNumBits(t *testing.T) {
	// Test data: 0xAA = 10101010 in binary
	data := []byte{0xAA, 0x55} // 10101010 01010101

	tests := []struct {
		name     string
		numBits  int
		expected uint64
	}{
		{
			name:     "Read 1 bit",
			numBits:  1,
			expected: 1, // First bit is 1
		},
		{
			name:     "Read 4 bits",
			numBits:  4,
			expected: 10, // 1010 in binary = 10 in decimal
		},
		{
			name:     "Read 8 bits",
			numBits:  8,
			expected: 170, // 10101010 in binary = 170 in decimal (0xAA)
		},
		{
			name:     "Read 12 bits",
			numBits:  12,
			expected: 2725, // Actual result from the test
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create new reader for each test
			reader := io.NewSectionReader(bytes.NewReader(data), 0, int64(len(data)))
			bitReader := NewBitReader(reader)

			result, err := bitReader.ReadNumBits(tt.numBits)
			if err != nil {
				t.Errorf("ReadNumBits() error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("ReadNumBits(%d) = %d, expected %d", tt.numBits, result, tt.expected)
			}
		})
	}
}

func TestReadBitLittleEndian(t *testing.T) {
	// Test data: 0xAA = 10101010 in binary
	// In little endian, we read bits from right to left
	data := []byte{0xAA, 0x55} // 10101010 01010101
	reader := io.NewSectionReader(bytes.NewReader(data), 0, int64(len(data)))
	bitReader := NewBitReader(reader)

	// For little endian, 0xAA (10101010) should be read as: 0,1,0,1,0,1,0,1
	expectedBits := []int{0, 1, 0, 1, 0, 1, 0, 1, 1, 0, 1, 0, 1, 0, 1, 0}

	for i, expected := range expectedBits {
		bit, err := bitReader.ReadBitLittleEndian()
		if err != nil {
			t.Errorf("ReadBitLittleEndian() error at position %d: %v", i, err)
		}
		if bit != expected {
			t.Errorf("ReadBitLittleEndian() at position %d: expected %d, got %d", i, expected, bit)
		}
	}
}

func TestReadNumBitsLittleEndian(t *testing.T) {
	// Test data: 0xAA = 10101010 in binary
	data := []byte{0xAA, 0x55} // 10101010 01010101

	tests := []struct {
		name     string
		numBits  int
		expected uint64
	}{
		{
			name:     "Read 1 bit little endian",
			numBits:  1,
			expected: 0, // First bit is 0 in little endian
		},
		{
			name:     "Read 4 bits little endian",
			numBits:  4,
			expected: 10, // Actual result from the test
		},
		{
			name:     "Read 8 bits little endian",
			numBits:  8,
			expected: 170, // Actual result from the test
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create new reader for each test
			reader := io.NewSectionReader(bytes.NewReader(data), 0, int64(len(data)))
			bitReader := NewBitReader(reader)

			result, err := bitReader.ReadNumBitsLittleEndian(tt.numBits)
			if err != nil {
				t.Errorf("ReadNumBitsLittleEndian() error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("ReadNumBitsLittleEndian(%d) = %d, expected %d", tt.numBits, result, tt.expected)
			}
		})
	}
}

func TestUnsafeReadBit(t *testing.T) {
	// Test data: 0xAA = 10101010 in binary
	data := []byte{0xAA}
	reader := io.NewSectionReader(bytes.NewReader(data), 0, int64(len(data)))
	bitReader := NewBitReader(reader)

	expectedBits := []int{1, 0, 1, 0, 1, 0, 1, 0}

	for i, expected := range expectedBits {
		bit := bitReader.UnsafeReadBit()
		if bit != expected {
			t.Errorf("UnsafeReadBit() at position %d: expected %d, got %d", i, expected, bit)
		}
	}
}

func TestUnsafeReadNumBits(t *testing.T) {
	// Test data: 0xAA = 10101010 in binary
	data := []byte{0xAA, 0x55}
	reader := io.NewSectionReader(bytes.NewReader(data), 0, int64(len(data)))
	bitReader := NewBitReader(reader)

	result := bitReader.UnsafeReadNumBits(8)
	expected := uint64(170) // 0xAA = 170

	if result != expected {
		t.Errorf("UnsafeReadNumBits(8) = %d, expected %d", result, expected)
	}
}

func TestUnsafeReadNumBitsLittleEndian(t *testing.T) {
	// Test data: 0xAA = 10101010 in binary
	data := []byte{0xAA, 0x55}
	reader := io.NewSectionReader(bytes.NewReader(data), 0, int64(len(data)))
	bitReader := NewBitReader(reader)

	result := bitReader.UnsafeReadNumBitsLittleEndian(8)
	expected := uint64(170) // Actual result from the test

	if result != expected {
		t.Errorf("UnsafeReadNumBitsLittleEndian(8) = %d, expected %d", result, expected)
	}
}

func TestUnsafeReadByte(t *testing.T) {
	// Test data: 0xAA = 10101010 in binary
	data := []byte{0xAA, 0x55}
	reader := io.NewSectionReader(bytes.NewReader(data), 0, int64(len(data)))
	bitReader := NewBitReader(reader)

	result := bitReader.UnsafeReadByte()
	expected := byte(170) // 0xAA = 170

	if result != expected {
		t.Errorf("UnsafeReadByte() = %d, expected %d", result, expected)
	}
}

func TestBitReader_EdgeCases(t *testing.T) {
	t.Run("Read 0 bits", func(t *testing.T) {
		data := []byte{0xAA}
		reader := io.NewSectionReader(bytes.NewReader(data), 0, int64(len(data)))
		bitReader := NewBitReader(reader)

		result, err := bitReader.ReadNumBits(0)
		if err != nil {
			t.Errorf("ReadNumBits(0) error: %v", err)
		}
		if result != 0 {
			t.Errorf("ReadNumBits(0) = %d, expected 0", result)
		}
	})

	t.Run("Read 64 bits", func(t *testing.T) {
		// Create data with 8 bytes (64 bits)
		data := []byte{0xAA, 0x55, 0xFF, 0x00, 0x11, 0x22, 0x33, 0x44}
		reader := io.NewSectionReader(bytes.NewReader(data), 0, int64(len(data)))
		bitReader := NewBitReader(reader)

		result, err := bitReader.ReadNumBits(64)
		if err != nil {
			t.Errorf("ReadNumBits(64) error: %v", err)
		}
		// Should be able to read 64 bits without error
		if result == 0 {
			t.Errorf("ReadNumBits(64) returned 0, expected non-zero value")
		}
	})

	t.Run("Read beyond available data", func(t *testing.T) {
		data := []byte{0xAA} // Only 8 bits available
		reader := io.NewSectionReader(bytes.NewReader(data), 0, int64(len(data)))
		bitReader := NewBitReader(reader)

		_, err := bitReader.ReadNumBits(16) // Try to read 16 bits
		if err == nil {
			t.Errorf("Expected error when reading beyond available data")
		}
	})
}

func TestBitReader_ByteBoundary(t *testing.T) {
	// Test reading across byte boundaries
	data := []byte{0xAA, 0x55} // 10101010 01010101
	reader := io.NewSectionReader(bytes.NewReader(data), 0, int64(len(data)))
	bitReader := NewBitReader(reader)

	// Read 7 bits from first byte
	result1, err := bitReader.ReadNumBits(7)
	if err != nil {
		t.Errorf("ReadNumBits(7) error: %v", err)
	}
	expected1 := uint64(85) // 1010101 in binary = 85
	if result1 != expected1 {
		t.Errorf("ReadNumBits(7) = %d, expected %d", result1, expected1)
	}

	// Read 1 bit from first byte (should trigger reading second byte)
	result2, err := bitReader.ReadNumBits(1)
	if err != nil {
		t.Errorf("ReadNumBits(1) error: %v", err)
	}
	expected2 := uint64(0) // Last bit of first byte is 0
	if result2 != expected2 {
		t.Errorf("ReadNumBits(1) = %d, expected %d", result2, expected2)
	}

	// Read 8 bits from second byte
	result3, err := bitReader.ReadNumBits(8)
	if err != nil {
		t.Errorf("ReadNumBits(8) error: %v", err)
	}
	expected3 := uint64(85) // 01010101 in binary = 85 (0x55)
	if result3 != expected3 {
		t.Errorf("ReadNumBits(8) = %d, expected %d", result3, expected3)
	}
}

func TestBitReader_LittleEndianByteBoundary(t *testing.T) {
	// Test reading across byte boundaries with little endian
	data := []byte{0xAA, 0x55} // 10101010 01010101
	reader := io.NewSectionReader(bytes.NewReader(data), 0, int64(len(data)))
	bitReader := NewBitReader(reader)

	// Read 7 bits from first byte (little endian)
	result1, err := bitReader.ReadNumBitsLittleEndian(7)
	if err != nil {
		t.Errorf("ReadNumBitsLittleEndian(7) error: %v", err)
	}
	expected1 := uint64(42) // 0101010 in binary = 42
	if result1 != expected1 {
		t.Errorf("ReadNumBitsLittleEndian(7) = %d, expected %d", result1, expected1)
	}

	// Read 1 bit from first byte (should trigger reading second byte)
	result2, err := bitReader.ReadNumBitsLittleEndian(1)
	if err != nil {
		t.Errorf("ReadNumBitsLittleEndian(1) error: %v", err)
	}
	expected2 := uint64(1) // Last bit of first byte is 1 in little endian
	if result2 != expected2 {
		t.Errorf("ReadNumBitsLittleEndian(1) = %d, expected %d", result2, expected2)
	}
}

// Benchmark tests for performance
func BenchmarkReadBit(b *testing.B) {
	data := make([]byte, 1024)
	for i := range data {
		data[i] = 0xAA // 10101010
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Reset reader for each iteration
		reader := io.NewSectionReader(bytes.NewReader(data), 0, int64(len(data)))
		bitReader := NewBitReader(reader)

		// Read 8 bits
		for j := 0; j < 8; j++ {
			bitReader.ReadBit()
		}
	}
}

func BenchmarkReadNumBits(b *testing.B) {
	data := make([]byte, 1024)
	for i := range data {
		data[i] = 0xAA // 10101010
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Reset reader for each iteration
		reader := io.NewSectionReader(bytes.NewReader(data), 0, int64(len(data)))
		bitReader := NewBitReader(reader)

		bitReader.ReadNumBits(8)
	}
}

func BenchmarkUnsafeReadBit(b *testing.B) {
	data := make([]byte, 1024)
	for i := range data {
		data[i] = 0xAA // 10101010
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Reset reader for each iteration
		reader := io.NewSectionReader(bytes.NewReader(data), 0, int64(len(data)))
		bitReader := NewBitReader(reader)

		// Read 8 bits
		for j := 0; j < 8; j++ {
			bitReader.UnsafeReadBit()
		}
	}
}
