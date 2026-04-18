package teacherportal

import (
	"time"

	"doan/internal/entities"
)

type LeaveRequestStudentItem struct {
	ID       string `json:"id"`
	Code     string `json:"code"`
	FullName string `json:"full_name"`
}

type LeaveRequestClassItem struct {
	ID   string `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

type LeaveRequestLessonItem struct {
	ID        string    `json:"id"`
	DateStart time.Time `json:"date_start"`
	DateEnd   time.Time `json:"date_end"`
}

type LeaveRequestItem struct {
	ID              string                  `json:"id"`
	Student         LeaveRequestStudentItem `json:"student"`
	LeaveType       string                  `json:"leave_type"`
	ApplyDate       time.Time               `json:"apply_date"`
	LateMinutes     int                     `json:"late_minutes"`
	EarlyMinutes    int                     `json:"early_minutes"`
	Reason          string                  `json:"reason"`
	Documents       []string                `json:"documents"`
	Class           *LeaveRequestClassItem  `json:"class,omitempty"`
	Lesson          *LeaveRequestLessonItem `json:"lesson,omitempty"`
	Subject         string                  `json:"subject"`
	Status          string                  `json:"status"`
	ApprovedByID    *string                 `json:"approved_by_id,omitempty"`
	ApprovedAt      *time.Time              `json:"approved_at,omitempty"`
	RejectionReason string                  `json:"rejection_reason"`
	CreatedAt       time.Time               `json:"created_at"`
	UpdatedAt       time.Time               `json:"updated_at"`
}

func buildLeaveRequestItem(request entities.LeaveRequest) LeaveRequestItem {
	item := LeaveRequestItem{
		ID: request.ID,
		Student: LeaveRequestStudentItem{
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
		item.Class = &LeaveRequestClassItem{
			ID:   request.Class.ID,
			Code: request.Class.Code,
			Name: request.Class.Name,
		}
	}
	if request.LessonID != nil {
		item.Lesson = &LeaveRequestLessonItem{
			ID:        request.Lesson.ID,
			DateStart: request.Lesson.DateStart,
			DateEnd:   request.Lesson.DateEnd,
		}
	}
	return item
}
