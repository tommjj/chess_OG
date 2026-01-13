package domain

import (
	"time"

	"github.com/google/uuid"
)

type Role string

const (
	RoleAdmin Role = "ADMIN"
	RoleUser  Role = "USER"
	RoleMod   Role = "MODERATOR"
	RoleGuest Role = "GUEST"
)

type User struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"` // unique
	Email    string    `json:"email"`    // for common identification
	Avatar   string    `json:"avatar"`
	Role     Role      `json:"role"`

	CreatedAt time.Time `json:"created_at"`
	UpdateAt  time.Time `json:"updated_at"`
}
