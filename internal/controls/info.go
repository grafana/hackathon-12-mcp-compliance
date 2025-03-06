package controls

import (
	"context"
	"fmt"

	"github.com/grafana/mcp-compliance/internal/programs"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// RegisterControlTools registers control information tools with the MCP server
func RegisterControlTools(s *server.MCPServer, registry *programs.Registry) {
	// Register tool for getting control information
	registerGetControlInfoTool(s, registry)

	// Register tool for getting control family information
	registerGetControlFamilyTool(s, registry)
}

// Register a tool for getting control information
func registerGetControlInfoTool(s *server.MCPServer, registry *programs.Registry) {
	tool := mcp.NewTool("get_control_info",
		mcp.WithDescription("Get detailed information about a specific security control"),
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

		// Get the control
		control, err := program.GetControl(controlID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get control: %v", err)), nil
		}

		// Format the control as text
		text := fmt.Sprintf("Control Information: %s - %s\n\n", control.ID, control.Title)
		text += fmt.Sprintf("Family: %s\n", control.Family)
		text += fmt.Sprintf("Impact: %s\n\n", control.Impact)

		text += "Description:\n"
		text += fmt.Sprintf("%s\n\n", control.Description)

		if control.Guidance != "" {
			text += "Implementation Guidance:\n"
			text += fmt.Sprintf("%s\n\n", control.Guidance)
		}

		if control.AssessmentInfo != nil && len(control.AssessmentInfo.Objectives) > 0 {
			text += "Assessment Objectives:\n"
			for _, objective := range control.AssessmentInfo.Objectives {
				text += fmt.Sprintf("- %s\n", objective)
			}
			text += "\n"
		}

		if control.AssessmentInfo != nil && len(control.AssessmentInfo.Methods) > 0 {
			text += "Assessment Methods:\n"
			for _, method := range control.AssessmentInfo.Methods {
				text += fmt.Sprintf("- %s\n", method)
			}
			text += "\n"
		}

		if len(control.Enhancements) > 0 {
			text += "Control Enhancements:\n"
			for _, enhancement := range control.Enhancements {
				text += fmt.Sprintf("- %s: %s\n", enhancement.ID, enhancement.Title)
			}
		}

		return mcp.NewToolResultText(text), nil
	})
}

// Register a tool for getting control family information
func registerGetControlFamilyTool(s *server.MCPServer, registry *programs.Registry) {
	tool := mcp.NewTool("get_control_family",
		mcp.WithDescription("Get information about a control family/group"),
		mcp.WithString("familyId",
			mcp.Required(),
			mcp.Description("The ID of the control family (e.g., AC, IA)"),
		),
		mcp.WithString("program",
			mcp.Description("The compliance program (e.g., FedRAMP High)"),
		),
	)

	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Get family ID from parameters
		familyID, ok := request.Params.Arguments["familyId"].(string)
		if !ok || familyID == "" {
			return mcp.NewToolResultError("Family ID is required"), nil
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

		// Get all control families
		families, err := program.GetControlFamilies()
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get control families: %v", err)), nil
		}

		// Find the requested family
		var targetFamily *programs.ControlFamily
		for _, family := range families {
			if family.ID == familyID {
				targetFamily = family
				break
			}
		}

		if targetFamily == nil {
			return mcp.NewToolResultError(fmt.Sprintf("Control family not found: %s", familyID)), nil
		}

		// Format the family as text
		text := fmt.Sprintf("Control Family: %s - %s\n\n", targetFamily.ID, targetFamily.Name)
		text += fmt.Sprintf("Description: %s\n\n", targetFamily.Description)

		text += fmt.Sprintf("Controls in this family (%d):\n", len(targetFamily.Controls))
		for _, control := range targetFamily.Controls {
			text += fmt.Sprintf("- %s: %s\n", control.ID, control.Title)
		}

		return mcp.NewToolResultText(text), nil
	})
}
