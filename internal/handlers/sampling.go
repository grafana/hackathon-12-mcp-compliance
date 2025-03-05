package handlers

import (
	"net/http"

	"github.com/rbrady/example-mcp/internal/models"
	"github.com/rbrady/example-mcp/internal/utils"
)

// HandleCreateMessage processes AI query requests
func HandleCreateMessage(w http.ResponseWriter, request models.MCPRequest) {
	params, ok := request.Params.(map[string]interface{})
	if !ok {
		utils.SendErrorResponse(w, request.ID, -32602, "Invalid parameters")
		return
	}

	messages, exists := params["messages"].([]interface{})
	if !exists || len(messages) == 0 {
		utils.SendErrorResponse(w, request.ID, -32602, "Missing messages")
		return
	}

	lastMessage := messages[len(messages)-1].(map[string]interface{})
	question := lastMessage["content"].(map[string]interface{})["text"].(string)

	// Placeholder response logic
	answer := "I don't have an answer for that yet."
	if question == "Why is SR-5 important?" {
		answer = "SR-5 ensures supply chain risk management, critical for FedRAMP compliance."
	}

	response := models.MCPResponse{
		ID:      request.ID,
		JSONRPC: "2.0",
		Result: map[string]interface{}{
			"content": []map[string]interface{}{
				{"type": "text", "text": answer},
			},
			"model": "mcp-ai-model",
			"role":  "assistant",
		},
	}
	utils.SendJSONResponse(w, response)
}
