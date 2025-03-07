package server

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/grafana/hackathon-12-mcp-compliance/internal/registry"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// NewServer creates a new MCP server with FedRAMP tools
func NewServer(r *registry.Registry) *server.MCPServer {
	s := server.NewMCPServer(
		"FedRAMP Compliance Server",
		"1.0.0",
	)

	registerTools(s, r)
	return s
}

// registerTools adds all FedRAMP control management tools to the server
func registerTools(s *server.MCPServer, r *registry.Registry) {
	// Get control details
	s.AddTool(mcp.NewTool("get_control",
		mcp.WithDescription("Get detailed information about a specific FedRAMP control"),
		mcp.WithString("program",
			mcp.Required(),
			mcp.Description("The FedRAMP program (High or Moderate)"),
			mcp.Enum("FedRAMP High", "FedRAMP Moderate"),
		),
		mcp.WithString("controlId",
			mcp.Required(),
			mcp.Description("The ID of the control (e.g., AC-1, IA-2)"),
		),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		program, err := r.GetProgram(request.Params.Arguments["program"].(string))
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		control, err := program.GetControl(request.Params.Arguments["controlId"].(string))
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		jsonBytes, err := json.Marshal(control)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to marshal control: %v", err)), nil
		}

		return mcp.NewToolResultText(string(jsonBytes)), nil
	})

	// Get control family
	s.AddTool(mcp.NewTool("get_control_family",
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
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		program, err := r.GetProgram(request.Params.Arguments["program"].(string))
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		controls, err := program.GetControlFamily(request.Params.Arguments["family"].(string))
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		jsonBytes, err := json.Marshal(controls)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to marshal controls: %v", err)), nil
		}

		return mcp.NewToolResultText(string(jsonBytes)), nil
	})

	// List control families
	s.AddTool(mcp.NewTool("list_control_families",
		mcp.WithDescription("List all control families in a program"),
		mcp.WithString("program",
			mcp.Required(),
			mcp.Description("The FedRAMP program (High or Moderate)"),
			mcp.Enum("FedRAMP High", "FedRAMP Moderate"),
		),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		program, err := r.GetProgram(request.Params.Arguments["program"].(string))
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		families, err := program.ListControlFamilies()
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		jsonBytes, err := json.Marshal(families)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to marshal families: %v", err)), nil
		}

		return mcp.NewToolResultText(string(jsonBytes)), nil
	})

	// Search controls
	s.AddTool(mcp.NewTool("search_controls",
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
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		program, err := r.GetProgram(request.Params.Arguments["program"].(string))
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		controls, err := program.SearchControls(request.Params.Arguments["query"].(string))
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		jsonBytes, err := json.Marshal(controls)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to marshal controls: %v", err)), nil
		}

		return mcp.NewToolResultText(string(jsonBytes)), nil
	})

	// Get controls by status
	s.AddTool(mcp.NewTool("get_controls_by_status",
		mcp.WithDescription("Get all controls with a specific implementation status"),
		mcp.WithString("program",
			mcp.Required(),
			mcp.Description("The FedRAMP program (High or Moderate)"),
			mcp.Enum("FedRAMP High", "FedRAMP Moderate"),
		),
		mcp.WithString("status",
			mcp.Required(),
			mcp.Description("The implementation status to filter by"),
		),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		program, err := r.GetProgram(request.Params.Arguments["program"].(string))
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		controls, err := program.GetControlsByStatus(request.Params.Arguments["status"].(string))
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		jsonBytes, err := json.Marshal(controls)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to marshal controls: %v", err)), nil
		}

		return mcp.NewToolResultText(string(jsonBytes)), nil
	})

	// Get control parameters
	s.AddTool(mcp.NewTool("get_control_parameters",
		mcp.WithDescription("Get parameters for a specific control"),
		mcp.WithString("program",
			mcp.Required(),
			mcp.Description("The FedRAMP program (High or Moderate)"),
			mcp.Enum("FedRAMP High", "FedRAMP Moderate"),
		),
		mcp.WithString("controlId",
			mcp.Required(),
			mcp.Description("The ID of the control (e.g., AC-1, IA-2)"),
		),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		program, err := r.GetProgram(request.Params.Arguments["program"].(string))
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		params, err := program.GetControlParameters(request.Params.Arguments["controlId"].(string))
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		jsonBytes, err := json.Marshal(params)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to marshal parameters: %v", err)), nil
		}

		return mcp.NewToolResultText(string(jsonBytes)), nil
	})
}
