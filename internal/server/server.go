package server

import (
	"encoding/json"
	"net/http"

	"github.com/rbrady/example-mcp/internal/handlers"
	"github.com/rbrady/example-mcp/internal/models"
	"github.com/rbrady/example-mcp/internal/utils"
)

// HandleMCPRequest processes incoming JSON-RPC 2.0 requests
func HandleMCPRequest(w http.ResponseWriter, r *http.Request) {
	var request models.MCPRequest

	// Decode JSON request
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		utils.SendErrorResponse(w, request.ID, -32700, "Invalid JSON format")
		return
	}

	// Validate JSON-RPC version
	if request.JSONRPC != "2.0" {
		utils.SendErrorResponse(w, request.ID, -32600, "Invalid JSON-RPC version")
		return
	}

	// Route request to appropriate handler
	switch request.Method {
	case "initialize":
		handlers.HandleInitialize(w, request)
	case "sampling/createMessage":
		handlers.HandleCreateMessage(w, request)
	default:
		utils.SendErrorResponse(w, request.ID, -32601, "Method not found")
	}
}
