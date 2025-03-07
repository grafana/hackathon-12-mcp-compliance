package programs

// Program defines the interface that all compliance programs must implement
type Program interface {
	// Basic information
	Name() string
	Description() string

	// Control management
	GetControl(id string) (map[string]interface{}, error)
	GetControlFamily(family string) ([]map[string]interface{}, error)
	ListControlFamilies() ([]string, error)
	SearchControls(query string) ([]map[string]interface{}, error)

	// Additional control operations
	GetControlsByStatus(status string) ([]map[string]interface{}, error)
	GetControlParameters(id string) (map[string]interface{}, error)
}
