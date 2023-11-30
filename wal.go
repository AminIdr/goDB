// wal/wal.go
package main

import (
	"errors"
	"os"
)

// appendToWAL appends a key-value entry to the Write-Ahead Logging file.
// It serializes the key and value, appends the entry to the WAL file, and returns any encountered error.
func (mem *fileDB) appendToWAL(key string, val value) error {

	// Check if the WAL exists. Otherwise, create it in append-only mode.*
	if _, err := os.Stat(walFileName); os.IsNotExist(err) {
		mem.wal, _ = os.OpenFile(walFileName, os.O_APPEND|os.O_CREATE, 0755)
	}

	res := kvToEntry(key, val)
	if _, err := mem.wal.Write(res); err != nil {
		return err
	}
	return nil
}

// recoverWAL reads the contents of the Write-Ahead Log file and write its entries to the Memtable.
// If the Memtable size exceeds a limit, it triggers a flush to disk.
// Returns any encountered error during the recovery process.
func (mem *fileDB) recoverWAL() error {
	// Read the WAL file
	wal, err := os.ReadFile(walFileName)
	if err != nil {
		return err
	}
	if len(wal) > 0 {
		position := 0
		for position < len(wal) {

			flag, keyBytes, valueBytes := entryToKv(wal, &position)

			mem.values.Put(string(keyBytes), value{
				flag,
				valueBytes,
			})
			// Check if the number
			if mem.values.Size() == memLimit {
				if err := mem.flush(); err != nil {
					return errors.New("Error while flushing Memtable to disk.")
				}
			}
		}
		return nil
	}
	return errors.New("Empty WAL file")
}
