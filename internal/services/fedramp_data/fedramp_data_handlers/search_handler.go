package fedramp_data_handlers

import (
	"strings"

	"github.com/grafana/hackathon-12-mcp-compliance/internal/domain/fedramp"
)

// SearchHandler handles search-related operations
type SearchHandler struct {
	// No ports needed for search operations as they operate on in-memory data
}

// NewSearchHandler creates a new search handler
func NewSearchHandler() *SearchHandler {
	return &SearchHandler{}
}

// HandleSearchControls searches for controls by keyword
func (h *SearchHandler) HandleSearchControls(cmd fedramp.SearchControlsCommand) []fedramp.Control {
	query := cmd.Query
	var results []fedramp.Control

	for _, family := range cmd.Program.Families {
		for _, control := range family.Controls {
			if contains(control.SearchIndex, query) {
				results = append(results, control)
			}
		}
	}

	return results
}

// Helper function to check if a string contains another string, ignoring case
func contains(s, substr string) bool {
	s, substr = toLower(s), toLower(substr)
	return s != "" && substr != "" && strings.Contains(s, substr)
}

// Helper function to convert a string to lowercase
func toLower(s string) string {
	return strings.ToLower(s)
}
