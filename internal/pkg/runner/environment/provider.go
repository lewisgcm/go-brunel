package environment

import "go-brunel/internal/pkg/shared"

type Provider interface {
	GetSecret(name string) (string, error)
	GetValue(name string) (string, error)
}

type Factory interface {
	Create(id shared.JobID) Provider
}
