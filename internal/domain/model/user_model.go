package model

import "bytes"

type UserResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	Status    bool   `json:"status"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
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

type RequestReset struct {
	Email string `json:"email" validate:"required,email,lte=100"`
}

type ResponseReset struct {
	Status string `json:"status"`
}

type ResetPassword struct {
	Token       string `param:"token" validate:"required"`
	Email       string `json:"email" validate:"required,email,lte=100"`
	NewPassword string `json:"new_password" validate:"required,min=8,lte=255"`
}

type VerifyEmail struct {
	Token string `param:"token" validate:"required"`
}

type SendEmail struct {
	EmailTo   string `json:"email_to" validate:"required,email,lte=100"`
	EmailFrom string `json:"email_from" validate:"required,email,lte=100"`
	Subject   string `json:"subject" validate:"required,lte=100"`
	Body      bytes.Buffer `json:"body" validate:"required"`
}