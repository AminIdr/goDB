package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"hash/crc32"

	"github.com/emirpasic/gods/maps/treemap"
)

// kvToEntry converts a key-value pair to a byte slice for storage.
// Returns the serialized byte slice.
func kvToEntry(key string, val value) []byte {
	keyLength := make([]byte, 4) // 4 bytes for the key length
	binary.LittleEndian.PutUint32(keyLength, uint32(len([]byte(key))))
	entry := append(keyLength, []byte(key)...)

	flag := val.flag
	// If it's a delete operation, store: flag, key length, key
	if flag == del {
		res := append([]byte{flag}, entry...)
		return res
	}
	// Otherwise, store: flag, key length, key, value length, value
	valueLength := make([]byte, 4) // 4 bytes for the value length
	binary.LittleEndian.PutUint32(valueLength, uint32(len(val.val)))
	entry = append(entry, valueLength...)
	entry = append(entry, val.val...)
	res := append([]byte{flag}, entry...) //	1 byte for the flag
	return res
}

// entryToKv converts a byte slice to a key-value pair.
// It extracts the flag, key, and value from the serialized byte slice.
// The position parameter is used to keep track of the parsing position.
// Returns the flag, key, and value.
func entryToKv(entry []byte, position *int) (flag byte, key, val []byte) {
	flag = entry[*position]
	*position++

	keyLength := binary.LittleEndian.Uint32(entry[*position : *position+4])
	*position += 4
	key = entry[*position : *position+int(keyLength)]
	*position += int(keyLength)

	if flag == del {
		return
	}

	valueLength := binary.LittleEndian.Uint32(entry[*position : *position+4])
	*position += 4
	val = entry[*position : *position+int(valueLength)]
	*position += int(valueLength)
	return
}

// calculateChecksum calculates the CRC32 checksum for a given byte slice.
// Returns a 4-byte checksum.
func calculateChecksum(data []byte) []byte {
	crc := crc32.ChecksumIEEE(data)
	checksumBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(checksumBytes, crc)
	return checksumBytes
}

// writeToBuffer creates a byte buffer from the contents of the treemap.
// It includes metadata such as magic number, entry count, smallest and largest keys, version, key-value tuples, and checksum.
// The isFlush parameter indicates whether the function is called from flush() or compact().
// Returns the resulting byte buffer.
func writeToBuffer(treemap *treemap.Map, isFlush bool) (buffer bytes.Buffer) {
	// 4 bytes for the Magic Number
	mag := make([]byte, 4)
	binary.LittleEndian.PutUint32(mag, magicNumber)
	buffer.Write(mag)

	// 4 bytes for the Entry Count
	ent := make([]byte, 4)
	binary.LittleEndian.PutUint32(ent, uint32(treemap.Size()))
	buffer.Write(ent)

	// 4 bytes for the smallest key's length
	sKeyLength := make([]byte, 4)
	sKey, _ := treemap.Min()
	binary.LittleEndian.PutUint32(sKeyLength, uint32(len(fmt.Sprint(sKey))))
	buffer.Write(sKeyLength)

	// Smallest Key
	buffer.Write([]byte(fmt.Sprint(sKey)))

	// 4 bytes for the largest key's length
	lKeyLength := make([]byte, 4)
	lKey, _ := treemap.Max()
	binary.LittleEndian.PutUint32(lKeyLength, uint32(len(fmt.Sprint(lKey))))
	buffer.Write(lKeyLength)

	// Largest Key
	buffer.Write([]byte(fmt.Sprint(lKey)))

	// 2 bytes for the version
	ver := make([]byte, 2)
	binary.LittleEndian.PutUint16(ver, version) // Version 1
	buffer.Write(ver)

	// Key-value entries
	iterator := treemap.Iterator()
	for iterator.Next() {
		key := iterator.Key().(string)
		val := iterator.Value().(value)
		if isFlush || val.flag != del {
			entry := kvToEntry(key, val)
			buffer.Write(entry)
		}
	}

	// Checksum
	checksum := calculateChecksum(buffer.Bytes()) // 4 bytes
	buffer.Write(checksum)

	return buffer
}
