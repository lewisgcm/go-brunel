// +build storeIntegrationTests !unit

package store

import (
	"errors"
	"go-brunel/internal/pkg/server/store"
	"go-brunel/internal/pkg/shared"
	"go-brunel/test"
	"testing"
	"time"
)

func TestAddEnvironment(t *testing.T) {
	suite := setup(t)

	for _, environmentStore := range suite.environmentStores {
		now := time.Now()

		env, err := environmentStore.AddOrUpdate(store.Environment{
			Name: "testing",
			Variables: []store.EnvironmentVariable{
				{
					Name:  "var",
					Value: "val",
					Type:  store.EnvironmentVariableTypeText,
				},
			},
		})

		if err != nil {
			t.Fatalf("error saving environment %s", err.Error())
		}

		if e := environmentStore.Delete(env.ID, true); e != nil {
			t.Fatalf("could not delete old environment %e", e)
		}

		if env.ID == "" {
			t.Errorf("no environment id returned")
		}

		if env.Name != "testing" {
			t.Errorf("name was not returned correctly: %s != %s", env.Name, "testing")
		}

		if env.CreatedAt.Before(now) {
			t.Errorf("created at date not returned correctly: %s before %s", env.CreatedAt.String(), now.String())
		}

		if env.UpdatedAt.Before(now) {
			t.Errorf("updated at date not returned correctly: %s before %s", env.UpdatedAt.String(), now.String())
		}

		if env.DeletedAt != nil {
			t.Errorf("deleted at should not be set")
		}

		if len(env.Variables) != 1 || env.Variables[0].Name != "var" || env.Variables[0].Value != "val" || env.Variables[0].Type != store.EnvironmentVariableTypeText {
			t.Errorf("environment variables not returned")
		}
	}
}

func TestAddAndUpdateEnvironment(t *testing.T) {
	suite := setup(t)

	for _, environmentStore := range suite.environmentStores {
		env, err := environmentStore.AddOrUpdate(store.Environment{
			Name: "testing",
			Variables: []store.EnvironmentVariable{
				{
					Name:  "var",
					Value: "val",
					Type:  store.EnvironmentVariableTypeText,
				},
			},
		})

		if err != nil {
			t.Fatalf("error saving environment %s", err.Error())
		}

		now := time.Now()
		env.Name = "testing2"
		env.Variables = []store.EnvironmentVariable{
			{
				Name:  "var1",
				Value: "val1",
				Type:  store.EnvironmentVariableTypePassword,
			},
		}
		env.CreatedAt = time.Now()
		updatedEnv, err := environmentStore.AddOrUpdate(*env)

		if err != nil {
			t.Fatalf("error updating environment %s", err.Error())
		}

		if e := environmentStore.Delete(env.ID, true); e != nil {
			t.Fatalf("could not delete old environment %e", e)
		}

		if updatedEnv.ID != env.ID {
			t.Errorf("returned id does not match old id")
		}

		if updatedEnv.Name != "testing2" {
			t.Errorf("name was not returned correctly")
		}

		if updatedEnv.CreatedAt.Equal(env.CreatedAt) {
			t.Errorf("created at date should not have updated")
		}

		if updatedEnv.UpdatedAt.Before(now) {
			t.Errorf("updated at date not returned correctly %s", env.UpdatedAt.String())
		}

		if env.DeletedAt != nil {
			t.Errorf("deleted at should not be set")
		}

		if len(updatedEnv.Variables) != 1 || updatedEnv.Variables[0].Name != "var1" || updatedEnv.Variables[0].Value != "val1" || updatedEnv.Variables[0].Type != store.EnvironmentVariableTypePassword {
			t.Errorf("environment variables not updated")
		}
	}
}

func TestUniqueNameEnvironment(t *testing.T) {
	suite := setup(t)

	for _, environmentStore := range suite.environmentStores {
		env, err := environmentStore.AddOrUpdate(store.Environment{
			Name: "testing",
			Variables: []store.EnvironmentVariable{
				{
					Name:  "var",
					Value: "val",
					Type:  store.EnvironmentVariableTypeText,
				},
			},
		})

		if err != nil {
			t.Fatalf("error saving environment %s", err.Error())
		}

		_, err = environmentStore.AddOrUpdate(store.Environment{
			Name: "testing",
			Variables: []store.EnvironmentVariable{
				{
					Name:  "var",
					Value: "val",
					Type:  store.EnvironmentVariableTypeText,
				},
			},
		})

		if e := environmentStore.Delete(env.ID, true); e != nil {
			t.Fatalf("could not delete old environment %e", e)
		}

		test.ExpectError(t, errors.New("environment name must be unique"), err)
	}
}

func TestGetEnvironment(t *testing.T) {
	suite := setup(t)

	for _, environmentStore := range suite.environmentStores {
		env, err := environmentStore.AddOrUpdate(store.Environment{
			Name: "testing",
			Variables: []store.EnvironmentVariable{
				{
					Name:  "var",
					Value: "val",
					Type:  store.EnvironmentVariableTypeText,
				},
			},
		})

		if err != nil {
			t.Fatalf("error saving environment %s", err.Error())
		}

		getEnv, err := environmentStore.Get(env.ID)

		if e := environmentStore.Delete(env.ID, true); e != nil {
			t.Fatalf("could not delete old environment %e", e)
		}

		if getEnv == nil || err != nil {
			t.Fatalf("error getting environment")
		}

		if env.ID != getEnv.ID {
			t.Errorf("environment ids do not match")
		}

		if getEnv.Name != "testing" {
			t.Errorf("environment name was not fetched correctly")
		}

		if !getEnv.CreatedAt.Equal(env.CreatedAt) {
			t.Errorf("created at dates do not match")
		}

		if !getEnv.UpdatedAt.Equal(env.UpdatedAt) {
			t.Errorf("updated at dates do not match")
		}

		if len(getEnv.Variables) != 1 || getEnv.Variables[0].Name != "var" || getEnv.Variables[0].Value != "val" || getEnv.Variables[0].Type != store.EnvironmentVariableTypeText {
			t.Errorf("environment variables not fetched correctly")
		}
	}
}

func TestGetEnvironmentNotFound(t *testing.T) {
	suite := setup(t)

	for _, environmentStore := range suite.environmentStores {
		_, err := environmentStore.Get(shared.EnvironmentID("5eb1d158a610b1d1024f0d59"))
		test.ExpectError(t, store.ErrorNotFound, err)
	}
}

func TestGetEnvironmentVariableNotFound(t *testing.T) {
	suite := setup(t)

	for _, environmentStore := range suite.environmentStores {
		_, err := environmentStore.GetVariable(shared.EnvironmentID("5eb1d158a610b1d1024f0d59"), "tesy")
		test.ExpectErrorLike(t, errors.New("error getting environment"), err)
	}
}

func TestGetEnvironmentVariable(t *testing.T) {
	suite := setup(t)

	for _, environmentStore := range suite.environmentStores {
		env, err := environmentStore.AddOrUpdate(store.Environment{
			Name: "testing",
			Variables: []store.EnvironmentVariable{
				{
					Name:  "var",
					Value: "val",
					Type:  store.EnvironmentVariableTypeText,
				},
			},
		})

		if err != nil {
			t.Fatalf("error saving environment %s", err.Error())
		}

		variable, err := environmentStore.GetVariable(env.ID, "var")
		if variable == nil || err != nil || *variable != "val" {
			t.Errorf("environment variable was not expected value")
		}

		_, err = environmentStore.GetVariable(env.ID, "vaz")
		test.ExpectErrorLike(t, errors.New("environment variable not found"), err)

		if e := environmentStore.Delete(env.ID, true); e != nil {
			t.Fatalf("could not delete old environment %e", e)
		}
	}
}

func TestFilterEnvironments(t *testing.T) {
	suite := setup(t)

	for _, environmentStore := range suite.environmentStores {
		env, err := environmentStore.AddOrUpdate(store.Environment{
			Name: "testing",
			Variables: []store.EnvironmentVariable{
				{
					Name:  "var",
					Value: "val",
					Type:  store.EnvironmentVariableTypeText,
				},
			},
		})

		if err != nil {
			t.Fatalf("error saving environment %s", err.Error())
		}

		environments, err := environmentStore.Filter("TESTING")
		if len(environments) == 0 || err != nil {
			t.Errorf("incorrect environments returned for 'TESTING'")
		}

		environments, err = environmentStore.Filter("TesTiNG")
		if len(environments) == 0 || err != nil {
			t.Errorf("incorrect environments returned for 'TesTiNG'")
		}

		environments, err = environmentStore.Filter("")
		if len(environments) == 0 || err != nil {
			t.Errorf("incorrect environments returned for ''")
		}

		environments, err = environmentStore.Filter("LOL")
		if len(environments) != 0 || err != nil {
			t.Errorf("incorrect environments returned for 'LOL'")
		}

		if e := environmentStore.Delete(env.ID, true); e != nil {
			t.Fatalf("could not delete old environment %e", e)
		}
	}
}
