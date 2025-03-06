# FedRAMP Compliance MCP Server

An MCP (Model Context Protocol) server that provides information about FedRAMP High security controls and guidance on collecting evidence for those controls.

## Features

- Provides up-to-date information about FedRAMP High security controls
- Offers guidance on how to collect evidence for security controls
- Supports searching for controls and evidence guidance
- Implements the Model Context Protocol for seamless integration with LLM applications

## Installation

### Prerequisites

- Go 1.24.1 or newer

### Building from Source

1. Clone the repository:
   ```
   git clone https://github.com/grafana/mcp-compliance.git
   cd mcp-compliance
   ```

2. Build the server:
   ```
   make build
   ```

3. Run the server:
   ```
   make run
   ```

## Usage

The server implements the Model Context Protocol and can be used with any MCP client. It exposes the following resources and tools:

### Resources

- `fedramp://controls/families` - List of all FedRAMP security control families
- `fedramp://controls/{id}` - Information about a specific FedRAMP security control
- `fedramp://controls/search/{query}` - Search for FedRAMP security controls by keyword

### Tools

- `get_evidence_guidance` - Get guidance on collecting evidence for a specific FedRAMP security control
- `search_evidence_guidance` - Search for evidence guidance by keyword

## Example Client Usage

Here's an example of how to use the server with an MCP client:

```go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
)

func main() {
	// Create a new MCP client connected to our server
	c, err := client.NewStdioMCPClient(
		"./bin/mcp-compliance", // Path to the server binary
		[]string{},             // Empty ENV
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

	// Get information about a specific control
	fmt.Println("Getting information about AC-1...")
	readRequest := mcp.ReadResourceRequest{}
	readRequest.Params.URI = "fedramp://controls/AC-1"

	resource, err := c.ReadResource(ctx, readRequest)
	if err != nil {
		log.Fatalf("Failed to read resource: %v", err)
	}

	for _, content := range resource.Contents {
		if textContent, ok := content.(mcp.TextContent); ok {
			fmt.Println(textContent.Text)
		}
	}
	fmt.Println()

	// Get evidence guidance for a control
	fmt.Println("Getting evidence guidance for AC-1...")
	toolRequest := mcp.CallToolRequest{}
	toolRequest.Params.Name = "get_evidence_guidance"
	toolRequest.Params.Arguments = map[string]interface{}{
		"controlId": "AC-1",
	}

	result, err := c.CallTool(ctx, toolRequest)
	if err != nil {
		log.Fatalf("Failed to call tool: %v", err)
	}

	for _, content := range result.Content {
		if textContent, ok := content.(mcp.TextContent); ok {
			fmt.Println(textContent.Text)
		}
	}
}
```

## License

This project is licensed under the MIT License - see the LICENSE file for details. 