package evidence

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// EvidenceGuidance represents guidance for collecting evidence for a security control
type EvidenceGuidance struct {
	ControlID       string   `json:"controlId"`
	Description     string   `json:"description"`
	EvidenceTypes   []string `json:"evidenceTypes"`
	CollectionSteps []string `json:"collectionSteps"`
	Examples        []string `json:"examples"`
	CommonPitfalls  []string `json:"commonPitfalls"`
}

// Map of evidence guidance by control ID
var evidenceGuidanceMap = map[string]EvidenceGuidance{
	"AC-1": {
		ControlID:   "AC-1",
		Description: "Guidance for collecting evidence for Access Control Policy and Procedures",
		EvidenceTypes: []string{
			"Policy Documents",
			"Procedure Documents",
			"Implementation Plans",
			"Security Assessment Reports",
		},
		CollectionSteps: []string{
			"Identify and collect the organization's access control policy",
			"Identify and collect the organization's access control procedures",
			"Collect evidence of regular reviews and updates to the policy and procedures",
			"Collect evidence of dissemination to relevant personnel",
		},
		Examples: []string{
			"Access Control Policy document",
			"Access Control Procedures manual",
			"Email notifications of policy updates",
			"Training records for access control procedures",
		},
		CommonPitfalls: []string{
			"Outdated policy documents",
			"Missing approval signatures",
			"Lack of evidence for regular reviews",
			"Incomplete procedure documentation",
		},
	},
	"AC-2": {
		ControlID:   "AC-2",
		Description: "Guidance for collecting evidence for Account Management",
		EvidenceTypes: []string{
			"Account Management Procedures",
			"System Configuration Screenshots",
			"Account Review Logs",
			"User Access Lists",
		},
		CollectionSteps: []string{
			"Collect account management procedures",
			"Gather evidence of account creation, modification, and termination processes",
			"Collect evidence of regular account reviews",
			"Document system configurations for account management",
		},
		Examples: []string{
			"Account request and approval forms",
			"Screenshots of account management interfaces",
			"Logs of account reviews",
			"User access recertification documentation",
		},
		CommonPitfalls: []string{
			"Missing evidence of account terminations",
			"Incomplete account review documentation",
			"Lack of evidence for privileged account management",
			"Inconsistent application of account management procedures",
		},
	},
	"AU-1": {
		ControlID:   "AU-1",
		Description: "Guidance for collecting evidence for Audit and Accountability Policy and Procedures",
		EvidenceTypes: []string{
			"Audit Policy Documents",
			"Audit Procedure Documents",
			"Implementation Plans",
			"Security Assessment Reports",
		},
		CollectionSteps: []string{
			"Identify and collect the organization's audit and accountability policy",
			"Identify and collect the organization's audit procedures",
			"Collect evidence of regular reviews and updates to the policy and procedures",
			"Collect evidence of dissemination to relevant personnel",
		},
		Examples: []string{
			"Audit and Accountability Policy document",
			"Audit Procedures manual",
			"Email notifications of policy updates",
			"Training records for audit procedures",
		},
		CommonPitfalls: []string{
			"Outdated policy documents",
			"Missing approval signatures",
			"Lack of evidence for regular reviews",
			"Incomplete procedure documentation",
		},
	},
	"AU-2": {
		ControlID:   "AU-2",
		Description: "Guidance for collecting evidence for Audit Events",
		EvidenceTypes: []string{
			"Audit Configuration Documents",
			"System Configuration Screenshots",
			"Audit Logs",
			"Audit Event Selection Documentation",
		},
		CollectionSteps: []string{
			"Document the organization's rationale for selecting auditable events",
			"Collect evidence of coordination with organizational entities",
			"Gather system configuration showing enabled audit events",
			"Document the review and update process for auditable events",
		},
		Examples: []string{
			"List of auditable events with rationale",
			"Screenshots of audit configuration settings",
			"Meeting minutes discussing audit event selection",
			"Documentation of periodic reviews of audit events",
		},
		CommonPitfalls: []string{
			"Missing rationale for selected audit events",
			"Incomplete coverage of required audit events",
			"Lack of evidence for coordination with stakeholders",
			"Outdated audit event configurations",
		},
	},
	"IA-1": {
		ControlID:   "IA-1",
		Description: "Guidance for collecting evidence for Identification and Authentication Policy and Procedures",
		EvidenceTypes: []string{
			"Policy Documents",
			"Procedure Documents",
			"Implementation Plans",
			"Security Assessment Reports",
		},
		CollectionSteps: []string{
			"Identify and collect the organization's identification and authentication policy",
			"Identify and collect the organization's identification and authentication procedures",
			"Collect evidence of regular reviews and updates to the policy and procedures",
			"Collect evidence of dissemination to relevant personnel",
		},
		Examples: []string{
			"Identification and Authentication Policy document",
			"Identification and Authentication Procedures manual",
			"Email notifications of policy updates",
			"Training records for identification and authentication procedures",
		},
		CommonPitfalls: []string{
			"Outdated policy documents",
			"Missing approval signatures",
			"Lack of evidence for regular reviews",
			"Incomplete procedure documentation",
		},
	},
	"IA-2": {
		ControlID:   "IA-2",
		Description: "Guidance for collecting evidence for Identification and Authentication (Organizational Users)",
		EvidenceTypes: []string{
			"Authentication Configuration Documents",
			"System Configuration Screenshots",
			"MFA Implementation Evidence",
			"Authentication Logs",
		},
		CollectionSteps: []string{
			"Document system configurations for user identification and authentication",
			"Collect evidence of multifactor authentication implementation",
			"Gather evidence of unique user identification",
			"Document authentication mechanisms for different access types",
		},
		Examples: []string{
			"Screenshots of authentication settings",
			"MFA configuration documentation",
			"User account listings showing unique identifiers",
			"Authentication logs showing MFA usage",
		},
		CommonPitfalls: []string{
			"Incomplete MFA implementation",
			"Lack of evidence for privileged account authentication",
			"Missing documentation for remote access authentication",
			"Inconsistent application of authentication requirements",
		},
	},
}

// RegisterTools registers all evidence-related tools with the MCP server
func RegisterTools(s *server.MCPServer) {
	// Register tool for getting evidence guidance for a specific control
	registerEvidenceGuidanceTool(s)

	// Register tool for searching evidence guidance
	registerEvidenceSearchTool(s)
}

// Register a tool for getting evidence guidance for a specific control
func registerEvidenceGuidanceTool(s *server.MCPServer) {
	tool := mcp.NewTool("get_evidence_guidance",
		mcp.WithDescription("Get guidance on collecting evidence for a specific FedRAMP security control"),
		mcp.WithString("controlId",
			mcp.Required(),
			mcp.Description("The ID of the security control (e.g., AC-1, IA-2)"),
		),
	)

	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Get control ID from parameters
		controlID, ok := request.Params.Arguments["controlId"].(string)
		if !ok || controlID == "" {
			return mcp.NewToolResultError("Control ID is required"), nil
		}

		// Look up evidence guidance
		guidance, exists := evidenceGuidanceMap[controlID]
		if !exists {
			return mcp.NewToolResultError(fmt.Sprintf("Evidence guidance not found for control: %s", controlID)), nil
		}

		// Convert guidance to JSON
		data, err := json.MarshalIndent(guidance, "", "  ")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal evidence guidance: %v", err)), nil
		}

		return mcp.NewToolResultText(string(data)), nil
	})
}

// Register a tool for searching evidence guidance
func registerEvidenceSearchTool(s *server.MCPServer) {
	tool := mcp.NewTool("search_evidence_guidance",
		mcp.WithDescription("Search for evidence guidance by keyword"),
		mcp.WithString("query",
			mcp.Required(),
			mcp.Description("The search query (e.g., 'policy', 'authentication')"),
		),
	)

	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Get search query from parameters
		query, ok := request.Params.Arguments["query"].(string)
		if !ok || query == "" {
			return mcp.NewToolResultError("Search query is required"), nil
		}

		query = strings.ToLower(query)

		// Search for guidance matching the query
		var results []EvidenceGuidance
		for _, guidance := range evidenceGuidanceMap {
			if strings.Contains(strings.ToLower(guidance.ControlID), query) ||
				strings.Contains(strings.ToLower(guidance.Description), query) {
				results = append(results, guidance)
				continue
			}

			// Search in evidence types
			for _, evidenceType := range guidance.EvidenceTypes {
				if strings.Contains(strings.ToLower(evidenceType), query) {
					results = append(results, guidance)
					break
				}
			}

			// Search in collection steps
			for _, step := range guidance.CollectionSteps {
				if strings.Contains(strings.ToLower(step), query) {
					results = append(results, guidance)
					break
				}
			}
		}

		// Convert results to JSON
		data, err := json.MarshalIndent(results, "", "  ")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal search results: %v", err)), nil
		}

		return mcp.NewToolResultText(string(data)), nil
	})
}
