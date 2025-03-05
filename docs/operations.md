# FedRAMP Compliance Assistant: How It Works

This document explains how the FedRAMP Compliance Assistant MCP server works, including the flow of communication between AI agents and the server.

## Overview of the Machine Chat Protocol (MCP)

The Machine Chat Protocol (MCP) is designed to enable AI agents to interact with external tools. It provides a standardized way for agents to discover and call tools, and for tools to return results to agents.

## Server Architecture

The FedRAMP Compliance Assistant is built using the `mcp-go` library, which provides a Go implementation of the Machine Chat Protocol. The server exposes tools that help with understanding, implementing, and collecting evidence for FedRAMP security controls.

## Communication Flow

### 1. Server Initialization

When you start the server with `./fedramp-compliance-assistant`, the following happens:

1. The `main()` function in `cmd/server/main.go` creates a new MCP server using `server.NewMCPServer()` with the name "FedRAMP Compliance Assistant" and version "1.0.0".
2. It configures the server with resource capabilities and logging.
3. It calls `registerComplianceTools()` to register all the tools that the server will provide.
4. Finally, it starts the server with `server.ServeStdio(s)`, which listens for MCP requests on standard input/output.

### 2. Tool Registration

During initialization, the `registerComplianceTools()` function registers several tools:

1. `getControlInfo`: Gets information about a specific security control
2. `getImplementationGuidance`: Gets guidance on implementing a control
3. `explainControlImportance`: Explains why a control is important
4. `getEvidenceGuidance`: Provides guidance on collecting evidence
5. Plus placeholder tools for future functionality (evidence collection, SSP generation)

Each tool is defined with:
- A name
- A description
- Parameters (with types, descriptions, and validation rules)
- A handler function that processes requests for that tool

### 3. MCP Protocol Flow

When an agent wants to interact with your server, the following happens:

1. **Initialization**: The agent sends an `initialize` request to the server.
   - The server responds with its capabilities and available tools.

2. **Tool Discovery**: The agent sends a `tools/list` request.
   - The server responds with a list of available tools and their schemas.

3. **Tool Invocation**: The agent sends a `tools/call` request with:
   - The name of the tool to call
   - Arguments for the tool parameters

4. **Tool Execution**: The server:
   - Validates the request
   - Routes it to the appropriate handler function
   - The handler processes the request and generates a response
   - The server sends the response back to the agent

### 4. Handler Execution

Let's trace through what happens when an agent calls the `getControlInfo` tool:

1. The agent sends a request like:
   ```json
   {
     "jsonrpc": "2.0",
     "id": "123",
     "method": "tools/call",
     "params": {
       "name": "getControlInfo",
       "arguments": {
         "controlId": "AC-1"
       }
     }
   }
   ```

2. The mcp-go library:
   - Parses the request
   - Validates the parameters against the tool schema
   - Calls the `HandleGetControlInfo` function with the request

3. The `HandleGetControlInfo` function in `internal/compliance/handlers.go`:
   - Extracts the `controlId` parameter
   - Calls `GetControlInfo(controlId)` from `service.go`
   - `GetControlInfo` loads the FedRAMP profile if not already loaded
   - It looks up the control in the profile
   - It returns the control data or an error

4. The handler formats the control information into a readable text format

5. The handler returns a `mcp.CallToolResult` with the formatted text

6. The mcp-go library converts this to a JSON-RPC response and sends it back to the agent

### 5. Data Flow

The data flow for control information is:

1. FedRAMP YAML file → Loaded by `LoadFedRAMPProfile()` → Stored in memory
2. Agent request → Parsed by mcp-go → Routed to handler
3. Handler → Calls service function → Gets control data
4. Handler → Formats data → Returns result
5. Result → Converted to JSON-RPC → Sent to agent

## Code Structure

The codebase is organized as follows:

- `cmd/server/main.go`: Entry point for the server, initializes the MCP server and registers tools
- `internal/compliance/`:
  - `handlers.go`: Contains handler functions for each tool
  - `models.go`: Defines data structures for FedRAMP controls and evidence
  - `service.go`: Provides services for working with controls and evidence
- `data/FedRAMP_rev5_HIGH-baseline_profile.yaml`: Contains the FedRAMP High baseline profile

## Handler Implementation

Each handler follows a similar pattern:

1. Extract parameters from the request
2. Call service functions to get or process data
3. Format the data into a readable response
4. Return the response as a `mcp.CallToolResult`

For example, the `HandleGetControlInfo` function:
- Extracts the `controlId` parameter
- Calls `GetControlInfo(controlId)` to get the control data
- Formats the control data into a readable text format with Markdown
- Returns the formatted text as a `mcp.CallToolResult`

## Service Implementation

The service layer provides functions for working with FedRAMP controls and evidence:

- `LoadFedRAMPProfile(filePath)`: Loads the FedRAMP profile from a YAML file
- `GetControlInfo(controlID)`: Gets information about a specific control
- `GetEvidenceForControl(controlID)`: Gets evidence for a specific control
- `SaveEvidence(evidence)`: Saves evidence for a control
- `GenerateSSP()`: Generates a System Security Plan

The service layer uses a singleton pattern to ensure the FedRAMP profile is loaded only once and shared across all requests.

## Testing the Server

To test the server, you can:

1. Start the server with `./fedramp-compliance-assistant`
2. Use a test client to send MCP requests to the server
3. Verify that the server responds correctly to each request

A test client can be implemented using the `mcp-go` library's client functionality, or by manually constructing JSON-RPC requests and sending them to the server's standard input.

## Extending the Server

To add new functionality to the server:

1. Define a new tool in `registerComplianceTools()` with a name, description, and parameters
2. Implement a handler function for the tool
3. Add any necessary service functions to support the handler
4. Register the tool with the server using `s.AddTool()`

For example, to add a tool for collecting evidence:

1. Define a `collectEvidence` tool with parameters for the control ID, evidence type, etc.
2. Implement a `HandleCollectEvidence` function that processes the request
3. Add a `SaveEvidence` function to the service layer
4. Register the tool with `s.AddTool(collectEvidenceTool, HandleCollectEvidence)`
