package utils

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"strings"
)

func Hash256(src string) string {
	h := sha256.New()
	h.Write([]byte(src))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func HashData(data []byte) string {
	h := sha256.New()
	h.Write(data)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func HashItems(src ...string) string {
	b := strings.Builder{}
	for i := range src {
		b.WriteString(src[i])
	}
	return Hash256(b.String())
}

func GetFileHash(filePath string) (string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("unable to open file: %w", err)
	}
	defer f.Close()

	h := sha256.New()
	if _, err = io.Copy(h, f); err != nil {
		return "", fmt.Errorf("unable to copy file: %w", err)
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
