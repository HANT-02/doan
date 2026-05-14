package scheduling

import (
	"context"
	"fmt"
	"math"

	"doan/internal/entities"
	repositoryinterface "doan/internal/repositories/interface"
)

const (
	defaultMinimumEnrollmentRatio  = 0.80
	defaultMaxManualAdjustmentRate = 0.35
	minimumManualAdjustmentLimit   = 3
)

func filterClassesByEnrollment(
	ctx context.Context,
	enrollmentRepo repositoryinterface.EnrollmentRepository,
	classes []entities.Class,
) ([]entities.Class, []PreviewConflict, error) {
	if len(classes) == 0 {
		return nil, nil, nil
	}

	eligible := make([]entities.Class, 0, len(classes))
	conflicts := make([]PreviewConflict, 0)
	for _, classEntity := range classes {
		enrollments, err := enrollmentRepo.ListByClassID(ctx, classEntity.ID)
		if err != nil {
			return nil, nil, err
		}

		currentCount := countApprovedEnrollments(enrollments)
		requiredCount := requiredEnrollmentCount(classEntity.MaxStudents, defaultMinimumEnrollmentRatio)
		if requiredCount > 0 && currentCount < requiredCount {
			conflicts = append(conflicts, PreviewConflict{
				VariableID: classEntity.ID,
				ClassID:    classEntity.ID,
				ClassCode:  classEntity.Code,
				ClassName:  classEntity.Name,
				Type:       "INSUFFICIENT_ENROLLMENT",
				Message: fmt.Sprintf(
					"Lớp hiện chỉ có %d/%d học sinh đạt điều kiện mở lớp. Cần tối thiểu %d học sinh (%.0f%% sĩ số mục tiêu) thì mới được đưa vào xếp lịch.",
					currentCount,
					classEntity.MaxStudents,
					requiredCount,
					defaultMinimumEnrollmentRatio*100,
				),
			})
			continue
		}

		eligible = append(eligible, classEntity)
	}

	return eligible, conflicts, nil
}

func countApprovedEnrollments(enrollments []entities.Enrollment) int {
	count := 0
	for _, enrollment := range enrollments {
		if enrollment.ApprovedAt != nil || enrollment.Status == "APPROVED" {
			count++
		}
	}
	return count
}

func requiredEnrollmentCount(maxStudents int, ratio float64) int {
	if maxStudents <= 0 || ratio <= 0 {
		return 0
	}
	return int(math.Ceil(float64(maxStudents) * ratio))
}

func allowedManualAdjustmentLimit(totalSessions int) int {
	if totalSessions <= 0 {
		return minimumManualAdjustmentLimit
	}

	derived := int(math.Ceil(float64(totalSessions) * defaultMaxManualAdjustmentRate))
	if derived < minimumManualAdjustmentLimit {
		return minimumManualAdjustmentLimit
	}
	return derived
}

func appendOperationalConflict(preview PreviewResult, conflict PreviewConflict) PreviewResult {
	preview.Conflicts = append(preview.Conflicts, conflict)
	preview.Summary.ConflictCount = len(preview.Conflicts)
	if preview.Status == "COMPLETED" {
		preview.Status = "PARTIAL"
	}
	return preview
}

func maybeAppendConflictDensityConflict(preview PreviewResult) PreviewResult {
	affectedSessions := make(map[string]struct{})
	for _, conflict := range preview.Conflicts {
		if conflict.VariableID == "" {
			continue
		}
		affectedSessions[conflict.VariableID] = struct{}{}
	}

	if len(affectedSessions) == 0 {
		return preview
	}

	limit := allowedManualAdjustmentLimit(preview.Summary.RequestedSessions)
	if len(affectedSessions) <= limit {
		return preview
	}

	return appendOperationalConflict(preview, PreviewConflict{
		Type: "EXCESSIVE_MANUAL_ADJUSTMENT",
		Message: fmt.Sprintf(
			"Preview hiện có %d ca học cần xử lý tay, vượt ngưỡng cho phép %d ca. Mật độ xung đột này không phù hợp để tiếp tục vá tay; hãy nới khoảng ngày, thêm ca/phòng hoặc điều chỉnh lịch tuần lớp trước khi chạy lại.",
			len(affectedSessions),
			limit,
		),
	})
}
