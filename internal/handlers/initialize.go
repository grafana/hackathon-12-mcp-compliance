package handlers

import (
	"net/http"

	"github.com/rbrady/example-mcp/internal/models"
	"github.com/rbrady/example-mcp/internal/utils"
)

// HandleInitialize processes the "initialize" request
func HandleInitialize(w http.ResponseWriter, request models.MCPRequest) {
	response := models.MCPResponse{
		ID:      request.ID,
		JSONRPC: "2.0",
		Result: map[string]interface{}{
			"protocolVersion": "1.0",
			"serverInfo": map[string]string{
				"name":    "MCP-Server",
				"version": "1.0",
			},
			"capabilities": map[string]interface{}{
				"sampling":  true,
				"tools":     true,
				"resources": true,
			},
		},
	}
	utils.SendJSONResponse(w, response)
}
