package leaveflow

import (
	"context"
	"strings"

	"doan/internal/entities"
	"doan/internal/repositories"
	repointerface "doan/internal/repositories/interface"
)

func resolveStudentByEmail(ctx context.Context, studentRepo repointerface.StudentRepository, email string) (*entities.Student, error) {
	condition := repositories.NewCommonCondition()
	condition.AddCondition("email", strings.TrimSpace(email), repositories.Equal)
	condition.SetPaging(1, 1)
	result, err := studentRepo.GetByCondition(ctx, condition)
	if err != nil {
		return nil, err
	}
	if result == nil || len(result.Data) == 0 || result.Data[0] == nil {
		return nil, ErrStudentNotFound
	}
	return result.Data[0], nil
}

func resolveTeacherByEmail(ctx context.Context, teacherRepo repointerface.TeacherRepository, email string) (*entities.Teacher, error) {
	condition := repositories.NewCommonCondition()
	condition.AddCondition("email", strings.TrimSpace(email), repositories.Equal)
	condition.SetPaging(1, 1)
	result, err := teacherRepo.GetByCondition(ctx, condition)
	if err != nil {
		return nil, err
	}
	if result == nil || len(result.Data) == 0 || result.Data[0] == nil {
		return nil, ErrTeacherNotFound
	}
	return result.Data[0], nil
}

func isValidLeaveType(value string) bool {
	switch strings.ToUpper(strings.TrimSpace(value)) {
	case "LEAVE", "LATE", "EARLY":
		return true
	default:
		return false
	}
}
