package utils

import (
	"crypto/sha256"
	"fmt"
	"strings"
)

func Hash256(src string) string {
	h := sha256.New()
	h.Write([]byte(src))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func HashItems(src ...string) string {
	b := strings.Builder{}
	for i := range src {
		b.WriteString(src[i])
	}
	return Hash256(b.String())
}
