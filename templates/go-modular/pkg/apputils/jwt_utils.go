package apputils

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
)

// ContextKey is a typed key for context values to avoid collisions and staticcheck SA1029.
type ContextKey string

const (
	// HeadersContextKey is used to store request headers map in context (map[string]string).
	HeadersContextKey ContextKey = "{{ package_name | kebab_case }}.headers"

	// JWTClaimsContextKey is used to store parsed JWT claims in context (map[string]any).
	JWTClaimsContextKey ContextKey = "{{ package_name | kebab_case }}.jwt_claims"
)

// JWTConfig holds configuration for JWT generation and validation.
type JWTConfig struct {
	SecretKey          []byte                 // Secret key for signing JWTs
	AccessTokenExpiry  time.Duration          // Access token expiration duration
	RefreshTokenExpiry time.Duration          // Refresh token expiration duration
	SigningAlg         jwa.SignatureAlgorithm // Signing algorithm (e.g. jwa.HS256)
	Issuer             string                 // JWT issuer claim
}

// JWTGenerator provides methods to sign and validate JWTs using HMAC.
type JWTGenerator struct {
	config JWTConfig
}

// NewJWTGenerator creates a new JWTGenerator with the given configuration.
// If no signing algorithm is provided, it defaults to HS256.
func NewJWTGenerator(config JWTConfig) *JWTGenerator {
	if config.SigningAlg == "" {
		config.SigningAlg = jwa.HS256
	}
	return &JWTGenerator{config: config}
}

// Sign generates a JWT string with the given payload (struct or map) and optional subject.
// The payload is flattened into the JWT claims. The "typ" claim is set to "access".
func (j *JWTGenerator) Sign(ctx context.Context, payload any, subject string) (string, error) {
	if len(j.config.SecretKey) == 0 {
		return "", errors.New("secret key is required")
	}
	token := jwt.New()
	now := time.Now()
	_ = token.Set(jwt.IssuerKey, j.config.Issuer)
	_ = token.Set(jwt.IssuedAtKey, now)
	_ = token.Set(jwt.ExpirationKey, now.Add(j.config.AccessTokenExpiry))
	_ = token.Set("typ", "access")
	if subject != "" {
		_ = token.Set(jwt.SubjectKey, subject)
	}

	// Marshal the payload to map[string]any and flatten into JWT claims.
	b, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("marshal payload: %w", err)
	}
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		return "", fmt.Errorf("unmarshal payload: %w", err)
	}
	for k, v := range m {
		_ = token.Set(k, v)
	}

	signed, err := jwt.Sign(token, j.config.SigningAlg, j.config.SecretKey)
	if err != nil {
		return "", fmt.Errorf("sign jwt: %w", err)
	}
	return string(signed), nil
}

// GenerateRefreshTokenJWT generates a JWT as a refresh token with a simple payload.
// The "typ" claim is set to "refresh" and "jti" is set to the refresh token ID.
func (j *JWTGenerator) GenerateRefreshTokenJWT(ctx context.Context, uid, audience, refreshTokenID string) (string, error) {
	if len(j.config.SecretKey) == 0 {
		return "", errors.New("secret key is required")
	}
	token := jwt.New()
	now := time.Now()
	_ = token.Set(jwt.IssuerKey, j.config.Issuer)
	_ = token.Set(jwt.IssuedAtKey, now)
	_ = token.Set(jwt.ExpirationKey, now.Add(j.config.RefreshTokenExpiry))
	_ = token.Set(jwt.SubjectKey, uid)
	_ = token.Set(jwt.AudienceKey, audience)
	_ = token.Set("typ", "refresh")
	_ = token.Set(jwt.JwtIDKey, refreshTokenID)

	signed, err := jwt.Sign(token, j.config.SigningAlg, j.config.SecretKey)
	if err != nil {
		return "", fmt.Errorf("sign refresh jwt: %w", err)
	}
	return string(signed), nil
}

// ParseAndValidate parses and validates a JWT string, returning the claims as a map if valid.
// It verifies the signature and validates standard claims (exp, nbf, etc).
func (j *JWTGenerator) ParseAndValidate(ctx context.Context, tokenString string) (map[string]any, error) {
	token, err := jwt.Parse(
		[]byte(tokenString),
		jwt.WithVerify(j.config.SigningAlg, j.config.SecretKey),
		jwt.WithValidate(true),
	)
	if err != nil {
		return nil, fmt.Errorf("parse/validate jwt: %w", err)
	}
	claims := make(map[string]any)
	for it := token.Iterate(ctx); it.Next(ctx); {
		pair := it.Pair()
		claims[pair.Key.(string)] = pair.Value
	}
	return claims, nil
}

// ParseAndUnmarshal parses a JWT and unmarshals claims into the provided struct (for custom claims).
func (j *JWTGenerator) ParseAndUnmarshal(ctx context.Context, tokenString string, out any) error {
	claims, err := j.ParseAndValidate(ctx, tokenString)
	if err != nil {
		return err
	}
	b, err := json.Marshal(claims)
	if err != nil {
		return fmt.Errorf("marshal claims: %w", err)
	}
	return json.Unmarshal(b, out)
}

// GetHash returns the hex-encoded SHA256 hash of the input string.
func (j *JWTGenerator) GetHash(s string) string {
	hash := sha256.Sum256([]byte(s))
	return hex.EncodeToString(hash[:])
}

// GetSigningKey returns the secret key used for signing JWTs.
func (j *JWTGenerator) GetSigningKey() []byte {
	return j.config.SecretKey
}

// AccessTokenExpiry returns the configured access token expiry duration.
func (j *JWTGenerator) AccessTokenExpiry() time.Duration {
	return j.config.AccessTokenExpiry
}

// RefreshTokenExpiry returns the configured refresh token expiry duration.
func (j *JWTGenerator) RefreshTokenExpiry() time.Duration {
	return j.config.RefreshTokenExpiry
}
