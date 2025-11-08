package randomly

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"strings"
)

const (
	charSet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()_+-=[]{}|;:,.<>?"
)

func GeneratePassword(k int) (string, error) {
	if k <= 0 {
		return "", errors.New("just positive number accepted")
	}

	randomBytes := make([]byte, k)

	if _, err := io.ReadFull(rand.Reader, randomBytes); err != nil {
		return "", fmt.Errorf("error in read randomly bytes, error: %s", err.Error())
	}

	var builder strings.Builder
	builder.Grow(k)

	charSetLength := byte(len(charSet))

	for _, b := range randomBytes {
		idx := b % charSetLength
		builder.WriteByte(charSet[idx])
	}

	return builder.String(), nil
}
