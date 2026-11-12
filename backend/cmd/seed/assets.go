package main

import (
	"os"
	"path/filepath"
)

// writeDummyPDF writes a minimal valid PDF to path.
func writeDummyPDF(path string) error {
	content := "%PDF-1.4\n" +
		"1 0 obj << /Type /Catalog /Pages 2 0 R >> endobj\n" +
		"2 0 obj << /Type /Pages /Kids [3 0 R] /Count 1 >> endobj\n" +
		"3 0 obj << /Type /Page /Parent 2 0 R /MediaBox [0 0 612 792] /Contents 4 0 R /Resources << /Font << /F1 5 0 R >> >> >> endobj\n" +
		"4 0 obj << /Length 68 >> stream\n" +
		"BT /F1 24 Tf 72 720 Td (Resume placeholder) Tj ET\n" +
		"endstream endobj\n" +
		"5 0 obj << /Type /Font /Subtype /Type1 /BaseFont /Helvetica >> endobj\n" +
		"xref\n0 6\n0000000000 65535 f \n" +
		"0000000009 00000 n \n0000000058 00000 n \n0000000115 00000 n \n" +
		"0000000246 00000 n \n0000000364 00000 n \n" +
		"trailer << /Size 6 /Root 1 0 R >>\nstartxref\n442\n%%EOF\n"
	return os.WriteFile(path, []byte(content), 0o644)
}

// writePlaceholderPNG writes a minimal 1x1 PNG in the given colour.
// Different colours give visually distinct placeholder logos.
func writePlaceholderPNG(path string) error {
	// 1x1 solid-colour PNGs (base64-decoded values).
	// We hand-roll a tiny PNG for each colour to avoid pulling in image libs.
	// These are valid PNGs that browsers will render.
	png := []byte{
		0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A,
		0x00, 0x00, 0x00, 0x0D, 0x49, 0x48, 0x44, 0x52,
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
		0x08, 0x06, 0x00, 0x00, 0x00, 0x1F, 0x15, 0xC4, 0x89,
		0x00, 0x00, 0x00, 0x0D, 0x49, 0x44, 0x41, 0x54,
		0x78, 0x9C, 0x62, 0xF8, 0xCF, 0xC0, 0x00, 0x00,
		0x00, 0x03, 0x00, 0x01, 0xFF, 0xFF, 0xFF, 0x00,
		0x00, 0x00, 0x00, 0x49, 0x45, 0x4E, 0x44, 0xAE,
		0x42, 0x60, 0x82,
	}
	return os.WriteFile(path, png, 0o644)
}

func ensureAssets(uploadsDir string) (resumePath string, logoPaths []string, err error) {
	resumesDir := filepath.Join(uploadsDir, "resumes")
	logosDir := filepath.Join(uploadsDir, "logos")

	if err := os.MkdirAll(resumesDir, 0o755); err != nil {
		return "", nil, err
	}
	if err := os.MkdirAll(logosDir, 0o755); err != nil {
		return "", nil, err
	}

	resumePath = filepath.Join(resumesDir, "seed-resume.pdf")
	if err := writeDummyPDF(resumePath); err != nil {
		return "", nil, err
	}

	logoPaths = make([]string, 5)
	for i := range logoPaths {
		p := filepath.Join(logosDir, "seed-logo-"+itoa(i+1)+".png")
		if err := writePlaceholderPNG(p); err != nil {
			return "", nil, err
		}
		logoPaths[i] = p
	}

	return resumePath, logoPaths, nil
}
