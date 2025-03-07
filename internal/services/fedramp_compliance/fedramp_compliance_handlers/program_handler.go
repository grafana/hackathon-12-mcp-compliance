package fedramp_compliance_handlers

import (
	"github.com/grafana/hackathon-12-mcp-compliance/internal/domain/fedramp"
	"github.com/grafana/hackathon-12-mcp-compliance/internal/ports"
)

// ProgramHandler handles program-related operations
type ProgramHandler struct {
	complianceRepo ports.ComplianceRepository
}

// NewProgramHandler creates a new program handler
func NewProgramHandler(complianceRepo ports.ComplianceRepository) *ProgramHandler {
	return &ProgramHandler{
		complianceRepo: complianceRepo,
	}
}

// HandleListCompliancePrograms lists available compliance programs
func (h *ProgramHandler) HandleListCompliancePrograms(cmd fedramp.ListComplianceProgramsCommand) ([]string, error) {
	return h.complianceRepo.ListPrograms()
}

// HandleGetProgram loads a specific compliance program
func (h *ProgramHandler) HandleGetProgram(cmd fedramp.GetProgramCommand) (fedramp.Program, error) {
	return h.complianceRepo.LoadProgram(cmd.ProgramName)
}
