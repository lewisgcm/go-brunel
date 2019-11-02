package environment

import (
	"bufio"
	"fmt"
	"github.com/pkg/errors"
	"go-brunel/internal/pkg/shared/util"
	"os"
	"regexp"
)

type localEnvironment struct {
	DotEnvPath string
}

type LocalEnvironmentFactory struct {
	DotEnvPath string
}

var envRegex = regexp.MustCompile(`^([^=]+)=(.+)$`)

func (envFactory *LocalEnvironmentFactory) Create(searchPath []string) Provider {
	return &localEnvironment{
		DotEnvPath: envFactory.DotEnvPath,
	}
}

func (e *localEnvironment) GetSecret(name string) (string, error) {
	return e.resolve(name)
}

func (e *localEnvironment) GetValue(name string) (string, error) {
	return e.resolve(name)
}

func (e *localEnvironment) resolve(name string) (string, error) {
	if e.DotEnvPath != "" {
		file, err := os.Open(e.DotEnvPath)
		if err != nil {
			return "", errors.Wrap(err, "error opening .env file")
		}

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			matches := envRegex.FindStringSubmatch(scanner.Text())
			if len(matches) != 3 {
				return "", util.ErrorAppend(errors.New("error parsing .env file, invalid format."), file.Close())
			}
			if name == matches[1] {
				return matches[2], errors.Wrap(file.Close(), "error closing .env file")
			}
		}
		err = file.Close()
		if err != nil {
			return "", errors.Wrap(err, "error closing .env file")
		}
	}

	if v, ok := os.LookupEnv(name); ok {
		return v, nil
	}
	return "", errors.New(fmt.Sprintf("error getting environment variable %s", name))
}
