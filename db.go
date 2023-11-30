package main

import (
	"os"

	"github.com/emirpasic/gods/maps/treemap"
)

// value represents a key-value pair with a flag indicating the operation type (set or del).
type value struct {
	flag byte
	val  []byte
}

// DB is an interface that defines basic operations for a key-value store.
type DB interface {
	Set(key string, value []byte) error

	Get(key string) ([]byte, error)

	Del(key string) ([]byte, error)
}

// fileDB is a key-value store that uses a TreeMap for in-memory storage and
// maintains a Write-Ahead Log (WAL) file for durability.
type fileDB struct {
	values *treemap.Map
	wal    *os.File
}

// newDB creates a new fileDB instance with an empty Memtable and an open Write-Ahead Log file.
// Returns the initialized fileDB and any encountered error.
func newDB() (*fileDB, error) {
	// Create a TreeMap with a string comparator for in-memory storage
	values := treemap.NewWithStringComparator()

	// Create the WAL file in append-only mode if it does not exist
	wal, err := os.OpenFile(walFileName, os.O_APPEND|os.O_CREATE, 0755)
	if err != nil {
		return nil, err
	}

	// Return the initialized fileDB instance
	return &fileDB{
		values: values,
		wal:    wal,
	}, nil
}
