package model

// LoginResponse represents the structure for login responses
type LoginResponse struct {
	BaseAPIResponse
	Error *ErrorResponse `json:"error,omitempty"`
	Data  LoginData      `json:"data"`
}

// LoginData contains specific login information
type LoginData struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	Role         string `json:"role"`
	User         *User  `json:"user,omitempty"`
}
