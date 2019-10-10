package store

import (
	"time"

	"github.com/pkg/errors"
)

var ErrorNotFound = errors.New("not found")

type Repository struct {
	ID        string `bson:"-"`
	Project   string
	Name      string
	URI       string
	CreatedAt time.Time `bson:"created_at"`
}

type RepositoryStore interface {
	AddOrUpdate(repository Repository) (Repository, error)

	Get(id string) (Repository, error)

	Filter(filter string) ([]Repository, error)
}
