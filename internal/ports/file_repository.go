package ports

import "io"

// FileRepository defines the interface for file operations
type FileRepository interface {
	// ReadFile reads a file from the given path and returns its contents
	ReadFile(path string) ([]byte, error)

	// WriteFile writes data to a file at the given path
	WriteFile(path string, data []byte) error

	// FileExists checks if a file exists at the given path
	FileExists(path string) bool

	// CreateDirectory creates a directory at the given path
	CreateDirectory(path string) error

	// OpenFile opens a file for reading and returns a ReadCloser
	OpenFile(path string) (io.ReadCloser, error)

	// CreateFile creates a file for writing and returns a WriteCloser
	CreateFile(path string) (io.WriteCloser, error)
}
