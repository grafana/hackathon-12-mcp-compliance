# FedRAMP Compliance MCP Server

An MCP (Model Context Protocol) server that provides information about FedRAMP security controls and guidance on collecting evidence for those controls.

## The Compliance Journey

The FedRAMP Compliance MCP Server is designed to support users throughout their compliance journey, which consists of three main phases:

1. **Understanding** - Learning about security controls, their requirements, and how they apply to your system
2. **Implementing** - Designing and implementing controls in your system to meet compliance requirements
3. **Evidencing** - Collecting and documenting evidence to demonstrate compliance with controls

![Compliance Journey](docs/images/compliance-journey.png)

## Features

- Provides up-to-date information about FedRAMP security controls
- Offers guidance on how to implement controls in your system
- Provides both general and company-specific evidence collection guidance
- Supports searching for controls and evidence guidance
- Implements the Model Context Protocol for seamless integration with LLM applications

## Installation

### Prerequisites

- Go 1.24.1 or newer

### Building from Source

1. Clone the repository:
   ```
   git clone https://github.com/grafana/hackathon-12-mcp-compliance.git
   cd hackathon-12-mcp-compliance
   ```

2. Build the server:
   ```
   make build
   ```

3. Run the server:
   ```
   make run
   ```

### Local Deployment (macOS)

To deploy the server locally on macOS:

```
make deploy-local
```

This will build the macOS version of the server and copy it to `~/.mcp-compliance/bin/`.

## Usage

The server implements the Model Context Protocol and can be used with any MCP client. It exposes the following resources and tools:

### Resources

- `compliance://programs` - List of all available compliance programs
- `compliance://{program}/families` - List of all control families for a compliance program
- `compliance://{program}/controls/{id}` - Information about a specific security control
- `compliance://{program}/search/{query}` - Search for security controls by keyword

### Tools

#### Understanding Controls

- `get_control_info` - Get detailed information about a specific security control
- `get_control_family` - Get information about a control family/group

#### Evidencing Controls

- `get_evidence_guidance` - Get general guidance on collecting evidence for a control
- `search_evidence_guidance` - Search for evidence guidance by keyword
- `get_company_evidence_practice` - Get company-specific evidence collection practices for a control
- `search_company_evidence_practices` - Search for company-specific evidence collection practices by keyword

## Customizing Company-Specific Evidence Practices

The server supports company-specific evidence collection practices through a YAML configuration file. This file can be placed in one of the following locations:

- `./evidence_practices.yaml`
- `./data/evidence_practices.yaml`
- `~/.mcp-compliance/evidence_practices.yaml`

The file should follow this format:

```yaml
AC-1:
  practice: "Description of how your company collects evidence for AC-1"
  responsible_team: "Team responsible for this control"
  artifacts:
    - "List of artifacts to collect"
    - "Another artifact"
  review_frequency: "How often to review"
  notes: "Additional notes"

AC-2:
  practice: "Description of how your company collects evidence for AC-2"
  # ...
```

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
	toolRequest := mcp.CallToolRequest{}
	toolRequest.Params.Name = "get_control_info"
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