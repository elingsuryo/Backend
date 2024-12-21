package model

type UserResponse struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Email     string  `json:"email"`
	Role      string  `json:"role"`
	Status    bool    `json:"status"`
	LastLogin *string `json:"last_login,omitempty"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email,lte=100"`
	Password string `json:"password" validate:"required,min=8,lte=255"`
	Name     string `json:"name" validate:"required,lte=100"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email,lte=100"`
	Password string `json:"password" validate:"required,min=8,lte=255"`
}

type UpdateRequest struct {
	Email    string `json:"email" validate:"omitempty,email,lte=100"`
	Password string `json:"password" validate:"omitempty,min=8,lte=255"`
	Name     string `json:"name" validate:"omitempty,lte=100"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type ReqResetPasswordRequest struct {
	Email string `json:"email" validate:"required,email,lte=100"`
}

type VerifyResponse struct {
	Status string `json:"status"`
}

type ResetPasswordRequest struct {
	Token       string `param:"token" validate:"required"`
	Email       string `json:"email" validate:"required,email,lte=100"`
	NewPassword string `json:"new_password" validate:"required,min=8,lte=255"`
}

type VerifyRequest struct {
	Token string `param:"token" validate:"required"`
}
