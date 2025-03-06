package controls

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/grafana/mcp-compliance/internal/programs"
	"github.com/grafana/mcp-compliance/internal/programs/fedramp/high"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// Global registry of compliance programs
var programRegistry *programs.Registry

// Initialize the registry and register programs
func init() {
	programRegistry = programs.NewRegistry()

	// Register FedRAMP High program
	fedRAMPHigh, err := high.NewFedRAMPHighProgram()
	if err != nil {
		// Log error but continue
		fmt.Printf("Failed to initialize FedRAMP High program: %v\n", err)
	} else {
		programRegistry.Register(fedRAMPHigh)
	}

	// Register other programs here
}

// RegisterResources registers all control-related resources with the MCP server
func RegisterResources(s *server.MCPServer) {
	// Register resource for getting all available compliance programs
	registerProgramsResource(s)

	// Register resource for getting all control families for a program
	registerControlFamiliesResource(s)

	// Register resource for getting a specific control by ID
	registerControlResource(s)

	// Register resource for searching controls
	registerControlSearchResource(s)
}

// Register a resource for getting all available compliance programs
func registerProgramsResource(s *server.MCPServer) {
	resource := mcp.NewResource(
		"compliance://programs",
		"Available Compliance Programs",
		mcp.WithResourceDescription("List of all available compliance programs"),
		mcp.WithMIMEType("application/json"),
	)

	s.AddResource(resource, func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		// Get all programs
		allPrograms := programRegistry.GetAllPrograms()

		// Convert to a simpler structure for JSON
		type ProgramInfo struct {
			Name        string `json:"name"`
			Version     string `json:"version"`
			Description string `json:"description"`
		}

		programInfos := make([]ProgramInfo, 0, len(allPrograms))
		for _, program := range allPrograms {
			programInfos = append(programInfos, ProgramInfo{
				Name:        program.GetName(),
				Version:     program.GetVersion(),
				Description: program.GetDescription(),
			})
		}

		// Convert to JSON
		data, err := json.MarshalIndent(programInfos, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("failed to marshal program info: %v", err)
		}

		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      "compliance://programs",
				MIMEType: "application/json",
				Text:     string(data),
			},
		}, nil
	})
}

// Register a resource for getting all control families for a program
func registerControlFamiliesResource(s *server.MCPServer) {
	template := mcp.NewResourceTemplate(
		"compliance://{program}/families",
		"Control Families",
		mcp.WithTemplateDescription("List of all control families for a compliance program"),
		mcp.WithTemplateMIMEType("application/json"),
	)

	s.AddResourceTemplate(template, func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		// Extract program name from URI
		uri := request.Params.URI
		parts := strings.Split(uri, "/")
		if len(parts) != 3 {
			return nil, fmt.Errorf("invalid URI format: %s", uri)
		}

		programName := parts[1]

		// Get the program
		program, exists := programRegistry.GetProgram(programName)
		if !exists {
			return nil, fmt.Errorf("compliance program not found: %s", programName)
		}

		// Get control families
		families, err := program.GetControlFamilies()
		if err != nil {
			return nil, fmt.Errorf("failed to get control families: %v", err)
		}

		// Convert to JSON
		data, err := json.MarshalIndent(families, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("failed to marshal control families: %v", err)
		}

		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      uri,
				MIMEType: "application/json",
				Text:     string(data),
			},
		}, nil
	})
}

// Register a resource for getting a specific control by ID
func registerControlResource(s *server.MCPServer) {
	template := mcp.NewResourceTemplate(
		"compliance://{program}/controls/{id}",
		"Security Control",
		mcp.WithTemplateDescription("Information about a specific security control"),
		mcp.WithTemplateMIMEType("application/json"),
	)

	s.AddResourceTemplate(template, func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		// Extract program name and control ID from URI
		uri := request.Params.URI
		parts := strings.Split(uri, "/")
		if len(parts) != 4 {
			return nil, fmt.Errorf("invalid URI format: %s", uri)
		}

		programName := parts[1]
		controlID := parts[3]

		// Get the program
		program, exists := programRegistry.GetProgram(programName)
		if !exists {
			return nil, fmt.Errorf("compliance program not found: %s", programName)
		}

		// Get the control
		control, err := program.GetControl(controlID)
		if err != nil {
			return nil, err
		}

		// Convert to JSON
		data, err := json.MarshalIndent(control, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("failed to marshal control: %v", err)
		}

		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      uri,
				MIMEType: "application/json",
				Text:     string(data),
			},
		}, nil
	})
}

// Register a resource for searching controls
func registerControlSearchResource(s *server.MCPServer) {
	template := mcp.NewResourceTemplate(
		"compliance://{program}/search/{query}",
		"Search Controls",
		mcp.WithTemplateDescription("Search for controls by keyword"),
		mcp.WithTemplateMIMEType("application/json"),
	)

	s.AddResourceTemplate(template, func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		// Extract program name and search query from URI
		uri := request.Params.URI
		parts := strings.Split(uri, "/")
		if len(parts) != 4 {
			return nil, fmt.Errorf("invalid URI format: %s", uri)
		}

		programName := parts[1]
		query := parts[3]

		// Get the program
		program, exists := programRegistry.GetProgram(programName)
		if !exists {
			return nil, fmt.Errorf("compliance program not found: %s", programName)
		}

		// Search for controls
		controls, err := program.SearchControls(query)
		if err != nil {
			return nil, err
		}

		// Convert to JSON
		data, err := json.MarshalIndent(controls, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("failed to marshal search results: %v", err)
		}

		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      uri,
				MIMEType: "application/json",
				Text:     string(data),
			},
		}, nil
	})
}
