package room

import (
	"context"

	"doan/internal/entities"
	repointerface "doan/internal/repositories/interface"
	"doan/pkg/logger"
)

type UpdateRoomInput struct {
	ID       string
	Name     string
	Capacity int
	Location *string
	Status   string
}

type UpdateRoomOutput struct {
	Room *entities.Room
}

type UpdateRoomUseCase interface {
	Execute(ctx context.Context, input UpdateRoomInput) (*UpdateRoomOutput, error)
}

type updateRoomUseCase struct {
	roomRepo repointerface.RoomRepository
}

func NewUpdateRoomUseCase(roomRepo repointerface.RoomRepository) UpdateRoomUseCase {
	return &updateRoomUseCase{
		roomRepo: roomRepo,
	}
}

func (uc *updateRoomUseCase) Execute(ctx context.Context, input UpdateRoomInput) (*UpdateRoomOutput, error) {
	ctxLogger := logger.NewLogger(ctx)

	// Fetch existing room to update
	room, err := uc.roomRepo.GetByID(ctx, input.ID)
	if err != nil {
		ctxLogger.Errorf("Room not found: %v", err)
		return nil, err
	}

	updateData := map[string]interface{}{
		"name":     input.Name,
		"capacity": input.Capacity,
	}

	err = uc.roomRepo.Update(ctx, input.ID, updateData)
	if err != nil {
		ctxLogger.Errorf("Failed to update room: %v", err)
		return nil, err
	}

	room.Name = input.Name
	room.Capacity = input.Capacity

	return &UpdateRoomOutput{Room: room}, nil
}
