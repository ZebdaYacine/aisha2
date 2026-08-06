package auth

import application "github.com/aisha-platform/aisha/apps/api/internal/features/auth/application"

func hashToken(token string) string { return application.HashToken(token) }
