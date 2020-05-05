package environment

import (
	"errors"
	"go-brunel/internal/pkg/runner/remote"
	"go-brunel/internal/pkg/shared"
)

type remoteEnvironment struct {
	remote remote.Remote
	id     *shared.EnvironmentID
}

type RemoteEnvironmentFactory struct {
	Remote remote.Remote
}

func (envFactory *RemoteEnvironmentFactory) Create(id *shared.EnvironmentID) Provider {
	return &remoteEnvironment{
		id:     id,
		remote: envFactory.Remote,
	}
}

func (env *remoteEnvironment) GetSecret(name string) (string, error) {
	if env.id == nil {
		return "", errors.New("no environment has been configured")
	}
	return env.remote.GetEnvironmentSecret(*env.id, name)
}

func (env *remoteEnvironment) GetValue(name string) (string, error) {
	if env.id == nil {
		return "", errors.New("no environment has been configured")
	}
	return env.remote.GetEnvironmentValue(*env.id, name)
}
