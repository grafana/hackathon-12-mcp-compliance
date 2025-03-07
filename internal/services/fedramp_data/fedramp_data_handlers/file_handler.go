package fedramp_data_handlers

import (
	"fmt"

	"github.com/grafana/hackathon-12-mcp-compliance/internal/domain/fedramp"
	"github.com/grafana/hackathon-12-mcp-compliance/internal/ports"
)

// FileHandler handles file processing operations
type FileHandler struct {
	fileRepo  ports.FileRepository
	oscalRepo ports.OSCALRepository
}

// NewFileHandler creates a new file handler
func NewFileHandler(fileRepo ports.FileRepository, oscalRepo ports.OSCALRepository) *FileHandler {
	return &FileHandler{
		fileRepo:  fileRepo,
		oscalRepo: oscalRepo,
	}
}

// HandleProcessFile processes a FedRAMP baseline file and returns a Program
func (h *FileHandler) HandleProcessFile(cmd fedramp.ProcessFileCommand) (fedramp.Program, error) {
	// Read the file
	data, err := h.fileRepo.ReadFile(cmd.InputPath)
	if err != nil {
		return fedramp.Program{}, fmt.Errorf("failed to read input file: %v", err)
	}

	// Parse the OSCAL catalog
	catalog, err := h.oscalRepo.ParseOSCALCatalog(data)
	if err != nil {
		return fedramp.Program{}, err
	}

	// Process the catalog into a Program
	return h.oscalRepo.ProcessOSCALCatalog(catalog, cmd.ProgramName)
}

// HandleWriteOutput writes a Program to a JSON file
func (h *FileHandler) HandleWriteOutput(cmd fedramp.WriteOutputCommand) error {
	// Serialize the program to JSON
	data, err := h.oscalRepo.SerializeProgram(cmd.Program)
	if err != nil {
		return err
	}

	// Write the data to the output file
	return h.fileRepo.WriteFile(cmd.OutputPath, data)
}
