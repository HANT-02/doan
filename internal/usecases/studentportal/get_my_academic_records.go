package studentportal

import (
	"context"
	"sort"
	"strings"
	"time"

	"doan/internal/entities"
	repointerface "doan/internal/repositories/interface"
	"doan/pkg/logger"
)

type StudentAcademicRecordLessonSummary struct {
	ID       string `json:"id"`
	Topic    string `json:"topic"`
	Homework string `json:"homework"`
}

type StudentAcademicRecordLesson struct {
	ID        string                             `json:"id"`
	ClassID   string                             `json:"class_id"`
	ClassName string                             `json:"class_name"`
	ClassCode string                             `json:"class_code"`
	Teacher   StudentAttendanceTeacher           `json:"teacher"`
	RoomID    *string                            `json:"room_id,omitempty"`
	RoomName  *string                            `json:"room_name,omitempty"`
	DateStart time.Time                          `json:"date_start"`
	DateEnd   time.Time                          `json:"date_end"`
	Notes     string                             `json:"notes"`
	Shift     *StudentAttendanceShift            `json:"shift,omitempty"`
	Summary   StudentAcademicRecordLessonSummary `json:"summary"`
}

type StudentAcademicRecordItem struct {
	RecordID           string                      `json:"record_id"`
	LessonSummaryID    string                      `json:"lesson_summary_id"`
	Lesson             StudentAcademicRecordLesson `json:"lesson"`
	HomeworkCompleted  bool                        `json:"homework_completed"`
	HomeworkScore      float64                     `json:"homework_score"`
	AttitudeRating     int                         `json:"attitude_rating"`
	ParticipationScore float64                     `json:"participation_score"`
	PersonalComment    string                      `json:"personal_comment"`
	TotalScore         float64                     `json:"total_score"`
	IsCompleted        bool                        `json:"is_completed"`
	CreatedAt          time.Time                   `json:"created_at"`
	UpdatedAt          time.Time                   `json:"updated_at"`
}

type StudentAcademicClassSummary struct {
	ClassID           string  `json:"class_id"`
	ClassName         string  `json:"class_name"`
	ClassCode         string  `json:"class_code"`
	RecordsCount      int     `json:"records_count"`
	CompletedCount    int     `json:"completed_count"`
	AverageTotalScore float64 `json:"average_total_score"`
}

type GetMyAcademicRecordsInput struct {
	Actor    Actor
	ClassID  string
	DateFrom *time.Time
	DateTo   *time.Time
}

type GetMyAcademicRecordsOutput struct {
	StudentID      string                        `json:"student_id"`
	ClassID        string                        `json:"class_id,omitempty"`
	ClassSummaries []StudentAcademicClassSummary `json:"class_summaries"`
	Records        []StudentAcademicRecordItem   `json:"records"`
}

type GetMyAcademicRecordsUseCase interface {
	Execute(ctx context.Context, input GetMyAcademicRecordsInput) (*GetMyAcademicRecordsOutput, error)
}

type getMyAcademicRecordsUseCase struct {
	studentRepo       repointerface.StudentRepository
	recordRepo        repointerface.AcademicRecordRepository
	classScheduleRepo repointerface.ClassScheduleRepository
}

func NewGetMyAcademicRecordsUseCase(
	studentRepo repointerface.StudentRepository,
	recordRepo repointerface.AcademicRecordRepository,
	classScheduleRepo repointerface.ClassScheduleRepository,
) GetMyAcademicRecordsUseCase {
	return &getMyAcademicRecordsUseCase{
		studentRepo:       studentRepo,
		recordRepo:        recordRepo,
		classScheduleRepo: classScheduleRepo,
	}
}

func (uc *getMyAcademicRecordsUseCase) Execute(ctx context.Context, input GetMyAcademicRecordsInput) (*GetMyAcademicRecordsOutput, error) {
	ctxLogger := logger.NewLogger(ctx)

	if strings.TrimSpace(input.Actor.Role) != "STUDENT" {
		return nil, ErrStudentAccessDenied
	}

	student, err := resolveStudentByEmail(ctx, uc.studentRepo, input.Actor.Email)
	if err != nil {
		ctxLogger.Errorf("Failed to resolve student from actor email %s: %v", input.Actor.Email, err)
		return nil, err
	}

	records, err := uc.recordRepo.ListByStudentID(ctx, student.ID)
	if err != nil {
		ctxLogger.Errorf("Failed to load academic records for student %s: %v", student.ID, err)
		return nil, err
	}

	classScheduleMap := make(map[string][]entities.ClassSchedule)
	filtered := make([]StudentAcademicRecordItem, 0, len(records))
	classSummaryMap := make(map[string]*StudentAcademicClassSummary)

	for _, record := range records {
		lesson := record.LessonSummary.Lesson
		if input.ClassID != "" && lesson.ClassID != input.ClassID {
			continue
		}
		if input.DateFrom != nil && lesson.DateEnd.Before(input.DateFrom.UTC()) {
			continue
		}
		if input.DateTo != nil && lesson.DateStart.After(input.DateTo.UTC()) {
			continue
		}

		schedules, exists := classScheduleMap[lesson.ClassID]
		if !exists {
			schedules, err = uc.classScheduleRepo.GetSchedulesByClassID(ctx, lesson.ClassID)
			if err != nil {
				ctxLogger.Errorf("Failed to load class schedules for class %s: %v", lesson.ClassID, err)
				return nil, err
			}
			classScheduleMap[lesson.ClassID] = schedules
		}

		item := StudentAcademicRecordItem{
			RecordID:        record.ID,
			LessonSummaryID: record.LessonSummaryID,
			Lesson: StudentAcademicRecordLesson{
				ID:        lesson.ID,
				ClassID:   lesson.ClassID,
				ClassName: lesson.Class.Name,
				ClassCode: lesson.Class.Code,
				Teacher:   StudentAttendanceTeacher{},
				RoomID:    lesson.RoomID,
				DateStart: lesson.DateStart,
				DateEnd:   lesson.DateEnd,
				Notes:     lesson.Notes,
				Summary: StudentAcademicRecordLessonSummary{
					ID:       record.LessonSummary.ID,
					Topic:    record.LessonSummary.Topic,
					Homework: record.LessonSummary.Homework,
				},
			},
			HomeworkCompleted:  record.HomeworkCompleted,
			HomeworkScore:      record.HomeworkScore,
			AttitudeRating:     record.AttitudeRating,
			ParticipationScore: record.ParticipationScore,
			PersonalComment:    record.PersonalComment,
			TotalScore:         record.TotalScore,
			IsCompleted:        record.IsCompleted,
			CreatedAt:          record.CreatedAt,
			UpdatedAt:          record.UpdatedAt,
		}

		if lesson.TeacherID != nil {
			item.Lesson.Teacher.ID = lesson.TeacherID
		}
		if lesson.Teacher.Code != "" {
			code := lesson.Teacher.Code
			item.Lesson.Teacher.Code = &code
		}
		if lesson.Teacher.FullName != "" {
			name := lesson.Teacher.FullName
			item.Lesson.Teacher.FullName = &name
		}
		if lesson.Room.Name != "" {
			roomName := lesson.Room.Name
			item.Lesson.RoomName = &roomName
		}
		if matched := matchShiftForLesson(lesson, schedules); matched != nil {
			item.Lesson.Shift = &StudentAttendanceShift{
				ID:              matched.Shift.ID,
				Code:            matched.Shift.Code,
				Name:            matched.Shift.Name,
				StartTime:       matched.Shift.StartTime,
				EndTime:         matched.Shift.EndTime,
				DurationMinutes: matched.Shift.DurationMinutes,
				SessionType:     matched.Shift.SessionType,
			}
			if item.Lesson.RoomID == nil && matched.RoomID != nil {
				item.Lesson.RoomID = matched.RoomID
			}
			if item.Lesson.RoomName == nil && matched.Room.Name != "" {
				roomName := matched.Room.Name
				item.Lesson.RoomName = &roomName
			}
		}

		summary := classSummaryMap[lesson.ClassID]
		if summary == nil {
			summary = &StudentAcademicClassSummary{
				ClassID:   lesson.ClassID,
				ClassName: lesson.Class.Name,
				ClassCode: lesson.Class.Code,
			}
			classSummaryMap[lesson.ClassID] = summary
		}
		summary.RecordsCount++
		summary.AverageTotalScore += record.TotalScore
		if record.IsCompleted {
			summary.CompletedCount++
		}

		filtered = append(filtered, item)
	}

	classSummaries := make([]StudentAcademicClassSummary, 0, len(classSummaryMap))
	for _, summary := range classSummaryMap {
		if summary.RecordsCount > 0 {
			summary.AverageTotalScore = summary.AverageTotalScore / float64(summary.RecordsCount)
		}
		classSummaries = append(classSummaries, *summary)
	}

	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].Lesson.DateStart.Before(filtered[j].Lesson.DateStart)
	})
	sort.Slice(classSummaries, func(i, j int) bool {
		return classSummaries[i].ClassName < classSummaries[j].ClassName
	})

	return &GetMyAcademicRecordsOutput{
		StudentID:      student.ID,
		ClassID:        input.ClassID,
		ClassSummaries: classSummaries,
		Records:        filtered,
	}, nil
}
