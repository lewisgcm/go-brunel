package store

import (
	"go-brunel/internal/pkg/server/security"
	"time"
)

type User struct {
	Username  string
	Email     string
	Name      string
	AvatarURL string
	Role      security.UserRole
	CreatedAt time.Time
}

type UserList struct {
	Username string
	Role     security.UserRole
}

type UserStore interface {
	Filter(filter string) ([]UserList, error)

	AddOrUpdate(user User) (*User, error)

	GetByUsername(username string) (*User, error)

	Delete(username string, hard bool) error
}
