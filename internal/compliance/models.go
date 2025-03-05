package compliance

// FedRAMPProfile represents the structure of a FedRAMP profile
type FedRAMPProfile struct {
	Profile ProfileMetadata `yaml:"profile"`
}

// ProfileMetadata contains metadata about the FedRAMP profile
type ProfileMetadata struct {
	Title       string                 `yaml:"title"`
	Description string                 `yaml:"description"`
	Controls    map[string]ControlData `yaml:"controls"`
}

// ControlData represents a security control in the FedRAMP profile
type ControlData struct {
	ID          string                 `yaml:"id"`
	Title       string                 `yaml:"title"`
	Description string                 `yaml:"description"`
	Parameters  map[string]Parameter   `yaml:"parameters,omitempty"`
	Parts       map[string]ControlPart `yaml:"parts,omitempty"`
}

// Parameter represents a parameter for a security control
type Parameter struct {
	ID          string   `yaml:"id"`
	Label       string   `yaml:"label"`
	Description string   `yaml:"description"`
	Values      []string `yaml:"values,omitempty"`
}

// ControlPart represents a part of a security control
type ControlPart struct {
	ID    string `yaml:"id"`
	Name  string `yaml:"name"`
	Prose string `yaml:"prose"`
}

// EvidenceStatus represents the status of evidence collection for a control
type EvidenceStatus string

const (
	EvidenceNotCollected EvidenceStatus = "not_collected"
	EvidencePartial      EvidenceStatus = "partial"
	EvidenceComplete     EvidenceStatus = "complete"
)

// ControlEvidence represents evidence collected for a security control
type ControlEvidence struct {
	ControlID   string         `json:"controlId"`
	Status      EvidenceStatus `json:"status"`
	Evidence    []Evidence     `json:"evidence"`
	LastUpdated string         `json:"lastUpdated"`
}

// Evidence represents a single piece of evidence
type Evidence struct {
	Type          string `json:"type"` // document, screenshot, log, etc.
	Description   string `json:"description"`
	Location      string `json:"location"` // file path, URL, etc.
	DateCollected string `json:"dateCollected"`
}
