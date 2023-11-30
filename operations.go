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
// If the Memtable size exceeds a limit, it triggers a flush to disk.
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
// It first checks the Memtable, then searches SST files from newest to oldest.
// Returns the value associated with the key and any encountered error.
func (mem *fileDB) Get(key string) ([]byte, error) {
	if v, ok := mem.values.Get(key); ok {
		if v.(value).flag == del {
			return nil, errors.New("Key was deleted")
		}
		return v.(value).val, nil
	}

	pattern := "db_*.sst"
	matchingFiles, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}

	for i := len(matchingFiles) - 1; i >= 0; i-- { // Newest to oldest
		file := matchingFiles[i]
		fileContent, err := os.ReadFile(file)
		// Unlike os.Open(file), this method reads the whole file once and returns a slice
		// When using os.Open(file), we have to go to disk each time to read a part from the file
		if err != nil {
			return nil, err
		}

		// Compression Marker
		fileContent, _ = decompress(fileContent)

		position := 0
		retrievedMag := binary.LittleEndian.Uint32(fileContent[position : position+4])
		position += 4
		if retrievedMag != magicNumber {
			fmt.Println("This file should not be considered")
			continue // Move to the next file
		}

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

		if key < string(sKey) || key > string(lKey) { // If key is not in the range of the current file
			continue // Move to the next file
			// Here the fact that our SST files are storted lexicographically by keys helped us
		}

		retrievedVersion := binary.LittleEndian.Uint16(fileContent[position : position+2])
		position += 2
		if retrievedVersion != version {
			fmt.Println("This file is of an old version")
			continue // Old version
		}

		for i := 0; i < int(entryCount); i++ {

			flag, keyBytes, valueBytes := entryToKv(fileContent, &position)
			// TODO: key may be deleted and set again
			// 1 SST file corresponds to 1 treemap. Thus, is the key exists, there is exactly one entry in the SST file with that key.
			// So once found, return value
			if string(keyBytes) == key {
				if flag == set {
					return valueBytes, nil
				}
				return nil, errors.New("Key was deleted")
			}
		}
	}
	return nil, errors.New("Key not found")
}

// Del deletes a key-value pair from the Memtable and appends the operation to the Write-Ahead Log.
// If the Memtable size exceeds a limit, it triggers a flush to disk.
// Returns the deleted value and any encountered error during the process.
func (mem *fileDB) Del(key string) ([]byte, error) {
	if val, err := mem.Get(key); err != nil {
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
