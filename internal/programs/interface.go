package programs

// Control represents a security control
type Control struct {
	ID             string
	Title          string
	Description    string
	Guidance       string
	Family         string
	Enhancements   []*Control
	Impact         string // High, Moderate, Low
	AssessmentInfo *AssessmentInfo
}

// AssessmentInfo contains information about how to assess a control
type AssessmentInfo struct {
	Objectives []string
	Methods    []string
}

// ControlFamily represents a group of related controls
type ControlFamily struct {
	ID          string
	Name        string
	Description string
	Controls    []*Control
}

// EvidenceGuidance provides guidance on collecting evidence
type EvidenceGuidance struct {
	ControlID       string
	Description     string
	EvidenceTypes   []string
	CollectionSteps []string
	Examples        []string
	CommonPitfalls  []string
}

// ComplianceProgram defines the interface that all compliance programs must implement
type ComplianceProgram interface {
	// GetName returns the name of the compliance program
	GetName() string

	// GetVersion returns the version of the compliance program
	GetVersion() string

	// GetDescription returns a description of the compliance program
	GetDescription() string

	// GetControl returns a control by ID
	GetControl(controlID string) (*Control, error)

	// GetControlFamilies returns all control families
	GetControlFamilies() ([]*ControlFamily, error)

	// SearchControls searches for controls by keyword
	SearchControls(query string) ([]*Control, error)

	// GetEvidenceGuidance returns evidence guidance for a control
	GetEvidenceGuidance(controlID string) (*EvidenceGuidance, error)
}

// Registry maintains a registry of all available compliance programs
type Registry struct {
	programs map[string]ComplianceProgram
}

// NewRegistry creates a new registry
func NewRegistry() *Registry {
	return &Registry{
		programs: make(map[string]ComplianceProgram),
	}
}

// Register registers a compliance program
func (r *Registry) Register(program ComplianceProgram) {
	r.programs[program.GetName()] = program
}

// GetProgram returns a compliance program by name
func (r *Registry) GetProgram(name string) (ComplianceProgram, bool) {
	program, exists := r.programs[name]
	return program, exists
}

// GetAllPrograms returns all registered compliance programs
func (r *Registry) GetAllPrograms() map[string]ComplianceProgram {
	return r.programs
}

// GetProgramNames returns the names of all registered compliance programs
func (r *Registry) GetProgramNames() []string {
	var names []string
	for name := range r.programs {
		names = append(names, name)
	}
	return names
}
