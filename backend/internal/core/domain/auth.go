package domain

import (
	"time"

	"github.com/google/uuid"
)

// Auth provider constants
const (
	ProviderGoogle = "google"
	ProviderGitHub = "github"
)

// OAuthUser struct for oauth user data
type OAuthUser struct {
	ID      string `json:"id"`
	Name    string `json:"name"` // provider username
	Email   string `json:"email"`
	Picture string `json:"picture"`
}

// OAuthToken struct for oauth token
type OAuthToken struct {
	AccessToken  string    `json:"access_token,omitempty"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	Expiry       time.Time `json:"expiry,omitempty"`
	Provider     string    `json:"provider,omitempty"`
}

// Token struct for access and refresh tokens
type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// AuthSession struct for authentication session
type AuthSession struct {
	ID           uuid.UUID `json:"id"`
	UserID       uuid.UUID `json:"user_id"`
	SessionToken string    `json:"session_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}

// Account struct for user account linked to OAuth provider
type Account struct {
	UserID uuid.UUID `json:"user_id"`

	Provider   string `json:"provider"`    // e.g., "google", "github"
	ProviderID string `json:"provider_id"` // unique id from provider

	AccessToken  string    `json:"-"`
	RefreshToken string    `json:"-"`
	ExpiresAt    time.Time `json:"expires_at"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// PasswordCredential struct for storing user password credentials
// Note: Remove in future if use passwordless auth only
type PasswordCredential struct {
	UserID             uuid.UUID `json:"user_id"`
	Email              string    `json:"email"` // unique
	PasswordHash       string    `json:"-"`
	LastPasswordChange time.Time `json:"last_password_change"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TokenPayload struct for JWT token payload
type TokenPayload struct {
	UserID   uuid.UUID `json:"user_id"`
	Username string    `json:"username"`
	Avatar   string    `json:"avatar"`
	Role     Role      `json:"role"`
}
