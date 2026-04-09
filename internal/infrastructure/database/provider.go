package database

import (
	"context"
	"doan/internal/infrastructure/database/postgres"
	"doan/internal/infrastructure/database/postgres/implement"
	"doan/pkg/config"
	"doan/pkg/logger"

	"github.com/google/wire"
	"gorm.io/gorm"
)

var DBProvider = wire.NewSet(
	ProvideDB,
	postgres.NewMigration,
	postgres.NewUnitOfWork,
	implement.NewUserRepository,
	implement.NewPasswordResetRepository,
	implement.NewTeacherRepository,
	implement.NewRoomRepository,
	implement.NewShiftRepository,
	implement.NewClassRepository,
	implement.NewStudentRepository,
	implement.NewLessonRepository,
	implement.NewCourseRepository,
	implement.NewProgramRepository,
	implement.NewEnrollmentRepository,
	implement.NewMaterialRepository,
	implement.NewLabelRepository,
	implement.NewAuditLogRepository,
	implement.NewApprovalDecisionRepository,
)

// ProvideDB wraps GetDBContext and panics on error (for Wire)
func ProvideDB(ctx context.Context, log logger.Logger, cfg config.Manager) *gorm.DB {
	db, err := postgres.GetDBContext(ctx, log, cfg)
	if err != nil {
		panic(err)
	}
	return db
}
