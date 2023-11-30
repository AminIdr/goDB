package main

import (
	"os"
	"testing"
)

func TestAppendToWAL(t *testing.T) {
	db, err := newDB()
	if err != nil {
		t.Fatal("Error creating a new DB:", err)
	}
	defer db.wal.Close()

	// Normal case: Append to WAL
	key := "wal_key"
	value := value{
		flag: set,
		val:  []byte("wal_value"),
	}

	if err := db.appendToWAL(key, value); err != nil {
		t.Fatalf("Error appending to WAL: %s", err)
	}

	// Edge case: Append to non-existent WAL file
	db.wal.Close()
	os.Remove(walFileName)

	if err := db.appendToWAL(key, value); err != nil {
		t.Fatalf("Error appending to non-existent WAL file: %s", err)
	}
}
