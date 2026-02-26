package program

import (
	"doan/cmd/http/controllers/course"
	"time"
)

// CreateProgramRequest represents the request body for creating a program
type CreateProgramRequest struct {
	Code          string     `json:"code" binding:"required"`
	Name          string     `json:"name" binding:"required"`
	Track         string     `json:"track"`
	EffectiveFrom *time.Time `json:"effective_from"`
	EffectiveTo   *time.Time `json:"effective_to"`
	ApprovalNote  string     `json:"approval_note"`
}

// UpdateProgramRequest represents the request body for updating a program
type UpdateProgramRequest struct {
	Code          *string    `json:"code"`
	Name          *string    `json:"name"`
	Track         *string    `json:"track"`
	EffectiveFrom *time.Time `json:"effective_from"`
	EffectiveTo   *time.Time `json:"effective_to"`
	ApprovalNote  *string    `json:"approval_note"`
}

// ProgramResponse represents a program in the response
type ProgramResponse struct {
	ID            string                  `json:"id"`
	Code          string                  `json:"code"`
	Name          string                  `json:"name"`
	Track         string                  `json:"track"`
	EffectiveFrom *time.Time              `json:"effective_from"`
	EffectiveTo   *time.Time              `json:"effective_to"`
	ApprovalNote  string                  `json:"approval_note"`
	CreatedAt     time.Time               `json:"created_at"`
	UpdatedAt     time.Time               `json:"updated_at"`
	Courses       []course.CourseResponse `json:"courses,omitempty"`
}

// ListProgramsResponse represents the response for listing programs
type ListProgramsResponse struct {
	Programs   []ProgramResponse `json:"programs"`
	Pagination PaginationMeta    `json:"pagination"`
}

// PaginationMeta represents pagination metadata
type PaginationMeta struct {
	ItemsPerPage uint64 `json:"items_per_page"`
	TotalItems   uint64 `json:"total_items"`
	CurrentPage  uint64 `json:"current_page"`
	TotalPages   uint64 `json:"total_pages"`
}

// AddRemoveCoursesRequest represents the specific courses to link/unlink
type AddRemoveCoursesRequest struct {
	CourseIDs []string `json:"course_ids" binding:"required"`
}

// MessageResponse represents a simple message response
type MessageResponse struct {
	Message string `json:"message"`
}
