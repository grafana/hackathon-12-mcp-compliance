package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
)

func main() {
	// Create a new MCP client connected to our server
	c, err := client.NewStdioMCPClient(
		"../bin/mcp-compliance", // Path to the server binary
		[]string{},              // Empty ENV
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer c.Close()

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Initialize the client
	fmt.Println("Initializing client...")
	initRequest := mcp.InitializeRequest{}
	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = mcp.Implementation{
		Name:    "example-client",
		Version: "1.0.0",
	}

	initResult, err := c.Initialize(ctx, initRequest)
	if err != nil {
		log.Fatalf("Failed to initialize: %v", err)
	}
	fmt.Printf(
		"Initialized with server: %s %s\n\n",
		initResult.ServerInfo.Name,
		initResult.ServerInfo.Version,
	)

	// List available tools
	fmt.Println("Listing available tools...")
	toolsRequest := mcp.ListToolsRequest{}
	tools, err := c.ListTools(ctx, toolsRequest)
	if err != nil {
		log.Fatalf("Failed to list tools: %v", err)
	}
	for _, tool := range tools.Tools {
		fmt.Printf("- %s: %s\n", tool.Name, tool.Description)
	}
	fmt.Println()

	// Get information about control families
	fmt.Println("Getting information about control families...")
	readFamiliesRequest := mcp.ReadResourceRequest{}
	readFamiliesRequest.Params.URI = "fedramp://controls/families"

	familiesResource, err := c.ReadResource(ctx, readFamiliesRequest)
	if err != nil {
		log.Fatalf("Failed to read resource: %v", err)
	}

	for _, content := range familiesResource.Contents {
		if textContent, ok := content.(mcp.TextResourceContents); ok {
			// Parse the JSON to pretty print it
			var data interface{}
			if err := json.Unmarshal([]byte(textContent.Text), &data); err != nil {
				fmt.Println(textContent.Text)
			} else {
				prettyJSON, _ := json.MarshalIndent(data, "", "  ")
				fmt.Println(string(prettyJSON))
			}
		}
	}
	fmt.Println()

	// Get information about a specific control
	fmt.Println("Getting information about AC-1...")
	readControlRequest := mcp.ReadResourceRequest{}
	readControlRequest.Params.URI = "fedramp://controls/AC-1"

	controlResource, err := c.ReadResource(ctx, readControlRequest)
	if err != nil {
		log.Fatalf("Failed to read resource: %v", err)
	}

	for _, content := range controlResource.Contents {
		if textContent, ok := content.(mcp.TextResourceContents); ok {
			// Parse the JSON to pretty print it
			var data interface{}
			if err := json.Unmarshal([]byte(textContent.Text), &data); err != nil {
				fmt.Println(textContent.Text)
			} else {
				prettyJSON, _ := json.MarshalIndent(data, "", "  ")
				fmt.Println(string(prettyJSON))
			}
		}
	}
	fmt.Println()

	// Search for controls
	fmt.Println("Searching for controls with 'authentication'...")
	searchControlsRequest := mcp.ReadResourceRequest{}
	searchControlsRequest.Params.URI = "fedramp://controls/search/authentication"

	searchResults, err := c.ReadResource(ctx, searchControlsRequest)
	if err != nil {
		log.Fatalf("Failed to read resource: %v", err)
	}

	for _, content := range searchResults.Contents {
		if textContent, ok := content.(mcp.TextResourceContents); ok {
			// Parse the JSON to pretty print it
			var data interface{}
			if err := json.Unmarshal([]byte(textContent.Text), &data); err != nil {
				fmt.Println(textContent.Text)
			} else {
				prettyJSON, _ := json.MarshalIndent(data, "", "  ")
				fmt.Println(string(prettyJSON))
			}
		}
	}
	fmt.Println()

	// Get evidence guidance for a control
	fmt.Println("Getting evidence guidance for AC-1...")
	evidenceRequest := mcp.CallToolRequest{}
	evidenceRequest.Params.Name = "get_evidence_guidance"
	evidenceRequest.Params.Arguments = map[string]interface{}{
		"controlId": "AC-1",
	}

	evidenceResult, err := c.CallTool(ctx, evidenceRequest)
	if err != nil {
		log.Fatalf("Failed to call tool: %v", err)
	}

	for _, content := range evidenceResult.Content {
		if textContent, ok := content.(mcp.TextContent); ok {
			// Parse the JSON to pretty print it
			var data interface{}
			if err := json.Unmarshal([]byte(textContent.Text), &data); err != nil {
				fmt.Println(textContent.Text)
			} else {
				prettyJSON, _ := json.MarshalIndent(data, "", "  ")
				fmt.Println(string(prettyJSON))
			}
		}
	}
	fmt.Println()

	// Search for evidence guidance
	fmt.Println("Searching for evidence guidance with 'policy'...")
	searchEvidenceRequest := mcp.CallToolRequest{}
	searchEvidenceRequest.Params.Name = "search_evidence_guidance"
	searchEvidenceRequest.Params.Arguments = map[string]interface{}{
		"query": "policy",
	}

	searchEvidenceResult, err := c.CallTool(ctx, searchEvidenceRequest)
	if err != nil {
		log.Fatalf("Failed to call tool: %v", err)
	}

	for _, content := range searchEvidenceResult.Content {
		if textContent, ok := content.(mcp.TextContent); ok {
			// Parse the JSON to pretty print it
			var data interface{}
			if err := json.Unmarshal([]byte(textContent.Text), &data); err != nil {
				fmt.Println(textContent.Text)
			} else {
				prettyJSON, _ := json.MarshalIndent(data, "", "  ")
				fmt.Println(string(prettyJSON))
			}
		}
	}
}
