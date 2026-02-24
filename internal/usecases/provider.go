package usecases

import (
	"doan/internal/usecases/class"
	"doan/internal/usecases/room"
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

var RoomUseCaseProviders = wire.NewSet(
	room.NewCreateRoomUseCase,
	room.NewGetRoomUseCase,
	room.NewUpdateRoomUseCase,
	room.NewDeleteRoomUseCase,
	room.NewListRoomsUseCase,
)

var ClassUseCaseProviders = wire.NewSet(
	class.NewCreateClassUseCase,
	class.NewGetClassUseCase,
	class.NewUpdateClassUseCase,
	class.NewDeleteClassUseCase,
	class.NewListClassesUseCase,
)

var UseCaseProviders = wire.NewSet(
	UserUseCaseProviders,
	TeacherUseCaseProviders,
	RoomUseCaseProviders,
	ClassUseCaseProviders,
)
