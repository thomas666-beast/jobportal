package storage

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// LocalStorage writes files to a base directory and serves them at publicPath.
type LocalStorage struct {
	baseDir    string
	publicPath string
}

func NewLocalStorage(baseDir, publicPath string) (*LocalStorage, error) {
	if err := os.MkdirAll(baseDir, 0o755); err != nil {
		return nil, fmt.Errorf("create base dir: %w", err)
	}
	return &LocalStorage{baseDir: baseDir, publicPath: publicPath}, nil
}

// Save writes r to <baseDir>/<category>/<name> and returns the public URL path
// (e.g. /uploads/resumes/abc.pdf).
func (s *LocalStorage) Save(category, name string, r io.Reader) (string, error) {
	dir := filepath.Join(s.baseDir, category)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("create category dir: %w", err)
	}

	dst := filepath.Join(dir, name)
	f, err := os.Create(dst)
	if err != nil {
		return "", fmt.Errorf("create file: %w", err)
	}
	defer f.Close()

	if _, err := io.Copy(f, r); err != nil {
		_ = os.Remove(dst)
		return "", fmt.Errorf("write file: %w", err)
	}

	return s.publicPath + "/" + category + "/" + name, nil
}

// Delete removes the file referenced by urlPath. Missing files are ignored.
func (s *LocalStorage) Delete(urlPath string) error {
	if !strings.HasPrefix(urlPath, s.publicPath+"/") {
		return nil // not ours
	}
	rel := strings.TrimPrefix(urlPath, s.publicPath+"/")

	// Reject any path traversal attempts.
	for _, part := range strings.Split(rel, "/") {
		if part == ".." || part == "" {
			return errors.New("invalid path")
		}
	}

	full := filepath.Join(s.baseDir, rel)
	if err := os.Remove(full); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("delete file: %w", err)
	}
	return nil
}
