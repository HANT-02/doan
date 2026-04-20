package studentportal

import "doan/internal/entities"

func buildStudentLeaveRequestItem(request entities.LeaveRequest) StudentLeaveRequestItem {
	item := StudentLeaveRequestItem{
		ID: request.ID,
		Student: StudentLeaveRequestStudentItem{
			ID:       request.Student.ID,
			Code:     request.Student.Code,
			FullName: request.Student.FullName,
		},
		LeaveType:       request.LeaveType,
		ApplyDate:       request.ApplyDate,
		LateMinutes:     request.LateMinutes,
		EarlyMinutes:    request.EarlyMinutes,
		Reason:          request.Reason,
		Documents:       request.Documents,
		Subject:         request.Subject,
		Status:          request.Status,
		ApprovedByID:    request.ApprovedByID,
		ApprovedAt:      request.ApprovedAt,
		RejectionReason: request.RejectionReason,
		CreatedAt:       request.CreatedAt,
		UpdatedAt:       request.UpdatedAt,
	}
	if request.ClassID != nil {
		item.Class = &StudentLeaveRequestClassItem{
			ID:   request.Class.ID,
			Code: request.Class.Code,
			Name: request.Class.Name,
		}
	}
	if request.LessonID != nil {
		item.Lesson = &StudentLeaveRequestLessonItem{
			ID:        request.Lesson.ID,
			DateStart: request.Lesson.DateStart,
			DateEnd:   request.Lesson.DateEnd,
		}
	}
	return item
}
