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

### Get Control Details

Tool: `get_control`

Parameters:
- `program` (required): The FedRAMP program (High or Moderate)
- `controlId` (required): The ID of the control (e.g., AC-1, IA-2)

Returns detailed information about a specific FedRAMP security control.

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

### Get Control Family

Tool: `get_control_family`

Parameters:
- `program` (required): The FedRAMP program (High or Moderate)
- `family` (required): The control family ID (e.g., AC, IA)

Returns all controls in the specified family.

Example:
```json
[
  {
    "id": "AC-1",
    "family": "AC",
    "name": "Access Control Policy and Procedures",
    "description": "The organization develops, documents, and disseminates an access control policy and procedures."
  },
  {
    "id": "AC-2",
    "family": "AC",
    "name": "Account Management",
    "description": "The organization manages information system accounts..."
  }
]
```

### List Control Families

Tool: `list_control_families`

Parameters:
- `program` (required): The FedRAMP program (High or Moderate)

Returns a list of all control families in the program.

Example:
```json
[
  "AC",
  "AU",
  "CM",
  "IA",
  "IR",
  "SC"
]
```

### Search Controls

Tool: `search_controls`

Parameters:
- `program` (required): The FedRAMP program (High or Moderate)
- `query` (required): The search query

Searches for controls matching the query in their title or description.

Example:
```json
[
  {
    "id": "IA-2",
    "family": "IA",
    "name": "Identification and Authentication",
    "description": "The information system uniquely identifies and authenticates organizational users."
  }
]
```

### Get Controls by Status

Tool: `get_controls_by_status`

Parameters:
- `program` (required): The FedRAMP program (High or Moderate)
- `status` (required): The implementation status to filter by

Returns all controls with the specified implementation status.

Example:
```json
[
  {
    "id": "AC-1",
    "family": "AC",
    "name": "Access Control Policy and Procedures",
    "implementation-status": "implemented"
  }
]
```

### Get Control Parameters

Tool: `get_control_parameters`

Parameters:
- `program` (required): The FedRAMP program (High or Moderate)
- `controlId` (required): The ID of the control (e.g., AC-1, IA-2)

Returns the parameters for a specific control.

Example:
```json
{
  "ac-2_prm_1": {
    "label": "time period",
    "value": "90 days"
  },
  "ac-2_prm_2": {
    "label": "frequency",
    "value": "quarterly"
  }
}
```

## Using with an MCP Client

To use the server with an MCP client, you need to:

1. Start the server using `make run`
2. Connect to the server using an MCP client
3. Initialize the connection
4. Call tools as needed

See the example client in the `examples` directory for a complete example of how to use the server with an MCP client.

## Available Programs

The server supports the following compliance programs:
- FedRAMP High
- FedRAMP Moderate

## Error Handling

All tools return errors in a consistent format:

```json
{
  "error": "Error message describing what went wrong"
}
```

Common error cases:
- Program not found
- Control not found
- Invalid control ID format
- Invalid family ID
- Missing required parameters

## Common Questions

Here are examples of questions you can ask about FedRAMP controls:

### General Program Questions

1. "What security controls are required for FedRAMP High?"
2. "What security controls are required for FedRAMP Moderate?"
3. "What are the differences between FedRAMP High and Moderate for access control requirements?"
4. "Which program should I use for my cloud service?"

### Control-Specific Questions

1. "What are the requirements for access control in FedRAMP High?"
2. "How does FedRAMP Moderate handle incident response?"
3. "What is control AC-2 in FedRAMP High?"
4. "What parameters are required for AC-2 in FedRAMP Moderate?"
5. "Show me all the controls in the AC family for FedRAMP High"
6. "What are the audit requirements in FedRAMP Moderate?"

### Implementation Questions

1. "How do I implement AC-2 for FedRAMP High?"
2. "What evidence do I need to collect for AC-1 in FedRAMP Moderate?"
3. "What are the password requirements for FedRAMP High?"
4. "How often do I need to review access in FedRAMP Moderate?"
5. "What documentation is required for incident response in FedRAMP High?"

### Comparison Questions

1. "How do the access control requirements differ between FedRAMP High and Moderate?"
2. "Which program has stricter password requirements?"
3. "What additional controls does FedRAMP High require compared to Moderate?"
4. "Are there any controls in Moderate that aren't in High?"

### Search and Discovery

1. "Find all controls related to authentication"
2. "Which controls mention encryption?"
3. "Show me controls about password requirements"
4. "List all controls that require monthly reviews"
5. "Find controls related to incident response"

## Tips for Asking Questions

1. **Specify the Program**: Always mention which program you're asking about (High or Moderate) if your question is program-specific.

2. **Be Specific**: When asking about controls, use specific control IDs (e.g., AC-2) or clear topics (e.g., "password requirements").

3. **Context Matters**: If you're comparing programs or looking for differences, make that clear in your question.

4. **Implementation Details**: When asking about implementation, specify whether you're looking for:
   - Technical requirements
   - Documentation requirements
   - Evidence collection requirements
   - Review/audit requirements

## Example Interactions

```
Q: "What are the password requirements for FedRAMP High?"
A: The server will provide password-related controls and requirements specific to FedRAMP High.

Q: "How does this compare to FedRAMP Moderate?"
A: The server will highlight the differences in password requirements between High and Moderate.

Q: "What evidence do I need for AC-2?"
A: The server will ask which program (High or Moderate) you're interested in before providing evidence requirements.
```

## Getting Help

If you're not sure how to phrase your question or aren't getting the information you need:

1. Start with general questions about the program or control family
2. Use the search functionality to find relevant controls
3. Ask for clarification about specific requirements
4. Request examples or implementation guidance when needed

Remember that the server is designed to help you understand and implement FedRAMP requirements. Don't hesitate to ask follow-up questions or request clarification when needed.

What is the Access Control (AC) family in FedRAMP? What controls does it include? 

What FedRAMP controls are related to authentication? 

What compliance programs are available in the MCP server? Can you give me an overview of FedRAMP High? 

How should I implement AC-2 (Account Management) in my cloud application? 

I'm building a web application that uses OAuth for authentication. How should I implement IA-2 (Identification and Authentication)? 

What's the best way to implement the audit controls (AU family) together in a cohesive way? 

What evidence should I collect to demonstrate compliance with AC-1? 

How does our company collect evidence for AC-2? 

What evidence do we need to collect for the Identification and Authentication controls? 

What controls require quarterly reviews according to our company practices? 

I'm starting a new cloud project that will need FedRAMP High compliance. What controls should I focus on first, and how should I approach implementation and evidence collection? 

I have implemented basic authentication with username and password. What additional controls do I need to implement to meet FedRAMP High requirements for identification and authentication? 

We have an upcoming FedRAMP audit. What are the most common pitfalls when collecting evidence for the Access Control family, and how can we avoid them? 

How does AC-2 (Account Management) relate to other controls in the FedRAMP framework? Are there dependencies I should be aware of? 