package environment

import "go-brunel/internal/pkg/runner/remote"

type remoteEnvironment struct {
	remote     remote.Remote
	searchPath []string
}

type RemoteEnvironmentFactory struct {
	Remote remote.Remote
}

func (envFactory *RemoteEnvironmentFactory) Create(searchPath []string) Provider {
	return &remoteEnvironment{
		searchPath: searchPath,
		remote:     envFactory.Remote,
	}
}

func (env *remoteEnvironment) GetSecret(name string) (string, error) {
	return env.remote.SearchForSecret(env.searchPath, name)
}

func (env *remoteEnvironment) GetValue(name string) (string, error) {
	return env.remote.SearchForValue(env.searchPath, name)
}
