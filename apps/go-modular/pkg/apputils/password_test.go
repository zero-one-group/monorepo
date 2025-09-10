package apputils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPasswordHasher(t *testing.T) {
	t.Run("Hash_And_Validate_Success", func(t *testing.T) {
		hasher := NewPasswordHasher()
		pass := "myS3cretP@ss"
		hash, err := hasher.Hash(pass)
		require.NoError(t, err)
		require.NotEmpty(t, hash)

		ok, err := hasher.Validate(pass, hash)
		require.NoError(t, err)
		assert.True(t, ok, "expected password validation to succeed")
	})

	t.Run("WrongPassword_FailsValidation", func(t *testing.T) {
		hasher := NewPasswordHasher()
		pass := "correct-password"
		hash, err := hasher.Hash(pass)
		require.NoError(t, err)
		require.NotEmpty(t, hash)

		ok, err := hasher.Validate("incorrect-password", hash)
		require.NoError(t, err)
		assert.False(t, ok, "expected validation to fail for wrong password")
	})

	t.Run("InvalidHashFormat_ReturnsError", func(t *testing.T) {
		hasher := NewPasswordHasher()
		_, err := hasher.Validate("any", "not-a-valid-phc")
		require.Error(t, err)
	})

	t.Run("CustomParams_Hash_And_Validate", func(t *testing.T) {
		params := Argon2Params{
			Memory:      32768,
			Iterations:  3,
			Parallelism: 1,
			SaltLength:  16,
			KeyLength:   32,
		}
		hasher := NewPasswordHasherWithParams(params)
		pass := "anotherSecret!"
		hash, err := hasher.Hash(pass)
		require.NoError(t, err)
		require.NotEmpty(t, hash)

		ok, err := hasher.Validate(pass, hash)
		require.NoError(t, err)
		assert.True(t, ok, "expected validation to succeed with custom params")
	})
}
