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

// flush writes the contents of the Memtable to an SST file on disk.
// It also triggers compaction if the number of SST files exceeds a specified limit.
// Returns any encountered error during the process.
func (mem *fileDB) flush() error {
	buffer := writeToBuffer(mem.values, true)
	compressedData, _ := compress(buffer.Bytes())
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
	// Thus when starting the program, the first thing to check is the existence of the WAL
	// If it is the case, we call recover()
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
	// Since insertion in a sorted key-value map is in O(log(n)), the complexity of this compaction is O(nlog(n))

	// var tmp map[string]value
	tmp := treemap.NewWithStringComparator()

	for i := 0; i < len(matchingFiles); i++ { // From the oldest SST to the newest
		file := matchingFiles[i]
		fileContent, err := os.ReadFile(file)
		if err != nil {
			fmt.Println("Error in reading file")
			return err
		}

		// Compression Marker
		fileContent, _ = decompress(fileContent)

		position := 0
		retrievedMag := binary.LittleEndian.Uint32(fileContent[position : position+4])
		position += 4
		if retrievedMag != magicNumber {
			fmt.Println("This file should not be considered")
			// continue // Move to the next file
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
		position += 4 + int(sKeyLength) // Skip the smallest key
		lKeyLength := binary.LittleEndian.Uint32(fileContent[position : position+4])
		position += 4 + int(lKeyLength) // Skip the largest key

		retrievedVersion := binary.LittleEndian.Uint16(fileContent[position : position+2])
		position += 2
		if retrievedVersion != version {
			fmt.Println("This file is of an old version")
			continue // Old version
		}

		for i := 0; i < int(entryCount); i++ {

			flag, keyBytes, valueBytes := entryToKv(fileContent, &position)

			if flag == del {
				tmp.Remove(string(keyBytes))
			} else {
				tmp.Put(string(keyBytes), value{
					flag, // flag = set
					valueBytes,
				})
			}

		}
	}

	buffer := writeToBuffer(tmp, false)

	sstFileName := fmt.Sprintf(sstFileName, strconv.FormatInt(time.Now().UnixNano(), 10))
	sstFile, err := os.Create(sstFileName)
	if err != nil {
		fmt.Println("Error in creating sst file")
		return err
	}
	defer sstFile.Close()
	compressedData, _ := compress(buffer.Bytes())
	// Write the entire buffer to the file in a single operation
	if _, err := sstFile.Write(compressedData); err != nil {
		return err
	}

	for _, file := range matchingFiles {
		if err := os.Remove(file); err != nil {
			return err
		}
	}

	return nil
}
