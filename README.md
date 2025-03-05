
```bash
example-mcp/
│── go.mod                    # Go module file
│── go.sum                    # Go dependencies (auto-generated)
│── cmd/
│   └── server/
│       └── main.go            # Entry point for the MCP server
│── internal/
│   ├── server/
│   │   └── server.go          # Core MCP server logic
│   ├── handlers/              # JSON-RPC 2.0 handlers
│   ├── models/                # Data models for requests/responses
│   ├── utils/                 # Helper functions

```