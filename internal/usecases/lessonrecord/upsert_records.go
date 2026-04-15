package lessonrecord

import (
	"context"
	"strings"
	"time"

	"doan/internal/entities"
	repointerface "doan/internal/repositories/interface"
	"doan/pkg/logger"
)

type UpsertLessonAcademicRecordRow struct {
	StudentID          string  `json:"student_id"`
	HomeworkCompleted  bool    `json:"homework_completed"`
	HomeworkScore      float64 `json:"homework_score"`
	AttitudeRating     int     `json:"attitude_rating"`
	ParticipationScore float64 `json:"participation_score"`
	PersonalComment    string  `json:"personal_comment"`
}

type UpsertLessonAcademicRecordsInput struct {
	LessonID string
	Actor    LessonActor
	Records  []UpsertLessonAcademicRecordRow
}

type UpsertLessonAcademicRecordsOutput struct {
	SavedCount int `json:"saved_count"`
}

type UpsertLessonAcademicRecordsUseCase interface {
	Execute(ctx context.Context, input UpsertLessonAcademicRecordsInput) (*UpsertLessonAcademicRecordsOutput, error)
}

type upsertLessonAcademicRecordsUseCase struct {
	lessonRepo     repointerface.LessonRepository
	teacherRepo    repointerface.TeacherRepository
	enrollmentRepo repointerface.EnrollmentRepository
	summaryRepo    repointerface.LessonSummaryRepository
	recordRepo     repointerface.AcademicRecordRepository
}

func NewUpsertLessonAcademicRecordsUseCase(
	lessonRepo repointerface.LessonRepository,
	teacherRepo repointerface.TeacherRepository,
	enrollmentRepo repointerface.EnrollmentRepository,
	summaryRepo repointerface.LessonSummaryRepository,
	recordRepo repointerface.AcademicRecordRepository,
) UpsertLessonAcademicRecordsUseCase {
	return &upsertLessonAcademicRecordsUseCase{
		lessonRepo:     lessonRepo,
		teacherRepo:    teacherRepo,
		enrollmentRepo: enrollmentRepo,
		summaryRepo:    summaryRepo,
		recordRepo:     recordRepo,
	}
}

func (uc *upsertLessonAcademicRecordsUseCase) Execute(ctx context.Context, input UpsertLessonAcademicRecordsInput) (*UpsertLessonAcademicRecordsOutput, error) {
	ctxLogger := logger.NewLogger(ctx)

	lesson, err := authorizeLessonAccess(ctx, uc.lessonRepo, uc.teacherRepo, input.LessonID, input.Actor)
	if err != nil {
		ctxLogger.Errorf("Failed to authorize lesson academic record update %s: %v", input.LessonID, err)
		return nil, err
	}

	enrollments, err := uc.enrollmentRepo.ListByClassID(ctx, lesson.ClassID)
	if err != nil {
		return nil, err
	}

	allowedStudentIDs := map[string]struct{}{}
	for _, enrollment := range enrollments {
		allowedStudentIDs[enrollment.StudentID] = struct{}{}
	}

	summary, err := ensureLessonSummary(ctx, uc.summaryRepo, input.LessonID, input.Actor)
	if err != nil {
		return nil, err
	}

	savedCount := 0
	for _, row := range input.Records {
		if _, ok := allowedStudentIDs[row.StudentID]; !ok {
			return nil, ErrInvalidRecordRow
		}
		existing, err := uc.recordRepo.GetByLessonSummaryAndStudent(ctx, summary.ID, row.StudentID)
		if err != nil {
			return nil, err
		}

		updateData := map[string]interface{}{
			"homework_completed":  row.HomeworkCompleted,
			"homework_score":      row.HomeworkScore,
			"attitude_rating":     row.AttitudeRating,
			"participation_score": row.ParticipationScore,
			"personal_comment":    strings.TrimSpace(row.PersonalComment),
			"total_score":         calculateTotalScore(row.HomeworkScore, row.ParticipationScore, row.AttitudeRating),
			"is_completed":        false,
		}

		if existing == nil {
			_, err = uc.recordRepo.Create(ctx, &entities.AcademicRecord{
				LessonSummaryID:    summary.ID,
				StudentID:          row.StudentID,
				HomeworkCompleted:  row.HomeworkCompleted,
				HomeworkScore:      row.HomeworkScore,
				AttitudeRating:     row.AttitudeRating,
				ParticipationScore: row.ParticipationScore,
				PersonalComment:    strings.TrimSpace(row.PersonalComment),
				TotalScore:         calculateTotalScore(row.HomeworkScore, row.ParticipationScore, row.AttitudeRating),
				IsCompleted:        false,
				CreatedAt:          time.Now(),
				UpdatedAt:          time.Now(),
			})
		} else {
			err = uc.recordRepo.Update(ctx, existing.ID, updateData)
		}
		if err != nil {
			ctxLogger.Errorf("Failed to save academic record lesson %s student %s: %v", input.LessonID, row.StudentID, err)
			return nil, err
		}
		savedCount++
	}

	return &UpsertLessonAcademicRecordsOutput{SavedCount: savedCount}, nil
}
