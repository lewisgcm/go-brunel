package environment

import "go-brunel/internal/pkg/shared"

type Provider interface {
	GetVariable(name string) (string, error)
}

type Factory interface {
	Create(id *shared.EnvironmentID) Provider
}
