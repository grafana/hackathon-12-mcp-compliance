package fedramp_data

import (
	"fmt"

	"github.com/grafana/hackathon-12-mcp-compliance/internal/adapters"
	"github.com/grafana/hackathon-12-mcp-compliance/internal/domain/fedramp"
	"github.com/grafana/hackathon-12-mcp-compliance/internal/services/fedramp_data/fedramp_data_handlers"
)

// Service provides methods for processing FedRAMP data
type Service struct {
	fileHandler    *fedramp_data_handlers.FileHandler
	searchHandler  *fedramp_data_handlers.SearchHandler
	controlHandler *fedramp_data_handlers.ControlHandler
}

// NewService creates a new FedRAMP service
func NewService() *Service {
	// Create adapters
	fileRepo := adapters.NewLocalFileRepository()
	oscalRepo := adapters.NewLocalOSCALRepository()

	// Create handlers with appropriate adapters
	fileHandler := fedramp_data_handlers.NewFileHandler(fileRepo, oscalRepo)
	searchHandler := fedramp_data_handlers.NewSearchHandler()
	controlHandler := fedramp_data_handlers.NewControlHandler()

	return &Service{
		fileHandler:    fileHandler,
		searchHandler:  searchHandler,
		controlHandler: controlHandler,
	}
}

// ProcessFile processes a FedRAMP baseline file and returns a Program
func (s *Service) ProcessFile(inputPath, programName string) (fedramp.Program, error) {
	// Validate arguments
	if inputPath == "" {
		return fedramp.Program{}, fmt.Errorf("input path cannot be empty")
	}
	if programName == "" {
		return fedramp.Program{}, fmt.Errorf("program name cannot be empty")
	}

	// Create command
	cmd := fedramp.ProcessFileCommand{
		InputPath:   inputPath,
		ProgramName: programName,
	}

	// Delegate to file handler
	return s.fileHandler.HandleProcessFile(cmd)
}

// WriteOutput writes a Program to a JSON file
func (s *Service) WriteOutput(program fedramp.Program, outputPath string) error {
	// Validate arguments
	if outputPath == "" {
		return fmt.Errorf("output path cannot be empty")
	}

	// Create command
	cmd := fedramp.WriteOutputCommand{
		Program:    program,
		OutputPath: outputPath,
	}

	// Delegate to file handler
	return s.fileHandler.HandleWriteOutput(cmd)
}

// SearchControls searches for controls by keyword
func (s *Service) SearchControls(program fedramp.Program, query string) []fedramp.Control {
	// Validate arguments
	if query == "" {
		return []fedramp.Control{}
	}

	// Create command
	cmd := fedramp.SearchControlsCommand{
		Program: program,
		Query:   query,
	}

	// Delegate to search handler
	return s.searchHandler.HandleSearchControls(cmd)
}

// GetControl returns a control by ID
func (s *Service) GetControl(program fedramp.Program, controlID string) (fedramp.Control, bool) {
	// Validate arguments
	if controlID == "" {
		return fedramp.Control{}, false
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
func (s *Service) GetControlFamily(program fedramp.Program, familyID string) (fedramp.ControlFamily, bool) {
	// Validate arguments
	if familyID == "" {
		return fedramp.ControlFamily{}, false
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
func (s *Service) ListControlFamilies(program fedramp.Program) []fedramp.ControlFamily {
	// Create command
	cmd := fedramp.ListControlFamiliesCommand{
		Program: program,
	}

	// Delegate to control handler
	return s.controlHandler.HandleListControlFamilies(cmd)
}

// GetControlEvidenceGuidance returns evidence guidance for a control
func (s *Service) GetControlEvidenceGuidance(program fedramp.Program, controlID string) (string, bool) {
	// Validate arguments
	if controlID == "" {
		return "", false
	}

	// Create command
	cmd := fedramp.GetControlEvidenceGuidanceCommand{
		Program:   program,
		ControlID: controlID,
	}

	// Delegate to control handler
	return s.controlHandler.HandleGetControlEvidenceGuidance(cmd)
}
