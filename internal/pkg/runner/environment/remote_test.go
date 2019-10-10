package environment_test

import (
	"github.com/golang/mock/gomock"
	"go-brunel/internal/pkg/runner/environment"
	"go-brunel/test/mocks/go-brunel/pkg/runner/remote"
	"testing"
)

func TestRemoteEnvironment_GetValue(t *testing.T) {
	var envName = "a"
	var paths = []string{"a", "b"}

	controller := gomock.NewController(t)
	mockRemote := remote.NewMockRemote(controller)

	// We expect to send a request for the env variable, with the search paths we provided on creation
	mockRemote.EXPECT().
		SearchForValue(gomock.Eq(paths), gomock.Eq(envName)).
		Times(1)

	factory := environment.RemoteEnvironmentFactory{
		Remote: mockRemote,
	}
	provider := factory.Create(paths)

	_, _ = provider.GetValue(envName)
}

func TestRemoteEnvironment_GetSecret(t *testing.T) {
	var envName = "a"
	var paths = []string{"a", "b"}

	controller := gomock.NewController(t)
	mockRemote := remote.NewMockRemote(controller)

	// We expect to send a request for the secret variable, with the search paths we provided on creation
	mockRemote.EXPECT().
		SearchForSecret(gomock.Eq(paths), gomock.Eq(envName)).
		Times(1)

	factory := environment.RemoteEnvironmentFactory{
		Remote: mockRemote,
	}
	provider := factory.Create(paths)

	_, _ = provider.GetSecret(envName)
}
