package evidence

import (
	"context"
	"fmt"

	"github.com/grafana/mcp-compliance/internal/company"
	"github.com/grafana/mcp-compliance/internal/programs"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// RegisterCompanyTools registers company-specific evidence tools with the MCP server
func RegisterCompanyTools(s *server.MCPServer, registry *programs.Registry) error {
	// Load company evidence practices
	practices, err := company.NewEvidencePractices()
	if err != nil {
		return fmt.Errorf("failed to load company evidence practices: %v", err)
	}

	// Register tool for getting company-specific evidence practices
	registerCompanyEvidencePracticeTool(s, registry, practices)

	// Register tool for searching company-specific evidence practices
	registerSearchCompanyEvidenceTool(s, registry, practices)

	return nil
}

// Register a tool for getting company-specific evidence practices
func registerCompanyEvidencePracticeTool(s *server.MCPServer, registry *programs.Registry, practices *company.EvidencePractices) {
	tool := mcp.NewTool("get_company_evidence_practice",
		mcp.WithDescription("Get company-specific evidence collection practices for a control"),
		mcp.WithString("controlId",
			mcp.Required(),
			mcp.Description("The ID of the security control (e.g., AC-1, IA-2)"),
		),
	)

	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Get control ID from parameters
		controlID, ok := request.Params.Arguments["controlId"].(string)
		if !ok || controlID == "" {
			return mcp.NewToolResultError("Control ID is required"), nil
		}

		// Get the company evidence practice
		practice, err := practices.GetPractice(controlID)
		if err != nil {
			// If no company practice is found, try to get the general evidence guidance
			program, exists := registry.GetProgram("FedRAMP High")
			if !exists {
				return mcp.NewToolResultError(fmt.Sprintf("No company evidence practice found for control: %s", controlID)), nil
			}

			// Get evidence guidance from the program
			guidance, err := program.GetEvidenceGuidance(controlID)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("No evidence guidance found for control: %s", controlID)), nil
			}

			// Format a response that indicates this is general guidance, not company-specific
			text := fmt.Sprintf("No company-specific evidence practice found for %s. Here is general guidance:\n\n", controlID)
			text += fmt.Sprintf("Description: %s\n\n", guidance.Description)

			text += "Recommended Evidence Types:\n"
			for _, evidenceType := range guidance.EvidenceTypes {
				text += fmt.Sprintf("- %s\n", evidenceType)
			}
			text += "\n"

			text += "Suggested Collection Steps:\n"
			for i, step := range guidance.CollectionSteps {
				text += fmt.Sprintf("%d. %s\n", i+1, step)
			}
			text += "\n"

			text += "Examples:\n"
			for _, example := range guidance.Examples {
				text += fmt.Sprintf("- %s\n", example)
			}

			return mcp.NewToolResultText(text), nil
		}

		// Format the practice as text
		text := company.FormatPractice(practice)

		return mcp.NewToolResultText(text), nil
	})
}

// Register a tool for searching company-specific evidence practices
func registerSearchCompanyEvidenceTool(s *server.MCPServer, registry *programs.Registry, practices *company.EvidencePractices) {
	tool := mcp.NewTool("search_company_evidence_practices",
		mcp.WithDescription("Search for company-specific evidence collection practices by keyword"),
		mcp.WithString("query",
			mcp.Required(),
			mcp.Description("The search query (e.g., 'policy', 'authentication', 'quarterly')"),
		),
	)

	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Get search query from parameters
		query, ok := request.Params.Arguments["query"].(string)
		if !ok || query == "" {
			return mcp.NewToolResultError("Search query is required"), nil
		}

		// Search for practices matching the query
		results := practices.SearchPractices(query)

		// Format the results as text
		text := fmt.Sprintf("Company Evidence Practices Search Results for '%s'\n\n", query)

		if len(results) == 0 {
			text += "No company evidence practices found matching the query.\n"
			return mcp.NewToolResultText(text), nil
		}

		for i, practice := range results {
			if i > 0 {
				text += "\n---\n\n"
			}

			text += fmt.Sprintf("Control: %s\n", practice.ControlID)
			text += fmt.Sprintf("Responsible Team: %s\n", practice.ResponsibleTeam)
			text += fmt.Sprintf("Review Frequency: %s\n", practice.ReviewFrequency)
			text += fmt.Sprintf("Summary: %s\n", practice.Practice)
		}

		return mcp.NewToolResultText(text), nil
	})
}
