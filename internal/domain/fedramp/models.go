package fedramp

// ControlParameter represents a parameter for a control
type ControlParameter struct {
	ID         string   `json:"id"`
	Label      string   `json:"label,omitempty"`
	Guidelines []string `json:"guidelines,omitempty"`
}

// ControlStatement represents a statement or requirement in a control
type ControlStatement struct {
	ID       string             `json:"id"`
	Name     string             `json:"name"`
	Prose    string             `json:"prose,omitempty"`
	Parts    []ControlStatement `json:"parts,omitempty"`
	Label    string             `json:"label,omitempty"`
	SubParts []ControlStatement `json:"subParts,omitempty"`
}

// AssessmentMethod represents a method for assessing a control
type AssessmentMethod struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// AssessmentObjective represents an objective for assessing a control
type AssessmentObjective struct {
	ID      string                `json:"id"`
	Name    string                `json:"name"`
	Prose   string                `json:"prose,omitempty"`
	Methods []AssessmentMethod    `json:"methods,omitempty"`
	Parts   []AssessmentObjective `json:"parts,omitempty"`
}

// Control represents a security control
type Control struct {
	ID                   string                `json:"id"`
	Title                string                `json:"title"`
	Parameters           []ControlParameter    `json:"parameters,omitempty"`
	Statements           []ControlStatement    `json:"statements,omitempty"`
	Guidance             string                `json:"guidance,omitempty"`
	AssessmentObjectives []AssessmentObjective `json:"assessmentObjectives,omitempty"`
	FullText             string                `json:"fullText,omitempty"`         // Combined prose text of the control
	EvidenceGuidance     string                `json:"evidenceGuidance,omitempty"` // Guidance for evidence collection
	SearchIndex          string                `json:"-"`                          // Combined text for searching (not included in JSON output)
}

// ControlFamily represents a family of controls
type ControlFamily struct {
	ID       string    `json:"id"`
	Title    string    `json:"title"`
	Controls []Control `json:"controls"`
}

// Program represents a compliance program
type Program struct {
	Name     string          `json:"name"`
	Families []ControlFamily `json:"families"`
}

// OSCALCatalog represents the structure of the OSCAL catalog
type OSCALCatalog struct {
	Catalog struct {
		UUID     string `json:"uuid"`
		Metadata struct {
			Title string `json:"title"`
		} `json:"metadata"`
		Groups []struct {
			ID       string `json:"id"`
			Class    string `json:"class"`
			Title    string `json:"title"`
			Controls []struct {
				ID     string `json:"id"`
				Class  string `json:"class"`
				Title  string `json:"title"`
				Params []struct {
					ID         string `json:"id"`
					Label      string `json:"label,omitempty"`
					Guidelines []struct {
						Prose string `json:"prose,omitempty"`
					} `json:"guidelines,omitempty"`
				} `json:"params,omitempty"`
				Parts []struct {
					ID    string `json:"id"`
					Name  string `json:"name"`
					Prose string `json:"prose,omitempty"`
					Parts []struct {
						ID    string `json:"id"`
						Name  string `json:"name"`
						Prose string `json:"prose,omitempty"`
						Props []struct {
							Name  string `json:"name"`
							Value string `json:"value"`
						} `json:"props,omitempty"`
						Parts []struct {
							ID    string `json:"id"`
							Name  string `json:"name"`
							Prose string `json:"prose,omitempty"`
							Parts []struct {
								ID    string `json:"id"`
								Name  string `json:"name"`
								Prose string `json:"prose,omitempty"`
							} `json:"parts,omitempty"`
						} `json:"parts,omitempty"`
					} `json:"parts,omitempty"`
				} `json:"parts,omitempty"`
			} `json:"controls"`
		} `json:"groups"`
	} `json:"catalog"`
}
