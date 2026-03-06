package scheduling

import (
	"context"
	"errors"

	schedulingstore "doan/internal/services/scheduling"
)

type CommitPreviewInput struct {
	RunID string
}

type CommitPreviewOutput struct {
	Message          string `json:"message"`
	ScheduledLessons int    `json:"scheduled_lessons"`
	Status           string `json:"status"`
}

type CommitPreviewUseCase interface {
	Execute(ctx context.Context, input CommitPreviewInput) (*CommitPreviewOutput, error)
}

type commitPreviewUseCase struct {
	store schedulingstore.PreviewStore[PreviewResult]
}

func NewCommitPreviewUseCase(store schedulingstore.PreviewStore[PreviewResult]) CommitPreviewUseCase {
	return &commitPreviewUseCase{store: store}
}

func (uc *commitPreviewUseCase) Execute(_ context.Context, input CommitPreviewInput) (*CommitPreviewOutput, error) {
	if input.RunID == "" {
		return nil, errors.New("run_id is required")
	}

	preview, ok := uc.store.Get(input.RunID)
	if !ok {
		return nil, errors.New("preview run not found")
	}

	return &CommitPreviewOutput{
		Message:          "Scheduling commit scaffold executed. Persisting lessons is TODO.",
		ScheduledLessons: len(preview.Assignments),
		Status:           preview.Status,
	}, nil
}
