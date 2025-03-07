package main

import (
	"log"

	"github.com/grafana/hackathon-12-mcp-compliance/internal/programs"
	"github.com/grafana/hackathon-12-mcp-compliance/internal/registry"
	internalserver "github.com/grafana/hackathon-12-mcp-compliance/internal/server"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// Create program registry
	r := registry.New()

	// Create and register FedRAMP programs
	high, err := programs.NewFedRAMPHigh()
	if err != nil {
		log.Fatalf("Failed to create FedRAMP High program: %v", err)
	}
	r.Register(high)

	moderate, err := programs.NewFedRAMPModerate()
	if err != nil {
		log.Fatalf("Failed to create FedRAMP Moderate program: %v", err)
	}
	r.Register(moderate)

	// Create and start MCP server
	s := internalserver.NewServer(r)
	if err := server.ServeStdio(s); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
