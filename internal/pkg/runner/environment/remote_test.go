package environment_test

import (
	"github.com/golang/mock/gomock"
	"go-brunel/internal/pkg/runner/environment"
	"go-brunel/internal/pkg/shared"
	"go-brunel/test/mocks/go-brunel/pkg/runner/remote"
	"testing"
)

func TestRemoteEnvironment_GetValue(t *testing.T) {
	var envName = "a"

	id := shared.EnvironmentID("testy")
	controller := gomock.NewController(t)
	mockRemote := remote.NewMockRemote(controller)

	mockRemote.EXPECT().
		GetEnvironmentVariable(gomock.Eq(id), gomock.Eq(envName)).
		Times(1)

	factory := environment.RemoteEnvironmentFactory{
		Remote: mockRemote,
	}
	provider := factory.Create(&id)

	_, _ = provider.GetVariable(envName)
}
