package teacher

import "time"

// CreateTeacherRequest represents the request body for creating a teacher
type CreateTeacherRequest struct {
	Code            string `json:"code"`
	FullName        string `json:"full_name" binding:"required"`
	Email           string `json:"email"`
	Phone           string `json:"phone"`
	IsSchoolTeacher bool   `json:"is_school_teacher"`
	SchoolName      string `json:"school_name"`
	EmploymentType  string `json:"employment_type"` // PART_TIME, FULL_TIME
	Status          string `json:"status"`          // ACTIVE, INACTIVE
	Notes           string `json:"notes"`
}

// UpdateTeacherRequest represents the request body for updating a teacher
type UpdateTeacherRequest struct {
	Code            *string `json:"code"`
	FullName        *string `json:"full_name"`
	Email           *string `json:"email"`
	Phone           *string `json:"phone"`
	IsSchoolTeacher *bool   `json:"is_school_teacher"`
	SchoolName      *string `json:"school_name"`
	EmploymentType  *string `json:"employment_type"`
	Status          *string `json:"status"`
	Notes           *string `json:"notes"`
}

// TeacherResponse represents a teacher in the response
type TeacherResponse struct {
	ID              string    `json:"id"`
	Code            string    `json:"code"`
	FullName        string    `json:"full_name"`
	Email           string    `json:"email"`
	Phone           string    `json:"phone"`
	IsSchoolTeacher bool      `json:"is_school_teacher"`
	SchoolName      string    `json:"school_name"`
	EmploymentType  string    `json:"employment_type"`
	Status          string    `json:"status"`
	Notes           string    `json:"notes"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// ListTeachersResponse represents the response for listing teachers
type ListTeachersResponse struct {
	Teachers   []TeacherResponse `json:"teachers"`
	Pagination PaginationMeta    `json:"pagination"`
}

// PaginationMeta represents pagination metadata
type PaginationMeta struct {
	ItemsPerPage uint64 `json:"items_per_page"`
	TotalItems   uint64 `json:"total_items"`
	CurrentPage  uint64 `json:"current_page"`
	TotalPages   uint64 `json:"total_pages"`
}

// TimetableLesson represents a lesson in the timetable
type TimetableLesson struct {
	ID        string    `json:"id"`
	ClassID   string    `json:"class_id"`
	ClassName string    `json:"class_name"`
	RoomID    *string   `json:"room_id"`
	RoomName  *string   `json:"room_name"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Notes     string    `json:"notes"`
}

// TimetableResponse represents the response for teacher timetable
type TimetableResponse struct {
	Lessons []TimetableLesson `json:"lessons"`
}

// TeachingHoursStat represents teaching hours for a period
type TeachingHoursStat struct {
	Period string  `json:"period"`
	Hours  float64 `json:"hours"`
}

// TeachingHoursStatsResponse represents the response for teaching hours statistics
type TeachingHoursStatsResponse struct {
	TotalHours float64             `json:"total_hours"`
	Breakdown  []TeachingHoursStat `json:"breakdown"`
}

// MessageResponse represents a simple message response
type MessageResponse struct {
	Message string `json:"message"`
}
