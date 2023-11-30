# goDB

goDB is a simple, persistent key/value store written in Golang for learning purposes.

## Overview

goDB is designed to provide a basic understanding of key concepts in database storage, such as LSM Tree principles, memtable, and SST files. It serves as a practical educational tool for exploring the fundamentals of building a storage engine.

## Features

- **Pure Go Implementation:** Written entirely in Go, making it easy to understand and modify for learning purposes.
- **Key/Value Storage:** Provides a straightforward interface for storing and retrieving key/value pairs.
- **Persistence:** Utilizes LSM Tree concepts to ensure data persistence between program runs.
- **Educational Focus:** Designed with simplicity and clarity to aid learning about storage engine fundamentals.

## Getting Started

Follow these steps to get started with goDB:


##  Overall Architecture

![goDB Architecture](https://github.com/AminIdr/goDB/blob/main/images/KV%20Architecture.png?raw=true)

The goDB architecture enables users to interact with the application API through HTTP GET and POST requests. When a user performs a set or delete operation on a key, the operation is logged in both the in-memory table (memtable) and the Write-Ahead Log (WAL) before returning an OK status to the user.

For get operations, the program first checks if the key exists in the memtable. If found, the corresponding value is returned. If not, the program searches the SST files, starting with the newest and moving to the oldest.

To manage memory and disk usage efficiently, the memtable is periodically flushed to disk as an SST file when it exceeds a specified limit. Additionally, a compaction process is triggered when the number of SST files reaches five. During compaction, these files are merged into a single, larger SST file, ensuring a duplicate-free and optimized storage structure while removing keys that have been deleted.

This architecture ensures both durability and efficient retrieval of key/value pairs in goDB, providing a robust foundation for a persistent key/value storage engine.

##  SST File Format

![goDB Architecture](https://github.com/AminIdr/goDB/blob/main/images/SST%20File%20Format.png)

##  Key-Value Entry Format

![goDB Architecture](https://github.com/AminIdr/goDB/blob/main/images/Entry%20Format.png?raw=true)

## Recovery Mechanism

In the event of a system crash or unexpected termination, goDB employs a robust recovery mechanism to ensure data consistency and integrity.

### Recovery during Memtable Flushing

When the system crashes during the process of flushing the memtable to an SST file, the Write-Ahead Logging (WAL) file remains intact. Upon restarting the program, the first check involves examining the existence of the WAL file. If present, it indicates a previous crash. The recovery process involves reading entries from the WAL and populating the memtable. Once the flushing operation is successfully completed, the WAL file is automatically deleted.

### Recovery during Compaction

In each flush operation, after writing to a new SST file, goDB monitors the total number of SST files. If the count surpasses a predefined threshold (referred to as `compactingSize`), a compaction process is triggered. This involves merging the corresponding SST files into a single, larger SST file, ensuring data integrity and reducing redundancy.

In the event of a system crash during the compaction process, the remaining SST files are guaranteed to be intact. The compaction process is designed as a simulation, creating a new treemap by merging the oldest to newest SST files. This new treemap is then written to a buffer and subsequently to the newly compacted SST file in an atomic operation. Finally, the old SST files are removed. This approach guarantees that if a crash occurs during compaction, the old SST files remain present, maintaining the overall consistency of the data store.

By incorporating these recovery mechanisms, goDB ensures resilience and consistency in the face of unexpected failures.


## Usage
To set a key to a value, you can run the following command in Windows command line:

`curl -X POST -H "Content-Type: application/json" -d "{\"key\": \"yourKey\", \"value\": \"yourValue\"}" http://localhost:8080/set`

To delete a key, you can run the following command in Windows command line:

`curl http://localhost:8080/del?key=yourKey`

To get a key, you can use the following command in Windows:

`curl http://localhost:8080/get?key=yourKey`
