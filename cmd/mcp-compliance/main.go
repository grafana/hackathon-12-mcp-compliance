package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/grafana/hackathon-12-mcp-compliance/internal/services/fedramp_compliance"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// Create the compliance service
	complianceService := fedramp_compliance.NewService()

	// Create the MCP server
	s := server.NewMCPServer(
		"MCP Compliance Server",
		"1.0.0",
	)

	// Add tools to the server
	addComplianceTools(s, complianceService)

	// Start the server using stdio
	log.Println("Starting MCP Compliance Server...")
	if err := server.ServeStdio(s); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

// addComplianceTools adds all compliance-related tools to the MCP server
func addComplianceTools(s *server.MCPServer, service *fedramp_compliance.Service) {
	// Tool: list_compliance_programs
	listProgramsTool := mcp.NewTool("list_compliance_programs",
		mcp.WithDescription("List all available compliance programs"),
	)
	s.AddTool(listProgramsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		programs, err := service.ListCompliancePrograms()
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to list compliance programs: %v", err)), nil
		}

		// Format the result as JSON
		programsJSON, err := json.MarshalIndent(programs, "", "  ")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal programs to JSON: %v", err)), nil
		}

		return mcp.NewToolResultText(string(programsJSON)), nil
	})

	// Tool: get_control
	getControlTool := mcp.NewTool("get_control",
		mcp.WithDescription("Get detailed information about a specific control"),
		mcp.WithString("program",
			mcp.Required(),
			mcp.Description("The FedRAMP program (High or Moderate)"),
			mcp.Enum("FedRAMP High", "FedRAMP Moderate"),
		),
		mcp.WithString("controlId",
			mcp.Required(),
			mcp.Description("The ID of the control (e.g., AC-1, IA-2)"),
		),
	)
	s.AddTool(getControlTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		program := request.Params.Arguments["program"].(string)
		controlID := request.Params.Arguments["controlId"].(string)

		control, found, err := service.GetControl(program, controlID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get control: %v", err)), nil
		}
		if !found {
			return mcp.NewToolResultError(fmt.Sprintf("Control %s not found in %s", controlID, program)), nil
		}

		// Format the result as JSON
		controlJSON, err := json.MarshalIndent(control, "", "  ")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal control to JSON: %v", err)), nil
		}

		return mcp.NewToolResultText(string(controlJSON)), nil
	})

	// Tool: get_control_family
	getControlFamilyTool := mcp.NewTool("get_control_family",
		mcp.WithDescription("Get all controls in a family (e.g., AC for Access Control)"),
		mcp.WithString("program",
			mcp.Required(),
			mcp.Description("The FedRAMP program (High or Moderate)"),
			mcp.Enum("FedRAMP High", "FedRAMP Moderate"),
		),
		mcp.WithString("family",
			mcp.Required(),
			mcp.Description("The control family ID (e.g., AC, IA)"),
		),
	)
	s.AddTool(getControlFamilyTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		program := request.Params.Arguments["program"].(string)
		familyID := request.Params.Arguments["family"].(string)

		family, found, err := service.GetControlFamily(program, familyID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get control family: %v", err)), nil
		}
		if !found {
			return mcp.NewToolResultError(fmt.Sprintf("Control family %s not found in %s", familyID, program)), nil
		}

		// Format the result as JSON
		familyJSON, err := json.MarshalIndent(family, "", "  ")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal control family to JSON: %v", err)), nil
		}

		return mcp.NewToolResultText(string(familyJSON)), nil
	})

	// Tool: list_control_families
	listControlFamiliesTool := mcp.NewTool("list_control_families",
		mcp.WithDescription("List all control families in a program"),
		mcp.WithString("program",
			mcp.Required(),
			mcp.Description("The FedRAMP program (High or Moderate)"),
			mcp.Enum("FedRAMP High", "FedRAMP Moderate"),
		),
	)
	s.AddTool(listControlFamiliesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		program := request.Params.Arguments["program"].(string)

		families, err := service.ListControlFamilies(program)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to list control families: %v", err)), nil
		}

		// Create a simplified response with just family ID, title, and control count
		type SimplifiedFamily struct {
			ID           string `json:"id"`
			Title        string `json:"title"`
			ControlCount int    `json:"controlCount"`
		}

		simplifiedFamilies := make([]SimplifiedFamily, 0, len(families))
		for _, family := range families {
			simplifiedFamilies = append(simplifiedFamilies, SimplifiedFamily{
				ID:           family.ID,
				Title:        family.Title,
				ControlCount: len(family.Controls),
			})
		}

		// Format the result as JSON
		familiesJSON, err := json.MarshalIndent(simplifiedFamilies, "", "  ")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal families to JSON: %v", err)), nil
		}

		return mcp.NewToolResultText(string(familiesJSON)), nil
	})

	// Tool: search_controls
	searchControlsTool := mcp.NewTool("search_controls",
		mcp.WithDescription("Search for controls by keyword"),
		mcp.WithString("program",
			mcp.Required(),
			mcp.Description("The FedRAMP program (High or Moderate)"),
			mcp.Enum("FedRAMP High", "FedRAMP Moderate"),
		),
		mcp.WithString("query",
			mcp.Required(),
			mcp.Description("The search query"),
		),
	)
	s.AddTool(searchControlsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		program := request.Params.Arguments["program"].(string)
		query := request.Params.Arguments["query"].(string)

		controls, err := service.SearchControls(program, query)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to search controls: %v", err)), nil
		}

		// Create a simplified response with just control ID, title
		type SimplifiedControl struct {
			ID    string `json:"id"`
			Title string `json:"title"`
		}

		simplifiedControls := make([]SimplifiedControl, 0, len(controls))
		for _, control := range controls {
			simplifiedControls = append(simplifiedControls, SimplifiedControl{
				ID:    control.ID,
				Title: control.Title,
			})
		}

		// Format the result as JSON
		controlsJSON, err := json.MarshalIndent(simplifiedControls, "", "  ")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal controls to JSON: %v", err)), nil
		}

		return mcp.NewToolResultText(string(controlsJSON)), nil
	})

	// Tool: get_control_evidence_guidance
	getControlEvidenceGuidanceTool := mcp.NewTool("get_control_evidence_guidance",
		mcp.WithDescription("Get detailed guidance for evidence about a specific control"),
		mcp.WithString("program",
			mcp.Required(),
			mcp.Description("The FedRAMP program (High or Moderate)"),
			mcp.Enum("FedRAMP High", "FedRAMP Moderate"),
		),
		mcp.WithString("controlId",
			mcp.Required(),
			mcp.Description("The ID of the control (e.g., AC-1, IA-2)"),
		),
	)
	s.AddTool(getControlEvidenceGuidanceTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		program := request.Params.Arguments["program"].(string)
		controlID := request.Params.Arguments["controlId"].(string)

		guidance, found, err := service.GetControlEvidenceGuidance(program, controlID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get control evidence guidance: %v", err)), nil
		}
		if !found {
			return mcp.NewToolResultError(fmt.Sprintf("Control %s not found in %s", controlID, program)), nil
		}

		// Create a response structure
		response := struct {
			ControlID string `json:"controlId"`
			Program   string `json:"program"`
			Guidance  string `json:"guidance"`
		}{
			ControlID: controlID,
			Program:   program,
			Guidance:  guidance,
		}

		// Format the result as JSON
		responseJSON, err := json.MarshalIndent(response, "", "  ")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal response to JSON: %v", err)), nil
		}

		return mcp.NewToolResultText(string(responseJSON)), nil
	})
}
