package user

// CreateAuthTokenInput is a struct that contains the input for CreateAuthToken method
type CreateAuthTokenInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// CreateAuthTokenOutput is a struct that contains the output for CreateAuthToken method
type CreateAuthTokenOutput struct {
	AccessToken  string     `json:"access_token"`
	RefreshToken string     `json:"refresh_token"`
	User         UserOutput `json:"user"`
}

// UserOutput is a struct that contains user information
type UserOutput struct {
	ID       string `json:"id"`
	Code     string `json:"code"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	IsActive bool   `json:"is_active"`
}

// TokenClaims is a struct that contains JWT token claims
type TokenClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
}
