package lessonrecord

import (
	"context"

	"doan/internal/entities"
	repointerface "doan/internal/repositories/interface"
	"doan/pkg/logger"
)

type GetLessonAcademicRecordsInput struct {
	LessonID string
	Actor    LessonActor
}

type GetLessonAcademicRecordsOutput struct {
	Lesson  entities.Lesson            `json:"lesson"`
	Records []LessonAcademicRecordItem `json:"records"`
}

type GetLessonAcademicRecordsUseCase interface {
	Execute(ctx context.Context, input GetLessonAcademicRecordsInput) (*GetLessonAcademicRecordsOutput, error)
}

type getLessonAcademicRecordsUseCase struct {
	lessonRepo     repointerface.LessonRepository
	teacherRepo    repointerface.TeacherRepository
	enrollmentRepo repointerface.EnrollmentRepository
	summaryRepo    repointerface.LessonSummaryRepository
	recordRepo     repointerface.AcademicRecordRepository
}

func NewGetLessonAcademicRecordsUseCase(
	lessonRepo repointerface.LessonRepository,
	teacherRepo repointerface.TeacherRepository,
	enrollmentRepo repointerface.EnrollmentRepository,
	summaryRepo repointerface.LessonSummaryRepository,
	recordRepo repointerface.AcademicRecordRepository,
) GetLessonAcademicRecordsUseCase {
	return &getLessonAcademicRecordsUseCase{
		lessonRepo:     lessonRepo,
		teacherRepo:    teacherRepo,
		enrollmentRepo: enrollmentRepo,
		summaryRepo:    summaryRepo,
		recordRepo:     recordRepo,
	}
}

func (uc *getLessonAcademicRecordsUseCase) Execute(ctx context.Context, input GetLessonAcademicRecordsInput) (*GetLessonAcademicRecordsOutput, error) {
	ctxLogger := logger.NewLogger(ctx)

	lesson, err := authorizeLessonAccess(ctx, uc.lessonRepo, uc.teacherRepo, input.LessonID, input.Actor)
	if err != nil {
		ctxLogger.Errorf("Failed to authorize lesson academic records %s: %v", input.LessonID, err)
		return nil, err
	}

	enrollments, err := uc.enrollmentRepo.ListByClassID(ctx, lesson.ClassID)
	if err != nil {
		ctxLogger.Errorf("Failed to load roster for lesson academic records %s: %v", input.LessonID, err)
		return nil, err
	}

	recordMap := map[string]entities.AcademicRecord{}
	summary, err := uc.summaryRepo.GetByLessonID(ctx, input.LessonID)
	if err != nil {
		return nil, err
	}
	if summary != nil {
		existingRecords, listErr := uc.recordRepo.ListByLessonSummaryID(ctx, summary.ID)
		if listErr != nil {
			ctxLogger.Errorf("Failed to list academic records for lesson summary %s: %v", summary.ID, listErr)
			return nil, listErr
		}
		for _, record := range existingRecords {
			recordMap[record.StudentID] = record
		}
	}

	items := make([]LessonAcademicRecordItem, 0, len(enrollments))
	for _, enrollment := range enrollments {
		item := LessonAcademicRecordItem{Student: enrollment.Student}
		if record, exists := recordMap[enrollment.StudentID]; exists {
			recordCopy := record
			item.Record = &recordCopy
		}
		items = append(items, item)
	}

	return &GetLessonAcademicRecordsOutput{
		Lesson:  *lesson,
		Records: items,
	}, nil
}
