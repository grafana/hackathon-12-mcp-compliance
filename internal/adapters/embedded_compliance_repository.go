package adapters

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/grafana/hackathon-12-mcp-compliance/internal/domain/fedramp"
	"github.com/grafana/hackathon-12-mcp-compliance/internal/resources"
)

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
	data, err := resources.Data.ReadFile(filePath)
	if err != nil {
		// If the file doesn't exist in the embedded FS, it might not have been processed yet
		// In this case, return a more helpful error message
		return fedramp.Program{}, fmt.Errorf("program data file not found: %s (run 'make run-fedramp-data-%s' to generate it)",
			filePath, strings.ToLower(strings.Split(programName, " ")[1]))
	}

	// Unmarshal the JSON data
	var program fedramp.Program
	if err := json.Unmarshal(data, &program); err != nil {
		return fedramp.Program{}, fmt.Errorf("failed to parse program data: %v", err)
	}

	return program, nil
}
