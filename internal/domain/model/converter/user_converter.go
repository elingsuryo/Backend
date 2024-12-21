package converter

import (
	"time"

	"github.com/TrinityKnights/Backend/internal/domain/entity"
	"github.com/TrinityKnights/Backend/internal/domain/model"
)

func UserToResponse(user *entity.User) *model.UserResponse {
	var lastLogin *string
	if user.LastLogin != nil {
		formatted := user.LastLogin.UTC().Add(time.Hour * 7).Format(time.RFC3339)
		lastLogin = &formatted
	}

	return &model.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		Role:      user.Role,
		Status:    user.Status,
		LastLogin: lastLogin,
		CreatedAt: user.CreatedAt.UTC().Add(time.Hour * 7).Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.UTC().Add(time.Hour * 7).Format(time.RFC3339),
	}
}

func LoginToTokenResponse(accessToken, refreshToken string) *model.TokenResponse {
	return &model.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
}
