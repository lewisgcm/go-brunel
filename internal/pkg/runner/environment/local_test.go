package environment_test

import (
	"github.com/pkg/errors"
	"go-brunel/internal/pkg/runner/environment"
	"os"
	"testing"
)

// TestLocalEnvironment_GetValue tests that environment variables can be retrieved locally
// using the LocalEnvironmentFactory implementation of the environment.Factory interface
func TestLocalEnvironment_GetValue(t *testing.T) {
	envName := "a"
	envVal := "val"

	if err := os.Setenv(envName, envVal); err != nil {
		t.Error(err)
		t.FailNow()
	}

	factory := environment.LocalEnvironmentFactory{}
	provider := factory.Create(nil)
	v, e := provider.GetValue(envName)
	if e != nil {
		t.Fail()
		t.Error(e)
	}
	if v != envVal {
		t.Fail()
		t.Errorf("expecting value '%s' but got '%s'", envVal, v)
	}

	_, e = provider.GetValue(envName + "aaaa")
	if e == nil {
		t.Fail()
		t.Error(errors.New("expecting error but no error occurred"))
	}

	if err := os.Unsetenv(envName); err != nil {
		t.Error(err)
		t.FailNow()
	}
}

// TestLocalEnvironment_GetValue tests that environment secrets can be retrieved locally
// using the LocalEnvironmentFactory implementation of the environment.Factory interface
func TestLocalEnvironment_GetSecret(t *testing.T) {
	envName := "a"
	envVal := "val"

	if err := os.Setenv(envName, envVal); err != nil {
		t.Error(err)
		t.FailNow()
	}

	factory := environment.LocalEnvironmentFactory{}
	provider := factory.Create(nil)
	v, e := provider.GetSecret(envName)
	if e != nil {
		t.Fail()
		t.Error(e)
	}
	if v != envVal {
		t.Fail()
		t.Errorf("expecting value '%s' but got '%s'", envVal, v)
	}

	_, e = provider.GetSecret(envName + "aaaa")
	if e == nil {
		t.Fail()
		t.Error(errors.New("expecting error but no error occurred"))
	}

	if err := os.Unsetenv(envName); err != nil {
		t.Error(err)
		t.FailNow()
	}
}
