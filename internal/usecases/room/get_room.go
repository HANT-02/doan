package room

import (
	"context"
	"errors"

	"doan/internal/entities"
	"doan/internal/repositories"
	repointerface "doan/internal/repositories/interface"
	"doan/pkg/logger"
)

type GetRoomInput struct {
	ID string
}

type GetRoomOutput struct {
	Room *entities.Room
}

type GetRoomUseCase interface {
	Execute(ctx context.Context, input GetRoomInput) (*GetRoomOutput, error)
}

type getRoomUseCase struct {
	roomRepo repointerface.RoomRepository
}

func NewGetRoomUseCase(roomRepo repointerface.RoomRepository) GetRoomUseCase {
	return &getRoomUseCase{
		roomRepo: roomRepo,
	}
}

func (uc *getRoomUseCase) Execute(ctx context.Context, input GetRoomInput) (*GetRoomOutput, error) {
	ctxLogger := logger.NewLogger(ctx)

	commonCond := repositories.NewCommonCondition()
	commonCond.AddCondition("id", input.ID, repositories.Equal)
	commonCond.SetPreload([]string{"Campus"})

	result, err := uc.roomRepo.GetByCondition(ctx, commonCond)
	if err != nil {
		ctxLogger.Errorf("Failed to get room: %v", err)
		return nil, err
	}
	if result == nil || len(result.Data) == 0 || result.Data[0] == nil {
		return nil, errors.New("room not found")
	}

	return &GetRoomOutput{Room: result.Data[0]}, nil
}
