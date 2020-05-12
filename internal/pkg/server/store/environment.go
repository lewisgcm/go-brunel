package store

import (
	"go-brunel/internal/pkg/shared"
	"strings"
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
	ID        shared.EnvironmentID
	Name      string
	Variables []EnvironmentVariable
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type EnvironmentVariable struct {
	Name  string
	Value string
	Type  EnvironmentVariableType
}

func (environment *Environment) Clean() {
	environment.Name = strings.TrimSpace(environment.Name)

	for _, variable := range environment.Variables {
		variable.Name = strings.TrimSpace(variable.Name)
	}
}

func (environment *Environment) IsValid() (string, bool) {
	environment.Clean()

	if len(environment.Name) == 0 {
		return "environment name cannot be empty", false
	}

	for i, variable := range environment.Variables {
		if len(variable.Name) == 0 {
			return "environment variable name cannot be empty", false
		}

		if variable.Type != EnvironmentVariableTypePassword && variable.Type != EnvironmentVariableTypeText {
			return "environment variable type must be either password or text", false
		}

		for j, other := range environment.Variables {
			if i != j && variable.Name == other.Name {
				return "environment variable names must be unique", false
			}
		}
	}

	return "", true
}

type EnvironmentStore interface {
	Filter(filter string) ([]EnvironmentList, error)

	Get(id shared.EnvironmentID) (*Environment, error)

	AddOrUpdate(environment Environment) (*Environment, error)

	GetVariable(id shared.EnvironmentID, name string) (*string, error)

	Delete(id shared.EnvironmentID, hard bool) error
}
