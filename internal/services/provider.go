package services

import (
	_interface "doan/internal/infrastructure/queue/interface"
	"doan/internal/infrastructure/queue/noop"
	"doan/internal/services/mailer"
	"doan/internal/services/security"
	"doan/internal/services/user"
	"doan/pkg/config"
	"doan/pkg/logger"
	"github.com/google/wire"
)

var UserServiceProvider = wire.NewSet(
	user.NewAuthService,
	NewPasswordCipher,
	NewPasswordHasher,
	ProvideQueue,
	NewMailer,
)

// Wrapper providers to keep wire_gen imports minimal
func NewPasswordCipher(cfg config.Manager) security.PasswordCipher {
	return security.NewPasswordCipher(cfg)
}
func NewPasswordHasher(cfg config.Manager) security.PasswordHasher {
	return security.NewPasswordHasher(cfg)
}
func ProvideQueue() _interface.Queue { return noop.New() }
func NewMailer(q _interface.Queue, log logger.Logger, cfg config.Manager) mailer.Mailer {
	return mailer.NewMailer(q, log, cfg)
}
