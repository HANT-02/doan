package room

import (
	"context"

	"doan/internal/entities"
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

	room, err := uc.roomRepo.GetByID(ctx, input.ID)
	if err != nil {
		ctxLogger.Errorf("Failed to get room: %v", err)
		return nil, err
	}

	return &GetRoomOutput{Room: room}, nil
}
