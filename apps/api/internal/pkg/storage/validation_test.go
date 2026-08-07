package storage

import (
	"errors"
	"testing"
)

func TestValidateFileChecksSignatureAndDeclaredType(t *testing.T) {
	png := []byte{0x89, 'P', 'N', 'G', '\r', '\n', 0x1a, '\n'}
	file, err := ValidateFile(png, "image/png", 1024, map[string]bool{"IMAGE": true})
	if err != nil {
		t.Fatal(err)
	}
	if file.MediaType != "image/png" || file.Kind != "IMAGE" || file.Checksum == "" {
		t.Fatalf("validated file=%#v", file)
	}

	_, err = ValidateFile(png, "application/pdf", 1024, map[string]bool{"IMAGE": true})
	if !errors.Is(err, ErrInvalidFileSignature) {
		t.Fatalf("err=%v", err)
	}
}

func TestValidateFileRejectsOversizeAndUnsupportedContent(t *testing.T) {
	_, err := ValidateFile([]byte("too large"), "text/plain", 2, map[string]bool{"IMAGE": true})
	if !errors.Is(err, ErrFileTooLarge) {
		t.Fatalf("oversize err=%v", err)
	}

	_, err = ValidateFile([]byte("not media"), "image/png", 1024, map[string]bool{"IMAGE": true})
	if !errors.Is(err, ErrUnsupportedFileType) {
		t.Fatalf("unsupported err=%v", err)
	}
}
