# MCP Compliance

A project for compliance that provides CLI tools and an MCP server for agents to interact with compliance data.

## Overview

This project provides tools for working with FedRAMP compliance data, including:

1. CLI tools for processing and querying FedRAMP baseline data
2. An MCP server that exposes compliance data to LLM agents

## Easy Button Quick Start

For the quickest setup:

```bash
# Create the directory for the MCP compliance binary
mkdir -p ~/.mcp-compliance/bin

# Add to your PATH (add this to your .bashrc or .zshrc for persistence)
export PATH=$PATH:~/.mcp-compliance/bin

# Clone the repository
git clone https://github.com/grafana/hackathon-12-mcp-compliance.git
cd hackathon-12-mcp-compliance

# Build and deploy locally
make deploy-local
```

Now you can configure your agent (Cursor or Claude Desktop) to use the MCP compliance server. See the [Getting Started Guide](docs/getting_started.md) for configuration details.

## Documentation

- [Getting Started Guide](docs/getting_started.md) - Instructions for setting up and using the project
- [Concept of Operations](docs/concept_of_operations.md) - Detailed explanation of system architecture and data flow

## MCP Server

The MCP server provides the following tools for LLM agents:

- `get_control`: Get detailed information about a specific control
- `get_control_family`: Get all controls in a specific family
- `list_control_families`: List all control families in a program
- `search_controls`: Search for controls by keyword
- `get_control_evidence_guidance`: Get detailed guidance for evidence about a specific control

## Data Sources

The FedRAMP baseline files are sourced from the official GSA FedRAMP Automation GitHub repository:
- [FedRAMP Rev 5 HIGH Baseline](https://github.com/GSA/fedramp-automation/blob/master/dist/content/rev5/baselines/json/FedRAMP_rev5_HIGH-baseline-resolved-profile_catalog.json)
- [FedRAMP Rev 5 MODERATE Baseline](https://github.com/GSA/fedramp-automation/blob/master/dist/content/rev5/baselines/json/FedRAMP_rev5_MODERATE-baseline-resolved-profile_catalog.json)
