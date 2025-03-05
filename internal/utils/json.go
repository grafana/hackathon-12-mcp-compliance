package utils

import (
	"encoding/json"
	"net/http"

	"github.com/rbrady/example-mcp/internal/models"
)

// sendErrorResponse sends a JSON-RPC error response
func SendErrorResponse(w http.ResponseWriter, id interface{}, code int, message string) {
	response := models.MCPResponse{
		ID:      id,
		JSONRPC: "2.0",
		Error:   &models.MCPError{Code: code, Message: message},
	}
	SendJSONResponse(w, response)
}

// SendJSONResponse sends a JSON response
func SendJSONResponse(w http.ResponseWriter, response models.MCPResponse) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
