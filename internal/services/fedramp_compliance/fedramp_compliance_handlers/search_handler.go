package fedramp_compliance_handlers

import (
	"strings"

	"github.com/grafana/hackathon-12-mcp-compliance/internal/domain/fedramp"
	"github.com/grafana/hackathon-12-mcp-compliance/internal/ports"
)

// SearchHandler handles search-related operations
type SearchHandler struct {
	complianceRepo ports.ComplianceRepository
}

// NewSearchHandler creates a new search handler
func NewSearchHandler(complianceRepo ports.ComplianceRepository) *SearchHandler {
	return &SearchHandler{
		complianceRepo: complianceRepo,
	}
}

// HandleSearchControls searches for controls by keyword
func (h *SearchHandler) HandleSearchControls(cmd fedramp.SearchControlsCommand) ([]fedramp.Control, error) {
	// Load the program
	program, err := h.complianceRepo.LoadProgram(cmd.Program.Name)
	if err != nil {
		return nil, err
	}

	// Search for controls
	query := cmd.Query
	var results []fedramp.Control

	for _, family := range program.Families {
		for _, control := range family.Controls {
			if contains(control.SearchIndex, query) {
				results = append(results, control)
			}
		}
	}

	return results, nil
}

// Helper function to check if a string contains another string, ignoring case
func contains(s, substr string) bool {
	s, substr = strings.ToLower(s), strings.ToLower(substr)
	return s != "" && substr != "" && strings.Contains(s, substr)
}
