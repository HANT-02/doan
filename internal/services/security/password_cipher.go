package security

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"doan/pkg/config"
)

// PasswordCipher defines decrypt contract for FE->BE encrypted password payloads.
// Phase 1: allow plaintext when security.accept_plain_password=true (default for dev).
// If plaintext is not accepted, this implementation currently supports a minimal JSON
// envelope and returns an error (AES-GCM can be added later without changing call sites).
//
// Expected JSON envelope (future-proof):
// {
//   "alg": "AES-256-GCM",
//   "nonce": "<base64>",
//   "ciphertext": "<base64>"
// }
//
// For now, if accept_plain_password=true, Decrypt returns the input as-is.
// If accept_plain_password=false and the input looks like a JSON envelope, an
// informative error is returned so clients know encryption is not yet enabled.

type PasswordCipher interface {
	Decrypt(enc string) (string, error)
}

type passwordCipher struct {
	acceptPlain bool
}

type encEnvelope struct {
	Alg        string `json:"alg"`
	Nonce      string `json:"nonce"`
	Ciphertext string `json:"ciphertext"`
}

func NewPasswordCipher(cfg config.Manager) PasswordCipher {
	// Default behavior: allow plaintext by default for developer convenience.
	var acceptPlain bool
	if err := cfg.UnmarshalKey("security.accept_plain_password", &acceptPlain); err != nil {
		acceptPlain = true
	}
	// If app.env is explicitly set to a development value, always allow plaintext.
	var appEnv string
	_ = cfg.UnmarshalKey("app.env", &appEnv)
	appEnv = strings.ToLower(strings.TrimSpace(appEnv))
	switch appEnv {
	case "dev", "development", "local":
		acceptPlain = true
	}
	return &passwordCipher{acceptPlain: acceptPlain}
}

func (p *passwordCipher) Decrypt(enc string) (string, error) {
	if p.acceptPlain {
		return enc, nil
	}
	// Try to parse as JSON envelope to provide friendly error
	trim := strings.TrimSpace(enc)
	if strings.HasPrefix(trim, "{") && strings.HasSuffix(trim, "}") {
		var env encEnvelope
		_ = json.Unmarshal([]byte(trim), &env)
		return "", errors.New("encrypted password received but AES-GCM decryption is not enabled yet on server")
	}
	return "", fmt.Errorf("plaintext password is not accepted; enable security.accept_plain_password or send encrypted payload")
}
