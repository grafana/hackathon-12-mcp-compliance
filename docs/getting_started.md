# Getting Started with FedRAMP Compliance Server

This guide will walk you through setting up and using the FedRAMP Compliance Server, which provides access to FedRAMP controls and requirements through the Model Context Protocol (MCP).

## Prerequisites

- [Go 1.24.1](https://golang.org/dl/) or later
- [Git](https://git-scm.com/downloads)
- [Cursor](https://cursor.sh/) or [Claude Desktop](https://claude.ai/desktop) for using the MCP server

## Clone the Repository

```bash
git clone https://github.com/grafana/hackathon-12-mcp-compliance.git
cd hackathon-12-mcp-compliance
```

## Build and Process FedRAMP Data

The project includes make targets to download and process FedRAMP data files:

```bash
# Download and process FedRAMP High baseline
make run-fedramp-data-high

# Download and process FedRAMP Moderate baseline
make run-fedramp-data-moderate

# Or run both at once
make run-fedramp-data-high run-fedramp-data-moderate
```

These commands will:
1. Download the FedRAMP baseline files from the official source
2. Process them into a more usable format
3. Store the processed files in the `data/` directory
4. Copy the processed files to the `internal/resources/data/` directory for embedding

## Build the MCP Compliance Server

Build the MCP compliance server with:

```bash
make build-mcp-compliance
```

This will create a binary at `bin/mcp-compliance`.

## Deploy Locally

This assumes you have a directory `~/.mcp-compliance/bin` and it's added to your `$PATH`.
To build everything and deploy the server locally:

```bash
make deploy-local
```

This command:
1. Builds the MCP compliance server
2. Processes both FedRAMP High and Moderate baselines
3. Copies the server binary to `~/.mcp-compliance/bin/mcp-compliance`

## Configure Cursor or Claude Desktop

### Cursor Configuration

1. Open Cursor
2. Go to Settings (gear icon) > AI > MCP Servers
3. Click "Add Server"
4. Enter the following details:
   - Name: FedRAMP Compliance
   - Command: `~/.mcp-compliance/bin/mcp-compliance`
5. Click "Save"
6. Enable the server by toggling it on

### Claude Desktop Configuration

1. Open Claude Desktop
2. Go to Settings > Advanced > MCP Servers
3. Click "Add Server"
4. Enter the following details:
   - Name: FedRAMP Compliance
   - Command: `~/.mcp-compliance/bin/mcp-compliance`
5. Click "Add"
6. Enable the server by toggling it on

## Using the Server

Once configured, you can ask Claude about FedRAMP controls and requirements. For example:

- "What is AC-1 in FedRAMP High?"
- "Compare AC-2 between FedRAMP High and Moderate"
- "What evidence should I collect for IA-2?"
- "Summarize the SC control family"

## Available Tools

The server provides several tools for interacting with FedRAMP data:

- `get_control`: Get detailed information about a specific control
- `get_control_family`: Get all controls in a specific family
- `list_control_families`: List all control families in a program
- `search_controls`: Search for controls by keyword
- `get_controls_by_status`: Get controls with a specific implementation status
- `get_control_parameters`: Get parameters for a specific control
- `get_evidence_guidance`: Get evidence guidance for a specific control

## Troubleshooting

If you encounter issues:

1. Ensure the FedRAMP data has been processed:
   ```bash
   make run-fedramp-data-high run-fedramp-data-moderate
   ```

2. Rebuild the server:
   ```bash
   make build-mcp-compliance
   ```

3. Check that the server is properly deployed:
   ```bash
   make deploy-local
   ```

4. Verify the server is running by checking the logs in Cursor or Claude Desktop

## Additional Commands

- `make clean`: Clean build artifacts
- `make clean-resources`: Clean processed resources
- `make clean-all`: Clean all artifacts
- `make help`: Show all available make targets 