# Concept of Operations: FedRAMP Compliance Server

This document explains the end-to-end flow of how the FedRAMP Compliance Server works, from when a user asks a question to when they receive an answer. It covers the architecture, communication protocols, and data flow through each layer of the system.

## System Architecture Overview

The FedRAMP Compliance Server follows a hexagonal architecture with domain-driven design principles:

```
┌───────────────────────────────────────────────────┐
│                                                   │
│  ┌─────────┐    ┌───────────┐    ┌────────────┐   │
│  │         │    │           │    │            │   │
│  │  Agent  │◄──►│ MCP Server│◄──►│  Services  │   │
│  │         │    │           │    │            │   │
│  └─────────┘    └───────────┘    └────────────┘   │
│                                                   │
└───────────────────────────────────────────────────┘
```

## End-to-End Flow

### 1. User Query in Agent (Cursor/Claude Desktop)

When a user types a question about FedRAMP controls (e.g., "What is AC-1 in FedRAMP High?"):

1. The agent (Claude in Cursor or Claude Desktop) recognizes that the query is related to FedRAMP
2. The agent determines it needs to use an external tool to answer the question
3. The agent formulates a tool call to the MCP server

### 2. MCP Protocol Communication

The agent communicates with the MCP server using the Model Context Protocol:

1. **Initialization**: The agent establishes a connection with the MCP server
   ```json
   {
     "method": "initialize",
     "params": {
       "protocolVersion": "0.1.0",
       "clientInfo": {
         "name": "Claude",
         "version": "3.0"
       }
     }
   }
   ```

2. **Tool Discovery**: The agent queries available tools
   ```json
   {
     "method": "tools/list"
   }
   ```

3. **Tool Call**: The agent calls the appropriate tool with parameters
   ```json
   {
     "method": "tools/call",
     "params": {
       "name": "get_control",
       "arguments": {
         "program": "FedRAMP High",
         "controlId": "AC-1"
       }
     }
   }
   ```

### 3. MCP Server Processing

The MCP server (in `cmd/mcp-compliance/main.go`) handles the incoming request:

1. The server receives the tool call via stdio
2. It parses the JSON request and validates the parameters
3. It routes the request to the appropriate tool handler in `internal/server/server.go`

```go
// Simplified example from server.go
s.AddTool(mcp.NewTool("get_control", ...), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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
```

### 4. Service Layer Processing

The service layer (in `internal/services/`) orchestrates the operation:

1. The service receives the request from the MCP server
2. It validates the input parameters
3. It creates a command object with the validated parameters
4. It passes the command to the appropriate handler

### 5. Domain Layer Processing

The domain layer (in `internal/domain/fedramp/`) contains the core business logic:

1. The handler processes the command
2. It uses ports (interfaces) to interact with external systems
3. It applies business rules and transformations
4. It returns the result to the service layer

### 6. Adapter Layer Data Access

The adapter layer (in `internal/adapters/`) implements the ports defined in the domain layer:

1. The `EmbeddedComplianceRepository` adapter is used to access FedRAMP data
2. It looks up the requested program (FedRAMP High) in its registry
3. It retrieves the embedded data file for that program from `internal/resources/data/`
4. It unmarshals the JSON data into domain objects
5. It returns the requested control (AC-1) to the domain layer

```go
// Simplified example from embedded_compliance_repository.go
func (r *EmbeddedComplianceRepository) LoadProgram(programName string) (fedramp.Program, error) {
    // Find the file path for the program
    filePath, ok := r.programFiles[programName]
    if !ok {
        // Try case-insensitive match...
        return fedramp.Program{}, fmt.Errorf("program not found: %s", programName)
    }

    // Read the embedded file
    data, err := resources.Data.ReadFile(filePath)
    if err != nil {
        return fedramp.Program{}, fmt.Errorf("failed to read program data: %v", err)
    }

    // Unmarshal the JSON data
    var program fedramp.Program
    if err := json.Unmarshal(data, &program); err != nil {
        return fedramp.Program{}, fmt.Errorf("failed to parse program data: %v", err)
    }

    return program, nil
}
```

### 7. Response Flow Back to User

The response flows back through the layers in reverse:

1. The adapter returns the control data to the domain layer
2. The domain layer processes and enriches the data if needed
3. The service layer formats the response
4. The MCP server sends the response back to the agent
   ```json
   {
     "result": {
       "content": [
         {
           "type": "text",
           "text": "{\"id\":\"AC-1\",\"title\":\"Policy and Procedures\",\"description\":\"...\"}"
         }
       ]
     }
   }
   ```
5. The agent (Claude) interprets the JSON response
6. The agent formulates a natural language answer for the user
7. The user sees the response about AC-1 in FedRAMP High

## Detailed Component Descriptions

### MCP Server (cmd/mcp-compliance)

The MCP server is the entry point for all requests. It:
- Implements the Model Context Protocol
- Registers available tools
- Routes requests to the appropriate service
- Formats responses according to the MCP specification

### Registry (internal/registry)

The registry manages the available compliance programs:
- Maintains a map of program names to program implementations
- Provides methods to register and retrieve programs
- Ensures program names are case-insensitive for better user experience

### FedRAMP Program (internal/domain/fedramp)

The FedRAMP program implementation:
- Defines the structure of FedRAMP controls and families
- Provides methods to search and retrieve controls
- Implements business logic for control relationships

### Embedded Compliance Repository (internal/adapters)

The embedded compliance repository:
- Uses Go's embed package to include data files in the binary
- Maps program names to embedded file paths
- Loads and parses JSON data into domain objects
- Provides a clean interface for accessing compliance data

### Resources (internal/resources)

The resources package:
- Embeds the processed FedRAMP data files
- Makes them available to the rest of the application
- Ensures the binary is self-contained

## Data Flow Sequence Diagram

```
┌─────┐          ┌─────────┐          ┌──────────┐          ┌────────┐
│User │          │  Agent  │          │MCP Server│          │Services│
└──┬──┘          └────┬────┘          └────┬─────┘          └───┬────┘
   │                  │                    │                    │
   │ Ask about AC-1   │                    │                    │
   │─────────────────>│                    │                    │
   │                  │                    │                    │
   │                  │ Initialize         │                    │
   │                  │───────────────────>│                    │
   │                  │                    │                    │
   │                  │ List tools         │                    │
   │                  │───────────────────>│                    │
   │                  │                    │                    │
   │                  │ Call get_control   │                    │
   │                  │───────────────────>│                    │
   │                  │                    │                    │
   │                  │                    │ Get program        │
   │                  │                    │───────────────────>│
   │                  │                    │                    │
   │                  │                    │<───────────────────│
   │                  │                    │ Return program     │
   │                  │                    │                    │
   │                  │                    │ Get control        │
   │                  │                    │───────────────────>│
   │                  │                    │                    │
   │                  │                    │<───────────────────│
   │                  │                    │ Return control     │
   │                  │                    │                    │
   │                  │<───────────────────│                    │
   │                  │ Return control     │                    │
   │                  │                    │                    │
   │<─────────────────│                    │                    │
   │ Display answer   │                    │                    │
   │                  │                    │                    │
```

## Error Handling

The system includes comprehensive error handling at each layer:

1. **Adapter Layer**: Handles file not found, parsing errors, and data validation
2. **Domain Layer**: Handles business rule violations and data integrity issues
3. **Service Layer**: Handles invalid input parameters and command validation
4. **MCP Server**: Formats errors according to the MCP specification
5. **Agent**: Presents errors to the user in a friendly, actionable way

## Performance Considerations

The system is designed for optimal performance:

1. **Embedded Data**: All data is embedded in the binary, eliminating disk I/O
2. **In-Memory Processing**: Once loaded, all data remains in memory for fast access
3. **Case-Insensitive Lookups**: Flexible matching of control IDs and program names
4. **Efficient JSON Parsing**: Minimal unmarshaling for better performance

## Security Considerations

The system includes several security features:

1. **Read-Only Data**: All embedded data is read-only, preventing modification
2. **Input Validation**: All user inputs are validated before processing
3. **Error Sanitization**: Error messages don't expose sensitive information
4. **No External Dependencies**: The binary is self-contained, reducing attack surface

## Conclusion

This concept of operations document provides a comprehensive overview of how the FedRAMP Compliance Server works, from user query to response. By following the hexagonal architecture and domain-driven design principles, the system achieves a clean separation of concerns, making it maintainable, extensible, and robust. 