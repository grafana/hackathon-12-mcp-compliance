# Adding Support for a New Compliance Program

This guide explains how to add support for a new compliance program to the MCP Compliance project. It focuses on the architecture and interfaces rather than specific implementation details, as each program may have different data formats and requirements.

## Overview

The MCP Compliance project is designed to be extensible, allowing support for various compliance programs. Each program implementation must:

1. Implement the Program interface
2. Provide access to its control catalog
3. Register itself with the program registry

## Step-by-Step Guide

### 1. Understand the Program Interface

All compliance programs must implement the following interface:

```go
type Program interface {
    // Name returns the program's display name
    Name() string

    // Description returns a human-readable description of the program
    Description() string

    // GetControl returns detailed information about a specific control
    GetControl(id string) (map[string]interface{}, error)

    // GetControlFamily returns all controls in a specific family
    GetControlFamily(family string) ([]map[string]interface{}, error)

    // ListControlFamilies returns all available control families
    ListControlFamilies() ([]string, error)

    // GetControlsByStatus returns controls with a specific implementation status
    GetControlsByStatus(status string) ([]map[string]interface{}, error)

    // SearchControls finds controls matching a search query
    SearchControls(query string) ([]map[string]interface{}, error)

    // GetControlParameters returns parameters for a specific control
    GetControlParameters(id string) (map[string]interface{}, error)
}
```

### 2. Create the Program Structure

1. Create a new directory under `internal/programs/` for your program
2. Create the basic program structure:
   ```
   internal/programs/yourprogram/
   ├── data/           # Program-specific data files
   ├── program.go      # Main program implementation
   └── README.md       # Program-specific documentation
   ```

### 3. Implement the Program

1. Create a struct that implements the Program interface
2. Decide how to load and store your program's control catalog
3. Implement each interface method based on your program's data structure

Example structure:
```go
package yourprogram

type Program struct {
    // Add fields needed for your implementation
    catalog    map[string]interface{}
    metadata   ProgramMetadata
}

func New() (*Program, error) {
    // Initialize your program
}

// Implement interface methods
func (p *Program) Name() string {
    return "Your Program Name"
}

// ... implement other interface methods ...
```

### 4. Data Management

Consider these aspects when implementing data management:

1. **Data Format**: Choose an appropriate format for your control catalog (YAML, JSON, etc.)
2. **Data Loading**: Implement a strategy for loading your data:
   - Embed static files using `//go:embed`
   - Load from external sources
   - Generate programmatically
3. **Data Structure**: Design an internal representation that makes it efficient to:
   - Look up controls by ID
   - Group controls by family
   - Search control content
   - Access control parameters

### 5. Register Your Program

1. Add your program to the registry initialization in internal/programs/fedramp/registry.go:
   ```go
   func RegisterPrograms(r *registry.Registry) error {
       // Create and register your program
       yourProgram, err := yourpackage.New()
       if err != nil {
           return err
       }
       r.Register(yourProgram)
       return nil
   }
   ```

### 6. Testing

Create tests that verify:

1. Your program correctly implements the interface
2. Control data is properly loaded and accessible
3. Search and filtering functions work as expected
4. Error cases are handled appropriately

## Best Practices

1. **Error Handling**
   - Return clear, specific error messages
   - Handle missing or malformed data gracefully
   - Validate input parameters

2. **Performance**
   - Cache data when appropriate
   - Use efficient data structures for lookups
   - Consider memory usage for large control sets

3. **Documentation**
   - Document any program-specific behaviors
   - Explain any deviations from standard patterns
   - Include examples of using your program

4. **Maintainability**
   - Keep the implementation modular
   - Separate data loading from business logic
   - Use clear, consistent naming conventions

## Example Use Cases

Here are some common ways your program might be used:

```go
// Get a specific control
control, err := program.GetControl("AC-1")

// List all controls in a family
controls, err := program.GetControlFamily("Access Control")

// Search for controls
results, err := program.SearchControls("authentication")

// Get control parameters
params, err := program.GetControlParameters("AC-2")
```

## Common Challenges

1. **Data Normalization**: Different programs may use different terms for similar concepts. Consider providing mapping or normalization functions.

2. **Version Management**: If your program has multiple versions, design your implementation to handle version differences.

3. **Control Relationships**: Some controls may reference or depend on others. Consider how to represent and navigate these relationships.

4. **Custom Extensions**: Your program may have unique features not covered by the base interface. Consider creating program-specific extensions while maintaining base compatibility.

## Need Help?

- Review existing program implementations for examples
- Check the project's issue tracker for known challenges
- Reach out to the maintainers for guidance

Remember: The goal is to provide a consistent interface for working with your compliance program while handling its unique characteristics appropriately. 