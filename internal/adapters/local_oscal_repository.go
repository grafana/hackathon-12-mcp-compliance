package adapters

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/grafana/hackathon-12-mcp-compliance/internal/domain/fedramp"
	"github.com/grafana/hackathon-12-mcp-compliance/internal/ports"
)

// LocalOSCALRepository implements the OSCALRepository interface for local OSCAL processing
type LocalOSCALRepository struct{}

// NewLocalOSCALRepository creates a new LocalOSCALRepository
func NewLocalOSCALRepository() *LocalOSCALRepository {
	return &LocalOSCALRepository{}
}

// ParseOSCALCatalog parses OSCAL catalog data and returns a structured representation
func (r *LocalOSCALRepository) ParseOSCALCatalog(data []byte) (fedramp.OSCALCatalog, error) {
	var oscalCatalog fedramp.OSCALCatalog
	if err := json.Unmarshal(data, &oscalCatalog); err != nil {
		return fedramp.OSCALCatalog{}, fmt.Errorf("failed to parse OSCAL catalog: %v", err)
	}
	return oscalCatalog, nil
}

// ProcessOSCALCatalog processes an OSCAL catalog into a Program
func (r *LocalOSCALRepository) ProcessOSCALCatalog(catalog fedramp.OSCALCatalog, programName string) (fedramp.Program, error) {
	// Create our simplified program structure
	program := fedramp.Program{
		Name:     programName,
		Families: []fedramp.ControlFamily{},
	}

	// Process each control family
	for _, group := range catalog.Catalog.Groups {
		if group.Class == "family" {
			family := fedramp.ControlFamily{
				ID:       group.ID,
				Title:    group.Title,
				Controls: []fedramp.Control{},
			}

			// Process each control in the family
			for _, oscalControl := range group.Controls {
				control := fedramp.Control{
					ID:                   oscalControl.ID,
					Title:                oscalControl.Title,
					Parameters:           []fedramp.ControlParameter{},
					Statements:           []fedramp.ControlStatement{},
					AssessmentObjectives: []fedramp.AssessmentObjective{},
				}

				// Extract parameters
				for _, param := range oscalControl.Params {
					parameter := fedramp.ControlParameter{
						ID:    param.ID,
						Label: param.Label,
					}

					// Extract guidelines
					for _, guideline := range param.Guidelines {
						if guideline.Prose != "" {
							parameter.Guidelines = append(parameter.Guidelines, guideline.Prose)
						}
					}

					control.Parameters = append(control.Parameters, parameter)
				}

				// Extract statements, guidance, and assessment objectives
				var statementText strings.Builder
				var evidenceGuidanceBuilder strings.Builder

				for _, part := range oscalControl.Parts {
					if part.Name == "statement" {
						statement := r.extractStatement(part)
						control.Statements = append(control.Statements, statement)

						// Build the full statement text
						if part.Prose != "" {
							statementText.WriteString(part.Prose)
							statementText.WriteString("\n\n")
						}

						// Add sub-parts prose to the full text
						for _, subPart := range part.Parts {
							if subPart.Prose != "" {
								statementText.WriteString(subPart.Prose)
								statementText.WriteString("\n")
							}

							// Add deeper nested parts
							for _, subSubPart := range subPart.Parts {
								if subSubPart.Prose != "" {
									statementText.WriteString("  " + subSubPart.Prose)
									statementText.WriteString("\n")
								}
							}
						}
					} else if part.Name == "guidance" {
						control.Guidance = part.Prose

						// Check if guidance contains evidence-related information
						if strings.Contains(strings.ToLower(part.Prose), "evidence") ||
							strings.Contains(strings.ToLower(part.Prose), "assess") ||
							strings.Contains(strings.ToLower(part.Prose), "audit") ||
							strings.Contains(strings.ToLower(part.Prose), "document") {
							evidenceGuidanceBuilder.WriteString("Guidance related to evidence:\n")
							evidenceGuidanceBuilder.WriteString(part.Prose)
							evidenceGuidanceBuilder.WriteString("\n\n")
						}
					} else if part.Name == "assessment-objective" {
						objective := r.extractAssessmentObjective(part)
						control.AssessmentObjectives = append(control.AssessmentObjectives, objective)

						// Add assessment objectives to evidence guidance
						evidenceGuidanceBuilder.WriteString("Assessment Objective:\n")
						evidenceGuidanceBuilder.WriteString(part.Prose)
						evidenceGuidanceBuilder.WriteString("\n")

						// Add assessment methods
						for _, subPart := range part.Parts {
							for _, prop := range subPart.Props {
								if prop.Name == "method" {
									evidenceGuidanceBuilder.WriteString("Assessment Method: ")
									evidenceGuidanceBuilder.WriteString(prop.Value)
									evidenceGuidanceBuilder.WriteString("\n")
								}
							}

							if subPart.Prose != "" {
								evidenceGuidanceBuilder.WriteString(subPart.Prose)
								evidenceGuidanceBuilder.WriteString("\n")
							}
						}
					}
				}

				// Set the full text of the control
				control.FullText = statementText.String()

				// Set the evidence guidance
				control.EvidenceGuidance = evidenceGuidanceBuilder.String()

				// Create a search index by combining all text fields
				var searchIndexBuilder strings.Builder
				searchIndexBuilder.WriteString(control.ID)
				searchIndexBuilder.WriteString(" ")
				searchIndexBuilder.WriteString(control.Title)
				searchIndexBuilder.WriteString(" ")
				searchIndexBuilder.WriteString(control.FullText)
				searchIndexBuilder.WriteString(" ")
				searchIndexBuilder.WriteString(control.Guidance)
				searchIndexBuilder.WriteString(" ")
				searchIndexBuilder.WriteString(control.EvidenceGuidance)

				// Add parameter labels and guidelines to search index
				for _, param := range control.Parameters {
					searchIndexBuilder.WriteString(" ")
					searchIndexBuilder.WriteString(param.Label)
					for _, guideline := range param.Guidelines {
						searchIndexBuilder.WriteString(" ")
						searchIndexBuilder.WriteString(guideline)
					}
				}

				control.SearchIndex = strings.ToLower(searchIndexBuilder.String())

				family.Controls = append(family.Controls, control)
			}

			program.Families = append(program.Families, family)
		}
	}

	return program, nil
}

// SerializeProgram serializes a Program to JSON
func (r *LocalOSCALRepository) SerializeProgram(program fedramp.Program) ([]byte, error) {
	data, err := json.MarshalIndent(program, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to serialize program: %v", err)
	}
	return data, nil
}

// Recursively extract control statements
func (r *LocalOSCALRepository) extractStatement(part struct {
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
}) fedramp.ControlStatement {
	statement := fedramp.ControlStatement{
		ID:    part.ID,
		Name:  part.Name,
		Prose: part.Prose,
	}

	// Extract label from props if available
	for _, subPart := range part.Parts {
		for _, prop := range subPart.Props {
			if prop.Name == "label" {
				statement.Label = prop.Value
				break
			}
		}

		// Recursively process sub-parts
		if len(subPart.Parts) > 0 {
			for _, subSubPart := range subPart.Parts {
				subStatement := fedramp.ControlStatement{
					ID:    subSubPart.ID,
					Name:  subSubPart.Name,
					Prose: subSubPart.Prose,
				}
				statement.SubParts = append(statement.SubParts, subStatement)
			}
		}

		subStatement := fedramp.ControlStatement{
			ID:    subPart.ID,
			Name:  subPart.Name,
			Prose: subPart.Prose,
		}
		statement.Parts = append(statement.Parts, subStatement)
	}

	return statement
}

// Extract assessment objectives
func (r *LocalOSCALRepository) extractAssessmentObjective(part struct {
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
}) fedramp.AssessmentObjective {
	objective := fedramp.AssessmentObjective{
		ID:    part.ID,
		Name:  part.Name,
		Prose: part.Prose,
	}

	// Extract methods from props if available
	for _, subPart := range part.Parts {
		for _, prop := range subPart.Props {
			if prop.Name == "method" {
				method := fedramp.AssessmentMethod{
					Name:  "method",
					Value: prop.Value,
				}
				objective.Methods = append(objective.Methods, method)
			}
		}

		// Recursively process sub-parts
		if len(subPart.Parts) > 0 {
			for _, subSubPart := range subPart.Parts {
				subObjective := fedramp.AssessmentObjective{
					ID:    subSubPart.ID,
					Name:  subSubPart.Name,
					Prose: subSubPart.Prose,
				}
				objective.Parts = append(objective.Parts, subObjective)
			}
		}
	}

	return objective
}

// Ensure LocalOSCALRepository implements OSCALRepository
var _ ports.OSCALRepository = (*LocalOSCALRepository)(nil)
