package usecases

import (
	"doan/internal/usecases/user"
	"github.com/google/wire"
)

var UserUseCaseProviders = wire.NewSet(
	user.NewGetUserByIdUseCase,
	user.NewLoginUseCase,
	user.NewLogoutUseCase,
	user.NewRefreshTokenUseCase,
	user.NewRegisterUseCase,
	user.NewForgotPasswordUseCase,
	user.NewResetPasswordUseCase,
	user.NewChangePasswordUseCase,
)
