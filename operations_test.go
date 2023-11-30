package main

import (
	"testing"
)

func TestSetGetDel(t *testing.T) {
	db, err := newDB()
	if err != nil {
		t.Fatal("Error creating a new DB:", err)
	}
	defer db.wal.Close()

	// Normal case: Set, Get, and Delete a key
	key := "normal_key"
	value := []byte("normal_value")

	if err := db.Set(key, value); err != nil {
		t.Fatalf("Error setting key: %s", err)
	}

	retrievedValue, err := db.Get(key)
	if err != nil {
		t.Fatalf("Error getting key: %s", err)
	}
	if string(retrievedValue) != string(value) {
		t.Fatalf("Expected value %s, got %s", value, retrievedValue)
	}

	deletedValue, err := db.Del(key)
	if err != nil {
		t.Fatalf("Error deleting key: %s", err)
	}
	if string(deletedValue) != string(value) {
		t.Fatalf("Expected deleted value %s, got %s", value, deletedValue)
	}

	// Edge case 1: Get an inexistant key
	emptyValue, err := db.Get("nonexistent_key")
	if err == nil || emptyValue != nil {
		t.Fatal("Expected error getting nonexistent key")
	}

	// Edge case 2: Delete nonexistent key
	deletedValue, err = db.Del("nonexistent_key")
	if err == nil || deletedValue != nil {
		t.Fatal("Expected error deleting nonexistent key")
	}
}
