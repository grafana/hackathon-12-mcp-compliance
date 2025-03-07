package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/grafana/hackathon-12-mcp-compliance/internal/services/fedramp_data"
)

func main() {
	// Define command-line flags
	inputFile := flag.String("input", "", "Path to the FedRAMP baseline JSON file")
	outputFile := flag.String("output", "", "Path to the output JSON file")
	programName := flag.String("program", "FedRAMP High", "Program name (e.g., FedRAMP High, FedRAMP Moderate)")
	searchQuery := flag.String("search", "", "Search for controls by keyword (optional)")
	flag.Parse()

	// Validate flags
	if *inputFile == "" {
		fmt.Println("Usage: fedramp-data -input <input-file> -output <output-file> [-program <program-name>] [-search <keyword>]")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// If we're just searching, we don't need an output file
	if *searchQuery != "" && (*outputFile == "" || *outputFile == "/dev/null") {
		// This is fine, we're just searching
	} else if *outputFile == "" {
		fmt.Println("Error: Output file is required unless searching with -search flag")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Ensure input file exists
	if _, err := os.Stat(*inputFile); os.IsNotExist(err) {
		log.Fatalf("Input file does not exist: %s", *inputFile)
	}

	// Create output directory if it doesn't exist and we're writing to a file
	if *outputFile != "" && *outputFile != "/dev/null" {
		outputDir := filepath.Dir(*outputFile)
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			log.Fatalf("Failed to create output directory: %v", err)
		}
	}

	// Create a new FedRAMP service
	service := fedramp_data.NewService()

	// Process the file
	fmt.Printf("Processing %s...\n", *inputFile)
	programData, err := service.ProcessFile(*inputFile, *programName)
	if err != nil {
		log.Fatalf("Failed to process file: %v", err)
	}

	// If search query is provided, search for controls
	if *searchQuery != "" {
		fmt.Printf("Searching for controls matching '%s'...\n", *searchQuery)
		results := service.SearchControls(programData, *searchQuery)
		fmt.Printf("Found %d matching controls\n", len(results))
		for _, control := range results {
			fmt.Printf("- %s: %s\n", control.ID, control.Title)
		}
	}

	// Write the output if an output file is specified and it's not /dev/null
	if *outputFile != "" && *outputFile != "/dev/null" {
		fmt.Printf("Writing output to %s...\n", *outputFile)
		if err := service.WriteOutput(programData, *outputFile); err != nil {
			log.Fatalf("Failed to write output: %v", err)
		}
		fmt.Printf("Output written to %s\n", *outputFile)
	}
}
