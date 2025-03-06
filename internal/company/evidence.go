package company

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

//go:embed data/*.yaml
var dataFiles embed.FS

// EvidencePractice represents a company-specific evidence collection practice
type EvidencePractice struct {
	ControlID       string   `yaml:"-"` // Set from the map key
	Practice        string   `yaml:"practice"`
	ResponsibleTeam string   `yaml:"responsible_team"`
	Artifacts       []string `yaml:"artifacts"`
	ReviewFrequency string   `yaml:"review_frequency"`
	Notes           string   `yaml:"notes"`
}

// EvidencePractices stores company-specific evidence collection practices
type EvidencePractices struct {
	Practices map[string]*EvidencePractice
}

// NewEvidencePractices creates a new EvidencePractices instance
func NewEvidencePractices() (*EvidencePractices, error) {
	practices := &EvidencePractices{
		Practices: make(map[string]*EvidencePractice),
	}

	// Try to load from embedded files first
	err := practices.loadFromEmbedded()
	if err != nil {
		// If that fails, try to load from external file
		err = practices.loadFromExternal()
		if err != nil {
			return nil, fmt.Errorf("failed to load evidence practices: %v", err)
		}
	}

	return practices, nil
}

// loadFromEmbedded loads evidence practices from embedded files
func (ep *EvidencePractices) loadFromEmbedded() error {
	// Read the embedded file
	data, err := dataFiles.ReadFile("data/evidence_practices.yaml")
	if err != nil {
		return err
	}

	return ep.parseYAML(data)
}

// loadFromExternal loads evidence practices from an external file
func (ep *EvidencePractices) loadFromExternal() error {
	// Check common locations for the file
	locations := []string{
		"./evidence_practices.yaml",
		"./data/evidence_practices.yaml",
		filepath.Join(os.Getenv("HOME"), ".mcp-compliance/evidence_practices.yaml"),
	}

	for _, location := range locations {
		data, err := os.ReadFile(location)
		if err == nil {
			return ep.parseYAML(data)
		}
	}

	return fmt.Errorf("evidence practices file not found in any of the expected locations")
}

// parseYAML parses the YAML data into the EvidencePractices struct
func (ep *EvidencePractices) parseYAML(data []byte) error {
	// Parse the YAML into a map
	var practicesMap map[string]*EvidencePractice
	if err := yaml.Unmarshal(data, &practicesMap); err != nil {
		return fmt.Errorf("failed to parse evidence practices YAML: %v", err)
	}

	// Set the ControlID field from the map key
	for id, practice := range practicesMap {
		practice.ControlID = id
		ep.Practices[strings.ToUpper(id)] = practice
	}

	return nil
}

// GetPractice returns the evidence practice for a specific control
func (ep *EvidencePractices) GetPractice(controlID string) (*EvidencePractice, error) {
	practice, exists := ep.Practices[strings.ToUpper(controlID)]
	if !exists {
		return nil, fmt.Errorf("no evidence practice found for control: %s", controlID)
	}
	return practice, nil
}

// SearchPractices searches for evidence practices by keyword
func (ep *EvidencePractices) SearchPractices(query string) []*EvidencePractice {
	query = strings.ToLower(query)
	var results []*EvidencePractice

	for _, practice := range ep.Practices {
		// Search in control ID
		if strings.Contains(strings.ToLower(practice.ControlID), query) {
			results = append(results, practice)
			continue
		}

		// Search in practice description
		if strings.Contains(strings.ToLower(practice.Practice), query) {
			results = append(results, practice)
			continue
		}

		// Search in responsible team
		if strings.Contains(strings.ToLower(practice.ResponsibleTeam), query) {
			results = append(results, practice)
			continue
		}

		// Search in artifacts
		for _, artifact := range practice.Artifacts {
			if strings.Contains(strings.ToLower(artifact), query) {
				results = append(results, practice)
				break
			}
		}

		// Search in notes
		if strings.Contains(strings.ToLower(practice.Notes), query) {
			results = append(results, practice)
		}
	}

	return results
}

// FormatPractice formats an evidence practice as a human-readable string
func FormatPractice(practice *EvidencePractice) string {
	text := fmt.Sprintf("Company Evidence Practice for %s\n\n", practice.ControlID)
	text += fmt.Sprintf("%s\n\n", practice.Practice)

	text += fmt.Sprintf("Responsible Team: %s\n", practice.ResponsibleTeam)
	text += fmt.Sprintf("Review Frequency: %s\n\n", practice.ReviewFrequency)

	text += "Required Artifacts:\n"
	for _, artifact := range practice.Artifacts {
		text += fmt.Sprintf("- %s\n", artifact)
	}
	text += "\n"

	if practice.Notes != "" {
		text += fmt.Sprintf("Notes: %s\n", practice.Notes)
	}

	return text
}
