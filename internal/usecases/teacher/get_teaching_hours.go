package teacher

import (
	"context"
	repointerface "doan/internal/repositories/interface"
	"doan/pkg/logger"
	"errors"
	"time"
)

// GetTeachingHoursStatsInput represents the input for getting teaching hours statistics
type GetTeachingHoursStatsInput struct {
	TeacherID string    `json:"teacher_id"`
	From      time.Time `json:"from"`
	To        time.Time `json:"to"`
	GroupBy   string    `json:"group_by"` // day, week, month
}

// TeachingHoursStat represents teaching hours for a period
type TeachingHoursStat struct {
	Period string  `json:"period"`
	Hours  float64 `json:"hours"`
}

// GetTeachingHoursStatsOutput represents the output after getting teaching hours statistics
type GetTeachingHoursStatsOutput struct {
	TotalHours float64             `json:"total_hours"`
	Breakdown  []TeachingHoursStat `json:"breakdown"`
}

// GetTeachingHoursStatsUseCase defines the interface for getting teaching hours statistics
type GetTeachingHoursStatsUseCase interface {
	Execute(ctx context.Context, input GetTeachingHoursStatsInput) (*GetTeachingHoursStatsOutput, error)
}

type getTeachingHoursStatsUseCase struct {
	teacherRepo repointerface.TeacherRepository
}

// NewGetTeachingHoursStatsUseCase creates a new instance of GetTeachingHoursStatsUseCase
func NewGetTeachingHoursStatsUseCase(teacherRepo repointerface.TeacherRepository) GetTeachingHoursStatsUseCase {
	return &getTeachingHoursStatsUseCase{
		teacherRepo: teacherRepo,
	}
}

func (uc *getTeachingHoursStatsUseCase) Execute(ctx context.Context, input GetTeachingHoursStatsInput) (*GetTeachingHoursStatsOutput, error) {
	ctxLogger := logger.NewLogger(ctx)

	if input.TeacherID == "" {
		return nil, errors.New("teacher ID is required")
	}

	// Validate groupBy
	if input.GroupBy == "" {
		input.GroupBy = "day" // Default to day
	}
	if input.GroupBy != "day" && input.GroupBy != "week" && input.GroupBy != "month" {
		return nil, errors.New("groupBy must be one of: day, week, month")
	}

	// Verify teacher exists
	_, err := uc.teacherRepo.GetByID(ctx, input.TeacherID)
	if err != nil {
		ctxLogger.Errorf("Failed to get teacher: %v", err)
		return nil, err
	}

	// Get stats
	stats, err := uc.teacherRepo.GetTeachingHoursStats(ctx, input.TeacherID, input.From, input.To, input.GroupBy)
	if err != nil {
		ctxLogger.Errorf("Failed to get teaching hours stats: %v", err)
		return nil, err
	}

	// Calculate total hours
	var totalHours float64
	breakdown := make([]TeachingHoursStat, 0, len(stats))
	for _, stat := range stats {
		totalHours += stat.Hours
		breakdown = append(breakdown, TeachingHoursStat{
			Period: stat.Period,
			Hours:  stat.Hours,
		})
	}

	return &GetTeachingHoursStatsOutput{
		TotalHours: totalHours,
		Breakdown:  breakdown,
	}, nil
}
