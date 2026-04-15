package lessonactivity

import (
	"context"

	"doan/internal/entities"
	repointerface "doan/internal/repositories/interface"
	"doan/pkg/logger"
)

type GetLessonAttendanceInput struct {
	LessonID string
	Actor    LessonActor
}

type GetLessonAttendanceOutput struct {
	Lesson  entities.Lesson        `json:"lesson"`
	Records []AttendanceRecordItem `json:"records"`
}

type GetLessonAttendanceUseCase interface {
	Execute(ctx context.Context, input GetLessonAttendanceInput) (*GetLessonAttendanceOutput, error)
}

type getLessonAttendanceUseCase struct {
	lessonRepo     repointerface.LessonRepository
	teacherRepo    repointerface.TeacherRepository
	enrollmentRepo repointerface.EnrollmentRepository
	attendanceRepo repointerface.AttendanceRepository
}

func NewGetLessonAttendanceUseCase(
	lessonRepo repointerface.LessonRepository,
	teacherRepo repointerface.TeacherRepository,
	enrollmentRepo repointerface.EnrollmentRepository,
	attendanceRepo repointerface.AttendanceRepository,
) GetLessonAttendanceUseCase {
	return &getLessonAttendanceUseCase{
		lessonRepo:     lessonRepo,
		teacherRepo:    teacherRepo,
		enrollmentRepo: enrollmentRepo,
		attendanceRepo: attendanceRepo,
	}
}

func (uc *getLessonAttendanceUseCase) Execute(ctx context.Context, input GetLessonAttendanceInput) (*GetLessonAttendanceOutput, error) {
	ctxLogger := logger.NewLogger(ctx)

	lesson, err := authorizeLessonAccess(ctx, uc.lessonRepo, uc.teacherRepo, input.LessonID, input.Actor)
	if err != nil {
		ctxLogger.Errorf("Failed to authorize lesson attendance access %s: %v", input.LessonID, err)
		return nil, err
	}

	enrollments, err := uc.enrollmentRepo.ListByClassID(ctx, lesson.ClassID)
	if err != nil {
		ctxLogger.Errorf("Failed to load class roster for lesson %s: %v", input.LessonID, err)
		return nil, err
	}

	existingRecords, err := uc.attendanceRepo.ListByLessonID(ctx, input.LessonID)
	if err != nil {
		ctxLogger.Errorf("Failed to load attendance for lesson %s: %v", input.LessonID, err)
		return nil, err
	}

	recordMap := make(map[string]entities.Attendance, len(existingRecords))
	for _, record := range existingRecords {
		recordMap[record.StudentID] = record
	}

	records := make([]AttendanceRecordItem, 0, len(enrollments))
	for _, enrollment := range enrollments {
		item := AttendanceRecordItem{
			Student: enrollment.Student,
			Status:  0,
			Note:    "",
		}

		if record, exists := recordMap[enrollment.StudentID]; exists {
			recordCopy := record
			item.Attendance = &recordCopy
			item.Status = record.Status
			item.Note = record.Note
		}

		records = append(records, item)
	}

	return &GetLessonAttendanceOutput{
		Lesson:  *lesson,
		Records: records,
	}, nil
}
