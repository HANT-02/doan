package user

// CreateAuthTokenInput is a struct that contains the input for CreateAuthToken method
type CreateAuthTokenInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// CreateAuthTokenOutput is a struct that contains the output for CreateAuthToken method
type CreateAuthTokenOutput struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
