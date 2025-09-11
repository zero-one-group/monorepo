package apputils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSummarizeUserAgent(t *testing.T) {
	t.Run("Chrome on macOS", func(t *testing.T) {
		ua := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.5790.170 Safari/537.36"
		got := SummarizeUserAgent(ua)
		require.NotEmpty(t, got)
		assert.Equal(t, "Chrome v115.0 on macOS 10_15_7", got)
	})

	t.Run("Mobile Safari on iOS", func(t *testing.T) {
		ua := "Mozilla/5.0 (iPhone; CPU iPhone OS 14_4 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0 Mobile/15E148 Safari/604.1"
		got := SummarizeUserAgent(ua)
		require.NotEmpty(t, got)
		assert.Equal(t, "Mobile Safari v14.0 on iOS 14_4", got)
	})

	t.Run("Chrome on Android", func(t *testing.T) {
		ua := "Mozilla/5.0 (Linux; Android 11; Pixel 4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.210 Mobile Safari/537.36"
		got := SummarizeUserAgent(ua)
		require.NotEmpty(t, got)
		assert.Equal(t, "Chrome v90.0 on Android 11", got)
	})

	t.Run("Chrome on Windows 10", func(t *testing.T) {
		ua := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.5790.170 Safari/537.36"
		got := SummarizeUserAgent(ua)
		require.NotEmpty(t, got)
		assert.Equal(t, "Chrome v115.0 on Windows 10", got)
	})

	t.Run("Unknown user agent", func(t *testing.T) {
		got := SummarizeUserAgent("")
		assert.Equal(t, "Unknown", got)
	})
}
