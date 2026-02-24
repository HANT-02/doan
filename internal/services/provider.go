package services

import (
	_interface "doan/internal/infrastructure/queue/interface"
	"doan/internal/services/mailer"
	"doan/internal/services/security"
	"doan/internal/services/user"
	"doan/pkg/config"
	"doan/pkg/logger"
	"github.com/google/wire"
)

// ServiceProviders provides all application services
// Including: Auth, Security, Mailer
var ServiceProviders = wire.NewSet(
	// Auth & User services
	user.NewAuthService,

	// Security services
	NewPasswordCipher,
	NewPasswordHasher,

	// Mailer service
	NewMailer,
)

// Wrapper providers to keep wire_gen imports minimal
func NewPasswordCipher(cfg config.Manager) security.PasswordCipher {
	return security.NewPasswordCipher(cfg)
}

func NewPasswordHasher(cfg config.Manager) security.PasswordHasher {
	return security.NewPasswordHasher(cfg)
}

func NewMailer(q _interface.Queue, log logger.Logger, cfg config.Manager) mailer.Mailer {
	return mailer.NewMailer(q, log, cfg)
}
