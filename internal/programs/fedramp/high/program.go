package high

import (
	"embed"
	"fmt"
	"path/filepath"
	"strings"
	"sync"

	"github.com/grafana/mcp-compliance/internal/programs"
	"github.com/grafana/mcp-compliance/internal/programs/fedramp/common"
	"gopkg.in/yaml.v3"
)

//go:embed data/*.yaml
var dataFiles embed.FS

// FedRAMPHighProgram implements the ComplianceProgram interface for FedRAMP High
type FedRAMPHighProgram struct {
	name        string
	version     string
	description string

	nistCatalog     *common.NISTCatalog
	controlFamilies []*programs.ControlFamily
	controlsMap     map[string]*programs.Control

	// Evidence guidance
	evidenceGuidance map[string]*programs.EvidenceGuidance

	// Search index
	searchIndex map[string][]string

	// Mutex for thread safety
	mu sync.RWMutex
}

// FedRAMPProfile represents the FedRAMP High baseline profile
type FedRAMPProfile struct {
	Profile struct {
		UUID     string `yaml:"uuid"`
		Metadata struct {
			Title        string `yaml:"title"`
			Published    string `yaml:"published"`
			LastModified string `yaml:"last-modified"`
			Version      string `yaml:"version"`
			OSCALVersion string `yaml:"oscal-version"`
		} `yaml:"metadata"`
		Imports []struct {
			Href string `yaml:"href"`
		} `yaml:"imports"`
		Modify struct {
			Alterations []struct {
				ControlID string `yaml:"control-id"`
			} `yaml:"alterations"`
		} `yaml:"modify"`
	} `yaml:"profile"`
}

// NewFedRAMPHighProgram creates a new FedRAMP High program
func NewFedRAMPHighProgram() (*FedRAMPHighProgram, error) {
	program := &FedRAMPHighProgram{
		name:             "FedRAMP High",
		version:          "Rev 5",
		description:      "Federal Risk and Authorization Management Program (FedRAMP) High Baseline",
		controlsMap:      make(map[string]*programs.Control),
		evidenceGuidance: make(map[string]*programs.EvidenceGuidance),
		searchIndex:      make(map[string][]string),
	}

	// Load NIST catalog
	catalogPath := filepath.Join("..", "common", "NIST_SP-800-53_rev5_catalog.yaml")
	nistCatalog, err := common.ParseNISTCatalog(catalogPath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse NIST catalog: %v", err)
	}
	program.nistCatalog = nistCatalog

	// Load FedRAMP High profile
	profileData, err := dataFiles.ReadFile("data/FedRAMP_rev5_HIGH-baseline_profile.yaml")
	if err != nil {
		return nil, fmt.Errorf("failed to read FedRAMP High profile: %v", err)
	}

	var profile FedRAMPProfile
	if err := yaml.Unmarshal(profileData, &profile); err != nil {
		return nil, fmt.Errorf("failed to parse FedRAMP High profile: %v", err)
	}

	// Extract control IDs from the profile
	controlIDs := make(map[string]bool)
	for _, alteration := range profile.Profile.Modify.Alterations {
		controlIDs[strings.ToLower(alteration.ControlID)] = true
	}

	// Build control families and controls
	program.buildControlFamilies(nistCatalog, controlIDs)

	// Build evidence guidance
	program.buildEvidenceGuidance()

	// Build search index
	program.buildSearchIndex()

	return program, nil
}

// buildControlFamilies builds the control families and controls from the NIST catalog
func (p *FedRAMPHighProgram) buildControlFamilies(nistCatalog *common.NISTCatalog, controlIDs map[string]bool) {
	p.controlFamilies = make([]*programs.ControlFamily, 0, len(nistCatalog.Catalog.Groups))

	for _, group := range nistCatalog.Catalog.Groups {
		family := common.ConvertToControlFamily(group)

		// Filter controls to only include those in the FedRAMP High baseline
		filteredControls := make([]*programs.Control, 0)
		for _, control := range family.Controls {
			if controlIDs[strings.ToLower(control.ID)] {
				control.Impact = "High" // Set impact level
				filteredControls = append(filteredControls, control)
				p.controlsMap[strings.ToLower(control.ID)] = control
			}
		}

		family.Controls = filteredControls
		if len(filteredControls) > 0 {
			p.controlFamilies = append(p.controlFamilies, family)
		}
	}
}

// buildEvidenceGuidance builds the evidence guidance for controls
func (p *FedRAMPHighProgram) buildEvidenceGuidance() {
	// For now, we'll create some sample evidence guidance
	// In a real implementation, this would be loaded from a data source

	// AC-1 evidence guidance
	p.evidenceGuidance["ac-1"] = &programs.EvidenceGuidance{
		ControlID:   "AC-1",
		Description: "Guidance for collecting evidence for Access Control Policy and Procedures",
		EvidenceTypes: []string{
			"Policy Documents",
			"Procedure Documents",
			"Implementation Plans",
			"Security Assessment Reports",
		},
		CollectionSteps: []string{
			"Identify and collect the organization's access control policy",
			"Identify and collect the organization's access control procedures",
			"Collect evidence of regular reviews and updates to the policy and procedures",
			"Collect evidence of dissemination to relevant personnel",
		},
		Examples: []string{
			"Access Control Policy document",
			"Access Control Procedures manual",
			"Email notifications of policy updates",
			"Training records for access control procedures",
		},
		CommonPitfalls: []string{
			"Outdated policy documents",
			"Missing approval signatures",
			"Lack of evidence for regular reviews",
			"Incomplete procedure documentation",
		},
	}

	// Add more evidence guidance as needed
}

// buildSearchIndex builds the search index for controls
func (p *FedRAMPHighProgram) buildSearchIndex() {
	for id, control := range p.controlsMap {
		// Index by ID
		p.addToSearchIndex(id, id)

		// Index by title words
		for _, word := range strings.Fields(strings.ToLower(control.Title)) {
			if len(word) > 3 { // Only index words longer than 3 characters
				p.addToSearchIndex(word, id)
			}
		}

		// Index by description words
		for _, word := range strings.Fields(strings.ToLower(control.Description)) {
			if len(word) > 3 { // Only index words longer than 3 characters
				p.addToSearchIndex(word, id)
			}
		}
	}
}

// addToSearchIndex adds a control ID to the search index for a keyword
func (p *FedRAMPHighProgram) addToSearchIndex(keyword, controlID string) {
	if _, exists := p.searchIndex[keyword]; !exists {
		p.searchIndex[keyword] = make([]string, 0)
	}

	// Check if the control ID is already in the list
	for _, id := range p.searchIndex[keyword] {
		if id == controlID {
			return
		}
	}

	p.searchIndex[keyword] = append(p.searchIndex[keyword], controlID)
}

// GetName returns the name of the compliance program
func (p *FedRAMPHighProgram) GetName() string {
	return p.name
}

// GetVersion returns the version of the compliance program
func (p *FedRAMPHighProgram) GetVersion() string {
	return p.version
}

// GetDescription returns a description of the compliance program
func (p *FedRAMPHighProgram) GetDescription() string {
	return p.description
}

// GetControl returns a control by ID
func (p *FedRAMPHighProgram) GetControl(controlID string) (*programs.Control, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	control, exists := p.controlsMap[strings.ToLower(controlID)]
	if !exists {
		return nil, fmt.Errorf("control not found: %s", controlID)
	}

	return control, nil
}

// GetControlFamilies returns all control families
func (p *FedRAMPHighProgram) GetControlFamilies() ([]*programs.ControlFamily, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.controlFamilies, nil
}

// SearchControls searches for controls by keyword
func (p *FedRAMPHighProgram) SearchControls(query string) ([]*programs.Control, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	query = strings.ToLower(query)
	var results []*programs.Control

	// Check if the query is in the search index
	if controlIDs, exists := p.searchIndex[query]; exists {
		for _, id := range controlIDs {
			results = append(results, p.controlsMap[id])
		}
		return results, nil
	}

	// If not in the index, do a more comprehensive search
	for _, control := range p.controlsMap {
		if strings.Contains(strings.ToLower(control.ID), query) ||
			strings.Contains(strings.ToLower(control.Title), query) ||
			strings.Contains(strings.ToLower(control.Description), query) {
			results = append(results, control)
		}
	}

	return results, nil
}

// GetEvidenceGuidance returns evidence guidance for a control
func (p *FedRAMPHighProgram) GetEvidenceGuidance(controlID string) (*programs.EvidenceGuidance, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	guidance, exists := p.evidenceGuidance[strings.ToLower(controlID)]
	if !exists {
		return nil, fmt.Errorf("evidence guidance not found for control: %s", controlID)
	}

	return guidance, nil
}
