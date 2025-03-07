package fedramp

// ListComplianceProgramsCommand is a command to list available compliance programs
type ListComplianceProgramsCommand struct {
	// No parameters needed for this command
}

// GetProgramCommand is a command to get a specific compliance program
type GetProgramCommand struct {
	ProgramName string
}

// Note: The following commands are already defined in commands.go:
// - GetControlCommand
// - GetControlFamilyCommand
// - ListControlFamiliesCommand
// - SearchControlsCommand
// - GetControlEvidenceGuidanceCommand
