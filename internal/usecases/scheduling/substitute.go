package scheduling

import (
	"context"
	"doan/internal/entities"
	"doan/internal/repositories"
	repositoryinterface "doan/internal/repositories/interface"
	schedulingservice "doan/internal/services/scheduling"
	skillservice "doan/internal/services/skills"
	"errors"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Actor struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
}

var (
	ErrSubstituteAccessDenied    = errors.New("substitute access denied")
	ErrSubstituteTeacherNotFound = errors.New("substitute teacher profile missing")
	ErrSubstituteNotEligible     = errors.New("substitute teacher is not eligible")
)

type SubstituteSuggestion struct {
	TeacherID    string   `json:"teacher_id"`
	TeacherName  string   `json:"teacher_name"`
	TeacherCode  string   `json:"teacher_code"`
	Score        int      `json:"score"`
	MatchReasons []string `json:"match_reasons"`
	IsAvailable  bool     `json:"is_available"`
}

type SubstituteUseCase interface {
	SuggestSubstituteTeachers(ctx context.Context, actor Actor, lessonID string) ([]SubstituteSuggestion, error)
	AssignSubstitute(ctx context.Context, actor Actor, lessonID string, newTeacherID string, reason string) error
}

type substituteUseCase struct {
	lessonRepo       repositoryinterface.LessonRepository
	teacherRepo      repositoryinterface.TeacherRepository
	classRepo        repositoryinterface.ClassRepository
	roomRepo         repositoryinterface.RoomRepository
	campusTravelRepo repositoryinterface.CampusTravelTimeRepository
}

func NewSubstituteUseCase(
	lessonRepo repositoryinterface.LessonRepository,
	teacherRepo repositoryinterface.TeacherRepository,
	classRepo repositoryinterface.ClassRepository,
	roomRepo repositoryinterface.RoomRepository,
	campusTravelRepo repositoryinterface.CampusTravelTimeRepository,
) SubstituteUseCase {
	return &substituteUseCase{
		lessonRepo:       lessonRepo,
		teacherRepo:      teacherRepo,
		classRepo:        classRepo,
		roomRepo:         roomRepo,
		campusTravelRepo: campusTravelRepo,
	}
}

func (uc *substituteUseCase) SuggestSubstituteTeachers(ctx context.Context, actor Actor, lessonID string) ([]SubstituteSuggestion, error) {
	lesson, err := uc.lessonRepo.GetLessonWithRelations(ctx, lessonID)
	if err != nil || lesson == nil {
		return nil, err
	}

	if err := uc.ensureSubstituteAccess(ctx, actor, lesson); err != nil {
		return nil, err
	}

	classCondition := repositories.NewCommonCondition()
	classCondition.AddCondition("id", lesson.ClassID, repositories.Equal)
	classCondition.SetPreload([]string{"Course", "Room"})
	classCondition.SetPaging(1, 1)
	classPage, err := uc.classRepo.GetByCondition(ctx, classCondition)
	if err != nil {
		return nil, err
	}
	if classPage == nil || len(classPage.Data) == 0 || classPage.Data[0] == nil {
		return nil, ErrSubstituteTeacherNotFound
	}
	class := classPage.Data[0]

	teacherCondition := repositories.NewCommonCondition().
		WithCondition("status", "ACTIVE", "eq").
		WithPaging(1000, 1)
	teacherPage, err := uc.teacherRepo.GetByCondition(ctx, teacherCondition)
	if err != nil {
		return nil, err
	}

	dayStart := time.Date(lesson.DateStart.Year(), lesson.DateStart.Month(), lesson.DateStart.Day(), 0, 0, 0, 0, lesson.DateStart.Location())
	dayEnd := dayStart.Add(24 * time.Hour)
	dayLessons, err := uc.lessonRepo.ListInRange(ctx, dayStart, dayEnd)
	if err != nil {
		return nil, err
	}

	teacherLessons := make(map[string][]entities.Lesson)
	for _, l := range dayLessons {
		if l.Status != entities.LessonStatusPublished && l.Status != entities.LessonStatusDraft {
			continue
		}
		if l.TeacherID != nil {
			teacherLessons[*l.TeacherID] = append(teacherLessons[*l.TeacherID], l)
		}
	}

	travelMap := map[string]int{}
	travelPage, err := uc.campusTravelRepo.GetByCondition(ctx, repositories.NewCommonCondition().WithPaging(1000, 1))
	if err == nil && travelPage != nil {
		travelTimes := make([]entities.CampusTravelTime, 0)
		for _, ptr := range travelPage.Data {
			if ptr != nil {
				travelTimes = append(travelTimes, *ptr)
			}
		}
		travelMap = schedulingservice.BuildCampusTravelTimeMap(travelTimes)
	}

	var suggestions []SubstituteSuggestion
	requiredSkills := []string(class.Course.RequiredSkills)

	for _, teacherPtr := range teacherPage.Data {
		teacher := *teacherPtr

		if lesson.TeacherID != nil && *lesson.TeacherID == teacher.ID {
			continue
		}

		score := 100
		var reasons []string
		isAvailable := true

		missingSkills := skillservice.MissingRequiredCodes([]string(teacher.Skills), requiredSkills)
		if len(missingSkills) > 0 {
			isAvailable = false
			reasons = append(reasons, "Thiếu kỹ năng: "+strings.Join(missingSkills, ", "))
		} else {
			reasons = append(reasons, "Đủ kỹ năng chuyên môn")
		}

		workloadToday := 0
		for _, teacherLesson := range teacherLessons[teacher.ID] {
			if teacherLesson.ID == lesson.ID {
				continue
			}
			workloadToday++
			if overlaps(teacherLesson.DateStart, teacherLesson.DateEnd, lesson.DateStart, lesson.DateEnd) {
				isAvailable = false
				reasons = append(reasons, "Bận lịch trong cùng khung giờ")
				break
			}
			if !hasTravelFeasibility(teacherLesson, *lesson, travelMap) {
				isAvailable = false
				reasons = append(reasons, "Không đủ thời gian di chuyển giữa 2 cơ sở")
				break
			}
		}

		if !isAvailable {
			continue
		}

		score -= workloadToday * 5
		if workloadToday == 0 {
			reasons = append(reasons, "Tải dạy trong ngày thấp")
		} else {
			reasons = append(reasons, "Tải dạy trong ngày: "+strconv.Itoa(workloadToday)+" ca")
		}
		reasons = append(reasons, "Trống lịch và đủ thời gian di chuyển")

		suggestions = append(suggestions, SubstituteSuggestion{
			TeacherID:    teacher.ID,
			TeacherName:  teacher.FullName,
			TeacherCode:  teacher.Code,
			Score:        score,
			MatchReasons: reasons,
			IsAvailable:  true,
		})
	}

	sort.Slice(suggestions, func(i, j int) bool {
		return suggestions[i].Score > suggestions[j].Score
	})

	return suggestions, nil
}

func (uc *substituteUseCase) AssignSubstitute(ctx context.Context, actor Actor, lessonID string, newTeacherID string, reason string) error {
	lesson, err := uc.lessonRepo.GetLessonWithRelations(ctx, lessonID)
	if err != nil || lesson == nil {
		return err
	}

	if err := uc.ensureSubstituteAccess(ctx, actor, lesson); err != nil {
		return err
	}

	suggestions, err := uc.SuggestSubstituteTeachers(ctx, actor, lessonID)
	if err != nil {
		return err
	}

	isEligible := false
	for _, suggestion := range suggestions {
		if suggestion.TeacherID == newTeacherID && suggestion.IsAvailable {
			isEligible = true
			break
		}
	}
	if !isEligible {
		return ErrSubstituteNotEligible
	}

	return uc.lessonRepo.Update(ctx, lessonID, map[string]interface{}{
		"teacher_id":    newTeacherID,
		"change_reason": reason,
		"status":        entities.LessonStatusPublished,
	})
}

func (uc *substituteUseCase) ensureSubstituteAccess(ctx context.Context, actor Actor, lesson *entities.Lesson) error {
	role := strings.TrimSpace(actor.Role)
	if role == "ADMIN" || role == "SUPER_ADMIN" {
		return nil
	}
	if role != "TEACHER" {
		return ErrSubstituteAccessDenied
	}

	teacher, err := resolveSchedulingTeacherByEmail(ctx, uc.teacherRepo, actor.Email)
	if err != nil {
		return err
	}
	if lesson.TeacherID == nil || *lesson.TeacherID != teacher.ID {
		return ErrSubstituteAccessDenied
	}
	return nil
}

func resolveSchedulingTeacherByEmail(ctx context.Context, teacherRepo repositoryinterface.TeacherRepository, email string) (*entities.Teacher, error) {
	condition := repositories.NewCommonCondition()
	condition.AddCondition("email", strings.TrimSpace(email), repositories.Equal)
	condition.SetPaging(1, 1)

	result, err := teacherRepo.GetByCondition(ctx, condition)
	if err != nil {
		return nil, err
	}
	if result == nil || len(result.Data) == 0 || result.Data[0] == nil {
		return nil, ErrSubstituteTeacherNotFound
	}
	return result.Data[0], nil
}

func hasTravelFeasibility(existingLesson entities.Lesson, candidateLesson entities.Lesson, travelMap map[string]int) bool {
	if existingLesson.RoomID == nil || candidateLesson.RoomID == nil {
		return true
	}

	if overlaps(existingLesson.DateStart, existingLesson.DateEnd, candidateLesson.DateStart, candidateLesson.DateEnd) {
		return false
	}

	if existingLesson.DateEnd.Before(candidateLesson.DateStart) || existingLesson.DateEnd.Equal(candidateLesson.DateStart) {
		return schedulingservice.HasSufficientTravelGap(existingLesson.DateEnd, candidateLesson.DateStart, &existingLesson.Room, &candidateLesson.Room, travelMap)
	}

	return schedulingservice.HasSufficientTravelGap(candidateLesson.DateEnd, existingLesson.DateStart, &candidateLesson.Room, &existingLesson.Room, travelMap)
}
