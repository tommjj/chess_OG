package schema

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/tommjj/chess_OG/backend/internal/core/domain"
	"gorm.io/gorm"
)

// RoleEnum defines the enum type for user roles in the database
type RoleEnum struct{}

func (r RoleEnum) Values() []string {
	return []string{
		string(domain.RoleAdmin),
		string(domain.RoleUser),
		string(domain.RoleMod),
		string(domain.RoleGuest),
	}
}
func (r RoleEnum) Name() string { return "role" }

// AppSchemas holds all the database schemas for the application
var AppSchemas = []any{
	&User{},
	&Account{},
	&Session{},
}

// WithDate adds created_at and updated_at timestamps to a schema
type WithDate struct {
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

// WithSoftDelete adds soft delete functionality to a schema
type WithSoftDelete struct {
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// User represents the database schema for the users table
type User struct {
	ID       uuid.UUID   `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Username string      `gorm:"uniqueIndex;not null"`
	Email    string      `gorm:"not null"`
	Avatar   string      `gorm:""`
	Role     domain.Role `gorm:"type:role;size:50;not null;default:'USER'"`

	WithDate
	WithSoftDelete

	Accounts []Account `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
	Sessions []Session `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
}

// Account represents the database schema for the accounts table
type Account struct {
	UserID       uuid.UUID    `gorm:"type:uuid;not null;index"`
	Provider     string       `gorm:"size:50;primaryKey;not null"`
	ProviderID   string       `gorm:"size:255;primaryKey;not null"`
	AccessToken  string       `gorm:""`
	RefreshToken string       `gorm:""`
	ExpiresAt    sql.NullTime `gorm:"default:null"`

	WithDate

	User User `gorm:"foreignKey:UserID;references:ID"`
}

// Session represents the database schema for the sessions table
type Session struct {
	ID           uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	UserID       uuid.UUID `gorm:"type:uuid;not null"`
	SessionToken string    `gorm:"uniqueIndex;not null"`
	ExpiresAt    time.Time `gorm:"not null"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`

	User User `gorm:"foreignKey:UserID;references:ID"`
}
