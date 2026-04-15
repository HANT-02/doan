package lessonactivity

import (
	"context"
	"strings"
	"time"

	"doan/internal/entities"
	repointerface "doan/internal/repositories/interface"
	"doan/pkg/logger"
)

type UpsertLessonAttendanceRecord struct {
	StudentID string `json:"student_id"`
	Status    int    `json:"status"`
	Note      string `json:"note"`
}

type UpsertLessonAttendanceInput struct {
	LessonID string
	Actor    LessonActor
	Records  []UpsertLessonAttendanceRecord
}

type UpsertLessonAttendanceOutput struct {
	SavedCount int `json:"saved_count"`
}

type UpsertLessonAttendanceUseCase interface {
	Execute(ctx context.Context, input UpsertLessonAttendanceInput) (*UpsertLessonAttendanceOutput, error)
}

type upsertLessonAttendanceUseCase struct {
	lessonRepo     repointerface.LessonRepository
	teacherRepo    repointerface.TeacherRepository
	enrollmentRepo repointerface.EnrollmentRepository
	attendanceRepo repointerface.AttendanceRepository
}

func NewUpsertLessonAttendanceUseCase(
	lessonRepo repointerface.LessonRepository,
	teacherRepo repointerface.TeacherRepository,
	enrollmentRepo repointerface.EnrollmentRepository,
	attendanceRepo repointerface.AttendanceRepository,
) UpsertLessonAttendanceUseCase {
	return &upsertLessonAttendanceUseCase{
		lessonRepo:     lessonRepo,
		teacherRepo:    teacherRepo,
		enrollmentRepo: enrollmentRepo,
		attendanceRepo: attendanceRepo,
	}
}

func (uc *upsertLessonAttendanceUseCase) Execute(ctx context.Context, input UpsertLessonAttendanceInput) (*UpsertLessonAttendanceOutput, error) {
	ctxLogger := logger.NewLogger(ctx)

	lesson, err := authorizeLessonAccess(ctx, uc.lessonRepo, uc.teacherRepo, input.LessonID, input.Actor)
	if err != nil {
		ctxLogger.Errorf("Failed to authorize attendance update for lesson %s: %v", input.LessonID, err)
		return nil, err
	}

	enrollments, err := uc.enrollmentRepo.ListByClassID(ctx, lesson.ClassID)
	if err != nil {
		ctxLogger.Errorf("Failed to load class roster for attendance save %s: %v", input.LessonID, err)
		return nil, err
	}

	allowedStudentIDs := make(map[string]struct{}, len(enrollments))
	for _, enrollment := range enrollments {
		allowedStudentIDs[enrollment.StudentID] = struct{}{}
	}

	savedCount := 0
	for _, record := range input.Records {
		if _, ok := allowedStudentIDs[record.StudentID]; !ok {
			return nil, ErrInvalidAttendanceRow
		}
		if !IsValidAttendanceStatus(record.Status) {
			return nil, ErrInvalidAttendanceRow
		}

		existing, err := uc.attendanceRepo.GetByLessonAndStudent(ctx, input.LessonID, record.StudentID)
		if err != nil {
			ctxLogger.Errorf("Failed to get attendance row for lesson %s student %s: %v", input.LessonID, record.StudentID, err)
			return nil, err
		}

		now := time.Now()
		if existing == nil {
			_, err = uc.attendanceRepo.Create(ctx, &entities.Attendance{
				LessonID:  input.LessonID,
				StudentID: record.StudentID,
				Status:    record.Status,
				Note:      strings.TrimSpace(record.Note),
				MarkedAt:  now,
				CreatedAt: now,
				UpdatedAt: now,
			})
		} else {
			err = uc.attendanceRepo.Update(ctx, existing.ID, map[string]interface{}{
				"status":    record.Status,
				"note":      strings.TrimSpace(record.Note),
				"marked_at": now,
			})
		}
		if err != nil {
			ctxLogger.Errorf("Failed to save attendance row for lesson %s student %s: %v", input.LessonID, record.StudentID, err)
			return nil, err
		}
		savedCount++
	}

	return &UpsertLessonAttendanceOutput{SavedCount: savedCount}, nil
}
