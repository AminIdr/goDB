package main

import (
	"bytes"
	"compress/gzip"
	"io"
)

// compress compresses a byte slice using gzip compression.
// Returns the compressed byte slice and any encountered error.
func compress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	writer := gzip.NewWriter(&buf)

	if _, err := writer.Write(data); err != nil {
		return nil, err
	}

	if err := writer.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// decompress decompresses a gzip-compressed byte slice.
// Returns the decompressed byte slice and any encountered error.
func decompress(data []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	decompressedData, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	return decompressedData, nil
}
