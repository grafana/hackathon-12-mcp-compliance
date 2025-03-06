package evidence

import (
	"context"
	"fmt"

	"github.com/grafana/mcp-compliance/internal/programs"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// RegisterTools registers all evidence-related tools with the MCP server
func RegisterTools(s *server.MCPServer, registry *programs.Registry) {
	// Register tool for getting evidence guidance for a specific control
	registerEvidenceGuidanceTool(s, registry)

	// Register tool for searching evidence guidance
	registerEvidenceSearchTool(s, registry)

	// Register company-specific evidence tools
	err := RegisterCompanyTools(s, registry)
	if err != nil {
		// Log the error but continue
		fmt.Printf("Warning: Failed to register company evidence tools: %v\n", err)
	}
}

// Register a tool for getting evidence guidance for a specific control
func registerEvidenceGuidanceTool(s *server.MCPServer, registry *programs.Registry) {
	tool := mcp.NewTool("get_evidence_guidance",
		mcp.WithDescription("Get guidance on collecting evidence for a specific FedRAMP security control"),
		mcp.WithString("controlId",
			mcp.Required(),
			mcp.Description("The ID of the security control (e.g., AC-1, IA-2)"),
		),
		mcp.WithString("program",
			mcp.Description("The compliance program (e.g., FedRAMP High)"),
		),
	)

	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Get control ID from parameters
		controlID, ok := request.Params.Arguments["controlId"].(string)
		if !ok || controlID == "" {
			return mcp.NewToolResultError("Control ID is required"), nil
		}

		// Get program name from parameters (default to FedRAMP High)
		programName, ok := request.Params.Arguments["program"].(string)
		if !ok || programName == "" {
			programName = "FedRAMP High"
		}

		// Get the program
		program, exists := registry.GetProgram(programName)
		if !exists {
			return mcp.NewToolResultError(fmt.Sprintf("Compliance program not found: %s", programName)), nil
		}

		// Get evidence guidance
		guidance, err := program.GetEvidenceGuidance(controlID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get evidence guidance: %v", err)), nil
		}

		// Format the guidance as text
		text := fmt.Sprintf("Evidence Guidance for %s in %s\n\n", controlID, programName)
		text += fmt.Sprintf("Description: %s\n\n", guidance.Description)

		text += "Evidence Types:\n"
		for _, evidenceType := range guidance.EvidenceTypes {
			text += fmt.Sprintf("- %s\n", evidenceType)
		}
		text += "\n"

		text += "Collection Steps:\n"
		for i, step := range guidance.CollectionSteps {
			text += fmt.Sprintf("%d. %s\n", i+1, step)
		}
		text += "\n"

		text += "Examples:\n"
		for _, example := range guidance.Examples {
			text += fmt.Sprintf("- %s\n", example)
		}
		text += "\n"

		text += "Common Pitfalls:\n"
		for _, pitfall := range guidance.CommonPitfalls {
			text += fmt.Sprintf("- %s\n", pitfall)
		}

		return mcp.NewToolResultText(text), nil
	})
}

// Register a tool for searching evidence guidance
func registerEvidenceSearchTool(s *server.MCPServer, registry *programs.Registry) {
	tool := mcp.NewTool("search_evidence_guidance",
		mcp.WithDescription("Search for evidence guidance by keyword"),
		mcp.WithString("query",
			mcp.Required(),
			mcp.Description("The search query (e.g., 'policy', 'authentication')"),
		),
		mcp.WithString("program",
			mcp.Description("The compliance program (e.g., FedRAMP High)"),
		),
	)

	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Get search query from parameters
		query, ok := request.Params.Arguments["query"].(string)
		if !ok || query == "" {
			return mcp.NewToolResultError("Search query is required"), nil
		}

		// Get program name from parameters (default to FedRAMP High)
		programName, ok := request.Params.Arguments["program"].(string)
		if !ok || programName == "" {
			programName = "FedRAMP High"
		}

		// Get the program
		program, exists := registry.GetProgram(programName)
		if !exists {
			return mcp.NewToolResultError(fmt.Sprintf("Compliance program not found: %s", programName)), nil
		}

		// Search for controls matching the query
		controls, err := program.SearchControls(query)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to search controls: %v", err)), nil
		}

		// Format the results as text
		text := fmt.Sprintf("Evidence Guidance Search Results for '%s' in %s\n\n", query, programName)

		if len(controls) == 0 {
			text += "No controls found matching the query.\n"
			return mcp.NewToolResultText(text), nil
		}

		for _, control := range controls {
			text += fmt.Sprintf("Control: %s - %s\n", control.ID, control.Title)

			// Get evidence guidance for this control
			guidance, err := program.GetEvidenceGuidance(control.ID)
			if err != nil {
				text += "  No evidence guidance available for this control.\n\n"
				continue
			}

			text += fmt.Sprintf("  Description: %s\n", guidance.Description)
			text += "  Evidence Types: "
			for i, evidenceType := range guidance.EvidenceTypes {
				if i > 0 {
					text += ", "
				}
				text += evidenceType
			}
			text += "\n\n"
		}

		return mcp.NewToolResultText(text), nil
	})
}
