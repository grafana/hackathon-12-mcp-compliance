package models

// MCPRequest represents a JSON-RPC request
type MCPRequest struct {
	ID      interface{} `json:"id"`
	JSONRPC string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
}
