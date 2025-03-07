.PHONY: build clean run-fedramp-data download-fedramp-files download-fedramp-high download-fedramp-moderate build-test-compliance run-test-compliance build-mcp-compliance run-mcp-compliance deploy-local

# Default target
all: build

# Build all binaries
build: build-fedramp-data build-test-compliance build-mcp-compliance

# Build the fedramp-data tool
build-fedramp-data:
	@echo "Building fedramp-data..."
	@go build -o bin/fedramp-data ./cmd/fedramp_data

# Build the test-compliance tool
build-test-compliance:
	@echo "Building test-compliance..."
	@go build -o bin/test-compliance ./cmd/test-compliance

# Build the mcp-compliance server
build-mcp-compliance:
	@echo "Building mcp-compliance..."
	@go build -o bin/mcp-compliance ./cmd/mcp-compliance

# Deploy the mcp-compliance server locally
deploy-local: build-mcp-compliance run-fedramp-data-high run-fedramp-data-moderate
	@echo "Deploying mcp-compliance to ~/.mcp-compliance/bin..."
	@mkdir -p ~/.mcp-compliance/bin
	@cp bin/mcp-compliance ~/.mcp-compliance/bin/
	@echo "Deployed mcp-compliance to ~/.mcp-compliance/bin/mcp-compliance"
	@echo "You can now use this server in Cursor"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/

# Download FedRAMP baseline files
download-fedramp-files: download-fedramp-high download-fedramp-moderate
	@echo "All FedRAMP baseline files downloaded to data/ directory"

# Download FedRAMP High baseline
download-fedramp-high:
	@echo "Downloading FedRAMP High baseline..."
	@mkdir -p data
	@curl -s -o data/FedRAMP_rev5_HIGH-baseline-resolved-profile_catalog.json \
		https://raw.githubusercontent.com/GSA/fedramp-automation/refs/heads/master/dist/content/rev5/baselines/json/FedRAMP_rev5_HIGH-baseline-resolved-profile_catalog.json
	@echo "FedRAMP High baseline downloaded to data/FedRAMP_rev5_HIGH-baseline-resolved-profile_catalog.json"

# Download FedRAMP Moderate baseline
download-fedramp-moderate:
	@echo "Downloading FedRAMP Moderate baseline..."
	@mkdir -p data
	@curl -s -o data/FedRAMP_rev5_MODERATE-baseline-resolved-profile_catalog.json \
		https://raw.githubusercontent.com/GSA/fedramp-automation/refs/heads/master/dist/content/rev5/baselines/json/FedRAMP_rev5_MODERATE-baseline-resolved-profile_catalog.json
	@echo "FedRAMP Moderate baseline downloaded to data/FedRAMP_rev5_MODERATE-baseline-resolved-profile_catalog.json"

# Run the fedramp-data tool with FedRAMP High baseline
run-fedramp-data-high: download-fedramp-high
	@echo "Processing FedRAMP High baseline..."
	@mkdir -p data/processed
	@bin/fedramp-data \
		-input data/FedRAMP_rev5_HIGH-baseline-resolved-profile_catalog.json \
		-output data/processed/fedramp-high.json \
		-program "FedRAMP High"

# Run the fedramp-data tool with FedRAMP Moderate baseline
run-fedramp-data-moderate: download-fedramp-moderate
	@echo "Processing FedRAMP Moderate baseline..."
	@mkdir -p data/processed
	@bin/fedramp-data \
		-input data/FedRAMP_rev5_MODERATE-baseline-resolved-profile_catalog.json \
		-output data/processed/fedramp-moderate.json \
		-program "FedRAMP Moderate"

# Search for controls in the FedRAMP High baseline
search-high: download-fedramp-high
	@echo "Searching FedRAMP High baseline..."
	@bin/fedramp-data \
		-input data/FedRAMP_rev5_HIGH-baseline-resolved-profile_catalog.json \
		-output /dev/null \
		-program "FedRAMP High" \
		-search $(QUERY)

# Search for controls in the FedRAMP Moderate baseline
search-moderate: download-fedramp-moderate
	@echo "Searching FedRAMP Moderate baseline..."
	@bin/fedramp-data \
		-input data/FedRAMP_rev5_MODERATE-baseline-resolved-profile_catalog.json \
		-output /dev/null \
		-program "FedRAMP Moderate" \
		-search $(QUERY)

# Run the test-compliance tool
run-test-compliance: build-test-compliance run-fedramp-data-high run-fedramp-data-moderate
	@echo "Running test-compliance..."
	@bin/test-compliance

# Run the mcp-compliance server
run-mcp-compliance: build-mcp-compliance run-fedramp-data-high run-fedramp-data-moderate
	@echo "Running mcp-compliance server..."
	@bin/mcp-compliance

# Help target
help:
	@echo "Available targets:"
	@echo "  build                - Build all binaries"
	@echo "  build-fedramp-data   - Build the fedramp-data tool"
	@echo "  clean                - Clean build artifacts"
	@echo "  download-fedramp-files - Download all FedRAMP baseline files"
	@echo "  download-fedramp-high - Download FedRAMP High baseline"
	@echo "  download-fedramp-moderate - Download FedRAMP Moderate baseline"
	@echo "  run-fedramp-data-high - Process FedRAMP High baseline (downloads if needed)"
	@echo "  run-fedramp-data-moderate - Process FedRAMP Moderate baseline (downloads if needed)"
	@echo "  search-high QUERY=<keyword> - Search for controls in FedRAMP High baseline (downloads if needed)"
	@echo "  search-moderate QUERY=<keyword> - Search for controls in FedRAMP Moderate baseline (downloads if needed)"
	@echo "  run-test-compliance    - Run test-compliance"
	@echo "  run-mcp-compliance     - Run mcp-compliance server"
	@echo "  deploy-local           - Deploy mcp-compliance server locally"
	@echo "  help                 - Show this help message" 