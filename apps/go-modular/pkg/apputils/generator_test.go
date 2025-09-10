package apputils

import (
	"regexp"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestGenerateURLSafeToken(t *testing.T) {
	t.Run("ValidLength", func(t *testing.T) {
		length := 32
		tok, err := GenerateURLSafeToken(length)
		require.NoError(t, err)
		require.Equal(t, length, len(tok), "token length should match requested length")

		re := regexp.MustCompile(`^[A-Za-z0-9]+[0-9]{10}$`)
		require.True(t, re.MatchString(tok), "token must be alphanumeric followed by 10-digit timestamp")

		tsStr := tok[len(tok)-10:]
		ts, err := strconv.ParseInt(tsStr, 10, 64)
		require.NoError(t, err)

		now := time.Now().Unix()
		diff := now - ts
		if diff < 0 {
			diff = -diff
		}
		require.LessOrEqual(t, diff, int64(5), "timestamp should be within 5 seconds of now")
	})

	t.Run("MinLength", func(t *testing.T) {
		length := 11 // smallest allowed: 1-char random + 10-digit timestamp
		tok, err := GenerateURLSafeToken(length)
		require.NoError(t, err)
		require.Equal(t, length, len(tok), "token length should match requested length")

		re := regexp.MustCompile(`^[A-Za-z0-9][0-9]{10}$`)
		require.True(t, re.MatchString(tok), "min-length token must be 1 alnum char followed by 10-digit timestamp")
	})

	t.Run("TooShort", func(t *testing.T) {
		length := 10 // too short, must return error
		tok, err := GenerateURLSafeToken(length)
		require.Error(t, err)
		require.Empty(t, tok, "token should be empty when generation fails")
	})
}
