package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/emirpasic/gods/maps/treemap"
)

// flush writes the Memtable content to an SST file on disk, clears the Memtable,
// removes the Write-Ahead Log, and triggers compaction if the SST file count exceeds compactingSize.
// Returns any encountered error during the process.
func (mem *fileDB) flush() error {
	buffer := writeToBuffer(mem.values, true)
	compressedData, _ := compress(buffer.Bytes()) // Compress the buffer
	// Create the SST file
	sstFileName := fmt.Sprintf(sstFileName, strconv.FormatInt(time.Now().UnixNano(), 10))
	sstFile, err := os.Create(sstFileName)
	if err != nil {
		return err
	}
	// Write the entire buffer to the file in a single operation
	if _, err := sstFile.Write(compressedData); err != nil {
		return err
	}
	sstFile.Close()
	// If the program crashes while writing to the SST file, the WAL won't be deleted
	// When starting the program, the first thing to check is the existence of the WAL
	// If it exists, we call recoverWAL()
	mem.wal.Close()
	if err := os.Remove(walFileName); err != nil {
		return nil
	}

	// Clear the MemTable
	mem.values.Clear()

	// Check if SST files need to be compacted
	pattern := "db_*.sst"
	matchingFiles, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}
	if len(matchingFiles) >= compactingSize {
		if err := compact(matchingFiles); err != nil {
			return err
		}
	}
	return nil
}

// compact merges multiple SST files into one, removing duplicates and deleted keys.
// It reads each SST file, builds a new treemap, and writes the compacted data to a new SST file.
// Returns any encountered error during the compaction process.
func compact(matchingFiles []string) error {
	// Since insertion in a sorted key-value treemap is in O(log(n)), the complexity of this compaction is O(nlog(n))

	// Create a new temporary map
	tmp := treemap.NewWithStringComparator()

	// Iterate through the SST files from the oldest SST to the newest one
	for i := 0; i < len(matchingFiles); i++ {
		file := matchingFiles[i]
		fileContent, err := os.ReadFile(file)
		if err != nil {
			fmt.Println("Error in reading file")
			return err
		}

		// Decompress the file content
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
			continue // Corrupted file
		}

		entryCount := binary.LittleEndian.Uint32(fileContent[position : position+4])
		position += 4
		sKeyLength := binary.LittleEndian.Uint32(fileContent[position : position+4])
		position += 4 + int(sKeyLength) // Skip the smallest key
		lKeyLength := binary.LittleEndian.Uint32(fileContent[position : position+4])
		position += 4 + int(lKeyLength) // Skip the largest key

		// Check the checksum
		retrievedVersion := binary.LittleEndian.Uint16(fileContent[position : position+2])
		position += 2
		if retrievedVersion != version {
			fmt.Println("This file is of an old version")
			continue // Old version
		}

		// Simulate the execution of the SST file in the temporary map
		for i := 0; i < int(entryCount); i++ {

			flag, keyBytes, valueBytes := entryToKv(fileContent, &position)

			if flag == del {
				tmp.Remove(string(keyBytes))
			} else {
				tmp.Put(string(keyBytes), value{
					set,
					valueBytes,
				})
			}

		}
	}

	// Write the temporary map to the buffer
	buffer := writeToBuffer(tmp, false)

	// Create a new compacted SST file
	sstFileName := fmt.Sprintf(sstFileName, strconv.FormatInt(time.Now().UnixNano(), 10))
	sstFile, err := os.Create(sstFileName)
	if err != nil {
		fmt.Println("Error in creating SST file")
		return err
	}
	// Compress the buffer
	compressedData, _ := compress(buffer.Bytes())
	// Write the entire buffer to the file in a single operation
	if _, err := sstFile.Write(compressedData); err != nil {
		return err
	}
	sstFile.Close()

	// Remove the compacted SST files at the end to ensure consistency if the system crashes
	for _, file := range matchingFiles {
		if err := os.Remove(file); err != nil {
			return err
		}
	}

	return nil
}
