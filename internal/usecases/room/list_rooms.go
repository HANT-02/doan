package room

import (
	"context"

	"doan/internal/entities"
	"doan/internal/repositories"
	repointerface "doan/internal/repositories/interface"
	"doan/pkg/logger"
)

type ListRoomsInput struct {
	Search    string
	Status    string
	Page      int
	Limit     int
	SortBy    string
	SortOrder string
}

type ListRoomsOutput struct {
	Rooms      []entities.Room
	Pagination struct {
		CurrentPage  int
		ItemsPerPage int
		TotalItems   int64
		TotalPages   int
	}
}

type ListRoomsUseCase interface {
	Execute(ctx context.Context, input ListRoomsInput) (*ListRoomsOutput, error)
}

type listRoomsUseCase struct {
	roomRepo repointerface.RoomRepository
}

func NewListRoomsUseCase(roomRepo repointerface.RoomRepository) ListRoomsUseCase {
	return &listRoomsUseCase{
		roomRepo: roomRepo,
	}
}

func (uc *listRoomsUseCase) Execute(ctx context.Context, input ListRoomsInput) (*ListRoomsOutput, error) {
	ctxLogger := logger.NewLogger(ctx)

	commonCond := repositories.NewCommonCondition()

	if input.Search != "" {
		commonCond.AddCondition("name ILIKE ?", "%"+input.Search+"%", repositories.Like)
	}

	if input.Status != "" {
		commonCond.AddCondition("status", input.Status, repositories.Equal)
	}

	if input.Page > 0 && input.Limit > 0 {
		commonCond.SetPaging(uint64(input.Limit), uint64(input.Page))
	}

	orderBy := "created_at DESC"
	if input.SortBy != "" {
		order := repositories.Asc
		if input.SortOrder == "desc" || input.SortOrder == "DESC" {
			order = repositories.Desc
		}
		orderBy = input.SortBy + " " + order
	}
	commonCond.AddSorting(orderBy, "")

	result, err := uc.roomRepo.GetByCondition(ctx, commonCond)
	if err != nil {
		ctxLogger.Errorf("Failed to list rooms: %v", err)
		return nil, err
	}

	var rooms []entities.Room
	total := int64(0)
	totalPages := 0
	if result != nil {
		for _, ptr := range result.Data {
			rooms = append(rooms, *ptr)
		}
		total = int64(result.Meta.TotalItems)
		totalPages = int(result.Meta.TotalPages)
	}

	return &ListRoomsOutput{
		Rooms: rooms,
		Pagination: struct {
			CurrentPage  int
			ItemsPerPage int
			TotalItems   int64
			TotalPages   int
		}{
			CurrentPage:  input.Page,
			ItemsPerPage: input.Limit,
			TotalItems:   total,
			TotalPages:   totalPages,
		},
	}, nil
}
