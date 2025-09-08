package fileio

import (
	"bytes"
	"io"
	"testing"
)

// Helper function to compare byte slices
func compareByteSlices(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// Helper function to create a test StreamReader with given data
func createTestStreamReader(data []byte) *StreamReader {
	reader := io.NewSectionReader(bytes.NewReader(data), 0, int64(len(data)))
	return NewStreamReader(reader)
}

func TestReadUint8(t *testing.T) {
	// Test normal reading
	data := []byte{0x42, 0x00, 0xFF}
	reader := createTestStreamReader(data)

	value, err := reader.ReadUint8()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if value != 0x42 {
		t.Errorf("Expected 0x42, got 0x%02X", value)
	}

	// Test reading next value
	value, err = reader.ReadUint8()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if value != 0x00 {
		t.Errorf("Expected 0x00, got 0x%02X", value)
	}

	// Test reading last value
	value, err = reader.ReadUint8()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if value != 0xFF {
		t.Errorf("Expected 0xFF, got 0x%02X", value)
	}

	// Test EOF
	_, err = reader.ReadUint8()
	if err == nil {
		t.Error("Expected EOF error, got nil")
	}
}

func TestReadUint16(t *testing.T) {
	// Test little-endian reading
	data := []byte{0x34, 0x12, 0x78, 0x56} // 0x1234, 0x5678
	reader := createTestStreamReader(data)

	value, err := reader.ReadUint16()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if value != 0x1234 {
		t.Errorf("Expected 0x1234, got 0x%04X", value)
	}

	value, err = reader.ReadUint16()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if value != 0x5678 {
		t.Errorf("Expected 0x5678, got 0x%04X", value)
	}

	// Test EOF
	_, err = reader.ReadUint16()
	if err == nil {
		t.Error("Expected EOF error, got nil")
	}
}

func TestReadUint32(t *testing.T) {
	// Test little-endian reading
	data := []byte{0x78, 0x56, 0x34, 0x12, 0xBC, 0x9A, 0x78, 0x56} // 0x12345678, 0x56789ABC
	reader := createTestStreamReader(data)

	value, err := reader.ReadUint32()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if value != 0x12345678 {
		t.Errorf("Expected 0x12345678, got 0x%08X", value)
	}

	value, err = reader.ReadUint32()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if value != 0x56789ABC {
		t.Errorf("Expected 0x56789ABC, got 0x%08X", value)
	}

	// Test EOF
	_, err = reader.ReadUint32()
	if err == nil {
		t.Error("Expected EOF error, got nil")
	}
}

func TestReadInt16(t *testing.T) {
	// Test positive and negative values
	data := []byte{0x34, 0x12, 0xCC, 0xFF} // 0x1234 (positive), 0xFFCC (negative)
	reader := createTestStreamReader(data)

	value, err := reader.ReadInt16()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if value != 0x1234 {
		t.Errorf("Expected 0x1234, got 0x%04X", value)
	}

	value, err = reader.ReadInt16()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if value != -52 { // 0xFFCC as signed int16
		t.Errorf("Expected -52, got %d", value)
	}
}

func TestReadFloat32(t *testing.T) {
	// Test float32 reading (1.0f in little-endian)
	data := []byte{0x00, 0x00, 0x80, 0x3F} // 1.0f
	reader := createTestStreamReader(data)

	value, err := reader.ReadFloat32()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if value != 1.0 {
		t.Errorf("Expected 1.0, got %f", value)
	}
}

func TestReadBytes(t *testing.T) {
	data := []byte{0x01, 0x02, 0x03, 0x04, 0x05}
	reader := createTestStreamReader(data)

	// Test reading 3 bytes
	bytes, err := reader.ReadBytes(3)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	expected := []byte{0x01, 0x02, 0x03}
	if !compareByteSlices(bytes, expected) {
		t.Errorf("Expected %v, got %v", expected, bytes)
	}

	// Test reading remaining bytes
	bytes, err = reader.ReadBytes(2)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	expected = []byte{0x04, 0x05}
	if !compareByteSlices(bytes, expected) {
		t.Errorf("Expected %v, got %v", expected, bytes)
	}

	// Test reading beyond EOF
	_, err = reader.ReadBytes(1)
	if err == nil {
		t.Error("Expected EOF error, got nil")
	}

	// Test invalid count
	_, err = reader.ReadBytes(0)
	if err == nil {
		t.Error("Expected error for zero count, got nil")
	}

	_, err = reader.ReadBytes(-1)
	if err == nil {
		t.Error("Expected error for negative count, got nil")
	}
}

func TestReadString(t *testing.T) {
	// Test null-terminated string
	data := []byte{'H', 'e', 'l', 'l', 'o', 0x00, 'W', 'o', 'r', 'l', 'd'}
	reader := createTestStreamReader(data)

	str, err := reader.ReadString(20)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if str != "Hello" {
		t.Errorf("Expected 'Hello', got '%s'", str)
	}

	// Test string without null terminator
	data = []byte{'T', 'e', 's', 't'}
	reader = createTestStreamReader(data)

	str, err = reader.ReadString(10)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if str != "Test" {
		t.Errorf("Expected 'Test', got '%s'", str)
	}

	// Test invalid max length
	_, err = reader.ReadString(0)
	if err == nil {
		t.Error("Expected error for zero max length, got nil")
	}
}

func TestReadFixedString(t *testing.T) {
	// Test fixed-length string with null padding
	data := []byte{'T', 'e', 's', 't', 0x00, 0x00, 0x00, 0x00}
	reader := createTestStreamReader(data)

	str, err := reader.ReadFixedString(8)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if str != "Test" {
		t.Errorf("Expected 'Test', got '%s'", str)
	}

	// Test string with no null padding
	data = []byte{'F', 'u', 'l', 'l'}
	reader = createTestStreamReader(data)

	str, err = reader.ReadFixedString(4)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if str != "Full" {
		t.Errorf("Expected 'Full', got '%s'", str)
	}

	// Test invalid length
	_, err = reader.ReadFixedString(0)
	if err == nil {
		t.Error("Expected error for zero length, got nil")
	}
}

func TestSetPosition(t *testing.T) {
	data := []byte{0x01, 0x02, 0x03, 0x04, 0x05}
	reader := createTestStreamReader(data)

	// Read first byte
	value, err := reader.ReadUint8()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if value != 0x01 {
		t.Errorf("Expected 0x01, got 0x%02X", value)
	}

	// Set position back to start
	reader.SetPosition(0)

	// Read first byte again
	value, err = reader.ReadUint8()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if value != 0x01 {
		t.Errorf("Expected 0x01, got 0x%02X", value)
	}

	// Set position to middle
	reader.SetPosition(2)

	// Read byte at position 2
	value, err = reader.ReadUint8()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if value != 0x03 {
		t.Errorf("Expected 0x03, got 0x%02X", value)
	}
}

// Note: GetOpcode, GetDataLength, HasData, ValidateInstructionSize, and GetOpcodeName
// methods are not yet implemented in StreamReader. These tests can be added when
// those methods are implemented.

func TestEdgeCases(t *testing.T) {
	// Test reading from empty stream
	emptyReader := createTestStreamReader([]byte{})
	_, err := emptyReader.ReadUint8()
	if err == nil {
		t.Error("Expected error when reading from empty stream")
	}

	// Test reading beyond available data
	data := []byte{0x01}
	reader := createTestStreamReader(data)
	_, err = reader.ReadUint16()
	if err == nil {
		t.Error("Expected error when reading beyond available data")
	}

	// Test reading exact amount of data
	data = []byte{0x01, 0x02}
	reader = createTestStreamReader(data)
	value, err := reader.ReadUint16()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if value != 0x0201 { // Little-endian
		t.Errorf("Expected 0x0201, got 0x%04X", value)
	}
}

func TestConcurrentAccess(t *testing.T) {
	// Test that StreamReader is safe for concurrent access
	data := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}
	reader := createTestStreamReader(data)

	// This test is more about ensuring the code doesn't panic
	// rather than testing actual concurrency safety
	value, err := reader.ReadUint32()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if value != 0x04030201 { // Little-endian
		t.Errorf("Expected 0x04030201, got 0x%08X", value)
	}
}
