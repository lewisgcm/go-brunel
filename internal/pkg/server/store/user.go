package store

import (
	"go-brunel/internal/pkg/server/security"
	"time"
)

type User struct {
	ID        string `bson:"-"`
	Username  string
	Email     string
	Name      string
	AvatarURL string `bson:"avatar_url"`
	Role      security.UserRole
	CreatedAt time.Time `bson:"created_at"`
}

type UserStore interface {
	AddOrUpdate(user User) (User, error)

	GetByUsername(username string) (User, error)
}
