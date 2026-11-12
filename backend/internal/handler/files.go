package handler

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"job-portal/internal/storage"
)

var (
	resumeAllowedMIME = []string{"application/pdf"}
	logoAllowedMIME   = []string{"image/png", "image/jpeg", "image/webp"}
)

func extFromMIME(mime string) string {
	switch mime {
	case "application/pdf":
		return "pdf"
	case "image/png":
		return "png"
	case "image/jpeg":
		return "jpg"
	case "image/webp":
		return "webp"
	}
	return ""
}

type uploadedFile struct {
	Reader      io.ReadSeeker
	Size        int64
	Original    string
	SniffedMIME string
	closer      io.Closer
}

func (u *uploadedFile) Close() error {
	if u.closer != nil {
		return u.closer.Close()
	}
	return nil
}

func parseUpload(c echo.Context, maxBytes int64, allowedMIME []string) (*uploadedFile, error) {
	fh, err := c.FormFile("file")
	if err != nil {
		return nil, fmt.Errorf("missing file field: %w", err)
	}
	if fh.Size <= 0 {
		return nil, errors.New("empty file")
	}
	if fh.Size > maxBytes {
		return nil, storage.ErrFileTooLarge
	}

	f, err := fh.Open()
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}

	header := make([]byte, 512)
	n, _ := io.ReadFull(f, header)
	header = header[:n]
	mime := http.DetectContentType(header)

	if !mimeAllowed(mime, allowedMIME) {
		f.Close()
		return nil, storage.ErrUnsupportedType
	}

	if _, err := f.Seek(0, io.SeekStart); err != nil {
		f.Close()
		return nil, fmt.Errorf("rewind file: %w", err)
	}

	original := filepath.Base(strings.TrimSpace(fh.Filename))
	if original == "" || original == "." {
		original = "file"
	}

	return &uploadedFile{
		Reader:      f,
		Size:        fh.Size,
		Original:    original,
		SniffedMIME: mime,
		closer:      f,
	}, nil
}

func mimeAllowed(got string, allowed []string) bool {
	if i := strings.Index(got, ";"); i != -1 {
		got = strings.TrimSpace(got[:i])
	}
	for _, a := range allowed {
		if a == got {
			return true
		}
	}
	return false
}

func mapUploadError(c echo.Context, err error) error {
	switch {
	case errors.Is(err, storage.ErrFileTooLarge):
		return c.JSON(http.StatusRequestEntityTooLarge, map[string]string{
			"error": "file is too large",
		})
	case errors.Is(err, storage.ErrUnsupportedType):
		return c.JSON(http.StatusUnsupportedMediaType, map[string]string{
			"error": "unsupported file type",
		})
	default:
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid upload: " + err.Error(),
		})
	}
}

func newStorageName(mime string) (string, error) {
	ext := extFromMIME(mime)
	if ext == "" {
		return "", storage.ErrUnsupportedType
	}
	return uuid.NewString() + "." + ext, nil
}

var _ = multipart.FileHeader{}
