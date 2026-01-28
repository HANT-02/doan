package security

import (
	"doan/pkg/config"
	"golang.org/x/crypto/bcrypt"
)

// PasswordHasher wraps bcrypt operations with configurable cost.

type PasswordHasher interface {
	Hash(plain string) (string, error)
	Compare(hash, plain string) error
}

type passwordHasher struct {
	cost int
}

func NewPasswordHasher(cfg config.Manager) PasswordHasher {
	cost := 12
	_ = cfg.UnmarshalKey("security.bcrypt_cost", &cost)
	if cost < bcrypt.MinCost {
		cost = bcrypt.DefaultCost
	}
	return &passwordHasher{cost: cost}
}

func (h *passwordHasher) Hash(plain string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(plain), h.cost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (h *passwordHasher) Compare(hash, plain string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain))
}
