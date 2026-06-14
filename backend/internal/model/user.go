package model

import "time"

type User struct {
	ID           int64     `db:"id"`
	FullName     string    `db:"full_name"`
	Email        string    `db:"email"`
	PasswordHash string    `db:"password_hash"`
	RoleID       int64     `db:"role_id"`
	IsActive     bool      `db:"is_active"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`

	Role *Role `db:"-"`
}
