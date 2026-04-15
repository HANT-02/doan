package repositoryinterface

import (
	"context"
	"time"

	"doan/internal/entities"
	"doan/internal/repositories"
)

type LeaveRequestFilter struct {
	StudentID string
	ClassID   string
	ClassIDs  []string
	Status    string
	DateFrom  *time.Time
	DateTo    *time.Time
}

type LeaveRequestRepository interface {
	repositories.BaseRepository[entities.LeaveRequest]
	GetWithRelations(ctx context.Context, id string) (*entities.LeaveRequest, error)
	ListWithRelations(ctx context.Context, filter LeaveRequestFilter) ([]entities.LeaveRequest, error)
}
