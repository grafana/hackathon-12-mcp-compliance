package ports

import (
	"github.com/grafana/hackathon-12-mcp-compliance/internal/domain/fedramp"
)

// OSCALRepository defines the interface for OSCAL data operations
type OSCALRepository interface {
	// ParseOSCALCatalog parses OSCAL catalog data and returns a structured representation
	ParseOSCALCatalog(data []byte) (fedramp.OSCALCatalog, error)

	// ProcessOSCALCatalog processes an OSCAL catalog into a Program
	ProcessOSCALCatalog(catalog fedramp.OSCALCatalog, programName string) (fedramp.Program, error)

	// SerializeProgram serializes a Program to JSON
	SerializeProgram(program fedramp.Program) ([]byte, error)
}
