package user

import (
	"context"
	_interface "doan/internal/repositories/interface"
	"time"
)

// GetUserByIdInput DTO Input for GetUserById UseCase
type GetUserByIdInput struct {
	ID int `json:"id"`
}

// GetUserByIdOutput DTO Output for GetUserById UseCase
type GetUserByIdOutput struct {
	ID        string     `json:"id"`
	UserName  string     `json:"user_name"`
	Status    string     `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

// GetUserByIdUseCase defines an interface for retrieving a user by their ID.
type GetUserByIdUseCase interface {
	// Execute retrieves a user based on the provided ID.
	// It takes a context and GetUserByIdInput containing the user ID,
	// and returns a GetUserByIdOutput containing user details and an error if the retrieval fails.
	Execute(ctx context.Context, input GetUserByIdInput) (*GetUserByIdOutput, error)
}

type getUserByIdUseCase struct {
	repo _interface.UserRepository
}

func NewGetUserByIdUseCase(repo _interface.UserRepository) GetUserByIdUseCase {
	return &getUserByIdUseCase{repo: repo}
}

func (u *getUserByIdUseCase) Execute(ctx context.Context, input GetUserByIdInput) (*GetUserByIdOutput, error) {
	user, err := u.repo.GetByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	result := &GetUserByIdOutput{
		ID:        user.ID,
		UserName:  user.FullName,
		Status:    "active",
		CreatedAt: user.CreatedAt,
		UpdatedAt: &user.UpdatedAt,
	}

	return result, nil
}
