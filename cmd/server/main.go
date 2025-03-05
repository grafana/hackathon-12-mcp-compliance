package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rbrady/example-mcp/internal/server"
)

func main() {
	// Create a new router
	r := mux.NewRouter()

	// Register MCP handler
	r.HandleFunc("/mcp", server.HandleMCPRequest).Methods("POST")

	// Start the HTTP server
	log.Println("MCP Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
