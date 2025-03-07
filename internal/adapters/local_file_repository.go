package adapters

import (
	"io"
	"os"
	"path/filepath"

	"github.com/grafana/hackathon-12-mcp-compliance/internal/ports"
)

// LocalFileRepository implements the FileRepository interface using the local filesystem
type LocalFileRepository struct{}

// NewLocalFileRepository creates a new LocalFileRepository
func NewLocalFileRepository() *LocalFileRepository {
	return &LocalFileRepository{}
}

// ReadFile reads a file from the given path and returns its contents
func (r *LocalFileRepository) ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

// WriteFile writes data to a file at the given path
func (r *LocalFileRepository) WriteFile(path string, data []byte) error {
	return os.WriteFile(path, data, 0644)
}

// FileExists checks if a file exists at the given path
func (r *LocalFileRepository) FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// CreateDirectory creates a directory at the given path
func (r *LocalFileRepository) CreateDirectory(path string) error {
	return os.MkdirAll(path, 0755)
}

// OpenFile opens a file for reading and returns a ReadCloser
func (r *LocalFileRepository) OpenFile(path string) (io.ReadCloser, error) {
	return os.Open(path)
}

// CreateFile creates a file for writing and returns a WriteCloser
func (r *LocalFileRepository) CreateFile(path string) (io.WriteCloser, error) {
	// Create parent directory if it doesn't exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	return os.Create(path)
}

// Ensure LocalFileRepository implements FileRepository
var _ ports.FileRepository = (*LocalFileRepository)(nil)
