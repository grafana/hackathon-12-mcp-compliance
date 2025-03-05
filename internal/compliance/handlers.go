package compliance

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
)

// HandleGetControlInfo handles requests for information about a specific security control
func HandleGetControlInfo(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	controlId := request.Params.Arguments["controlId"].(string)

	control, err := GetControlInfo(controlId)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Error retrieving control information: %v", err)), nil
	}

	// Format the control information
	controlInfo := fmt.Sprintf("# %s: %s\n\n", control.ID, control.Title)
	controlInfo += fmt.Sprintf("## Description\n\n%s\n\n", control.Description)

	// Add control parts if available
	if len(control.Parts) > 0 {
		controlInfo += "## Control Requirements\n\n"
		for _, part := range control.Parts {
			controlInfo += fmt.Sprintf("### %s\n\n%s\n\n", part.Name, part.Prose)
		}
	}

	// Add parameters if available
	if len(control.Parameters) > 0 {
		controlInfo += "## Parameters\n\n"
		for _, param := range control.Parameters {
			controlInfo += fmt.Sprintf("- **%s**: %s\n", param.Label, param.Description)
			if len(param.Values) > 0 {
				controlInfo += "  - Values: "
				for i, val := range param.Values {
					if i > 0 {
						controlInfo += ", "
					}
					controlInfo += val
				}
				controlInfo += "\n"
			}
		}
	}

	return mcp.NewToolResultText(controlInfo), nil
}

// HandleGetImplementationGuidance handles requests for guidance on implementing a security control
func HandleGetImplementationGuidance(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	controlId := request.Params.Arguments["controlId"].(string)

	// Check if context was provided
	var context string
	if contextVal, ok := request.Params.Arguments["context"]; ok {
		context = contextVal.(string)
	}

	control, err := GetControlInfo(controlId)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Error retrieving control information: %v", err)), nil
	}

	// Generate implementation guidance based on the control
	guidance := fmt.Sprintf("# Implementation Guidance for %s: %s\n\n", control.ID, control.Title)

	// Add context-specific information if provided
	if context != "" {
		guidance += fmt.Sprintf("## Context: %s\n\n", context)
	}

	guidance += "## General Implementation Approach\n\n"

	// Add control-specific guidance
	guidance += fmt.Sprintf("To implement %s, you should consider the following steps:\n\n", control.ID)

	// Generate guidance based on control parts
	if len(control.Parts) > 0 {
		for _, part := range control.Parts {
			guidance += fmt.Sprintf("### For requirement: %s\n\n", part.Name)
			guidance += fmt.Sprintf("- Requirement: %s\n", part.Prose)
			guidance += "- Implementation suggestions:\n"

			// TODO: Generate more specific implementation guidance based on the control part
			guidance += "  - Document your approach in policies and procedures\n"
			guidance += "  - Implement technical controls to enforce this requirement\n"
			guidance += "  - Establish monitoring and auditing to verify compliance\n\n"
		}
	} else {
		guidance += "- Document your approach in policies and procedures\n"
		guidance += "- Implement technical controls to enforce this requirement\n"
		guidance += "- Establish monitoring and auditing to verify compliance\n\n"
	}

	guidance += "## Documentation Requirements\n\n"
	guidance += "Ensure you document the following:\n\n"
	guidance += "- Your implementation approach\n"
	guidance += "- Responsible parties\n"
	guidance += "- Testing and validation procedures\n"
	guidance += "- Ongoing monitoring approach\n"

	return mcp.NewToolResultText(guidance), nil
}

// HandleExplainControlImportance handles requests to explain why a security control is important
func HandleExplainControlImportance(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	controlId := request.Params.Arguments["controlId"].(string)

	control, err := GetControlInfo(controlId)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Error retrieving control information: %v", err)), nil
	}

	// Generate explanation of importance
	importance := fmt.Sprintf("# Importance of %s: %s\n\n", control.ID, control.Title)
	importance += "## Security Risks Addressed\n\n"

	// TODO: Generate more specific risk information based on the control
	// For now, provide generic information based on the control family

	// Extract control family from ID (e.g., AC from AC-1)
	family := ""
	if len(controlId) >= 2 {
		family = controlId[:2]
	}

	switch family {
	case "AC":
		importance += "This control addresses risks related to unauthorized access to systems and data. Without proper access controls:\n\n"
		importance += "- Unauthorized users may gain access to sensitive information\n"
		importance += "- Insiders may abuse their privileges\n"
		importance += "- Systems may be vulnerable to lateral movement by attackers\n"
	case "AU":
		importance += "This control addresses risks related to detecting and investigating security incidents. Without proper auditing:\n\n"
		importance += "- Security incidents may go undetected\n"
		importance += "- Forensic investigations may lack necessary evidence\n"
		importance += "- Compliance requirements for monitoring may not be met\n"
	case "IA":
		importance += "This control addresses risks related to verifying the identity of users and systems. Without proper identification and authentication:\n\n"
		importance += "- Unauthorized users may impersonate legitimate users\n"
		importance += "- Accountability for actions may be compromised\n"
		importance += "- Multi-factor authentication protections may be absent\n"
	case "CM":
		importance += "This control addresses risks related to system configuration and change management. Without proper configuration management:\n\n"
		importance += "- Systems may drift from secure configurations\n"
		importance += "- Unauthorized changes may introduce vulnerabilities\n"
		importance += "- Security patches may not be applied consistently\n"
	default:
		importance += "This control addresses important security risks that could impact the confidentiality, integrity, or availability of your systems and data. Proper implementation helps protect against threats and vulnerabilities that could lead to security breaches or compliance violations.\n"
	}

	importance += "\n## Compliance Impact\n\n"
	importance += "FedRAMP High compliance requires this control to be implemented. Failure to properly implement this control:\n\n"
	importance += "- May result in non-compliance with FedRAMP High requirements\n"
	importance += "- Could prevent authorization to operate (ATO) approval\n"
	importance += "- May require remediation plans and additional oversight\n"

	return mcp.NewToolResultText(importance), nil
}

// HandleGetEvidenceGuidance handles requests for guidance on finding or collecting evidence
func HandleGetEvidenceGuidance(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	controlId := request.Params.Arguments["controlId"].(string)

	control, err := GetControlInfo(controlId)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Error retrieving control information: %v", err)), nil
	}

	// Generate evidence guidance
	guidance := fmt.Sprintf("# Evidence Collection Guidance for %s: %s\n\n", control.ID, control.Title)
	guidance += "## Types of Evidence to Collect\n\n"

	// Extract control family from ID (e.g., AC from AC-1)
	family := ""
	if len(controlId) >= 2 {
		family = controlId[:2]
	}

	// Provide evidence guidance based on control family
	switch family {
	case "AC":
		guidance += "For access control requirements, collect:\n\n"
		guidance += "- Access control policies and procedures\n"
		guidance += "- Screenshots of access control settings in systems\n"
		guidance += "- User access reviews and documentation\n"
		guidance += "- Role-based access control configurations\n"
	case "AU":
		guidance += "For audit requirements, collect:\n\n"
		guidance += "- Audit policies and procedures\n"
		guidance += "- Screenshots of audit configurations\n"
		guidance += "- Sample audit logs showing required events\n"
		guidance += "- Evidence of log review processes\n"
	case "IA":
		guidance += "For identification and authentication requirements, collect:\n\n"
		guidance += "- Authentication policies and procedures\n"
		guidance += "- Screenshots of authentication settings\n"
		guidance += "- Multi-factor authentication configurations\n"
		guidance += "- Password policy enforcement evidence\n"
	case "CM":
		guidance += "For configuration management requirements, collect:\n\n"
		guidance += "- Configuration management policies and procedures\n"
		guidance += "- Baseline configuration documentation\n"
		guidance += "- Change management records\n"
		guidance += "- Configuration scanning results\n"
	default:
		guidance += "For this control, consider collecting:\n\n"
		guidance += "- Relevant policies and procedures\n"
		guidance += "- System configuration evidence\n"
		guidance += "- Process documentation\n"
		guidance += "- Testing and validation results\n"
	}

	guidance += "\n## Documentation Format\n\n"
	guidance += "Evidence should be well-documented with:\n\n"
	guidance += "- Clear labels identifying the control being addressed\n"
	guidance += "- Timestamps showing when the evidence was collected\n"
	guidance += "- Explanations of how the evidence satisfies the control requirements\n"
	guidance += "- Contact information for the person who collected the evidence\n"

	guidance += "\n## Collection Methods\n\n"
	guidance += "Consider these methods for collecting evidence:\n\n"
	guidance += "- System screenshots with timestamps and user information visible\n"
	guidance += "- Export of configuration settings from systems\n"
	guidance += "- Policy and procedure documents with approval signatures\n"
	guidance += "- Automated scanning reports from security tools\n"
	guidance += "- Interview notes from discussions with system administrators\n"

	return mcp.NewToolResultText(guidance), nil
}
