package course

import "time"

// CreateCourseRequest represents the request body for creating a course
type CreateCourseRequest struct {
	Code                   string  `json:"code" binding:"required"`
	Name                   string  `json:"name" binding:"required"`
	Description            string  `json:"description"`
	GradeLevel             string  `json:"grade_level"`
	Subject                string  `json:"subject"`
	SessionCount           int     `json:"session_count"`
	SessionDurationMinutes int     `json:"session_duration_minutes"`
	TotalHours             float64 `json:"total_hours"`
	Price                  float64 `json:"price"`
	Status                 string  `json:"status"`
}

// UpdateCourseRequest represents the request body for updating a course
type UpdateCourseRequest struct {
	Code                   *string  `json:"code"`
	Name                   *string  `json:"name"`
	Description            *string  `json:"description"`
	GradeLevel             *string  `json:"grade_level"`
	Subject                *string  `json:"subject"`
	SessionCount           *int     `json:"session_count"`
	SessionDurationMinutes *int     `json:"session_duration_minutes"`
	TotalHours             *float64 `json:"total_hours"`
	Price                  *float64 `json:"price"`
	Status                 *string  `json:"status"`
}

// CourseResponse represents a course in the response
type CourseResponse struct {
	ID                     string    `json:"id"`
	Code                   string    `json:"code"`
	Name                   string    `json:"name"`
	Description            string    `json:"description"`
	GradeLevel             string    `json:"grade_level"`
	Subject                string    `json:"subject"`
	SessionCount           int       `json:"session_count"`
	SessionDurationMinutes int       `json:"session_duration_minutes"`
	TotalHours             float64   `json:"total_hours"`
	Price                  float64   `json:"price"`
	Status                 string    `json:"status"`
	CreatedAt              time.Time `json:"created_at"`
	UpdatedAt              time.Time `json:"updated_at"`
}

// ListCoursesResponse represents the response for listing courses
type ListCoursesResponse struct {
	Courses    []CourseResponse `json:"courses"`
	Pagination PaginationMeta   `json:"pagination"`
}

// PaginationMeta represents pagination metadata
type PaginationMeta struct {
	ItemsPerPage uint64 `json:"items_per_page"`
	TotalItems   uint64 `json:"total_items"`
	CurrentPage  uint64 `json:"current_page"`
	TotalPages   uint64 `json:"total_pages"`
}

// MessageResponse represents a simple message response
type MessageResponse struct {
	Message string `json:"message"`
}
