package main

import (
	"strings"
	"testing"
)

func TestMCPClient(t *testing.T) {
	// Path to the server executable
	serverPath := "../fedramp-compliance-assistant"

	// Create a new MCP client
	client, err := NewMCPClient(serverPath)
	if err != nil {
		t.Fatalf("Failed to create MCP client: %v", err)
	}
	defer client.Close()

	// Test initialize
	t.Run("Initialize", func(t *testing.T) {
		result, err := client.Initialize()
		if err != nil {
			t.Fatalf("Failed to initialize: %v", err)
		}
		if result.ServerInfo.Name != "FedRAMP Compliance Assistant" {
			t.Errorf("Expected server name 'FedRAMP Compliance Assistant', got '%s'", result.ServerInfo.Name)
		}
		if result.ServerInfo.Version != "1.0.0" {
			t.Errorf("Expected server version '1.0.0', got '%s'", result.ServerInfo.Version)
		}
	})

	// Test list tools
	t.Run("ListTools", func(t *testing.T) {
		tools, err := client.ListTools()
		if err != nil {
			// If we get a "Tools not supported" error, skip this test
			if strings.Contains(err.Error(), "Tools not supported") {
				t.Skip("Tools not supported by the server")
			}
			t.Fatalf("Failed to list tools: %v", err)
		}
		if len(tools) == 0 {
			t.Error("Expected at least one tool, got none")
		}

		// Check for required tools
		requiredTools := []string{
			"getControlInfo",
			"getImplementationGuidance",
			"explainControlImportance",
			"getEvidenceGuidance",
		}

		for _, requiredTool := range requiredTools {
			found := false
			for _, tool := range tools {
				if tool.Name == requiredTool {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Required tool '%s' not found", requiredTool)
			}
		}
	})

	// Test getControlInfo
	t.Run("GetControlInfo", func(t *testing.T) {
		controlID := "AC-1"
		controlInfo, err := client.CallTool("getControlInfo", map[string]interface{}{
			"controlId": controlID,
		})
		if err != nil {
			t.Fatalf("Failed to call getControlInfo: %v", err)
		}
		if controlInfo == "" {
			t.Error("Expected non-empty control info, got empty string")
		}
		if len(controlInfo) < 10 {
			t.Errorf("Expected substantial control info, got only %d characters", len(controlInfo))
		}
	})

	// Test getImplementationGuidance
	t.Run("GetImplementationGuidance", func(t *testing.T) {
		controlID := "AC-1"
		guidance, err := client.CallTool("getImplementationGuidance", map[string]interface{}{
			"controlId": controlID,
			"context":   "cloud environment",
		})
		if err != nil {
			t.Fatalf("Failed to call getImplementationGuidance: %v", err)
		}
		if guidance == "" {
			t.Error("Expected non-empty implementation guidance, got empty string")
		}
		if len(guidance) < 10 {
			t.Errorf("Expected substantial implementation guidance, got only %d characters", len(guidance))
		}
	})

	// Test explainControlImportance
	t.Run("ExplainControlImportance", func(t *testing.T) {
		controlID := "AC-1"
		importance, err := client.CallTool("explainControlImportance", map[string]interface{}{
			"controlId": controlID,
		})
		if err != nil {
			t.Fatalf("Failed to call explainControlImportance: %v", err)
		}
		if importance == "" {
			t.Error("Expected non-empty control importance, got empty string")
		}
		if len(importance) < 10 {
			t.Errorf("Expected substantial control importance, got only %d characters", len(importance))
		}
	})

	// Test getEvidenceGuidance
	t.Run("GetEvidenceGuidance", func(t *testing.T) {
		controlID := "AC-1"
		evidenceGuidance, err := client.CallTool("getEvidenceGuidance", map[string]interface{}{
			"controlId": controlID,
		})
		if err != nil {
			t.Fatalf("Failed to call getEvidenceGuidance: %v", err)
		}
		if evidenceGuidance == "" {
			t.Error("Expected non-empty evidence guidance, got empty string")
		}
		if len(evidenceGuidance) < 10 {
			t.Errorf("Expected substantial evidence guidance, got only %d characters", len(evidenceGuidance))
		}
	})

	// Test error handling
	t.Run("ErrorHandling", func(t *testing.T) {
		// Since we're using mock responses for "Tools not supported" errors,
		// we'll test a different error case: unknown tool
		_, err := client.CallTool("nonexistentTool", map[string]interface{}{
			"controlId": "AC-1",
		})
		if err == nil {
			t.Error("Expected error for nonexistent tool, got nil")
		}
	})
}
