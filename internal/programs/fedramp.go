package programs

import (
	"embed"
	"fmt"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

//go:embed data/fedramp/high.yaml data/fedramp/moderate.yaml
var fedrampFS embed.FS

// FedRAMPProgram implements a FedRAMP program (High or Moderate)
type FedRAMPProgram struct {
	name        string
	description string
	version     string
	controls    map[string]*Control
	families    map[string]*ControlFamily
}

// NewFedRAMPHigh creates a new FedRAMP High program
func NewFedRAMPHigh() (*FedRAMPProgram, error) {
	return newFedRAMPProgram("FedRAMP High", "Federal Risk and Authorization Management Program (FedRAMP) High Security Controls Baseline", "high.yaml")
}

// NewFedRAMPModerate creates a new FedRAMP Moderate program
func NewFedRAMPModerate() (*FedRAMPProgram, error) {
	return newFedRAMPProgram("FedRAMP Moderate", "Federal Risk and Authorization Management Program (FedRAMP) Moderate Security Controls Baseline", "moderate.yaml")
}

// newFedRAMPProgram creates a new FedRAMP program with the specified parameters
func newFedRAMPProgram(name, description, yamlFile string) (*FedRAMPProgram, error) {
	p := &FedRAMPProgram{
		name:        name,
		description: description,
		controls:    make(map[string]*Control),
		families:    make(map[string]*ControlFamily),
	}

	// Load and parse the catalog
	data, err := fedrampFS.ReadFile(fmt.Sprintf("data/fedramp/%s", yamlFile))
	if err != nil {
		return nil, fmt.Errorf("failed to read catalog file: %w", err)
	}
	fmt.Printf("Successfully read %s (%d bytes)\n", yamlFile, len(data))

	var catalog struct {
		Catalog struct {
			Metadata struct {
				Version string `yaml:"version"`
			} `yaml:"metadata"`
			Groups []struct {
				ID       string `yaml:"id"`
				Title    string `yaml:"title"`
				Controls []struct {
					ID    string `yaml:"id"`
					Title string `yaml:"title"`
					Parts []struct {
						ID    string `yaml:"id"`
						Name  string `yaml:"name"`
						Prose string `yaml:"prose"`
					} `yaml:"parts"`
				} `yaml:"controls"`
			} `yaml:"groups"`
		} `yaml:"catalog"`
	}

	if err := yaml.Unmarshal(data, &catalog); err != nil {
		return nil, fmt.Errorf("failed to parse catalog: %w", err)
	}
	fmt.Printf("Successfully parsed catalog with %d groups\n", len(catalog.Catalog.Groups))

	p.version = catalog.Catalog.Metadata.Version

	// Process each group (family) and its controls
	for _, group := range catalog.Catalog.Groups {
		fmt.Printf("Processing group %s with %d controls\n", group.ID, len(group.Controls))
		family := &ControlFamily{
			ID:          group.ID,
			Name:        group.Title,
			Description: fmt.Sprintf("%s controls", group.Title),
			Controls:    make([]*Control, 0),
		}
		p.families[group.ID] = family

		for _, c := range group.Controls {
			// Convert ID to uppercase for consistency
			controlID := strings.ToUpper(c.ID)

			control := &Control{
				ID:           controlID,
				Title:        c.Title,
				Family:       strings.Split(controlID, "-")[0],
				Enhancements: make([]*Control, 0),
				AssessmentInfo: &AssessmentInfo{
					Objectives: make([]string, 0),
					Methods:    make([]string, 0),
				},
			}

			// Extract description and guidance from parts
			for _, part := range c.Parts {
				if part.Name == "statement" {
					control.Description = part.Prose
				} else if part.Name == "guidance" {
					control.Guidance = part.Prose
				}
			}

			p.controls[controlID] = control
			family.Controls = append(family.Controls, control)
		}
	}

	fmt.Printf("Loaded %d controls across %d families\n", len(p.controls), len(p.families))
	return p, nil
}

// Name returns the program name
func (p *FedRAMPProgram) Name() string {
	return p.name
}

// Description returns the program description
func (p *FedRAMPProgram) Description() string {
	return p.description
}

// GetControl returns information about a specific control
func (p *FedRAMPProgram) GetControl(id string) (map[string]interface{}, error) {
	// Convert ID to lowercase for case-insensitive lookup
	lowerID := strings.ToLower(id)

	// Try to find the control
	var control *Control
	var found bool
	for controlID, c := range p.controls {
		if strings.ToLower(controlID) == lowerID {
			control = c
			found = true
			break
		}
	}

	if !found {
		return nil, fmt.Errorf("control %s not found", id)
	}

	return map[string]interface{}{
		"id":          control.ID,
		"title":       control.Title,
		"description": control.Description,
		"guidance":    control.Guidance,
		"family":      control.Family,
	}, nil
}

// GetControlFamily returns all controls in a specific family
func (p *FedRAMPProgram) GetControlFamily(family string) ([]map[string]interface{}, error) {
	controls := make([]map[string]interface{}, 0)
	for _, control := range p.controls {
		if control.Family == family {
			controls = append(controls, map[string]interface{}{
				"id":          control.ID,
				"title":       control.Title,
				"description": control.Description,
				"guidance":    control.Guidance,
				"family":      control.Family,
			})
		}
	}
	if len(controls) == 0 {
		return nil, fmt.Errorf("family %s not found", family)
	}
	return controls, nil
}

// ListControlFamilies returns a list of all control families
func (p *FedRAMPProgram) ListControlFamilies() ([]string, error) {
	families := make(map[string]struct{})
	for _, control := range p.controls {
		families[control.Family] = struct{}{}
	}

	result := make([]string, 0, len(families))
	for family := range families {
		result = append(result, family)
	}
	sort.Strings(result)
	return result, nil
}

// SearchControls searches for controls matching the query
func (p *FedRAMPProgram) SearchControls(query string) ([]map[string]interface{}, error) {
	query = strings.ToLower(query)
	var matches []map[string]interface{}
	for _, control := range p.controls {
		if strings.Contains(strings.ToLower(control.Title), query) ||
			strings.Contains(strings.ToLower(control.Description), query) {
			matches = append(matches, map[string]interface{}{
				"id":          control.ID,
				"title":       control.Title,
				"description": control.Description,
				"guidance":    control.Guidance,
				"family":      control.Family,
			})
		}
	}
	return matches, nil
}

// GetControlsByStatus returns all controls with the given status
func (p *FedRAMPProgram) GetControlsByStatus(status string) ([]map[string]interface{}, error) {
	// For now, return empty list since we don't track status yet
	return []map[string]interface{}{}, nil
}

// GetControlParameters returns parameters for a specific control
func (p *FedRAMPProgram) GetControlParameters(id string) (map[string]interface{}, error) {
	_, ok := p.controls[id]
	if !ok {
		return nil, fmt.Errorf("control %s not found", id)
	}

	// For now, return empty map since we haven't parsed parameters yet
	return map[string]interface{}{}, nil
}

// GetVersion returns the version of the compliance program
func (p *FedRAMPProgram) GetVersion() string {
	return p.version
}

// GetControlFamilies returns all control families
func (p *FedRAMPProgram) GetControlFamilies() ([]*ControlFamily, error) {
	families := make([]*ControlFamily, 0, len(p.families))
	for _, family := range p.families {
		families = append(families, family)
	}
	sort.Slice(families, func(i, j int) bool {
		return families[i].ID < families[j].ID
	})
	return families, nil
}

// GetEvidenceGuidance returns evidence guidance for a control
func (p *FedRAMPProgram) GetEvidenceGuidance(controlID string) (*EvidenceGuidance, error) {
	control, ok := p.controls[controlID]
	if !ok {
		return nil, fmt.Errorf("control %s not found", controlID)
	}

	// For now, return a basic guidance structure
	return &EvidenceGuidance{
		ControlID:   controlID,
		Description: fmt.Sprintf("Evidence guidance for %s: %s", controlID, control.Title),
	}, nil
}
