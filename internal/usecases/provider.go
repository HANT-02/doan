package usecases

import (
	"doan/internal/usecases/account"
	"doan/internal/usecases/class"
	"doan/internal/usecases/course"
	"doan/internal/usecases/leaveflow"
	"doan/internal/usecases/lesson"
	"doan/internal/usecases/lessonactivity"
	"doan/internal/usecases/lessonrecord"
	"doan/internal/usecases/material"
	"doan/internal/usecases/predictive"
	"doan/internal/usecases/program"
	"doan/internal/usecases/room"
	"doan/internal/usecases/scheduling"
	"doan/internal/usecases/shift"
	"doan/internal/usecases/student"
	"doan/internal/usecases/studentportal"
	"doan/internal/usecases/teacher"
	"doan/internal/usecases/teacherportal"
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
	user.NewConfirmChangePasswordOTPUseCase,
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

var TeacherPortalUseCaseProviders = wire.NewSet(
	teacherportal.NewGetTeacherLessonsUseCase,
	teacherportal.NewGetAttendanceByLessonUseCase,
	teacherportal.NewMarkAttendanceUseCase,
	teacherportal.NewSubmitLessonAttendanceUseCase,
	teacherportal.NewGetAttendanceSummaryByStudentUseCase,
	teacherportal.NewGetLessonSummaryUseCase,
	teacherportal.NewUpsertLessonSummaryUseCase,
	teacherportal.NewGetAcademicRecordsByLessonUseCase,
	teacherportal.NewUpsertAcademicRecordUseCase,
	teacherportal.NewFinalizeAcademicRecordUseCase,
	teacherportal.NewGetAcademicRecordsByStudentUseCase,
	teacherportal.NewListLeaveRequestsForTeacherUseCase,
	teacherportal.NewApproveLeaveRequestUseCase,
	teacherportal.NewRejectLeaveRequestUseCase,
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

var StudentPortalUseCaseProviders = wire.NewSet(
	studentportal.NewGetStudentTimetableUseCase,
	studentportal.NewGetMyAttendanceUseCase,
	studentportal.NewGetMyAcademicRecordsUseCase,
	studentportal.NewGetMyAtRiskPredictionUseCase,
	studentportal.NewListMyLeaveRequestsUseCase,
	studentportal.NewCreateMyLeaveRequestUseCase,
	studentportal.NewCancelMyLeaveRequestUseCase,
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
	predictive.NewTrainAtRiskFromDBUseCase,
)

var LessonUseCaseProviders = wire.NewSet(
	lesson.NewListLessonsUseCase,
	lesson.NewGetLessonUseCase,
)

var LessonActivityUseCaseProviders = wire.NewSet(
	lessonactivity.NewGetLessonAttendanceUseCase,
	lessonactivity.NewUpsertLessonAttendanceUseCase,
	lessonactivity.NewGetLessonSummaryUseCase,
	lessonactivity.NewUpsertLessonSummaryUseCase,
)

var LessonRecordUseCaseProviders = wire.NewSet(
	lessonrecord.NewGetLessonAcademicRecordsUseCase,
	lessonrecord.NewUpsertLessonAcademicRecordsUseCase,
	lessonrecord.NewFinalizeLessonAcademicRecordsUseCase,
	lessonrecord.NewListMyAcademicRecordsUseCase,
)

var LeaveFlowUseCaseProviders = wire.NewSet(
	leaveflow.NewListLeaveRequestsUseCase,
	leaveflow.NewCreateLeaveRequestUseCase,
	leaveflow.NewUpdateLeaveRequestStatusUseCase,
	leaveflow.NewCancelLeaveRequestUseCase,
)

var UseCaseProviders = wire.NewSet(
	UserUseCaseProviders,
	AccountUseCaseProviders,
	TeacherUseCaseProviders,
	TeacherPortalUseCaseProviders,
	RoomUseCaseProviders,
	ShiftUseCaseProviders,
	ClassUseCaseProviders,
	StudentUseCaseProviders,
	StudentPortalUseCaseProviders,
	SchedulingUseCaseProviders,
	MaterialUseCaseProviders,
	CourseUseCaseProviders,
	ProgramUseCaseProviders,
	PredictiveUseCaseProviders,
	LessonUseCaseProviders,
	LessonActivityUseCaseProviders,
	LessonRecordUseCaseProviders,
	LeaveFlowUseCaseProviders,
)
