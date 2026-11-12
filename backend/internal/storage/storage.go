package storage

import (
	"errors"
	"io"
)

var (
	ErrFileTooLarge    = errors.New("file too large")
	ErrUnsupportedType = errors.New("unsupported file type")
)

type FileStorage interface {
	Save(category, name string, r io.Reader) (urlPath string, err error)
	Delete(urlPath string) error
}
