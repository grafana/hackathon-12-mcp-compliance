# FedRAMP Compliance Assistant

An MCP (Machine Chat Protocol) server that helps with compliance programs, starting with FedRAMP High compliance. This tool assists users in understanding, implementing, and collecting evidence for security controls.

## Features

- Get information about specific security controls
- Get guidance on how to implement security controls
- Understand why specific security controls are important
- Get guidance on finding and collecting evidence for security controls
- (Coming soon) Evidence collection and management
- (Coming soon) System Security Plan (SSP) generation

## Project Structure

```
hackathon-12-mcp-compliance/
│── go.mod                    # Go module file
│── go.sum                    # Go dependencies (auto-generated)
│── cmd/
│   └── server/
│       └── main.go           # Entry point for the MCP server
│── internal/
│   └── compliance/           # Compliance-related code
│       ├── handlers.go       # MCP tool handlers
│       ├── models.go         # Data models for FedRAMP controls
│       └── service.go        # Services for working with controls
│── data/
│   └── FedRAMP_rev5_HIGH-baseline_profile.yaml  # FedRAMP profile data
```

## Getting Started

1. Clone the repository
2. Place the FedRAMP_rev5_HIGH-baseline_profile.yaml file in the data directory
3. Run the server:

```bash
go run cmd/server/main.go
```

## Usage

The server implements the Machine Chat Protocol (MCP) and can be used with any MCP-compatible client. It provides tools for:

- `getControlInfo`: Get information about a specific security control
- `getImplementationGuidance`: Get guidance on how to implement a security control
- `explainControlImportance`: Explain why a specific security control is important
- `getEvidenceGuidance`: Get guidance on how to find or collect evidence for a security control

## License

MIT