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

![goDB Architecture](https://github.com/AminIdr/goDB/blob/main/KV%20Architecture.png?raw=true)

The goDB architecture enables users to interact with the application API through HTTP GET and POST requests. When a user performs a set or delete operation on a key, the operation is logged in both the in-memory table (memtable) and the Write-Ahead Log (WAL) before returning an OK status to the user.

For get operations, the program first checks if the key exists in the memtable. If found, the corresponding value is returned. If not, the program searches the SST files, starting with the newest and moving to the oldest.

To manage memory and disk usage efficiently, the memtable is periodically flushed to disk as an SST file when it exceeds a specified limit. Additionally, a compaction process is triggered when the number of SST files reaches five. During compaction, these files are merged into a single, larger SST file, ensuring a duplicate-free and optimized storage structure while removing keys that have been deleted.

This architecture ensures both durability and efficient retrieval of key/value pairs in goDB, providing a robust foundation for a persistent key/value storage engine.

##  SST File Format

![goDB Architecture](https://github.com/AminIdr/goDB/blob/main/SST%20File%20Format.png?raw=true)

##  Key-Value Entry Format

![goDB Architecture](https://github.com/AminIdr/goDB/blob/main/Entry%20Format.png?raw=true)
## Usage
To set a key to a value, you can run the following command in Windows command line:

`curl -X POST -H "Content-Type: application/json" -d "{\"key\": \"yourKey\", \"value\": \"yourValue\"}" http://localhost:8080/set`

To delete a key, you can run the following command in Windows command line:

`curl http://localhost:8080/del?key=yourKey`

To get a key, you can use the following command in Windows:

`curl http://localhost:8080/get?key=yourKey`
