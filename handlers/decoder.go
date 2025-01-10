package handlers

import (
	"encoding/base64"
	"mime"
	"path/filepath"
)

func Base64Decoder(input string) (string, error) {
	decodedBytes, err := base64.StdEncoding.DecodeString(input)
	if err != nil {
		return "", err
	}
	return string(decodedBytes), nil
}

func DetectMimeType(fileName string) string {
	ext := filepath.Ext(fileName)
	mimeType := mime.TypeByExtension(ext)

	if mimeType == "" {
		return "application/octet-stream"
	} else {
		return mimeType
	}
}
