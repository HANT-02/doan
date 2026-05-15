package infrastructure

import (
	"doan/internal/infrastructure/database"
	"doan/internal/infrastructure/database/postgres"
	"doan/internal/infrastructure/database/postgres/implement"
	_interface "doan/internal/infrastructure/queue/interface"
	"doan/internal/infrastructure/queue/noop"

	"github.com/google/wire"
)

// InfrastructureProviders provides all infrastructure dependencies
// Including: Database, Queue, External services
var InfrastructureProviders = wire.NewSet(
	database.ProvideDB,
	postgres.NewMigration,
	postgres.NewUnitOfWork,
	implement.NewUserRepository,
	implement.NewPasswordResetRepository,
	implement.NewTeacherRepository,
	implement.NewCampusRepository,
	implement.NewCampusTravelTimeRepository,
	implement.NewRoomRepository,
	implement.NewShiftRepository,
	implement.NewClassRepository,
	implement.NewStudentRepository,
	implement.NewLessonRepository,
	implement.NewCourseRepository,
	implement.NewSkillRepository,
	implement.NewProgramRepository,
	implement.NewEnrollmentRepository,
	implement.NewAttendanceRepository,
	implement.NewLessonSummaryRepository,
	implement.NewAcademicRecordRepository,
	implement.NewLeaveRequestRepository,
	implement.NewMaterialRepository,
	implement.NewLabelRepository,
	implement.NewAuditLogRepository,
	implement.NewApprovalDecisionRepository,
	implement.NewClassScheduleRepository,

	// Queue infrastructure
	ProvideQueue,
)

// ProvideQueue provides queue implementation (noop for now)
func ProvideQueue() _interface.Queue {
	return noop.New()
}
