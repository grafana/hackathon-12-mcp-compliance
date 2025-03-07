package fedramp_compliance_handlers

import (
	"fmt"
	"strings"

	"github.com/grafana/hackathon-12-mcp-compliance/internal/domain/fedramp"
	"github.com/grafana/hackathon-12-mcp-compliance/internal/ports"
)

// ControlHandler handles control-related operations
type ControlHandler struct {
	complianceRepo ports.ComplianceRepository
}

// NewControlHandler creates a new control handler
func NewControlHandler(complianceRepo ports.ComplianceRepository) *ControlHandler {
	return &ControlHandler{
		complianceRepo: complianceRepo,
	}
}

// HandleGetControl returns a control by ID
func (h *ControlHandler) HandleGetControl(cmd fedramp.GetControlCommand) (fedramp.Control, bool, error) {
	// Load the program
	program, err := h.complianceRepo.LoadProgram(cmd.Program.Name)
	if err != nil {
		return fedramp.Control{}, false, err
	}

	// Find the control
	for _, family := range program.Families {
		for _, control := range family.Controls {
			if equalIgnoreCase(control.ID, cmd.ControlID) {
				return control, true, nil
			}
		}
	}

	return fedramp.Control{}, false, nil
}

// HandleGetControlFamily returns a control family by ID
func (h *ControlHandler) HandleGetControlFamily(cmd fedramp.GetControlFamilyCommand) (fedramp.ControlFamily, bool, error) {
	// Load the program
	program, err := h.complianceRepo.LoadProgram(cmd.Program.Name)
	if err != nil {
		return fedramp.ControlFamily{}, false, err
	}

	// Find the family
	for _, family := range program.Families {
		if equalIgnoreCase(family.ID, cmd.FamilyID) {
			return family, true, nil
		}
	}

	return fedramp.ControlFamily{}, false, nil
}

// HandleListControlFamilies returns a list of all control families
func (h *ControlHandler) HandleListControlFamilies(cmd fedramp.ListControlFamiliesCommand) ([]fedramp.ControlFamily, error) {
	// Load the program
	program, err := h.complianceRepo.LoadProgram(cmd.Program.Name)
	if err != nil {
		return nil, err
	}

	return program.Families, nil
}

// HandleGetControlEvidenceGuidance returns evidence guidance for a control
func (h *ControlHandler) HandleGetControlEvidenceGuidance(cmd fedramp.GetControlEvidenceGuidanceCommand) (string, bool, error) {
	// Create a GetControlCommand
	getControlCmd := fedramp.GetControlCommand{
		Program:   cmd.Program,
		ControlID: cmd.ControlID,
	}

	// Get the control
	control, found, err := h.HandleGetControl(getControlCmd)
	if err != nil {
		return "", false, err
	}
	if !found {
		return "", false, fmt.Errorf("control not found: %s", cmd.ControlID)
	}

	return control.EvidenceGuidance, true, nil
}

// Helper function to check if two strings are equal, ignoring case
func equalIgnoreCase(s1, s2 string) bool {
	return strings.EqualFold(s1, s2)
}
