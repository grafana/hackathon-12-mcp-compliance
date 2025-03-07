package adapters

import (
	"embed"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/grafana/hackathon-12-mcp-compliance/internal/domain/fedramp"
)

//go:embed data/fedramp-high.json data/fedramp-moderate.json
var embeddedData embed.FS

// EmbeddedComplianceRepository implements the ComplianceRepository interface using embedded data
type EmbeddedComplianceRepository struct {
	// Map of program names to file paths
	programFiles map[string]string
}

// NewEmbeddedComplianceRepository creates a new embedded compliance repository
func NewEmbeddedComplianceRepository() *EmbeddedComplianceRepository {
	return &EmbeddedComplianceRepository{
		programFiles: map[string]string{
			"FedRAMP High":     "data/fedramp-high.json",
			"FedRAMP Moderate": "data/fedramp-moderate.json",
		},
	}
}

// ListPrograms returns a list of available compliance programs
func (r *EmbeddedComplianceRepository) ListPrograms() ([]string, error) {
	programs := make([]string, 0, len(r.programFiles))
	for program := range r.programFiles {
		programs = append(programs, program)
	}
	return programs, nil
}

// LoadProgram loads a specific compliance program by name
func (r *EmbeddedComplianceRepository) LoadProgram(programName string) (fedramp.Program, error) {
	// Find the file path for the program
	filePath, ok := r.programFiles[programName]
	if !ok {
		// Try case-insensitive match
		for name, path := range r.programFiles {
			if strings.EqualFold(name, programName) {
				filePath = path
				ok = true
				break
			}
		}

		if !ok {
			return fedramp.Program{}, fmt.Errorf("program not found: %s", programName)
		}
	}

	// Read the embedded file
	data, err := embeddedData.ReadFile(filePath)
	if err != nil {
		return fedramp.Program{}, fmt.Errorf("failed to read program data: %v", err)
	}

	// Unmarshal the JSON data
	var program fedramp.Program
	if err := json.Unmarshal(data, &program); err != nil {
		return fedramp.Program{}, fmt.Errorf("failed to parse program data: %v", err)
	}

	return program, nil
}
