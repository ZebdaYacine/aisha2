package storage

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/http"
	"strings"
)

var (
	ErrUnsupportedFileType  = errors.New("unsupported file type")
	ErrFileTooLarge         = errors.New("file is too large")
	ErrInvalidFileSignature = errors.New("file signature is invalid")
)

type ValidatedFile struct {
	Bytes     []byte
	MediaType string
	Kind      string
	Checksum  string
}

func ValidateFile(data []byte, declaredType string, maxBytes int64, allowedKinds map[string]bool) (ValidatedFile, error) {
	if int64(len(data)) > maxBytes {
		return ValidatedFile{}, ErrFileTooLarge
	}
	mediaType, kind := detect(data)
	if mediaType == "" || !allowedKinds[kind] {
		return ValidatedFile{}, ErrUnsupportedFileType
	}
	if declaredType != "" && declaredType != "application/octet-stream" && !sameFamily(strings.ToLower(declaredType), mediaType) {
		return ValidatedFile{}, ErrInvalidFileSignature
	}
	hash := sha256.Sum256(data)
	return ValidatedFile{Bytes: data, MediaType: mediaType, Kind: kind, Checksum: hex.EncodeToString(hash[:])}, nil
}

func detect(data []byte) (string, string) {
	switch {
	case bytes.HasPrefix(data, []byte{0xff, 0xd8, 0xff}):
		return "image/jpeg", "IMAGE"
	case bytes.HasPrefix(data, []byte{0x89, 'P', 'N', 'G', '\r', '\n', 0x1a, '\n'}):
		return "image/png", "IMAGE"
	case len(data) >= 12 && string(data[:4]) == "RIFF" && string(data[8:12]) == "WEBP":
		return "image/webp", "IMAGE"
	case bytes.HasPrefix(data, []byte("%PDF-")):
		return "application/pdf", "DOCUMENT"
	case len(data) >= 12 && string(data[4:8]) == "ftyp":
		return "video/mp4", "VIDEO"
	default:
		mediaType := http.DetectContentType(data)
		if mediaType == "image/jpeg" || mediaType == "image/png" || mediaType == "image/webp" {
			return mediaType, "IMAGE"
		}
		return "", ""
	}
}

func sameFamily(declared, detected string) bool {
	if declared == detected {
		return true
	}
	return strings.HasPrefix(declared, "image/") && strings.HasPrefix(detected, "image/")
}
