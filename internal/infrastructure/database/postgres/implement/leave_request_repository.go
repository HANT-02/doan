package implement

import (
	"context"

	"doan/internal/entities"
	"doan/internal/infrastructure/database/postgres"
	"doan/internal/repositories"
	repointerface "doan/internal/repositories/interface"
	"doan/pkg/base_struct"
	"doan/pkg/config"
	"doan/pkg/logger"

	"gorm.io/gorm"
)

type leaveRequestRepository struct {
	base_struct.BaseDependency
	repositories.BaseRepository[entities.LeaveRequest]
	db *gorm.DB
}

func NewLeaveRequestRepository(
	db *gorm.DB,
	log logger.Logger,
	manager config.Manager,
) repointerface.LeaveRequestRepository {
	modelRepo := postgres.NewBaseRepository[entities.LeaveRequest](log, manager, db, "leave_requests")
	return &leaveRequestRepository{
		BaseDependency: base_struct.BaseDependency{
			Log:           log,
			ConfigManager: manager,
		},
		BaseRepository: modelRepo,
		db:             db,
	}
}

func (r *leaveRequestRepository) GetWithRelations(ctx context.Context, id string) (*entities.LeaveRequest, error) {
	var entity entities.LeaveRequest
	err := r.db.WithContext(ctx).
		Preload("Student").
		Preload("Class").
		Preload("Lesson").
		Preload("ApprovedBy").
		First(&entity, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &entity, nil
}

func (r *leaveRequestRepository) ListWithRelations(ctx context.Context, filter repointerface.LeaveRequestFilter) ([]entities.LeaveRequest, error) {
	records := make([]entities.LeaveRequest, 0)
	query := r.db.WithContext(ctx).
		Preload("Student").
		Preload("Class").
		Preload("Lesson").
		Preload("ApprovedBy")

	if filter.StudentID != "" {
		query = query.Where("student_id = ?", filter.StudentID)
	}
	if filter.ClassID != "" {
		query = query.Where("class_id = ?", filter.ClassID)
	}
	if len(filter.ClassIDs) > 0 {
		query = query.Where("class_id IN ?", filter.ClassIDs)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.DateFrom != nil {
		query = query.Where("apply_date >= ?", *filter.DateFrom)
	}
	if filter.DateTo != nil {
		query = query.Where("apply_date <= ?", *filter.DateTo)
	}

	err := query.Order("apply_date DESC").Find(&records).Error
	if err != nil {
		return nil, err
	}
	return records, nil
}
