package account

type UserResponse struct {
	ID        string `json:"id"`
	Code      string `json:"code"`
	FullName  string `json:"full_name"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	IsActive  bool   `json:"is_active"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

type CreateUserRequest struct {
	Code     string `json:"code"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Role     string `json:"role" binding:"required"`
	IsActive *bool  `json:"is_active"`
	Password string `json:"password" binding:"required,min=6"`
}

type UpdateUserRequest struct {
	FullName string `json:"full_name"`
	Role     string `json:"role"`
	IsActive *bool  `json:"is_active"`
}

type ResetPasswordRequest struct {
	NewPassword string `json:"new_password" binding:"required,min=6"`
}
