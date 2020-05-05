package store

import (
	"github.com/pkg/errors"
	"go-brunel/internal/pkg/shared"
	"regexp"
	"time"
)

type RepositoryID string
type RepositoryTriggerType int8

var ErrorNotFound = errors.New("entity not found")

const (
	RepositoryTriggerTypeTag    RepositoryTriggerType = 0
	RepositoryTriggerTypeBranch RepositoryTriggerType = 1
)

type RepositoryTrigger struct {
	Type          RepositoryTriggerType
	Pattern       string
	EnvironmentID *shared.EnvironmentID `bson:"environment_id"`
}

type Repository struct {
	ID        RepositoryID `bson:"-"`
	Project   string
	Name      string
	URI       string
	Triggers  []RepositoryTrigger `bson:",omitempty"`
	CreatedAt time.Time           `bson:"created_at"`
	UpdatedAt time.Time           `bson:"updated_at"`
	DeletedAt *time.Time          `bson:"deleted_at" json:",omitempty"`
}

func (repository *Repository) IsValid() bool {
	return repository.Name != "" && repository.Project != "" && repository.URI != ""
}

func (trigger *RepositoryTrigger) IsValid() bool {
	if trigger.Type != RepositoryTriggerTypeTag && trigger.Type != RepositoryTriggerTypeBranch {
		return false
	}

	if _, e := regexp.Compile(trigger.Pattern); e != nil {
		return false
	}

	return true
}

type RepositoryStore interface {
	AddOrUpdate(repository Repository) (Repository, error)

	SetTriggers(id RepositoryID, triggers []RepositoryTrigger) error

	Get(id RepositoryID) (Repository, error)

	Filter(filter string) ([]Repository, error)
}
