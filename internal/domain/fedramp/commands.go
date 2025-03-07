package fedramp

// ProcessFileCommand represents a command to process a FedRAMP baseline file
type ProcessFileCommand struct {
	InputPath   string
	ProgramName string
}

// WriteOutputCommand represents a command to write a Program to a JSON file
type WriteOutputCommand struct {
	Program    Program
	OutputPath string
}

// SearchControlsCommand represents a command to search for controls by keyword
type SearchControlsCommand struct {
	Program Program
	Query   string
}

// GetControlCommand represents a command to get a control by ID
type GetControlCommand struct {
	Program   Program
	ControlID string
}

// GetControlFamilyCommand represents a command to get a control family by ID
type GetControlFamilyCommand struct {
	Program  Program
	FamilyID string
}

// ListControlFamiliesCommand represents a command to list all control families
type ListControlFamiliesCommand struct {
	Program Program
}

// GetControlEvidenceGuidanceCommand represents a command to get evidence guidance for a control
type GetControlEvidenceGuidanceCommand struct {
	Program   Program
	ControlID string
}
