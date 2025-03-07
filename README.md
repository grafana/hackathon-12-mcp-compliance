# MCP Compliance

A project for compliance that provides CLI tools and an MCP server for agents to interact with compliance data.

## Overview

This project provides tools for working with FedRAMP compliance data, including:

1. CLI tools for processing and querying FedRAMP baseline data
2. An MCP server that exposes compliance data to LLM agents

## Getting Started

### Prerequisites

- Go 1.24 or later
- Make (optional, for using the Makefile)
- curl (for downloading FedRAMP files)

### Building

To build all tools:

```bash
make build
```

Or build individual tools:

```bash
make build-fedramp-data
```

### Downloading FedRAMP Data

The project can automatically download the FedRAMP baseline files from the official GitHub repository:

```bash
# Download all FedRAMP baseline files
make download-fedramp-files

# Or download specific baselines
make download-fedramp-high
make download-fedramp-moderate
```

The files will be downloaded to the `data/` directory.

### Processing FedRAMP Data

To process the FedRAMP High baseline data (will download if needed):

```bash
make run-fedramp-data-high
```

This will read the FedRAMP High baseline JSON file and output a simplified JSON file to `data/processed/fedramp-high.json`.

Similarly, for the Moderate baseline:

```bash
make run-fedramp-data-moderate
```

### Searching for Controls

To search for controls by keyword:

```bash
# Search in FedRAMP High baseline
make search-high QUERY=encryption

# Search in FedRAMP Moderate baseline
make search-moderate QUERY=encryption
```

These commands will download the baseline files if needed and search for controls matching the keyword without generating processed output files.

You can also use the CLI tool directly for more control:

```bash
# Search only
bin/fedramp-data -input data/FedRAMP_rev5_HIGH-baseline-resolved-profile_catalog.json -output /dev/null -program "FedRAMP High" -search encryption

# Search and generate processed file
bin/fedramp-data -input data/FedRAMP_rev5_HIGH-baseline-resolved-profile_catalog.json -output data/processed/fedramp-high.json -program "FedRAMP High" -search encryption
```

## MCP Server

The MCP server provides the following tools for LLM agents:

- `get_control`: Get detailed information about a specific control
- `get_control_family`: Get all controls in a specific family
- `list_control_families`: List all control families in a program
- `search_controls`: Search for controls by keyword
- `get_control_evidence_guidance`: Get detailed guidance for evidence about a specific control

## Project Structure

```
.
├── cmd/                    # Command-line tools
│   └── fedramp_data/       # Tool for processing FedRAMP baseline data
├── data/                   # Input data files
│   └── processed/          # Processed data files
├── internal/               # Internal packages
│   ├── adapters/           # Adapter implementations
│   ├── domain/             # Domain layer
│   │   └── fedramp/        # FedRAMP domain models and commands
│   ├── ports/              # Port interfaces
│   └── services/           # Service layer
│       └── fedramp_data/   # FedRAMP data processing service
│           └── fedramp_data_handlers/ # Specialized handlers for different operations
├── Makefile                # Build automation
└── README.md               # This file
```

## Architecture

The project follows a hexagonal architecture pattern with domain-driven design, command-handler pattern, and ports and adapters:

- **CLI Layer** (`cmd/fedramp_data/main.go`): Handles command-line arguments and user interaction
- **Domain Layer** (`internal/domain/fedramp/`): Contains domain-specific models and commands
- **Service Layer** (`internal/services/fedramp_data/`): Contains services that validate input and create commands
- **Handler Layer** (`internal/services/fedramp_data/fedramp_data_handlers/`): Contains specialized handlers for different types of operations
- **Ports Layer** (`internal/ports/`): Contains interfaces that define contracts for external dependencies
- **Adapters Layer** (`internal/adapters/`): Contains implementations of port interfaces

### Command-Handler Pattern

The project uses the command-handler pattern for better separation of concerns:

1. The CLI layer calls the service layer
2. The service layer validates input and creates command objects
3. The service passes commands to appropriate specialized handlers
4. Handlers contain the business logic to process commands and return results

This pattern provides several benefits:
- Clear separation between input validation and business logic
- Improved testability through command objects
- Better maintainability with single-responsibility components

### Specialized Handlers

The project uses specialized handlers for different types of operations:

1. **FileHandler**: Handles file processing operations (requires file and OSCAL repositories)
2. **SearchHandler**: Handles search-related operations (operates on in-memory data)
3. **ControlHandler**: Handles control-related operations (operates on in-memory data)

This approach provides several benefits:
- Each handler only receives the dependencies it needs
- Improved modularity and testability
- Smaller, more focused components with clear responsibilities

### Ports and Adapters Pattern (Hexagonal Architecture)

The project also implements the ports and adapters pattern (also known as hexagonal architecture):

1. **Ports**: Interfaces that define contracts for external dependencies (e.g., file system, OSCAL parsing)
2. **Adapters**: Implementations of port interfaces that interact with external systems
3. **Handlers**: Core business logic that depends on ports, not concrete implementations
4. **Service**: Creates and configures adapters, then passes them to appropriate handlers

This pattern provides several benefits:
- Decoupling business logic from external dependencies
- Improved testability through mock implementations of ports
- Flexibility to change implementations without affecting business logic
- Clear boundaries between the application core and external systems

## Data Sources

The FedRAMP baseline files are sourced from the official GSA FedRAMP Automation GitHub repository:
- [FedRAMP Rev 5 HIGH Baseline](https://github.com/GSA/fedramp-automation/blob/master/dist/content/rev5/baselines/json/FedRAMP_rev5_HIGH-baseline-resolved-profile_catalog.json)
- [FedRAMP Rev 5 MODERATE Baseline](https://github.com/GSA/fedramp-automation/blob/master/dist/content/rev5/baselines/json/FedRAMP_rev5_MODERATE-baseline-resolved-profile_catalog.json)

## License

[Add license information here] 