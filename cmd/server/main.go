package main

import (
	"context"
	"fmt"
	"log"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rbrady/example-mcp/internal/compliance"
)

func main() {
	// Create a new MCP server
	s := server.NewMCPServer(
		"FedRAMP Compliance Assistant",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
	)

	// Add compliance tools
	registerComplianceTools(s)

	// Start the server
	log.Println("Starting FedRAMP Compliance Assistant MCP Server...")
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}

func registerComplianceTools(s *server.MCPServer) {
	// Add tool for getting information about security controls
	controlInfoTool := mcp.NewTool("getControlInfo",
		mcp.WithDescription("Get information about a specific security control"),
		mcp.WithString("controlId",
			mcp.Required(),
			mcp.Description("The ID of the security control (e.g., AC-1, IA-2)"),
		),
	)
	s.AddTool(controlInfoTool, compliance.HandleGetControlInfo)

	// Add tool for getting implementation guidance
	implementationTool := mcp.NewTool("getImplementationGuidance",
		mcp.WithDescription("Get guidance on how to implement a specific security control"),
		mcp.WithString("controlId",
			mcp.Required(),
			mcp.Description("The ID of the security control (e.g., AC-1, IA-2)"),
		),
		mcp.WithString("context",
			mcp.Description("Additional context about the project or environment (optional)"),
		),
	)
	s.AddTool(implementationTool, compliance.HandleGetImplementationGuidance)

	// Add tool for explaining the importance of a control
	importanceTool := mcp.NewTool("explainControlImportance",
		mcp.WithDescription("Explain why a specific security control is important"),
		mcp.WithString("controlId",
			mcp.Required(),
			mcp.Description("The ID of the security control (e.g., AC-1, IA-2)"),
		),
	)
	s.AddTool(importanceTool, compliance.HandleExplainControlImportance)

	// Add tool for finding evidence for a control
	evidenceGuidanceTool := mcp.NewTool("getEvidenceGuidance",
		mcp.WithDescription("Get guidance on how to find or collect evidence for a security control"),
		mcp.WithString("controlId",
			mcp.Required(),
			mcp.Description("The ID of the security control (e.g., AC-1, IA-2)"),
		),
	)
	s.AddTool(evidenceGuidanceTool, compliance.HandleGetEvidenceGuidance)

	// Future tools (placeholders for now)

	// Add tool for getting current evidence for a control
	getCurrentEvidenceTool := mcp.NewTool("getCurrentEvidence",
		mcp.WithDescription("Get the current evidence for a specific security control"),
		mcp.WithString("controlId",
			mcp.Required(),
			mcp.Description("The ID of the security control (e.g., AC-1, IA-2)"),
		),
	)
	s.AddTool(getCurrentEvidenceTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return mcp.NewToolResultText("This functionality is coming soon."), nil
	})

	// Add tool for collecting evidence for a control
	collectEvidenceTool := mcp.NewTool("collectEvidence",
		mcp.WithDescription("Collect evidence for a specific security control"),
		mcp.WithString("controlId",
			mcp.Required(),
			mcp.Description("The ID of the security control (e.g., AC-1, IA-2)"),
		),
		mcp.WithString("evidenceType",
			mcp.Required(),
			mcp.Description("The type of evidence to collect (e.g., document, screenshot, log)"),
		),
		mcp.WithString("description",
			mcp.Required(),
			mcp.Description("Description of the evidence"),
		),
		mcp.WithString("location",
			mcp.Required(),
			mcp.Description("Location of the evidence (e.g., file path, URL)"),
		),
	)
	s.AddTool(collectEvidenceTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return mcp.NewToolResultText("This functionality is coming soon."), nil
	})

	// Add tool for generating an SSP
	generateSSPTool := mcp.NewTool("generateSSP",
		mcp.WithDescription("Generate a System Security Plan (SSP) based on collected evidence"),
		mcp.WithString("outputLocation",
			mcp.Description("Location to store the generated SSP (optional)"),
		),
	)
	s.AddTool(generateSSPTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return mcp.NewToolResultText("This functionality is coming soon."), nil
	})
}
