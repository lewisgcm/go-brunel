package store

import (
	"fmt"
	"github.com/pkg/errors"
	"go-brunel/internal/pkg/shared"
	"regexp"
	"strings"
	"time"
)

type RepositoryID string

type RepositoryTriggerType int8

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
	Triggers  []RepositoryTrigger
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func (repository *Repository) Clean() {
	repository.Name = strings.TrimSpace(repository.Name)
	repository.Project = strings.TrimSpace(repository.Project)
	repository.URI = strings.TrimSpace(repository.URI)
}

func (repository *Repository) IsValid() error {
	if len(repository.Name) == 0 || len(repository.Project) == 0 || len(repository.URI) != 0 {
		return errors.New("repository name, project, and uri are required")
	}
	return nil
}

func (trigger *RepositoryTrigger) IsValid() error {
	if trigger.Type != RepositoryTriggerTypeTag && trigger.Type != RepositoryTriggerTypeBranch {
		return fmt.Errorf("unknown trigger type: %d", trigger.Type)
	}

	if _, e := regexp.Compile(trigger.Pattern); e != nil {
		return errors.Wrap(e, "invalid trigger pattern")
	}

	return nil
}

type RepositoryStore interface {
	AddOrUpdate(repository Repository) (*Repository, error)

	SetTriggers(id RepositoryID, triggers []RepositoryTrigger) error

	Get(id RepositoryID) (*Repository, error)

	Filter(filter string) ([]Repository, error)

	Delete(id RepositoryID, hard bool) error
}
