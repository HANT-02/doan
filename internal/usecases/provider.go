package usecases

import (
	"doan/internal/usecases/teacher"
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
	user.NewVerifyOTPUseCase,
)

var TeacherUseCaseProviders = wire.NewSet(
	teacher.NewCreateTeacherUseCase,
	teacher.NewDeleteTeacherUseCase,
	teacher.NewGetTeacherUseCase,
	teacher.NewGetTeachingHoursStatsUseCase,
	teacher.NewGetTeacherTimetableUseCase,
	teacher.NewListTeachersUseCase,
	teacher.NewUpdateTeacherUseCase,
)
