package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// CORS_ORIGINS reaches us in three shapes: the generated .env.example writes the Go
// slice as "[*]", operators write comma lists, and some paste the bracketed form back.
// All of them must produce real origins — a literal "[*]" allows nothing and silently
// breaks every cross-origin client.
func TestLoad_CORSOriginsAcceptsEveryCommonShape(t *testing.T) {
	cases := map[string][]string{
		"*":   {"*"},
		"[*]": {"*"},
		"https://app.example.com,https://admin.example.com":    {"https://app.example.com", "https://admin.example.com"},
		"[https://app.example.com, https://admin.example.com]": {"https://app.example.com", "https://admin.example.com"},
		" https://app.example.com , ":                          {"https://app.example.com"},
	}
	for raw, want := range cases {
		t.Run(raw, func(t *testing.T) {
			t.Setenv("CORS_ORIGINS", raw)
			cfg, err := Load("")
			require.NoError(t, err)
			assert.Equal(t, want, cfg.App.CORSOrigins)
		})
	}
}

func TestGenerateExampleEnvFile_WritesSlicesAsCommaLists(t *testing.T) {
	path := t.TempDir() + "/.env.example"
	require.NoError(t, GenerateExampleEnvFile(path))
	b, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Contains(t, string(b), "CORS_ORIGINS=*\n", "no Go-slice brackets in the example file")
}
