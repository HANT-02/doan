package services

import (
	"doan/internal/services/user"
	"github.com/google/wire"
)

var UserServiceProvider = wire.NewSet(user.NewAuthService)
