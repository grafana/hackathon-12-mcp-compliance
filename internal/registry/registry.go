package registry

import (
	"fmt"

	"github.com/grafana/hackathon-12-mcp-compliance/internal/programs"
)

// Registry manages the set of available compliance programs
type Registry struct {
	programs map[string]programs.Program
}

// New creates a new program registry
func New() *Registry {
	return &Registry{
		programs: make(map[string]programs.Program),
	}
}

// Register adds a program to the registry
func (r *Registry) Register(p programs.Program) {
	r.programs[p.Name()] = p
}

// GetProgram retrieves a program by name
func (r *Registry) GetProgram(name string) (programs.Program, error) {
	p, exists := r.programs[name]
	if !exists {
		return nil, fmt.Errorf("program not found: %s", name)
	}
	return p, nil
}

// ListPrograms returns a list of all registered program names
func (r *Registry) ListPrograms() []string {
	names := make([]string, 0, len(r.programs))
	for name := range r.programs {
		names = append(names, name)
	}
	return names
}
