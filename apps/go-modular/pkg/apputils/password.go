package apputils

// Example usage:
//
// // Using default parameters
// hasher := apputils.NewPasswordHasher()
// hash, err := hasher.Hash("mySecretPassword")
// if err != nil {
//     // handle error
// }
// ok, err := hasher.Validate("mySecretPassword", hash)
// if ok {
//     // password is valid
// }
//
// // Using custom parameters
// params := apputils.Argon2Params{
//     Memory:      32768,
//     Iterations:  6,
//     Parallelism: 4,
//     SaltLength:  16,
//     KeyLength:   32,
// }
// customHasher := apputils.NewPasswordHasherWithParams(params)
// customHash, err := customHasher.Hash("mySecretPassword")
// if err != nil {
//     // handle error
// }
// ok, err = customHasher.Validate("mySecretPassword", customHash)
// if ok {
//     // password is valid
// }

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

// Argon2Params holds the parameters for argon2id hashing
type Argon2Params struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

// PasswordHasher provides configurable password hashing and validation
type PasswordHasher struct {
	params Argon2Params
}

// DefaultArgon2Params returns recommended default parameters
func DefaultArgon2Params() Argon2Params {
	return Argon2Params{
		Memory:      16384, // 16 MB
		Iterations:  4,     // 4 iterations
		Parallelism: 2,     // 2 threads
		SaltLength:  8,     // 8 bytes salt
		KeyLength:   32,    // 32 bytes key length
	}
}

// NewPasswordHasher creates a PasswordHasher with default params
func NewPasswordHasher() *PasswordHasher {
	return &PasswordHasher{params: DefaultArgon2Params()}
}

// WithParams creates a PasswordHasher with custom params
func NewPasswordHasherWithParams(params Argon2Params) *PasswordHasher {
	return &PasswordHasher{params: params}
}

// Hash hashes the password using argon2id and returns a PHC string
func (h *PasswordHasher) Hash(password string) (string, error) {
	salt := make([]byte, h.params.SaltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, h.params.Iterations, h.params.Memory, h.params.Parallelism, h.params.KeyLength)
	phc := fmt.Sprintf("$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s",
		h.params.Memory, h.params.Iterations, h.params.Parallelism,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(hash),
	)

	return phc, nil
}

// Validate compares a plain password with a PHC argon2id hash
func (h *PasswordHasher) Validate(password, hash string) (bool, error) {
	parts := strings.Split(hash, "$")
	if len(parts) != 6 {
		return false, errors.New("invalid hash format")
	}

	var memory uint32
	var iterations uint32
	var parallelism uint8
	_, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &iterations, &parallelism)
	if err != nil {
		return false, err
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, err
	}

	expectedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, err
	}

	keyLen := uint32(len(expectedHash))
	computedHash := argon2.IDKey([]byte(password), salt, iterations, memory, parallelism, keyLen)

	return subtleCompare(computedHash, expectedHash), nil
}

// subtleCompare does a constant-time comparison of two byte slices
func subtleCompare(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	var result byte
	for i := range a {
		result |= a[i] ^ b[i]
	}
	return result == 0
}
