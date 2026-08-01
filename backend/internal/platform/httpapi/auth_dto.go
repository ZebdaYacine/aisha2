package httpapi

import (
	"github.com/aisha-platform/aisha/backend/internal/auth"
	"time"
)

type RegisterRequest struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	DisplayName string `json:"displayName"`
}
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type RefreshRequest struct {
	RefreshToken string `json:"refreshToken"`
}
type LogoutRequest struct {
	RefreshToken string `json:"refreshToken"`
}
type ForgotPasswordRequest struct {
	Email string `json:"email"`
}
type ResetPasswordRequest struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}
type UserResponse struct {
	ID          string    `json:"id"`
	Email       string    `json:"email"`
	DisplayName string    `json:"displayName"`
	Status      string    `json:"status"`
	Roles       []string  `json:"roles"`
	CreatedAt   time.Time `json:"createdAt"`
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
	return UserResponse{ID: user.ID, Email: user.Email, DisplayName: user.DisplayName, Status: user.Status, Roles: user.Roles, CreatedAt: user.CreatedAt}
}
func tokenResponseFrom(tokens auth.Tokens) TokenResponse {
	return TokenResponse{AccessToken: tokens.AccessToken, RefreshToken: tokens.RefreshToken, TokenType: "Bearer", ExpiresAt: tokens.ExpiresAt}
}
