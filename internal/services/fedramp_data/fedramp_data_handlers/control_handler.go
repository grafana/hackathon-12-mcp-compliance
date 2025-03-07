package fedramp_data_handlers

import (
	"strings"

	"github.com/grafana/hackathon-12-mcp-compliance/internal/domain/fedramp"
)

// ControlHandler handles control-related operations
type ControlHandler struct {
	// No ports needed for control operations as they operate on in-memory data
}

// NewControlHandler creates a new control handler
func NewControlHandler() *ControlHandler {
	return &ControlHandler{}
}

// HandleGetControl returns a control by ID
func (h *ControlHandler) HandleGetControl(cmd fedramp.GetControlCommand) (fedramp.Control, bool) {
	controlID := cmd.ControlID

	for _, family := range cmd.Program.Families {
		for _, control := range family.Controls {
			if equalIgnoreCase(control.ID, controlID) {
				return control, true
			}
		}
	}

	return fedramp.Control{}, false
}

// HandleGetControlFamily returns a control family by ID
func (h *ControlHandler) HandleGetControlFamily(cmd fedramp.GetControlFamilyCommand) (fedramp.ControlFamily, bool) {
	familyID := cmd.FamilyID

	for _, family := range cmd.Program.Families {
		if equalIgnoreCase(family.ID, familyID) {
			return family, true
		}
	}

	return fedramp.ControlFamily{}, false
}

// HandleListControlFamilies returns a list of all control families
func (h *ControlHandler) HandleListControlFamilies(cmd fedramp.ListControlFamiliesCommand) []fedramp.ControlFamily {
	return cmd.Program.Families
}

// HandleGetControlEvidenceGuidance returns evidence guidance for a control
func (h *ControlHandler) HandleGetControlEvidenceGuidance(cmd fedramp.GetControlEvidenceGuidanceCommand) (string, bool) {
	getControlCmd := fedramp.GetControlCommand{
		Program:   cmd.Program,
		ControlID: cmd.ControlID,
	}
	control, found := h.HandleGetControl(getControlCmd)
	if !found {
		return "", false
	}

	return control.EvidenceGuidance, true
}

// Helper function to check if two strings are equal, ignoring case
func equalIgnoreCase(s1, s2 string) bool {
	return strings.ToLower(s1) == strings.ToLower(s2)
}
