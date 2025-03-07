package ports

import (
	"github.com/grafana/hackathon-12-mcp-compliance/internal/domain/fedramp"
)

// ComplianceRepository defines methods for accessing compliance program data
type ComplianceRepository interface {
	// ListPrograms returns a list of available compliance programs
	ListPrograms() ([]string, error)

	// LoadProgram loads a specific compliance program by name
	LoadProgram(programName string) (fedramp.Program, error)
}
