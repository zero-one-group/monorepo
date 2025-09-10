package apputils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"
	"time"
)

// GenerateURLSafeToken generates a cryptographically secure, URL-safe, alphanumeric token
// with a total length of 'length' characters, including a 10-digit unix timestamp appended at the end.
// Returns an error if the length is too short or random generation fails.
func GenerateURLSafeToken(length int) (string, error) {
	tokenLen := length - 10
	if tokenLen < 1 {
		return "", fmt.Errorf("token length too short")
	}
	// Calculate the number of random bytes needed to produce at least tokenLen base64 characters
	// base64.RawURLEncoding: 3 bytes = 4 chars, 1 char = 6 bits
	// chars = ceil(bytes * 4 / 3), bytes = ceil(chars * 3 / 4)
	byteLen := (tokenLen*3 + 3) / 4
	b := make([]byte, byteLen)
	_, err := rand.Read(b)
	if err != nil {
		return "", fmt.Errorf("failed to generate secure random token: %w", err)
	}
	token := base64.RawURLEncoding.EncodeToString(b)
	// Keep only alphanumeric characters (A-Z, a-z, 0-9)
	token = strings.Map(func(r rune) rune {
		if (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			return r
		}
		return -1
	}, token)
	// Truncate to the required length
	if len(token) > tokenLen {
		token = token[:tokenLen]
	}
	// Append the current unix timestamp (10 digits)
	return fmt.Sprintf("%s%d", token, time.Now().Unix()), nil
}
