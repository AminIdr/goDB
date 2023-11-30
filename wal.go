// wal/wal.go
package main

import (
	"errors"
	"os"
)

const WalFileName = "db.wal"

// appendToWAL appends a key-value entry to the Write-Ahead Logging file.
// It serializes the key and value, appends the entry to the WAL file, and returns any encountered error.
func (mem *fileDB) appendToWAL(key string, val value) error {

	if _, err := os.Stat(walFileName); os.IsNotExist(err) {
		mem.wal, _ = os.OpenFile(walFileName, os.O_APPEND|os.O_CREATE, 0755)
	}

	res := kvToEntry(key, val)
	if _, err := mem.wal.Write(res); err != nil {
		return err
	}
	return nil
}

// recoverWAL reads the contents of the Write-Ahead Log file and replays the operations to reconstruct the Memtable.
// Returns any encountered error during the recovery process.
func (mem *fileDB) recoverWAL() error {
	// Load WAL to memory and flush it
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
