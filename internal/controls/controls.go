package controls

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// SecurityControl represents a FedRAMP security control
type SecurityControl struct {
	ID          string `json:"id"`
	Family      string `json:"family"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Impact      string `json:"impact"` // Low, Moderate, High
	Details     string `json:"details"`
}

// SecurityControlFamily represents a group of related security controls
type SecurityControlFamily struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Controls    []SecurityControl `json:"controls"`
}

// Global map of all security controls for quick lookup
var controlsMap = make(map[string]*SecurityControl)

// List of control families
var controlFamilies = []SecurityControlFamily{
	{
		ID:          "AC",
		Name:        "Access Control",
		Description: "Controls related to system access, permissions, and authentication",
		Controls: []SecurityControl{
			{
				ID:          "AC-1",
				Family:      "AC",
				Name:        "Access Control Policy and Procedures",
				Description: "The organization develops, documents, and disseminates an access control policy and procedures.",
				Impact:      "High",
				Details:     "Requires formal, documented policies and procedures for access control that address purpose, scope, roles, responsibilities, management commitment, coordination, and compliance.",
			},
			{
				ID:          "AC-2",
				Family:      "AC",
				Name:        "Account Management",
				Description: "The organization manages information system accounts, including establishing, activating, modifying, reviewing, disabling, and removing accounts.",
				Impact:      "High",
				Details:     "Requires comprehensive account management including account types, group memberships, privileges, and regular review of accounts.",
			},
		},
	},
	{
		ID:          "AU",
		Name:        "Audit and Accountability",
		Description: "Controls related to auditing, logging, and accountability",
		Controls: []SecurityControl{
			{
				ID:          "AU-1",
				Family:      "AU",
				Name:        "Audit and Accountability Policy and Procedures",
				Description: "The organization develops, documents, and disseminates an audit and accountability policy and procedures.",
				Impact:      "High",
				Details:     "Requires formal, documented policies and procedures for audit and accountability that address purpose, scope, roles, responsibilities, management commitment, coordination, and compliance.",
			},
			{
				ID:          "AU-2",
				Family:      "AU",
				Name:        "Audit Events",
				Description: "The organization determines the events to be audited within the information system.",
				Impact:      "High",
				Details:     "Requires identification of auditable events, coordination with organizational entities, and providing a rationale for selected events.",
			},
		},
	},
	{
		ID:          "IA",
		Name:        "Identification and Authentication",
		Description: "Controls related to identifying and authenticating users and devices",
		Controls: []SecurityControl{
			{
				ID:          "IA-1",
				Family:      "IA",
				Name:        "Identification and Authentication Policy and Procedures",
				Description: "The organization develops, documents, and disseminates an identification and authentication policy and procedures.",
				Impact:      "High",
				Details:     "Requires formal, documented policies and procedures for identification and authentication that address purpose, scope, roles, responsibilities, management commitment, coordination, and compliance.",
			},
			{
				ID:          "IA-2",
				Family:      "IA",
				Name:        "Identification and Authentication (Organizational Users)",
				Description: "The information system uniquely identifies and authenticates organizational users.",
				Impact:      "High",
				Details:     "Requires multifactor authentication for network access to privileged and non-privileged accounts and for local access to privileged accounts.",
			},
		},
	},
}

// Initialize the controls map
func init() {
	// Populate the controls map for quick lookup
	for _, family := range controlFamilies {
		for i, control := range family.Controls {
			controlsMap[control.ID] = &family.Controls[i]
		}
	}
}

// RegisterResources registers all control-related resources with the MCP server
func RegisterResources(s *server.MCPServer) {
	// Register resource for getting all control families
	registerControlFamiliesResource(s)

	// Register resource for getting a specific control by ID
	registerControlResource(s)

	// Register resource for searching controls
	registerControlSearchResource(s)
}

// Register a resource for getting all control families
func registerControlFamiliesResource(s *server.MCPServer) {
	resource := mcp.NewResource(
		"fedramp://controls/families",
		"FedRAMP Control Families",
		mcp.WithResourceDescription("List of all FedRAMP security control families"),
		mcp.WithMIMEType("application/json"),
	)

	s.AddResource(resource, func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		// Convert control families to JSON
		data, err := json.MarshalIndent(controlFamilies, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("failed to marshal control families: %v", err)
		}

		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      "fedramp://controls/families",
				MIMEType: "application/json",
				Text:     string(data),
			},
		}, nil
	})
}

// Register a resource for getting a specific control by ID
func registerControlResource(s *server.MCPServer) {
	template := mcp.NewResourceTemplate(
		"fedramp://controls/{id}",
		"FedRAMP Security Control",
		mcp.WithTemplateDescription("Information about a specific FedRAMP security control"),
		mcp.WithTemplateMIMEType("application/json"),
	)

	s.AddResourceTemplate(template, func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		// Extract control ID from URI
		uri := request.Params.URI
		parts := strings.Split(uri, "/")
		if len(parts) != 3 {
			return nil, fmt.Errorf("invalid URI format: %s", uri)
		}

		controlID := parts[2]
		control, exists := controlsMap[controlID]
		if !exists {
			return nil, fmt.Errorf("control not found: %s", controlID)
		}

		// Convert control to JSON
		data, err := json.MarshalIndent(control, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("failed to marshal control: %v", err)
		}

		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      uri,
				MIMEType: "application/json",
				Text:     string(data),
			},
		}, nil
	})
}

// Register a resource for searching controls
func registerControlSearchResource(s *server.MCPServer) {
	template := mcp.NewResourceTemplate(
		"fedramp://controls/search/{query}",
		"Search FedRAMP Security Controls",
		mcp.WithTemplateDescription("Search for FedRAMP security controls by keyword"),
		mcp.WithTemplateMIMEType("application/json"),
	)

	s.AddResourceTemplate(template, func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		// Extract search query from URI
		uri := request.Params.URI
		parts := strings.Split(uri, "/")
		if len(parts) != 4 {
			return nil, fmt.Errorf("invalid URI format: %s", uri)
		}

		query := strings.ToLower(parts[3])

		// Search for controls matching the query
		var results []SecurityControl
		for _, control := range controlsMap {
			if strings.Contains(strings.ToLower(control.ID), query) ||
				strings.Contains(strings.ToLower(control.Name), query) ||
				strings.Contains(strings.ToLower(control.Description), query) {
				results = append(results, *control)
			}
		}

		// Convert results to JSON
		data, err := json.MarshalIndent(results, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("failed to marshal search results: %v", err)
		}

		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      uri,
				MIMEType: "application/json",
				Text:     string(data),
			},
		}, nil
	})
}
