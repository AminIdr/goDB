package main

import (
	"os"

	"github.com/emirpasic/gods/maps/treemap"
)

// value represents a key-value pair with a flag indicating the operation type.
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

type fileDB struct {
	values *treemap.Map
	wal    *os.File
}

// newDB creates a new fileDB instance with an empty Memtable and an open Write-Ahead Log file.
// Returns the initialized fileDB and any encountered error.
func newDB() (*fileDB, error) {
	values := treemap.NewWithStringComparator()

	wal, err := os.OpenFile(walFileName, os.O_APPEND|os.O_CREATE, 0755)
	if err != nil {
		return nil, err
	}

	return &fileDB{

		values: values,
		wal:    wal,
	}, nil
}
