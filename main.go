package main

import (
	"log/slog"
	"net/http"
)

/*
# Prompt
Create a golang http server for storing unstructured objects by ID.

# Parameters
Any coding tools, googling, etc are allowed. Best practices can be skipped as long as they re called out when they occur.

# Data Format

The Structured fields for an object should be:

ID string
CreatedAt time.Time
Labels []string
Object any

The data storage mechanism should be an in-memory data structure defined in your code.

# API endpoints

- Create: Add an object for ID
- Delete: Remove an object by ID
- List: list objects with query modifiers:
	Ordering:
		- CreatedAt (Default)
		- ID
	Filtering:
		- Labels
*/

func main() {

	// https://go.dev/blog/routing-enhancements
	mux := http.NewServeMux()

	mux.HandleFunc("POST /create", handleCreate)
	mux.HandleFunc("POST /delete/{id}", handleDelete)
	mux.HandleFunc("GET /list", handleList)

	addr := "localhost:8090"
	slog.Info("starting server", slog.Any("addr", addr))

	err := http.ListenAndServe(addr, mux)
	if err != nil {
		slog.Error("listen error", slog.Any("error", err))
	}
}
