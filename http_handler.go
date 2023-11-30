package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// handleFunction returns an http.HandlerFunc that routes requests to specific handler functions based on the URL path.
// Supported paths include "/get", "/set", and "/del".
func handleFunction(db *fileDB) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		switch req.URL.Path {
		case "/get":
			handleGet(resp, req, db)
		case "/set":
			handleSet(resp, req, db)
		case "/del":
			handleDelete(resp, req, db)
		default:
			http.Error(resp, "Not Found", http.StatusNotFound)
		}
	}
}

// handleGet is an HTTP handler function for the "/get" endpoint.
// Retrieves the value for a given key and writes it to the response.
func handleGet(resp http.ResponseWriter, req *http.Request, db *fileDB) {
	key := req.URL.Query().Get("key")
	if key == "" {
		http.Error(resp, "Key parameter is missing", http.StatusBadRequest)
		return
	}

	value, err := db.Get(key)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusNotFound)
		return
	}

	resp.Write(value)
}

// handleSet is an HTTP handler function for the "/set" endpoint.
// Parses JSON input, sets the key-value pair in the Memtable, and writes the result to the response.
func handleSet(resp http.ResponseWriter, req *http.Request, db *fileDB) {
	var data map[string]string
	if err := json.NewDecoder(req.Body).Decode(&data); err != nil {
		http.Error(resp, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	key, ok := data["key"]
	if !ok || key == "" {
		http.Error(resp, "Key parameter is missing", http.StatusBadRequest)
		return
	}

	value, ok := data["value"]
	if !ok || value == "" {
		http.Error(resp, "Value parameter is missing", http.StatusBadRequest)
		return
	}

	if err := db.Set(key, []byte(value)); err != nil {
		http.Error(resp, fmt.Sprintf("Error setting key: %s", err), http.StatusInternalServerError)
		return
	}

	resp.Write([]byte("Key set successfully"))
}

// handleDelete is an HTTP handler function for the "/del" endpoint.
// Deletes a key-value pair and writes the result to the response.
func handleDelete(resp http.ResponseWriter, req *http.Request, db *fileDB) {
	key := req.URL.Query().Get("key")
	if key == "" {
		http.Error(resp, "Key parameter is missing", http.StatusBadRequest)
		return
	}

	value, err := db.Del(key)
	if err != nil {
		http.Error(resp, "Key not found", http.StatusNotFound)
		return
	}

	resp.Write([]byte(fmt.Sprintf("Key deleted successfully. Value: %s", value)))
}
