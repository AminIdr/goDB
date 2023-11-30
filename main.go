package main

import (
	"fmt"
	"net/http"
	"os"
)

const (
	set            = byte(0)
	del            = byte(1)
	walFileName    = "db.wal"
	sstFileName    = "db_%s.sst"
	magicNumber    = 1234
	version        = uint16(1)
	memLimit       = 10
	compactingSize = 5
)

// main is the entry point of the application.
// Initializes a new fileDB, recovers the Memtable from the Write-Ahead Log, and starts an HTTP server.
func main() {
	db, err := newDB()
	if err != nil {
		fmt.Println("Error in the WAL:", err)
		return
	}
	defer db.wal.Close()
	defer db.flush()

	if _, err := os.Stat(walFileName); err == nil {
		db.recoverWAL()
	}

	http.HandleFunc("/", handleFunction(db))
	http.ListenAndServe(":8080", nil)
}
