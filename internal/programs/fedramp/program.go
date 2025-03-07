package fedramp

import (
	"embed"
	"fmt"
	"io"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

//go:embed high/data/*.yaml
var highFS embed.FS

//go:embed moderate/data/*.yaml
var moderateFS embed.FS

// Program implements a FedRAMP program (High or Moderate)
type Program struct {
	name        string
	description string
	fs          embed.FS
	dataFile    string
	catalog     map[string]interface{}
}

// NewHighProgram creates a new FedRAMP High program
func NewHighProgram() (*Program, error) {
	p := &Program{
		name:        "FedRAMP High",
		description: "Federal Risk and Authorization Management Program (FedRAMP) High Security Controls Baseline",
		fs:          highFS,
		dataFile:    "high/data/FedRAMP_rev5_HIGH-baseline-resolved-profile_catalog.yaml",
	}
	if err := p.loadCatalog(); err != nil {
		return nil, fmt.Errorf("failed to load High catalog: %w", err)
	}
	return p, nil
}

// NewModerateProgram creates a new FedRAMP Moderate program
func NewModerateProgram() (*Program, error) {
	p := &Program{
		name:        "FedRAMP Moderate",
		description: "Federal Risk and Authorization Management Program (FedRAMP) Moderate Security Controls Baseline",
		fs:          moderateFS,
		dataFile:    "moderate/data/FedRAMP_rev5_MODERATE-baseline-resolved-profile_catalog.yaml",
	}
	if err := p.loadCatalog(); err != nil {
		return nil, fmt.Errorf("failed to load Moderate catalog: %w", err)
	}
	return p, nil
}

// loadCatalog loads and parses the YAML catalog file
func (p *Program) loadCatalog() error {
	f, err := p.fs.Open(p.dataFile)
	if err != nil {
		return fmt.Errorf("failed to open catalog file: %w", err)
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return fmt.Errorf("failed to read catalog file: %w", err)
	}

	var catalog map[string]interface{}
	if err := yaml.Unmarshal(data, &catalog); err != nil {
		return fmt.Errorf("failed to parse catalog file: %w", err)
	}

	p.catalog = catalog
	return nil
}

// Name returns the program name
func (p *Program) Name() string {
	return p.name
}

// Description returns the program description
func (p *Program) Description() string {
	return p.description
}

// GetControl returns the control information for the given control ID
func (p *Program) GetControl(id string) (map[string]interface{}, error) {
	controls, ok := p.catalog["controls"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid catalog format: controls section not found")
	}

	control, ok := controls[id].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("control %s not found", id)
	}

	// Add the ID to the control data
	control["id"] = id
	return control, nil
}

// GetControlFamily returns all controls in the given family
func (p *Program) GetControlFamily(family string) ([]map[string]interface{}, error) {
	controls, ok := p.catalog["controls"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid catalog format: controls section not found")
	}

	var familyControls []map[string]interface{}
	for id, control := range controls {
		c, ok := control.(map[string]interface{})
		if !ok {
			continue
		}

		if fam, ok := c["family"].(string); ok && fam == family {
			c["id"] = id
			familyControls = append(familyControls, c)
		}
	}

	// Sort controls by ID
	sort.Slice(familyControls, func(i, j int) bool {
		return familyControls[i]["id"].(string) < familyControls[j]["id"].(string)
	})

	return familyControls, nil
}

// ListControlFamilies returns a list of all control families in the catalog
func (p *Program) ListControlFamilies() ([]string, error) {
	controls, ok := p.catalog["controls"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid catalog format: controls section not found")
	}

	families := make(map[string]struct{})
	for _, control := range controls {
		c, ok := control.(map[string]interface{})
		if !ok {
			continue
		}

		if family, ok := c["family"].(string); ok {
			families[family] = struct{}{}
		}
	}

	result := make([]string, 0, len(families))
	for family := range families {
		result = append(result, family)
	}
	sort.Strings(result)
	return result, nil
}

// GetControlsByStatus returns all controls with the given implementation status
func (p *Program) GetControlsByStatus(status string) ([]map[string]interface{}, error) {
	controls, ok := p.catalog["controls"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid catalog format: controls section not found")
	}

	var statusControls []map[string]interface{}
	for id, control := range controls {
		c, ok := control.(map[string]interface{})
		if !ok {
			continue
		}

		if impl, ok := c["implementation-status"].(string); ok && strings.EqualFold(impl, status) {
			c["id"] = id
			statusControls = append(statusControls, c)
		}
	}

	// Sort controls by ID
	sort.Slice(statusControls, func(i, j int) bool {
		return statusControls[i]["id"].(string) < statusControls[j]["id"].(string)
	})

	return statusControls, nil
}

// SearchControls searches for controls matching the given query in their title or description
func (p *Program) SearchControls(query string) ([]map[string]interface{}, error) {
	controls, ok := p.catalog["controls"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid catalog format: controls section not found")
	}

	query = strings.ToLower(query)
	var matches []map[string]interface{}
	for id, control := range controls {
		c, ok := control.(map[string]interface{})
		if !ok {
			continue
		}

		title, _ := c["title"].(string)
		desc, _ := c["description"].(string)

		if strings.Contains(strings.ToLower(title), query) ||
			strings.Contains(strings.ToLower(desc), query) {
			c["id"] = id
			matches = append(matches, c)
		}
	}

	// Sort matches by ID
	sort.Slice(matches, func(i, j int) bool {
		return matches[i]["id"].(string) < matches[j]["id"].(string)
	})

	return matches, nil
}

// GetControlParameters returns the parameters for a given control
func (p *Program) GetControlParameters(id string) (map[string]interface{}, error) {
	control, err := p.GetControl(id)
	if err != nil {
		return nil, err
	}

	params, ok := control["parameters"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("no parameters found for control %s", id)
	}

	return params, nil
}

// InitializePrograms creates and registers all FedRAMP programs
func InitializePrograms(r Registry) error {
	// Create and register FedRAMP High program
	high, err := NewHighProgram()
	if err != nil {
		return fmt.Errorf("failed to create FedRAMP High program: %w", err)
	}
	r.Register(high)

	// Create and register FedRAMP Moderate program
	moderate, err := NewModerateProgram()
	if err != nil {
		return fmt.Errorf("failed to create FedRAMP Moderate program: %w", err)
	}
	r.Register(moderate)

	return nil
}

// Registry defines the interface for registering programs
type Registry interface {
	Register(p interface{})
}
