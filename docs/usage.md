# FedRAMP Compliance MCP Server - Usage Guide

This document provides detailed instructions on how to use the FedRAMP Compliance MCP Server.

## Starting the Server

To start the server, use the following command:

```bash
make run
```

This will build the server and start it in stdio mode, ready to accept connections from MCP clients.

## Available Resources

The server provides the following resources:

### Control Families

URI: `fedramp://controls/families`

Returns a list of all FedRAMP security control families, including their IDs, names, descriptions, and the controls within each family.

Example:
```json
[
  {
    "id": "AC",
    "name": "Access Control",
    "description": "Controls related to system access, permissions, and authentication",
    "controls": [
      {
        "id": "AC-1",
        "family": "AC",
        "name": "Access Control Policy and Procedures",
        "description": "The organization develops, documents, and disseminates an access control policy and procedures.",
        "impact": "High",
        "details": "Requires formal, documented policies and procedures for access control that address purpose, scope, roles, responsibilities, management commitment, coordination, and compliance."
      },
      ...
    ]
  },
  ...
]
```

### Specific Control

URI: `fedramp://controls/{id}`

Returns detailed information about a specific FedRAMP security control, where `{id}` is the control ID (e.g., AC-1, IA-2).

Example:
```json
{
  "id": "AC-1",
  "family": "AC",
  "name": "Access Control Policy and Procedures",
  "description": "The organization develops, documents, and disseminates an access control policy and procedures.",
  "impact": "High",
  "details": "Requires formal, documented policies and procedures for access control that address purpose, scope, roles, responsibilities, management commitment, coordination, and compliance."
}
```

### Search Controls

URI: `fedramp://controls/search/{query}`

Searches for FedRAMP security controls matching the specified query, where `{query}` is the search term.

Example:
```json
[
  {
    "id": "IA-2",
    "family": "IA",
    "name": "Identification and Authentication (Organizational Users)",
    "description": "The information system uniquely identifies and authenticates organizational users.",
    "impact": "High",
    "details": "Requires multifactor authentication for network access to privileged and non-privileged accounts and for local access to privileged accounts."
  },
  ...
]
```

## Available Tools

The server provides the following tools:

### Get Evidence Guidance

Tool: `get_evidence_guidance`

Parameters:
- `controlId` (required): The ID of the security control (e.g., AC-1, IA-2)

Returns guidance on collecting evidence for the specified FedRAMP security control.

Example:
```json
{
  "controlId": "AC-1",
  "description": "Guidance for collecting evidence for Access Control Policy and Procedures",
  "evidenceTypes": [
    "Policy Documents",
    "Procedure Documents",
    "Implementation Plans",
    "Security Assessment Reports"
  ],
  "collectionSteps": [
    "Identify and collect the organization's access control policy",
    "Identify and collect the organization's access control procedures",
    "Collect evidence of regular reviews and updates to the policy and procedures",
    "Collect evidence of dissemination to relevant personnel"
  ],
  "examples": [
    "Access Control Policy document",
    "Access Control Procedures manual",
    "Email notifications of policy updates",
    "Training records for access control procedures"
  ],
  "commonPitfalls": [
    "Outdated policy documents",
    "Missing approval signatures",
    "Lack of evidence for regular reviews",
    "Incomplete procedure documentation"
  ]
}
```

### Search Evidence Guidance

Tool: `search_evidence_guidance`

Parameters:
- `query` (required): The search query (e.g., 'policy', 'authentication')

Searches for evidence guidance matching the specified query.

Example:
```json
[
  {
    "controlId": "AC-1",
    "description": "Guidance for collecting evidence for Access Control Policy and Procedures",
    "evidenceTypes": [
      "Policy Documents",
      "Procedure Documents",
      "Implementation Plans",
      "Security Assessment Reports"
    ],
    ...
  },
  {
    "controlId": "AU-1",
    "description": "Guidance for collecting evidence for Audit and Accountability Policy and Procedures",
    "evidenceTypes": [
      "Audit Policy Documents",
      "Audit Procedure Documents",
      "Implementation Plans",
      "Security Assessment Reports"
    ],
    ...
  },
  ...
]
```

## Using with an MCP Client

To use the server with an MCP client, you need to:

1. Start the server using `make run`
2. Connect to the server using an MCP client
3. Initialize the connection
4. Call resources and tools as needed

See the example client in the `examples` directory for a complete example of how to use the server with an MCP client. 