package common

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/grafana/mcp-compliance/internal/programs"
	"gopkg.in/yaml.v3"
)

// NISTCatalog represents the NIST SP-800-53 catalog
type NISTCatalog struct {
	Catalog struct {
		UUID     string `yaml:"uuid"`
		Metadata struct {
			Title        string `yaml:"title"`
			LastModified string `yaml:"last-modified"`
			Version      string `yaml:"version"`
			OSCALVersion string `yaml:"oscal-version"`
		} `yaml:"metadata"`
		Groups []struct {
			ID       string        `yaml:"id"`
			Class    string        `yaml:"class"`
			Title    string        `yaml:"title"`
			Controls []NISTControl `yaml:"controls"`
		} `yaml:"groups"`
	} `yaml:"catalog"`
}

// NISTControl represents a control in the NIST catalog
type NISTControl struct {
	ID    string `yaml:"id"`
	Class string `yaml:"class"`
	Title string `yaml:"title"`
	Parts []struct {
		ID    string `yaml:"id"`
		Name  string `yaml:"name"`
		Prose string `yaml:"prose,omitempty"`
		Parts []struct {
			ID    string `yaml:"id"`
			Name  string `yaml:"name"`
			Prose string `yaml:"prose,omitempty"`
		} `yaml:"parts,omitempty"`
	} `yaml:"parts"`
}

// ParseNISTCatalog parses the NIST SP-800-53 catalog from a file
func ParseNISTCatalog(filePath string) (*NISTCatalog, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read NIST catalog: %v", err)
	}

	var catalog NISTCatalog
	if err := yaml.Unmarshal(data, &catalog); err != nil {
		return nil, fmt.Errorf("failed to parse NIST catalog: %v", err)
	}

	return &catalog, nil
}

// ConvertToControlFamily converts a NIST group to a ControlFamily
func ConvertToControlFamily(group struct {
	ID       string        `yaml:"id"`
	Class    string        `yaml:"class"`
	Title    string        `yaml:"title"`
	Controls []NISTControl `yaml:"controls"`
}) *programs.ControlFamily {
	family := &programs.ControlFamily{
		ID:          group.ID,
		Name:        group.Title,
		Description: fmt.Sprintf("%s controls", group.Title),
		Controls:    make([]*programs.Control, 0, len(group.Controls)),
	}

	for _, nistControl := range group.Controls {
		control := ConvertToControl(nistControl)
		if control != nil {
			family.Controls = append(family.Controls, control)
		}
	}

	return family
}

// ConvertToControl converts a NIST control to a Control
func ConvertToControl(nistControl NISTControl) *programs.Control {
	control := &programs.Control{
		ID:           nistControl.ID,
		Title:        nistControl.Title,
		Family:       strings.Split(nistControl.ID, "-")[0],
		Enhancements: make([]*programs.Control, 0),
		AssessmentInfo: &programs.AssessmentInfo{
			Objectives: make([]string, 0),
			Methods:    make([]string, 0),
		},
	}

	// Extract description and guidance from parts
	for _, part := range nistControl.Parts {
		if part.Name == "statement" {
			control.Description = part.Prose
		} else if part.Name == "guidance" {
			control.Guidance = part.Prose
		}
	}

	return control
}

// FindControlByID finds a control by ID in the NIST catalog
func FindControlByID(catalog *NISTCatalog, controlID string) *NISTControl {
	controlID = strings.ToLower(controlID)

	for _, group := range catalog.Catalog.Groups {
		for _, control := range group.Controls {
			if strings.ToLower(control.ID) == controlID {
				return &control
			}

			// Check for enhancements
			if strings.HasPrefix(controlID, strings.ToLower(control.ID)+".") {
				// This is an enhancement, but we need to handle it differently
				// For now, return nil as we'll need to implement enhancement handling
				return nil
			}
		}
	}

	return nil
}
