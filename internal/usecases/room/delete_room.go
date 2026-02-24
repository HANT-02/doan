package room

import (
	"context"

	repointerface "doan/internal/repositories/interface"
	"doan/pkg/logger"
)

type DeleteRoomInput struct {
	ID string
}

type DeleteRoomOutput struct {
	Message string
}

type DeleteRoomUseCase interface {
	Execute(ctx context.Context, input DeleteRoomInput) (*DeleteRoomOutput, error)
}

type deleteRoomUseCase struct {
	roomRepo repointerface.RoomRepository
}

func NewDeleteRoomUseCase(roomRepo repointerface.RoomRepository) DeleteRoomUseCase {
	return &deleteRoomUseCase{
		roomRepo: roomRepo,
	}
}

func (uc *deleteRoomUseCase) Execute(ctx context.Context, input DeleteRoomInput) (*DeleteRoomOutput, error) {
	ctxLogger := logger.NewLogger(ctx)

	err := uc.roomRepo.SoftDelete(ctx, input.ID)
	if err != nil {
		ctxLogger.Errorf("Failed to delete room: %v", err)
		return nil, err
	}

	return &DeleteRoomOutput{Message: "Room deleted successfully"}, nil
}
