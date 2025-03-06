package controls

// Control represents a security control
type Control struct {
	ID             string
	Title          string
	Description    string
	Guidance       string
	Family         string
	Enhancements   []*Control
	FedRAMPImpact  string // High, Moderate, Low
	AssessmentInfo *AssessmentInfo
}

// AssessmentInfo contains information about how to assess a control
type AssessmentInfo struct {
	Objectives []string
	Methods    []string
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

// ComplianceFramework represents a compliance framework like FedRAMP
type ComplianceFramework struct {
	Name             string
	Version          string
	Controls         map[string]*Control // Indexed by ID
	ControlsByFamily map[string][]*Control
	SearchIndex      map[string][]string // Simple keyword index
}

// SecurityControl represents a FedRAMP security control
type SecurityControl struct {
	ID          string `json:"id"`
	Family      string `json:"family"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Impact      string `json:"impact"` // Low, Moderate, High
	Details     string `json:"details"`
}

// SecurityControlFamily represents a group of related security controls
type SecurityControlFamily struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Controls    []SecurityControl `json:"controls"`
}
