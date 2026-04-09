package shift

import "time"

type CreateShiftRequest struct {
	Code            string `json:"code" binding:"required"`
	Name            string `json:"name" binding:"required"`
	StartTime       string `json:"start_time" binding:"required"`
	EndTime         string `json:"end_time" binding:"required"`
	DurationMinutes int    `json:"duration_minutes" binding:"required,min=1"`
	SessionType     string `json:"session_type" binding:"required,oneof=MORNING AFTERNOON EVENING CUSTOM"`
	IsActive        bool   `json:"is_active"`
	Notes           string `json:"notes"`
}

type UpdateShiftRequest struct {
	Code            string `json:"code" binding:"required"`
	Name            string `json:"name" binding:"required"`
	StartTime       string `json:"start_time" binding:"required"`
	EndTime         string `json:"end_time" binding:"required"`
	DurationMinutes int    `json:"duration_minutes" binding:"required,min=1"`
	SessionType     string `json:"session_type" binding:"required,oneof=MORNING AFTERNOON EVENING CUSTOM"`
	IsActive        bool   `json:"is_active"`
	Notes           string `json:"notes"`
}

type ShiftResponse struct {
	ID              string    `json:"id"`
	Code            string    `json:"code"`
	Name            string    `json:"name"`
	StartTime       string    `json:"start_time"`
	EndTime         string    `json:"end_time"`
	DurationMinutes int       `json:"duration_minutes"`
	SessionType     string    `json:"session_type"`
	IsActive        bool      `json:"is_active"`
	Notes           string    `json:"notes"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type ListShiftsResponse struct {
	Shifts     []ShiftResponse `json:"shifts"`
	Pagination PaginationMeta  `json:"pagination"`
}

type PaginationMeta struct {
	ItemsPerPage int   `json:"items_per_page"`
	TotalItems   int64 `json:"total_items"`
	CurrentPage  int   `json:"current_page"`
	TotalPages   int   `json:"total_pages"`
}
