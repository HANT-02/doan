package teacherportal

import (
	"context"
	"time"

	"doan/internal/entities"
	repointerface "doan/internal/repositories/interface"
)

type GetAcademicRecordsByStudentInput struct {
	Actor     Actor
	ClassID   string
	StudentID string
}

type TeacherStudentAcademicRecordHistory struct {
	Student TeacherStudentReference `json:"student"`
	Records []TeacherAcademicRecord `json:"records"`
}

type GetAcademicRecordsByStudentOutput struct {
	TeacherID string                              `json:"teacher_id"`
	ClassID   string                              `json:"class_id"`
	StudentID string                              `json:"student_id"`
	History   TeacherStudentAcademicRecordHistory `json:"history"`
}

type GetAcademicRecordsByStudentUseCase interface {
	Execute(ctx context.Context, input GetAcademicRecordsByStudentInput) (*GetAcademicRecordsByStudentOutput, error)
}

type getAcademicRecordsByStudentUseCase struct {
	teacherRepo repointerface.TeacherRepository
	recordRepo  repointerface.AcademicRecordRepository
}

func NewGetAcademicRecordsByStudentUseCase(
	teacherRepo repointerface.TeacherRepository,
	recordRepo repointerface.AcademicRecordRepository,
) GetAcademicRecordsByStudentUseCase {
	return &getAcademicRecordsByStudentUseCase{
		teacherRepo: teacherRepo,
		recordRepo:  recordRepo,
	}
}

func (uc *getAcademicRecordsByStudentUseCase) Execute(ctx context.Context, input GetAcademicRecordsByStudentInput) (*GetAcademicRecordsByStudentOutput, error) {
	if input.Actor.Role != "TEACHER" {
		return nil, ErrTeacherAccessDenied
	}

	teacher, err := resolveTeacherByEmail(ctx, uc.teacherRepo, input.Actor.Email)
	if err != nil {
		return nil, err
	}

	teacherLessons, err := uc.teacherRepo.GetTeacherLessons(ctx, teacher.ID, time.Time{}, time.Time{})
	if err != nil {
		return nil, err
	}

	allowedLessonIDs := make(map[string]entities.Lesson)
	for _, lesson := range teacherLessons {
		if lesson.ClassID != input.ClassID {
			continue
		}
		allowedLessonIDs[lesson.ID] = lesson
	}
	if len(allowedLessonIDs) == 0 {
		return nil, ErrTeacherAccessDenied
	}

	records, err := uc.recordRepo.ListByStudentID(ctx, input.StudentID)
	if err != nil {
		return nil, err
	}

	history := make([]TeacherAcademicRecord, 0)
	var student TeacherStudentReference
	for _, record := range records {
		lesson := record.LessonSummary.Lesson
		if _, ok := allowedLessonIDs[lesson.ID]; !ok {
			continue
		}
		if student.ID == "" {
			student = TeacherStudentReference{
				ID:       record.Student.ID,
				Code:     record.Student.Code,
				FullName: record.Student.FullName,
			}
		}
		recordID := record.ID
		summaryID := record.LessonSummaryID
		createdAt := record.CreatedAt
		updatedAt := record.UpdatedAt
		history = append(history, TeacherAcademicRecord{
			RecordID:           &recordID,
			LessonSummaryID:    &summaryID,
			Student:            TeacherStudentReference{ID: record.Student.ID, Code: record.Student.Code, FullName: record.Student.FullName},
			HomeworkCompleted:  record.HomeworkCompleted,
			HomeworkScore:      record.HomeworkScore,
			AttitudeRating:     record.AttitudeRating,
			ParticipationScore: record.ParticipationScore,
			PersonalComment:    record.PersonalComment,
			TotalScore:         record.TotalScore,
			IsCompleted:        record.IsCompleted,
			CreatedAt:          &createdAt,
			UpdatedAt:          &updatedAt,
		})
	}

	if student.ID == "" {
		return &GetAcademicRecordsByStudentOutput{
			TeacherID: teacher.ID,
			ClassID:   input.ClassID,
			StudentID: input.StudentID,
			History: TeacherStudentAcademicRecordHistory{
				Student: TeacherStudentReference{ID: input.StudentID},
				Records: []TeacherAcademicRecord{},
			},
		}, nil
	}

	return &GetAcademicRecordsByStudentOutput{
		TeacherID: teacher.ID,
		ClassID:   input.ClassID,
		StudentID: input.StudentID,
		History: TeacherStudentAcademicRecordHistory{
			Student: student,
			Records: history,
		},
	}, nil
}
