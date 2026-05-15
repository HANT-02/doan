package class

import (
	"context"
	repointerface "doan/internal/repositories/interface"
	skillservice "doan/internal/services/skills"
	"fmt"
	"strings"
)

func validateTeacherCourseSkills(
	ctx context.Context,
	teacherRepo repointerface.TeacherRepository,
	courseRepo repointerface.CourseRepository,
	teacherID *string,
	courseID *string,
) error {
	if teacherID == nil || *teacherID == "" || courseID == nil || *courseID == "" {
		return nil
	}

	teacherEntity, err := teacherRepo.GetByID(ctx, *teacherID)
	if err != nil {
		return err
	}

	courseEntity, err := courseRepo.GetByID(ctx, *courseID)
	if err != nil {
		return err
	}

	missingSkills := skillservice.MissingRequiredCodes([]string(teacherEntity.Skills), []string(courseEntity.RequiredSkills))
	if len(missingSkills) == 0 {
		return nil
	}

	return fmt.Errorf(
		"giáo viên chưa đáp ứng kỹ năng/chứng chỉ bắt buộc của khóa học, thiếu: %s",
		strings.Join(missingSkills, ", "),
	)
}
