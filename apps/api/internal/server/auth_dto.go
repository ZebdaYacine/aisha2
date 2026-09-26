package httpapi

import (
	"github.com/aisha-platform/aisha/apps/api/internal/features/auth"
	customer "github.com/aisha-platform/aisha/apps/api/internal/features/user"
	"time"
)

type RegisterRequest struct {
	Email       string `json:"email" validate:"required,email,max=254"`
	Password    string `json:"password" validate:"required,min=12,max=128"`
	DisplayName string `json:"displayName" validate:"required,min=2,max=100"`
}
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email,max=254"`
	Password string `json:"password" validate:"required,max=128"`
}
type RefreshRequest struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}
type LogoutRequest struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}
type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email,max=254"`
}
type ResetPasswordRequest struct {
	Token    string `json:"token" validate:"required"`
	Password string `json:"password" validate:"required,min=12,max=128"`
}
type ChangePasswordRequest struct {
	CurrentPassword string `json:"currentPassword" validate:"required,max=128"`
	NewPassword     string `json:"newPassword" validate:"required,min=12,max=128"`
}
type UserResponse struct {
	ID              string    `json:"id"`
	Email           string    `json:"email"`
	DisplayName     string    `json:"displayName"`
	Status          string    `json:"status"`
	Roles           []string  `json:"roles"`
	CreatedAt       time.Time `json:"createdAt"`
	CustomerEnabled bool      `json:"customerEnabled"`
	ArtisanStatus   string    `json:"artisanStatus"`
	ArtisanEnabled  bool      `json:"artisanEnabled"`
	Capabilities    []string  `json:"capabilities"`
}
type TokenResponse struct {
	AccessToken  string    `json:"accessToken"`
	RefreshToken string    `json:"refreshToken"`
	TokenType    string    `json:"tokenType"`
	ExpiresAt    time.Time `json:"expiresAt"`
}
type AuthenticationResponse struct {
	User   UserResponse  `json:"user"`
	Tokens TokenResponse `json:"tokens"`
}

func userResponseFrom(user auth.User) UserResponse {
	response := UserResponse{ID: user.ID, Email: user.Email, DisplayName: user.DisplayName, Status: user.Status, Roles: user.Roles, CreatedAt: user.CreatedAt}
	response.CustomerEnabled = user.Status == "ACTIVE"
	response.ArtisanStatus = "NOT_STARTED"
	response.Capabilities = []string{}
	if response.CustomerEnabled {
		response.Capabilities = []string{
			"customer.account.read",
			"customer.profile.read",
			"customer.profile.write",
			"customer.addresses.read",
			"customer.addresses.write",
			"customer.purchases.read",
			"customer.tracking.read",
		}
		for _, role := range user.Roles {
			if role == "administrator" {
				response.Capabilities = append(response.Capabilities, "admin.artisan_applications.read")
				break
			}
		}
	}
	return response
}

func applyAccountSummary(response UserResponse, summary customer.AccountSummary) UserResponse {
	response.CustomerEnabled = summary.CustomerEnabled
	response.ArtisanStatus = summary.ArtisanStatus
	response.ArtisanEnabled = summary.ArtisanEnabled
	response.Capabilities = summary.Capabilities
	return response
}
func tokenResponseFrom(tokens auth.Tokens) TokenResponse {
	return TokenResponse{AccessToken: tokens.AccessToken, RefreshToken: tokens.RefreshToken, TokenType: "Bearer", ExpiresAt: tokens.ExpiresAt}
}
