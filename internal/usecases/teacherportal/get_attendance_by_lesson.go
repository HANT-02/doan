package teacherportal

import (
	"context"
	"time"

	"doan/internal/entities"
	lessonactivity "doan/internal/usecases/lessonactivity"
)

type GetAttendanceByLessonInput struct {
	Actor    Actor
	LessonID string
}

type TeacherAttendanceRecord struct {
	AttendanceID *string          `json:"attendance_id,omitempty"`
	Student      entities.Student `json:"student"`
	Status       int              `json:"status"`
	Note         string           `json:"note"`
	MarkedAt     *time.Time       `json:"marked_at,omitempty"`
}

type GetAttendanceByLessonOutput struct {
	Lesson  entities.Lesson           `json:"lesson"`
	Records []TeacherAttendanceRecord `json:"records"`
}

type GetAttendanceByLessonUseCase interface {
	Execute(ctx context.Context, input GetAttendanceByLessonInput) (*GetAttendanceByLessonOutput, error)
}

type getAttendanceByLessonUseCase struct {
	getLessonAttendanceUseCase lessonactivity.GetLessonAttendanceUseCase
}

func NewGetAttendanceByLessonUseCase(
	getLessonAttendanceUseCase lessonactivity.GetLessonAttendanceUseCase,
) GetAttendanceByLessonUseCase {
	return &getAttendanceByLessonUseCase{
		getLessonAttendanceUseCase: getLessonAttendanceUseCase,
	}
}

func (uc *getAttendanceByLessonUseCase) Execute(ctx context.Context, input GetAttendanceByLessonInput) (*GetAttendanceByLessonOutput, error) {
	output, err := uc.getLessonAttendanceUseCase.Execute(ctx, lessonactivity.GetLessonAttendanceInput{
		LessonID: input.LessonID,
		Actor:    buildLessonActor(input.Actor),
	})
	if err != nil {
		return nil, err
	}

	records := make([]TeacherAttendanceRecord, 0, len(output.Records))
	for _, item := range output.Records {
		record := TeacherAttendanceRecord{
			Student: item.Student,
			Status:  TeacherAttendanceStatusUnmarked,
			Note:    item.Note,
		}
		if item.Attendance != nil {
			attendanceID := item.Attendance.ID
			record.AttendanceID = &attendanceID
			record.MarkedAt = &item.Attendance.MarkedAt
			record.Status = mapInternalAttendanceStatusToTeacher(item.Attendance.Status)
			record.Note = item.Attendance.Note
		}
		records = append(records, record)
	}

	return &GetAttendanceByLessonOutput{
		Lesson:  output.Lesson,
		Records: records,
	}, nil
}
