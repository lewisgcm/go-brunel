package environment

import (
	"go-brunel/internal/pkg/runner/remote"
	"go-brunel/internal/pkg/shared"
)

type remoteEnvironment struct {
	remote remote.Remote
	id     shared.JobID
}

type RemoteEnvironmentFactory struct {
	Remote remote.Remote
}

func (envFactory *RemoteEnvironmentFactory) Create(id shared.JobID) Provider {
	return &remoteEnvironment{
		id:     id,
		remote: envFactory.Remote,
	}
}

func (env *remoteEnvironment) GetSecret(name string) (string, error) {
	return env.remote.SearchForSecret(env.id, name)
}

func (env *remoteEnvironment) GetValue(name string) (string, error) {
	return env.remote.SearchForValue(env.id, name)
}
