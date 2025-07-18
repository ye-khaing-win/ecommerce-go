package models

import (
	"time"
)

type Admin struct {
	ID                string     `json:"id,omitempty" db:"id"`
	FirstName         string     `json:"first_name,omitempty" db:"first_name"`
	LastName          string     `json:"last_name,omitempty" db:"last_name"`
	Email             string     `json:"email,omitempty" db:"email"`
	Password          string     `json:"password,omitempty" db:"password"`
	PasswordChangedAt *time.Time `json:"password_changed_at,omitempty" db:"password_changed_at"`
	Active            bool       `json:"active,omitempty" db:"active"`
	Role              string     `json:"role,omitempty" db:"role"`
	CreatedAt         *time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt         *time.Time `json:"updated_at,omitempty" db:"updated_at"`
}
