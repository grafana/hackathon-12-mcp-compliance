package models

// MCPResponse represents a JSON-RPC response
type MCPResponse struct {
	ID      interface{} `json:"id"`
	JSONRPC string      `json:"jsonrpc"`
	Result  interface{} `json:"result,omitempty"`
	Error   *MCPError   `json:"error,omitempty"`
}

// MCPError represents an error in JSON-RPC format
type MCPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
