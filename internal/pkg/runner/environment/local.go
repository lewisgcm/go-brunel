package environment

import (
	"fmt"
	"github.com/pkg/errors"
	"os"
)

type localEnvironment struct {
}

type LocalEnvironmentFactory struct {
}

func (envFactory *LocalEnvironmentFactory) Create(searchPath []string) Provider {
	return &localEnvironment{}
}

func (e *localEnvironment) GetSecret(name string) (string, error) {
	if v, ok := os.LookupEnv(name); ok {
		return v, nil
	}
	return "", errors.New(fmt.Sprintf("error getting environment variable secret %s", name))
}

func (e *localEnvironment) GetValue(name string) (string, error) {
	if v, ok := os.LookupEnv(name); ok {
		return v, nil
	}
	return "", errors.New(fmt.Sprintf("error getting environment variable %s", name))
}
