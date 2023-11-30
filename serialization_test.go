package main

import (
	"bytes"
	"testing"
)

func TestCalculateChecksum(t *testing.T) {
	// Test normal case
	data := []byte("test data")
	expectedChecksum := calculateChecksum(data)
	if expectedChecksum == nil || bytes.Equal(expectedChecksum, []byte{0, 0, 0, 0}) {
		t.Error("Expected non-zero checksum for normal case")
	}

	// Test edge case with empty data
	emptyData := []byte{}
	emptyChecksum := calculateChecksum(emptyData)
	if emptyChecksum != nil && !bytes.Equal(emptyChecksum, []byte{0, 0, 0, 0}) {
		t.Error("Expected zero checksum for empty data")
	}

	// Test edge case with nil data
	nilChecksum := calculateChecksum(nil)
	if nilChecksum != nil && !bytes.Equal(nilChecksum, []byte{0, 0, 0, 0}) {
		t.Error("Expected zero checksum for nil data")
	}
}
