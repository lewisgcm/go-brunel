package environment_test

import (
	"github.com/pkg/errors"
	"go-brunel/internal/pkg/runner/environment"
	"io/ioutil"
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
	v, e := provider.GetVariable(envName)
	if e != nil {
		t.Fail()
		t.Error(e)
	}
	if v != envVal {
		t.Fail()
		t.Errorf("expecting value '%s' but got '%s'", envVal, v)
	}

	_, e = provider.GetVariable(envName + "aaaa")
	if e == nil {
		t.Fail()
		t.Error(errors.New("expecting error but no error occurred"))
	}

	if err := os.Unsetenv(envName); err != nil {
		t.Error(err)
		t.FailNow()
	}
}

// TestLocalEnvironment_DotEnv_GetValue tests that environment variables can be retrieved locally
// using the LocalEnvironmentFactory implementation of the environment.Factory interface from a
// .env file.
func TestLocalEnvironment_DotEnv_GetValue(t *testing.T) {
	envName := "a"
	envVal := "val"

	file, err := ioutil.TempFile(".", ".env.*")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	defer os.Remove(file.Name())

	_, err = file.WriteString(`a=val
b=val2`)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	filePath := file.Name()
	factory := environment.LocalEnvironmentFactory{
		DotEnvPath: filePath,
	}
	provider := factory.Create(nil)
	v, e := provider.GetVariable(envName)
	if e != nil {
		t.Fail()
		t.Error(e)
	}
	if v != envVal {
		t.Fail()
		t.Errorf("expecting value '%s' but got '%s'", envVal, v)
	}

	_, e = provider.GetVariable(envName + "aaaa")
	if e == nil {
		t.Fail()
		t.Error(errors.New("expecting error but no error occurred"))
	}
}

// TestLocalEnvironment_DotEnv_GetValue_InvalidValue tests that invalid variables in .env files
// are correctly handled
func TestLocalEnvironment_DotEnv_GetValue_InvalidValue(t *testing.T) {
	file, err := ioutil.TempFile(".", ".env.*")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	defer os.Remove(file.Name())

	_, err = file.WriteString(`ssdsd`)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	filePath := file.Name()
	factory := environment.LocalEnvironmentFactory{
		DotEnvPath: filePath,
	}
	provider := factory.Create(nil)
	_, e := provider.GetVariable("asdasd")

	if e == nil {
		t.Error("expected error message but got none")
		t.FailNow()
	}
}

// TestLocalEnvironment_DotEnv_GetValue_Non_Existent_File tests that we fail if file couldn't be found
func TestLocalEnvironment_DotEnv_GetValue_Non_Existent_File(t *testing.T) {
	filePath := ".env.dont-exist"
	factory := environment.LocalEnvironmentFactory{
		DotEnvPath: filePath,
	}
	provider := factory.Create(nil)
	_, e := provider.GetVariable("asdasd")

	if e == nil {
		t.Error("expected error message but got none")
		t.FailNow()
	}
}
