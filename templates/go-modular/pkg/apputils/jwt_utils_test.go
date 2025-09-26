package apputils

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJWTGenerator(t *testing.T) {
	ctx := context.Background()

	t.Run("Sign_ParseAndValidate_ParseAndUnmarshal", func(t *testing.T) {
		secret := []byte("my-very-secret-key")
		cfg := JWTConfig{
			SecretKey:          secret,
			AccessTokenExpiry:  time.Minute * 5,
			RefreshTokenExpiry: time.Hour * 24,
			Issuer:             "test-issuer",
		}
		gen := NewJWTGenerator(cfg)

		payload := map[string]any{
			"name": "alice",
			"role": "admin",
		}

		// Sign an access token
		tokenStr, err := gen.Sign(ctx, payload, "user-123")
		require.NoError(t, err)
		require.NotEmpty(t, tokenStr)

		// Parse and validate
		claims, err := gen.ParseAndValidate(ctx, tokenStr)
		require.NoError(t, err)
		assert.Equal(t, "access", claims["typ"])
		assert.Equal(t, "test-issuer", claims["iss"])
		assert.Equal(t, "alice", claims["name"])
		assert.Equal(t, "admin", claims["role"])
		assert.Equal(t, "user-123", claims["sub"])

		// Ensure exp claim exists and is in the future (handle common types)
		expVal, ok := claims["exp"]
		require.True(t, ok, "exp claim must be present")
		switch v := expVal.(type) {
		case time.Time:
			assert.True(t, v.After(time.Now()), "exp (time.Time) should be in the future")
		case float64:
			assert.True(t, time.Unix(int64(v), 0).After(time.Now()), "exp (float64) should be in the future")
		case json.Number:
			n, err := v.Int64()
			require.NoError(t, err)
			assert.True(t, time.Unix(n, 0).After(time.Now()), "exp (json.Number) should be in the future")
		default:
			t.Fatalf("unexpected exp claim type: %T", v)
		}

		// Test ParseAndUnmarshal into a struct
		type outClaims struct {
			Name string `json:"name"`
			Role string `json:"role"`
			Sub  string `json:"sub"`
			Typ  string `json:"typ"`
			Iss  string `json:"iss"`
		}
		var out outClaims
		err = gen.ParseAndUnmarshal(ctx, tokenStr, &out)
		require.NoError(t, err)
		assert.Equal(t, "alice", out.Name)
		assert.Equal(t, "admin", out.Role)
		assert.Equal(t, "user-123", out.Sub)
		assert.Equal(t, "access", out.Typ)
		assert.Equal(t, "test-issuer", out.Iss)
	})

	t.Run("GenerateRefreshTokenJWT", func(t *testing.T) {
		secret := []byte("refresh-secret")
		cfg := JWTConfig{
			SecretKey:          secret,
			AccessTokenExpiry:  time.Minute * 5,
			RefreshTokenExpiry: time.Hour * 24 * 7,
			Issuer:             "refresh-issuer",
		}
		gen := NewJWTGenerator(cfg)

		rt, err := gen.GenerateRefreshTokenJWT(ctx, "uid-321", "my-audience", "refresh-id-999")
		require.NoError(t, err)
		require.NotEmpty(t, rt)

		claims, err := gen.ParseAndValidate(ctx, rt)
		require.NoError(t, err)
		assert.Equal(t, "refresh", claims["typ"])
		assert.Equal(t, "uid-321", claims["sub"])

		// aud claim may be a string or a slice depending on parser representation
		audVal, ok := claims["aud"]
		require.True(t, ok, "aud claim must be present")
		switch v := audVal.(type) {
		case string:
			assert.Equal(t, "my-audience", v)
		case []string:
			assert.Equal(t, []string{"my-audience"}, v)
		case []any:
			if len(v) == 0 {
				t.Fatalf("aud claim empty")
			}
			s, ok := v[0].(string)
			require.True(t, ok)
			assert.Equal(t, "my-audience", s)
		default:
			t.Fatalf("unexpected aud claim type: %T", v)
		}

		assert.Equal(t, "refresh-id-999", claims["jti"])
		assert.Equal(t, "refresh-issuer", claims["iss"])
	})

	t.Run("ErrorsAndHelpers", func(t *testing.T) {
		// Missing secret key should cause Sign and GenerateRefreshTokenJWT to error
		genNoKey := NewJWTGenerator(JWTConfig{
			SecretKey:          nil,
			AccessTokenExpiry:  time.Minute,
			RefreshTokenExpiry: time.Hour,
			Issuer:             "x",
		})
		_, err := genNoKey.Sign(ctx, map[string]any{"a": 1}, "")
		assert.Error(t, err)

		_, err = genNoKey.GenerateRefreshTokenJWT(ctx, "u", "aud", "jti")
		assert.Error(t, err)

		// Helpers: GetHash, GetSigningKey, expiry getters
		secret := []byte("helper-secret")
		cfg := JWTConfig{
			SecretKey:          secret,
			AccessTokenExpiry:  2 * time.Minute,
			RefreshTokenExpiry: 48 * time.Hour,
			Issuer:             "help-iss",
		}
		gen := NewJWTGenerator(cfg)

		// GetSigningKey
		assert.Equal(t, secret, gen.GetSigningKey())

		// Access/Refresh expiry
		assert.Equal(t, 2*time.Minute, gen.AccessTokenExpiry())
		assert.Equal(t, 48*time.Hour, gen.RefreshTokenExpiry())

		// GetHash should match manual sha256 hex
		input := "some-string-to-hash"
		sum := sha256.Sum256([]byte(input))
		expected := hex.EncodeToString(sum[:])
		assert.Equal(t, expected, gen.GetHash(input))
	})
}
