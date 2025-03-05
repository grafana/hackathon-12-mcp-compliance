package compliance

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"gopkg.in/yaml.v3"
)

var (
	profile     *FedRAMPProfile
	profileOnce sync.Once
	profileErr  error
)

// LoadFedRAMPProfile loads the FedRAMP profile from the specified YAML file
func LoadFedRAMPProfile(filePath string) (*FedRAMPProfile, error) {
	profileOnce.Do(func() {
		data, err := os.ReadFile(filePath)
		if err != nil {
			profileErr = fmt.Errorf("failed to read FedRAMP profile file: %w", err)
			return
		}

		var p FedRAMPProfile
		if err := yaml.Unmarshal(data, &p); err != nil {
			profileErr = fmt.Errorf("failed to parse FedRAMP profile: %w", err)
			return
		}

		profile = &p
	})

	return profile, profileErr
}

// GetControlInfo retrieves information about a specific control
func GetControlInfo(controlID string) (ControlData, error) {
	if profile == nil {
		// Try to load from a default location if not already loaded
		defaultPath := filepath.Join("data", "FedRAMP_rev5_HIGH-baseline_profile.yaml")
		if _, err := LoadFedRAMPProfile(defaultPath); err != nil {
			return ControlData{}, fmt.Errorf("profile not loaded and failed to load from default path: %w", err)
		}
	}

	if control, ok := profile.Profile.Controls[controlID]; ok {
		return control, nil
	}

	return ControlData{}, fmt.Errorf("control %s not found in FedRAMP profile", controlID)
}

// GetEvidenceForControl retrieves the current evidence for a specific control
// In a real implementation, this would likely read from a database or file system
func GetEvidenceForControl(controlID string) (ControlEvidence, error) {
	// TODO: Implement actual evidence retrieval
	// For now, return a placeholder
	return ControlEvidence{
		ControlID:   controlID,
		Status:      EvidenceNotCollected,
		Evidence:    []Evidence{},
		LastUpdated: "",
	}, nil
}

// SaveEvidence saves evidence for a control
// In a real implementation, this would likely save to a database or file system
func SaveEvidence(evidence ControlEvidence) error {
	// TODO: Implement actual evidence saving
	return nil
}

// GenerateSSP generates a System Security Plan based on the collected evidence
func GenerateSSP() (string, error) {
	// TODO: Implement actual SSP generation
	return "This is a placeholder for a generated SSP", nil
}
