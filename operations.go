package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// Set stores a key-value pair in the Memtable and appends the operation to the Write-Ahead Log.
// If the Memtable size exceeds memLimit, it triggers a flush to disk.
// Returns any encountered error during the process.
func (mem *fileDB) Set(key string, val []byte) error {
	v := value{
		set,
		val,
	}
	if err := mem.appendToWAL(key, v); err != nil {
		return err
	}
	mem.values.Put(key, v)

	if mem.values.Size() == memLimit {
		if err := mem.flush(); err != nil {
			return errors.New("Error while flushing Memtable to disk.")
		}
	}

	return nil
}

// Get retrieves the value for a given key.
// It first checks in the Memtable. If not found, it looks in SST files from the newest to the oldest one.
// Returns the value associated with the key and any encountered error.
func (mem *fileDB) Get(key string) ([]byte, error) {
	if v, ok := mem.values.Get(key); ok { // Check the existence in the Memtable
		if v.(value).flag == del { // Check if it was deleted
			return nil, errors.New("Key not found")
		}
		return v.(value).val, nil
	}
	// Not found. Check in SST files
	pattern := "db_*.sst"
	matchingFiles, err := filepath.Glob(pattern) // Glob() returns a sorted []string.
	if err != nil {
		return nil, err
	}

	// Iterate through the SST files from the newest to the oldest one
	for i := len(matchingFiles) - 1; i >= 0; i-- {
		file := matchingFiles[i]

		// Read the entire content of the current file into a byte slice.
		fileContent, err := os.ReadFile(file)

		// Unlike using os.Open(file), which requires disk access each time to read parts of the file,
		// os.ReadFile(file) efficiently reads the entire file into memory in a single operation.

		if err != nil {
			return nil, err
		}

		fileContent, _ = decompress(fileContent)

		position := 0
		// Check the magic number
		retrievedMag := binary.LittleEndian.Uint32(fileContent[position : position+4])
		position += 4
		if retrievedMag != magicNumber {
			fmt.Println("This file should not be considered")
			continue // Move to the next file
		}

		// Check the checksum
		checksum := calculateChecksum(fileContent[0 : len(fileContent)-4])
		retrievedChecksum := fileContent[len(fileContent)-4:]

		if !bytes.Equal(checksum, retrievedChecksum) {
			fmt.Println("This was was corrupted")
			continue
		}

		entryCount := binary.LittleEndian.Uint32(fileContent[position : position+4])
		position += 4

		sKeyLength := binary.LittleEndian.Uint32(fileContent[position : position+4])
		position += 4
		sKey := fileContent[position : position+int(sKeyLength)]
		position += int(sKeyLength)

		lKeyLength := binary.LittleEndian.Uint32(fileContent[position : position+4])
		position += 4
		lKey := fileContent[position : position+int(lKeyLength)]
		position += int(lKeyLength)

		// Check if the key falls within the range of the current file.
		if key < string(sKey) || key > string(lKey) {
			// If the key is outside the range of the current file, skip to the next file.
			continue
			// The lexicographical sorting of SST files ensures that we can efficiently
			// determine whether the key is within the file's range, leveraging the file order.
		}

		// Check the version
		retrievedVersion := binary.LittleEndian.Uint16(fileContent[position : position+2])
		position += 2
		if retrievedVersion != version {
			fmt.Println("This file is of an old version")
			continue // Old version
		}

		for i := 0; i < int(entryCount); i++ {
			// Each SST file corresponds to a single TreeMap.
			// Consequently, if the key exists, there is precisely one entry in the SST file associated with that key.
			// Therefore, upon finding the key, we can promptly return its corresponding value.
			flag, keyBytes, valueBytes := entryToKv(fileContent, &position)
			if string(keyBytes) == key {
				if flag == set {
					return valueBytes, nil
				}
				return nil, errors.New("Key not found")
			}
		}
	}
	return nil, errors.New("Key not found")
}

// Del writes a delete entry to Memtable and appends it to the Write-Ahead Log.
// If the Memtable size exceeds a limit, it triggers a flush to disk.
// Returns the deleted value and any encountered error during the process.
func (mem *fileDB) Del(key string) ([]byte, error) {
	if val, err := mem.Get(key); err != nil { // Check the existence of the key
		return nil, err
	} else {
		v := value{
			del,
			val,
		}
		if err := mem.appendToWAL(key, v); err != nil {
			return nil, err
		}
		mem.values.Put(key, v)
		if mem.values.Size() == memLimit {
			if err := mem.flush(); err != nil {
				return nil, errors.New("Error while flushing Memtable to disk.")
			}
		}
		return val, nil
	}
}
