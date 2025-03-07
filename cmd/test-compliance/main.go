package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/grafana/hackathon-12-mcp-compliance/internal/services/fedramp_compliance"
)

func main() {
	// Create the compliance service
	service := fedramp_compliance.NewService()

	// List available programs
	fmt.Println("Available compliance programs:")
	programs, err := service.ListCompliancePrograms()
	if err != nil {
		fmt.Printf("Error listing programs: %v\n", err)
		os.Exit(1)
	}
	for _, program := range programs {
		fmt.Printf("- %s\n", program)
	}
	fmt.Println()

	// Get a control
	programName := "FedRAMP High"
	controlID := "ac-1"
	fmt.Printf("Getting control %s from %s:\n", controlID, programName)
	control, found, err := service.GetControl(programName, controlID)
	if err != nil {
		fmt.Printf("Error getting control: %v\n", err)
		os.Exit(1)
	}
	if !found {
		fmt.Printf("Control %s not found in %s\n", controlID, programName)
		os.Exit(1)
	}
	printJSON("Control", control)
	fmt.Println()

	// Get a control family
	familyID := "ac"
	fmt.Printf("Getting control family %s from %s:\n", familyID, programName)
	family, found, err := service.GetControlFamily(programName, familyID)
	if err != nil {
		fmt.Printf("Error getting control family: %v\n", err)
		os.Exit(1)
	}
	if !found {
		fmt.Printf("Control family %s not found in %s\n", familyID, programName)
		os.Exit(1)
	}
	fmt.Printf("Family: %s - %s (%d controls)\n", family.ID, family.Title, len(family.Controls))
	fmt.Println()

	// List control families
	fmt.Printf("Listing control families in %s:\n", programName)
	families, err := service.ListControlFamilies(programName)
	if err != nil {
		fmt.Printf("Error listing control families: %v\n", err)
		os.Exit(1)
	}
	for _, family := range families {
		fmt.Printf("- %s: %s (%d controls)\n", family.ID, family.Title, len(family.Controls))
	}
	fmt.Println()

	// Search for controls
	query := "access"
	fmt.Printf("Searching for controls with keyword '%s' in %s:\n", query, programName)
	results, err := service.SearchControls(programName, query)
	if err != nil {
		fmt.Printf("Error searching controls: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Found %d controls:\n", len(results))
	for _, control := range results {
		fmt.Printf("- %s: %s\n", control.ID, control.Title)
	}
	fmt.Println()

	// Get control evidence guidance
	fmt.Printf("Getting evidence guidance for control %s in %s:\n", controlID, programName)
	guidance, found, err := service.GetControlEvidenceGuidance(programName, controlID)
	if err != nil {
		fmt.Printf("Error getting evidence guidance: %v\n", err)
		os.Exit(1)
	}
	if !found {
		fmt.Printf("Control %s not found in %s\n", controlID, programName)
		os.Exit(1)
	}
	fmt.Printf("Evidence Guidance:\n%s\n", guidance)
}

// Helper function to print JSON
func printJSON(label string, v interface{}) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling %s to JSON: %v\n", label, err)
		return
	}
	fmt.Printf("%s:\n%s\n", label, string(data))
}
