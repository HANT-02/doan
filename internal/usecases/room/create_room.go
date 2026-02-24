package room

import (
	"context"

	"doan/internal/entities"
	repointerface "doan/internal/repositories/interface"
	"doan/pkg/logger"
)

type CreateRoomInput struct {
	Name     string
	Capacity int
	Location *string
	Status   string
}

type CreateRoomOutput struct {
	Room *entities.Room
}

type CreateRoomUseCase interface {
	Execute(ctx context.Context, input CreateRoomInput) (*CreateRoomOutput, error)
}

type createRoomUseCase struct {
	roomRepo repointerface.RoomRepository
}

func NewCreateRoomUseCase(roomRepo repointerface.RoomRepository) CreateRoomUseCase {
	return &createRoomUseCase{
		roomRepo: roomRepo,
	}
}

func (uc *createRoomUseCase) Execute(ctx context.Context, input CreateRoomInput) (*CreateRoomOutput, error) {
	ctxLogger := logger.NewLogger(ctx)

	room := &entities.Room{
		Name:     input.Name,
		Capacity: input.Capacity,
	}

	createdRoom, err := uc.roomRepo.Create(ctx, room)
	if err != nil {
		ctxLogger.Errorf("Failed to create room: %v", err)
		return nil, err
	}

	return &CreateRoomOutput{Room: createdRoom}, nil
}
