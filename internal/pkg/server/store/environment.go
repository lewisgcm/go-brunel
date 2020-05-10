package store

import (
	"go-brunel/internal/pkg/shared"
	"time"
)

type EnvironmentVariableType uint8

const (
	EnvironmentVariableTypeText     EnvironmentVariableType = 0
	EnvironmentVariableTypePassword EnvironmentVariableType = 1
)

type EnvironmentList struct {
	ID   shared.EnvironmentID `bson:"-"`
	Name string
}

type Environment struct {
	ID        shared.EnvironmentID  `bson:"-"`
	Name      string                `bson:"name"`
	Variables []EnvironmentVariable `bson:"variables"`
	CreatedAt time.Time             `bson:"created_at"`
	UpdatedAt time.Time             `bson:"updated_at" json:",omitempty"`
	DeletedAt *time.Time            `bson:"deleted_at" json:",omitempty"`
}

type EnvironmentVariable struct {
	Name  string
	Value string
	Type  EnvironmentVariableType
}

func (environment *Environment) IsValid() bool {
	return true
}

type EnvironmentStore interface {
	Filter(filter string) ([]EnvironmentList, error)

	Get(id shared.EnvironmentID) (*Environment, error)

	AddOrUpdate(environment Environment) (*Environment, error)

	GetVariable(id shared.EnvironmentID, name string) (*string, error)
}
