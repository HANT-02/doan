package usecases

import (
	"doan/internal/usecases/account"
	"doan/internal/usecases/class"
	"doan/internal/usecases/course"
	"doan/internal/usecases/lesson"
	"doan/internal/usecases/material"
	"doan/internal/usecases/predictive"
	"doan/internal/usecases/program"
	"doan/internal/usecases/room"
	"doan/internal/usecases/scheduling"
	"doan/internal/usecases/shift"
	"doan/internal/usecases/student"
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

var AccountUseCaseProviders = wire.NewSet(
	account.NewListUsersUseCase,
	account.NewGetUserUseCase,
	account.NewCreateUserUseCase,
	account.NewUpdateUserUseCase,
	account.NewResetUserPasswordUseCase,
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

var ShiftUseCaseProviders = wire.NewSet(
	shift.NewCreateShiftUseCase,
	shift.NewGetShiftUseCase,
	shift.NewUpdateShiftUseCase,
	shift.NewDeleteShiftUseCase,
	shift.NewListShiftsUseCase,
)

var ClassUseCaseProviders = wire.NewSet(
	class.NewCreateClassUseCase,
	class.NewGetClassUseCase,
	class.NewGetClassRosterUseCase,
	class.NewUpdateClassUseCase,
	class.NewDeleteClassUseCase,
	class.NewListClassesUseCase,
	class.NewEnrollStudentsUseCase,
	class.NewRemoveStudentsUseCase,
	class.NewAssignTeacherUseCase,
	class.NewGetClassSchedulesUseCase,
	class.NewCreateClassScheduleUseCase,
	class.NewDeleteClassScheduleUseCase,
)

var StudentUseCaseProviders = wire.NewSet(
	student.NewCreateStudentUseCase,
	student.NewGetStudentUseCase,
	student.NewUpdateStudentUseCase,
	student.NewDeleteStudentUseCase,
	student.NewListStudentsUseCase,
)

var SchedulingUseCaseProviders = wire.NewSet(
	scheduling.NewPreviewUseCase,
	scheduling.NewBenchmarkUseCase,
	scheduling.NewGetPreviewUseCase,
	scheduling.NewCommitPreviewUseCase,
)

var MaterialUseCaseProviders = wire.NewSet(
	material.NewUploadMaterialUseCase,
	material.NewListMaterialsUseCase,
	material.NewGetMaterialUseCase,
	material.NewDownloadMaterialUseCase,
	material.NewReviewMaterialUseCase,
)

var CourseUseCaseProviders = wire.NewSet(
	course.NewCreateCourseUseCase,
	course.NewGetCourseUseCase,
	course.NewUpdateCourseUseCase,
	course.NewDeleteCourseUseCase,
	course.NewListCoursesUseCase,
)

var ProgramUseCaseProviders = wire.NewSet(
	program.NewCreateProgramUseCase,
	program.NewGetProgramUseCase,
	program.NewUpdateProgramUseCase,
	program.NewDeleteProgramUseCase,
	program.NewListProgramsUseCase,
	program.NewAddCoursesUseCase,
	program.NewRemoveCoursesUseCase,
)

var PredictiveUseCaseProviders = wire.NewSet(
	predictive.NewListStudentPredictionsUseCase,
	predictive.NewGetModelMetadataUseCase,
)

var LessonUseCaseProviders = wire.NewSet(
	lesson.NewListLessonsUseCase,
	lesson.NewGetLessonUseCase,
)

var UseCaseProviders = wire.NewSet(
	UserUseCaseProviders,
	AccountUseCaseProviders,
	TeacherUseCaseProviders,
	RoomUseCaseProviders,
	ShiftUseCaseProviders,
	ClassUseCaseProviders,
	StudentUseCaseProviders,
	SchedulingUseCaseProviders,
	MaterialUseCaseProviders,
	CourseUseCaseProviders,
	ProgramUseCaseProviders,
	PredictiveUseCaseProviders,
	LessonUseCaseProviders,
)
