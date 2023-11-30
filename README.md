# goDB

goDB is a simple, persistent key/value store written in Golang for learning purposes.

## Overview

goDB is a lightweight, pure-Go implementation of a persistent key-value storage engine inspired by the Log-Structured Merge (LSM) Tree concepts. It serves as a practical educational tool for exploring the fundamentals of building a storage engine.

## Features

- **LSM Tree Architecture:** Leveraging the principles of LSM Trees for efficient storage and retrieval of key-value pairs.
- **Write-Ahead Logging (WAL):** Implements a WAL mechanism to ensure durability and recoverability in the face of crashes.
- **SST File Management:** Stores data in SST (Sorted String Table) files, with automatic compaction to maintain optimal performance.
- **SST File Compression:** Used gzip compression for SST files, effectively saving storage space.
- **HTTP API:** Provides a basic HTTP API for interacting with the key-value store, supporting GET, SET, and DELETE operations.

## Project Structure

- **main.go:** The entry point of the application. Initializes the database, recovers from Write-Ahead Log if needed, and starts an HTTP server.
- **db.go:** Defines the core database structure (`fileDB`) implementing the `DB` interface. Manages the Memtable, Write-Ahead Log (WAL), and provides functions for basic database operations like `Set`, `Get`, and `Del`. Also includes the initialization of the database (`newDB` function).
- **operations.go:** Defines core database operations like `Set`, `Get`, and `Del`. Handles interactions with the Memtable and triggers flushing to disk when necessary.
- **serialization.go:** Provides functions for converting key-value pairs to byte slices and vice versa. Handles the serialization and deserialization of data for storage and retrieval.
- **wal.go:** Manages Write-Ahead Logging, including functions for appending key-value entries to the Write-Ahead Log and recovering from the log during startup.
- **data_maintenance.go:** Handles data maintenance tasks such as flushing Memtable to disk and compacting SST files.
- **compression.go:** Provides functions for compressing and decompressing data, using gzip compression for storage efficiency.
- **http_handler.go:** Defines HTTP handler functions for various endpoints (`/get`, `/set`, `/del`). Parses incoming requests, calls corresponding database operations, and sends responses.



##  Overall Architecture

![goDB Architecture](https://github.com/AminIdr/goDB/blob/main/images/KV%20Architecture.png?raw=true)

The goDB architecture enables users to interact with the application API through HTTP GET and POST requests. When a user performs a set or delete operation on a key, the operation is logged in both the in-memory table (memtable) and the Write-Ahead Log (WAL) before returning an OK status to the user.

For get operations, the program first checks if the key exists in the memtable. If found, the corresponding value is returned. If not, the program searches the SST files, starting with the newest and moving to the oldest.

To manage memory and disk usage efficiently, the memtable is periodically flushed to disk as an SST file when it exceeds a specified limit. Additionally, a compaction process is triggered when the number of SST files reaches five. During compaction, these files are merged into a single, larger SST file, ensuring a duplicate-free and optimized storage structure while removing keys that have been deleted.

This architecture ensures both durability and efficient retrieval of key/value pairs in goDB, providing a robust foundation for a persistent key/value storage engine.

##  SST File Format

![goDB Architecture](https://github.com/AminIdr/goDB/blob/main/images/SST%20File%20Format.png)

##  Set Entry Format

![goDB Architecture](https://github.com/AminIdr/goDB/blob/main/images/Set.png?raw=true)

##  Delete Entry Format

![goDB Architecture](https://github.com/AminIdr/goDB/blob/main/images/Del.png?raw=true)

## Recovery Mechanism

In the event of a system crash or unexpected termination, goDB employs a recovery mechanism to ensure data consistency and integrity.

### Recovery during Memtable Flushing

When the system crashes during the process of flushing the memtable to an SST file, the Write-Ahead Logging (WAL) file remains intact. Upon restarting the program, the first check involves examining the existence of the WAL file. If present, it indicates a previous crash. The recovery process involves reading entries from the WAL and populating the memtable.
Once the flushing operation is successfully completed, the WAL file is automatically deleted.

### Recovery during Compaction

In each flush operation, after writing to a new SST file, goDB monitors the total number of SST files. If the count surpasses a predefined threshold (referred to as `compactingSize`), a compaction process is triggered. This involves merging the corresponding SST files into a single, larger SST file, ensuring data integrity and reducing redundancy.

In the event of a system crash during the compaction process, the remaining SST files are guaranteed to be intact. The compaction process is designed as a simulation, creating a new treemap by merging the oldest to newest SST files. This new treemap is then written to a buffer and subsequently to the newly compacted SST file in an atomic operation. Finally, the old SST files are removed. This approach guarantees that if a crash occurs during compaction, the old SST files remain present, maintaining the overall consistency of the data store.

By incorporating these recovery mechanisms, goDB ensures resilience and consistency in the face of unexpected failures.

## Getting Started

Follow these steps to get started with goDB:

1. Clone the repository: `git clone https://github.com/AminIdr/goDB.git`
2. Build and run the project: `go run .`


## Interact with the Database using Windows cmd

### Set a Key-Value Pair
`curl -X POST -H "Content-Type: application/json" -d "{\"key\": \"yourKey\", \"value\": \"yourValue\"}" http://localhost:8080/set`

### Get the Value for a Key
`curl http://localhost:8080/get?key=yourKey`

### Delete a Key
`curl http://localhost:8080/del?key=yourKey`

## Testing the Program

To test the program, execute the commands in `commands.txt`. This file contains 200 queries, organized as follows:

- 50 queries for setting keys.
- 50 queries for getting keys.
- 50 queries for deleting keys.
- 50 queries for getting keys again.
