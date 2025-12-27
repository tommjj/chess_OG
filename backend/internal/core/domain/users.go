package domain

import (
	"time"

	"github.com/google/uuid"
)

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
	RoleMod   Role = "moderator"
	RoleGuest Role = "guest"
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
