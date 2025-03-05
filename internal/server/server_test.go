package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rbrady/example-mcp/internal/models"
)

// Helper function to make test requests
func makeTestRequest(t *testing.T, payload map[string]interface{}) *httptest.ResponseRecorder {
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("Failed to marshal request: %v", err)
	}

	req, err := http.NewRequest("POST", "/mcp", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	HandleMCPRequest(rr, req)
	return rr
}

// Test Initialize Request
func TestHandleInitialize(t *testing.T) {
	request := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "initialize",
		"params":  map[string]interface{}{},
	}

	rr := makeTestRequest(t, request)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}

	var response models.MCPResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if response.Result == nil {
		t.Errorf("Expected a result, got nil")
	}
}

// Test AI Query Handling
func TestHandleCreateMessage(t *testing.T) {
	request := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      2,
		"method":  "sampling/createMessage",
		"params": map[string]interface{}{
			"messages": []map[string]interface{}{
				{
					"role":    "user",
					"content": map[string]interface{}{"type": "text", "text": "Why is SR-5 important?"},
				},
			},
		},
	}

	rr := makeTestRequest(t, request)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}

	var response models.MCPResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if response.Result == nil {
		t.Errorf("Expected a result, got nil")
	}
}
