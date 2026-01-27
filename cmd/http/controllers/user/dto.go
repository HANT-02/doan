package user

// LoginRequest represents the login request body
type LoginRequest struct {
	Username string `json:"username" binding:"required" example:"user@example.com"`
	Password string `json:"password" binding:"required" example:"password123"`
}

// LoginResponse represents the login response
type LoginResponse struct {
	AccessToken  string       `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	RefreshToken string       `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	User         UserResponse `json:"user"`
}

// UserResponse represents the user information in response
type UserResponse struct {
	ID       string `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Code     string `json:"code" example:"USER001"`
	FullName string `json:"full_name" example:"John Doe"`
	Email    string `json:"email" example:"user@example.com"`
	Role     string `json:"role" example:"STUDENT"`
	IsActive bool   `json:"is_active" example:"true"`
}

// LogoutRequest represents the logout request body
type LogoutRequest struct {
	Token string `json:"token" binding:"required"`
}

// LogoutResponse represents the logout response
type LogoutResponse struct {
	Message string `json:"message" example:"Logged out successfully"`
}

// RefreshTokenRequest represents the refresh token request body
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// RefreshTokenResponse represents the refresh token response
type RefreshTokenResponse struct {
	AccessToken string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}
