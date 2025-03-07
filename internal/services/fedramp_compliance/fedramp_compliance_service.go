package fedramp_compliance

import (
	"fmt"

	"github.com/grafana/hackathon-12-mcp-compliance/internal/adapters"
	"github.com/grafana/hackathon-12-mcp-compliance/internal/domain/fedramp"
	"github.com/grafana/hackathon-12-mcp-compliance/internal/services/fedramp_compliance/fedramp_compliance_handlers"
)

// Service provides methods for accessing FedRAMP compliance data
type Service struct {
	programHandler *fedramp_compliance_handlers.ProgramHandler
	controlHandler *fedramp_compliance_handlers.ControlHandler
	searchHandler  *fedramp_compliance_handlers.SearchHandler
}

// NewService creates a new FedRAMP compliance service
func NewService() *Service {
	// Create the compliance repository
	complianceRepo := adapters.NewEmbeddedComplianceRepository()

	// Create handlers with the repository
	programHandler := fedramp_compliance_handlers.NewProgramHandler(complianceRepo)
	controlHandler := fedramp_compliance_handlers.NewControlHandler(complianceRepo)
	searchHandler := fedramp_compliance_handlers.NewSearchHandler(complianceRepo)

	return &Service{
		programHandler: programHandler,
		controlHandler: controlHandler,
		searchHandler:  searchHandler,
	}
}

// ListCompliancePrograms returns a list of available compliance programs
func (s *Service) ListCompliancePrograms() ([]string, error) {
	// Create command
	cmd := fedramp.ListComplianceProgramsCommand{}

	// Delegate to program handler
	return s.programHandler.HandleListCompliancePrograms(cmd)
}

// GetControl returns a control by ID
func (s *Service) GetControl(programName, controlID string) (fedramp.Control, bool, error) {
	// Validate arguments
	if programName == "" {
		return fedramp.Control{}, false, fmt.Errorf("program name cannot be empty")
	}
	if controlID == "" {
		return fedramp.Control{}, false, fmt.Errorf("control ID cannot be empty")
	}

	// Load the program
	program, err := s.loadProgram(programName)
	if err != nil {
		return fedramp.Control{}, false, err
	}

	// Create command
	cmd := fedramp.GetControlCommand{
		Program:   program,
		ControlID: controlID,
	}

	// Delegate to control handler
	return s.controlHandler.HandleGetControl(cmd)
}

// GetControlFamily returns a control family by ID
func (s *Service) GetControlFamily(programName, familyID string) (fedramp.ControlFamily, bool, error) {
	// Validate arguments
	if programName == "" {
		return fedramp.ControlFamily{}, false, fmt.Errorf("program name cannot be empty")
	}
	if familyID == "" {
		return fedramp.ControlFamily{}, false, fmt.Errorf("family ID cannot be empty")
	}

	// Load the program
	program, err := s.loadProgram(programName)
	if err != nil {
		return fedramp.ControlFamily{}, false, err
	}

	// Create command
	cmd := fedramp.GetControlFamilyCommand{
		Program:  program,
		FamilyID: familyID,
	}

	// Delegate to control handler
	return s.controlHandler.HandleGetControlFamily(cmd)
}

// ListControlFamilies returns a list of all control families
func (s *Service) ListControlFamilies(programName string) ([]fedramp.ControlFamily, error) {
	// Validate arguments
	if programName == "" {
		return nil, fmt.Errorf("program name cannot be empty")
	}

	// Load the program
	program, err := s.loadProgram(programName)
	if err != nil {
		return nil, err
	}

	// Create command
	cmd := fedramp.ListControlFamiliesCommand{
		Program: program,
	}

	// Delegate to control handler
	return s.controlHandler.HandleListControlFamilies(cmd)
}

// SearchControls searches for controls by keyword
func (s *Service) SearchControls(programName, query string) ([]fedramp.Control, error) {
	// Validate arguments
	if programName == "" {
		return nil, fmt.Errorf("program name cannot be empty")
	}
	if query == "" {
		return []fedramp.Control{}, nil
	}

	// Load the program
	program, err := s.loadProgram(programName)
	if err != nil {
		return nil, err
	}

	// Create command
	cmd := fedramp.SearchControlsCommand{
		Program: program,
		Query:   query,
	}

	// Delegate to search handler
	return s.searchHandler.HandleSearchControls(cmd)
}

// GetControlEvidenceGuidance returns evidence guidance for a control
func (s *Service) GetControlEvidenceGuidance(programName, controlID string) (string, bool, error) {
	// Validate arguments
	if programName == "" {
		return "", false, fmt.Errorf("program name cannot be empty")
	}
	if controlID == "" {
		return "", false, fmt.Errorf("control ID cannot be empty")
	}

	// Load the program
	program, err := s.loadProgram(programName)
	if err != nil {
		return "", false, err
	}

	// Create command
	cmd := fedramp.GetControlEvidenceGuidanceCommand{
		Program:   program,
		ControlID: controlID,
	}

	// Delegate to control handler
	return s.controlHandler.HandleGetControlEvidenceGuidance(cmd)
}

// Helper method to load a program
func (s *Service) loadProgram(programName string) (fedramp.Program, error) {
	cmd := fedramp.GetProgramCommand{
		ProgramName: programName,
	}
	return s.programHandler.HandleGetProgram(cmd)
}
