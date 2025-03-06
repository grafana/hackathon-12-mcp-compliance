package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/grafana/mcp-compliance/internal/controls"
	"github.com/grafana/mcp-compliance/internal/evidence"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// Create a new MCP server
	s := server.NewMCPServer(
		"FedRAMP Compliance Server",
		"1.0.0",
	)

	// Register resources and tools
	registerResources(s)
	registerTools(s)

	// Setup signal handling for graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigCh
		log.Println("Shutting down server...")
		os.Exit(0)
	}()

	// Start the server using stdio
	log.Println("Starting FedRAMP Compliance Server...")
	if err := server.ServeStdio(s); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

func registerResources(s *server.MCPServer) {
	// Register control information resources
	controls.RegisterResources(s)
}

func registerTools(s *server.MCPServer) {
	// Register evidence collection tools
	evidence.RegisterTools(s)
}
