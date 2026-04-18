package teacherportal

import (
	"context"
	"time"

	lessonrecord "doan/internal/usecases/lessonrecord"
)

type TeacherAcademicRecord struct {
	RecordID           *string                 `json:"record_id,omitempty"`
	LessonSummaryID    *string                 `json:"lesson_summary_id,omitempty"`
	Student            TeacherStudentReference `json:"student"`
	HomeworkCompleted  bool                    `json:"homework_completed"`
	HomeworkScore      float64                 `json:"homework_score"`
	AttitudeRating     int                     `json:"attitude_rating"`
	ParticipationScore float64                 `json:"participation_score"`
	PersonalComment    string                  `json:"personal_comment"`
	TotalScore         float64                 `json:"total_score"`
	IsCompleted        bool                    `json:"is_completed"`
	CreatedAt          *time.Time              `json:"created_at,omitempty"`
	UpdatedAt          *time.Time              `json:"updated_at,omitempty"`
}

type TeacherStudentReference struct {
	ID       string `json:"id"`
	Code     string `json:"code"`
	FullName string `json:"full_name"`
}

type GetAcademicRecordsByLessonInput struct {
	Actor    Actor
	LessonID string
}

type GetAcademicRecordsByLessonOutput struct {
	Lesson  TeacherLessonItem       `json:"lesson"`
	Records []TeacherAcademicRecord `json:"records"`
}

type GetAcademicRecordsByLessonUseCase interface {
	Execute(ctx context.Context, input GetAcademicRecordsByLessonInput) (*GetAcademicRecordsByLessonOutput, error)
}

type getAcademicRecordsByLessonUseCase struct {
	getLessonAcademicRecordsUseCase lessonrecord.GetLessonAcademicRecordsUseCase
}

func NewGetAcademicRecordsByLessonUseCase(
	getLessonAcademicRecordsUseCase lessonrecord.GetLessonAcademicRecordsUseCase,
) GetAcademicRecordsByLessonUseCase {
	return &getAcademicRecordsByLessonUseCase{
		getLessonAcademicRecordsUseCase: getLessonAcademicRecordsUseCase,
	}
}

func (uc *getAcademicRecordsByLessonUseCase) Execute(ctx context.Context, input GetAcademicRecordsByLessonInput) (*GetAcademicRecordsByLessonOutput, error) {
	output, err := uc.getLessonAcademicRecordsUseCase.Execute(ctx, lessonrecord.GetLessonAcademicRecordsInput{
		LessonID: input.LessonID,
		Actor:    buildLessonRecordActor(input.Actor),
	})
	if err != nil {
		return nil, err
	}

	records := make([]TeacherAcademicRecord, 0, len(output.Records))
	for _, item := range output.Records {
		record := TeacherAcademicRecord{
			Student: TeacherStudentReference{
				ID:       item.Student.ID,
				Code:     item.Student.Code,
				FullName: item.Student.FullName,
			},
		}
		if item.Record != nil {
			recordID := item.Record.ID
			summaryID := item.Record.LessonSummaryID
			createdAt := item.Record.CreatedAt
			updatedAt := item.Record.UpdatedAt
			record.RecordID = &recordID
			record.LessonSummaryID = &summaryID
			record.HomeworkCompleted = item.Record.HomeworkCompleted
			record.HomeworkScore = item.Record.HomeworkScore
			record.AttitudeRating = item.Record.AttitudeRating
			record.ParticipationScore = item.Record.ParticipationScore
			record.PersonalComment = item.Record.PersonalComment
			record.TotalScore = item.Record.TotalScore
			record.IsCompleted = item.Record.IsCompleted
			record.CreatedAt = &createdAt
			record.UpdatedAt = &updatedAt
		}
		records = append(records, record)
	}

	result := &GetAcademicRecordsByLessonOutput{
		Lesson: TeacherLessonItem{
			ID:        output.Lesson.ID,
			ClassID:   output.Lesson.ClassID,
			ClassName: output.Lesson.Class.Name,
			ClassCode: output.Lesson.Class.Code,
			DateStart: output.Lesson.DateStart,
			DateEnd:   output.Lesson.DateEnd,
			Notes:     output.Lesson.Notes,
		},
		Records: records,
	}
	if output.Lesson.RoomID != nil {
		result.Lesson.RoomID = output.Lesson.RoomID
		roomName := output.Lesson.Room.Name
		result.Lesson.RoomName = &roomName
	}

	return result, nil
}
