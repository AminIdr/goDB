package main

import (
	"bytes"
	"testing"
)

func TestCompressDecompress(t *testing.T) {
	testData := []byte("Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.")
	compressedData, err := compress(testData)
	if err != nil {
		t.Errorf("Error during compression: %v", err)
	}

	if len(compressedData) >= len(testData) {
		t.Error("Compression did not reduce size")
	}

	decompressedData, err := decompress(compressedData)
	if err != nil {
		t.Errorf("Error during decompression: %v", err)
	}

	if !bytes.Equal(decompressedData, testData) {
		t.Error("Decompressed data does not match original data")
	}
}
