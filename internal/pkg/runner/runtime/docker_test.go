// +build unit !integration

package runtime_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"go-brunel/internal/pkg/runner/runtime"
	"go-brunel/internal/pkg/shared"
	"go-brunel/test"
	mock_client "go-brunel/test/mocks/mock_docker"
	"io/ioutil"
	"testing"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/golang/mock/gomock"
)

const (
	dockerStatusRunning = "running"
	dockerStatusExited  = "exited"
)

func TestDockerRuntime_WaitForContainer(t *testing.T) {
	suites := []struct {
		inspectArgJSON  types.ContainerJSON
		inspectArgError error
		waitCondition   shared.ContainerWaitCondition
		expectedError   error
	}{
		{
			// wait for running container when container is running
			inspectArgJSON: types.ContainerJSON{
				ContainerJSONBase: &types.ContainerJSONBase{
					State: &types.ContainerState{
						Status: dockerStatusRunning,
					},
				},
			},
			inspectArgError: nil,
			expectedError:   nil,
			waitCondition: shared.ContainerWaitCondition{
				State: shared.ContainerWaitRunning,
			},
		},
		{
			// wait for running container when container is stopped
			inspectArgJSON: types.ContainerJSON{
				ContainerJSONBase: &types.ContainerJSONBase{
					State: &types.ContainerState{
						ExitCode: 0,
						Status:   dockerStatusExited,
					},
				},
			},
			inspectArgError: nil,
			expectedError:   errors.New("container has exited whilst waiting for it to be running"),
			waitCondition: shared.ContainerWaitCondition{
				State: shared.ContainerWaitRunning,
			},
		},
		{
			// wait for running or stopped container when container is running
			inspectArgJSON: types.ContainerJSON{
				ContainerJSONBase: &types.ContainerJSONBase{
					State: &types.ContainerState{
						Status: dockerStatusRunning,
					},
				},
			},
			inspectArgError: nil,
			expectedError:   nil,
			waitCondition: shared.ContainerWaitCondition{
				State: shared.ContainerWaitRunning | shared.ContainerWaitStopped,
			},
		},
		{
			// wait for running or stopped container when container is stopped
			inspectArgJSON: types.ContainerJSON{
				ContainerJSONBase: &types.ContainerJSONBase{
					State: &types.ContainerState{
						ExitCode: 0,
						Status:   dockerStatusExited,
					},
				},
			},
			inspectArgError: nil,
			expectedError:   nil,
			waitCondition: shared.ContainerWaitCondition{
				State: shared.ContainerWaitRunning | shared.ContainerWaitStopped,
			},
		},
		{
			// wait for stopped or running container when container is stopped
			inspectArgJSON: types.ContainerJSON{
				ContainerJSONBase: &types.ContainerJSONBase{
					State: &types.ContainerState{
						Status: dockerStatusExited,
					},
				},
			},
			inspectArgError: nil,
			expectedError:   nil,
			waitCondition: shared.ContainerWaitCondition{
				State: shared.ContainerWaitRunning | shared.ContainerWaitStopped,
			},
		},
		{
			// wait for stopped container when container is stopped
			inspectArgJSON: types.ContainerJSON{
				ContainerJSONBase: &types.ContainerJSONBase{
					State: &types.ContainerState{
						Status: dockerStatusExited,
					},
				},
			},
			inspectArgError: nil,
			expectedError:   nil,
			waitCondition: shared.ContainerWaitCondition{
				State: shared.ContainerWaitStopped,
			},
		},
		{
			// wait for stopped container when container is stopped with non-zero exit
			inspectArgJSON: types.ContainerJSON{
				ContainerJSONBase: &types.ContainerJSONBase{
					State: &types.ContainerState{
						ExitCode: -1,
						Status:   dockerStatusExited,
					},
				},
			},
			inspectArgError: nil,
			expectedError:   errors.New("container has exited with non zero exit status"),
			waitCondition: shared.ContainerWaitCondition{
				State: shared.ContainerWaitStopped,
			},
		},
		{
			// wait for stopped or running container when container is stopped with non-zero exit
			inspectArgJSON: types.ContainerJSON{
				ContainerJSONBase: &types.ContainerJSONBase{
					State: &types.ContainerState{
						ExitCode: -1,
						Status:   dockerStatusExited,
					},
				},
			},
			inspectArgError: nil,
			expectedError:   errors.New("container has exited with non zero exit status"),
			waitCondition: shared.ContainerWaitCondition{
				State: shared.ContainerWaitRunning | shared.ContainerWaitStopped,
			},
		},
		{
			// wait for container when error is returned by inspect
			inspectArgJSON: types.ContainerJSON{
				ContainerJSONBase: &types.ContainerJSONBase{
					State: &types.ContainerState{
						ExitCode: -1,
						Status:   dockerStatusExited,
					},
				},
			},
			inspectArgError: errors.New("bad_error"),
			expectedError:   errors.New("error inspecting container: bad_error"),
			waitCondition: shared.ContainerWaitCondition{
				State: shared.ContainerWaitRunning | shared.ContainerWaitStopped,
			},
		},
	}

	for i, suite := range suites {
		t.Run(
			fmt.Sprintf("suites[%d]", i),
			func(t *testing.T) {
				controller := gomock.NewController(t)
				client := mock_client.NewMockCommonAPIClient(controller)

				client.
					EXPECT().
					ContainerInspect(gomock.Any(), gomock.Any()).
					Return(suite.inspectArgJSON, suite.inspectArgError).
					AnyTimes()

				dockerRuntime := runtime.DockerRuntime{
					Client: client,
				}

				err := dockerRuntime.WaitForContainer(
					context.TODO(),
					shared.ContainerID(""),
					suite.waitCondition,
				)

				test.ExpectError(t, suite.expectedError, err)
			},
		)
	}
}

func TestDockerRuntime_WaitForContainer_ContextTimeout(t *testing.T) {
	controller := gomock.NewController(t)
	client := mock_client.NewMockCommonAPIClient(controller)

	json := types.ContainerJSON{
		ContainerJSONBase: &types.ContainerJSONBase{
			State: &types.ContainerState{
				Status: dockerStatusRunning,
			},
		},
	}

	client.
		EXPECT().
		ContainerInspect(gomock.Any(), gomock.Any()).
		Return(json, nil).
		AnyTimes()

	dockerRuntime := runtime.DockerRuntime{
		Client: client,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()
	err := dockerRuntime.WaitForContainer(
		ctx,
		shared.ContainerID(""),
		shared.ContainerWaitCondition{State: shared.ContainerWaitStopped},
	)

	test.ExpectErrorLike(t, errors.New("context cancelled"), err)
}

func TestDockerRuntime_DispatchContainer_ImagePullError(t *testing.T) {
	mockError := errors.New("error_pulling_image_container")
	controller := gomock.NewController(t)
	client := mock_client.NewMockCommonAPIClient(controller)
	dockerRuntime := runtime.DockerRuntime{
		Client: client,
	}

	// Test that an error on image pull is handled, and we get an error when dispatching the container
	client.
		EXPECT().
		ImagePull(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, mockError)

	containerID, err := dockerRuntime.DispatchContainer(context.TODO(), shared.JobID(""), shared.Container{})
	test.ExpectErrorLike(t, mockError, err)
	test.ExpectString(t, string(shared.EmptyContainerID), string(containerID))

	// Test that the reader will return an error if the image pull progress reader returns an error
	mockError = errors.New("error_pulling_image_reader_container")
	errorReader := test.NoOpReadCloser{Reader: &test.ErroringReader{Error: mockError}}
	client.
		EXPECT().
		ImagePull(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&errorReader, nil)

	containerID, err = dockerRuntime.DispatchContainer(context.TODO(), shared.JobID(""), shared.Container{})

	test.ExpectErrorLike(t, mockError, err)
	test.ExpectString(t, string(shared.EmptyContainerID), string(containerID))
}

func TestDockerRuntime_DispatchContainer_CreateError(t *testing.T) {
	mockError := errors.New("error_creating_container")
	controller := gomock.NewController(t)
	client := mock_client.NewMockCommonAPIClient(controller)
	dockerRuntime := runtime.DockerRuntime{
		Client: client,
	}

	client.
		EXPECT().
		ImagePull(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&test.NoOpReadCloser{Reader: bytes.NewReader([]byte(""))}, nil)

	client.
		EXPECT().
		ContainerCreate(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(container.ContainerCreateCreatedBody{}, mockError)

	containerID, err := dockerRuntime.DispatchContainer(context.TODO(), shared.JobID(""), shared.Container{})

	test.ExpectErrorLike(t, mockError, err)
	test.ExpectString(t, string(shared.EmptyContainerID), string(containerID))
}

func TestDockerRuntime_DispatchContainer_StartError(t *testing.T) {
	mockError := errors.New("error_starting_container")
	controller := gomock.NewController(t)
	client := mock_client.NewMockCommonAPIClient(controller)
	dockerRuntime := runtime.DockerRuntime{
		Client: client,
	}

	client.
		EXPECT().
		ImagePull(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&test.NoOpReadCloser{Reader: bytes.NewReader([]byte(""))}, nil)

	client.
		EXPECT().
		ContainerCreate(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(container.ContainerCreateCreatedBody{ID: "container_id"}, nil)

	client.
		EXPECT().
		ContainerStart(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(mockError)

	containerID, err := dockerRuntime.DispatchContainer(context.TODO(), shared.JobID(""), shared.Container{})

	test.ExpectErrorLike(t, mockError, err)
	test.ExpectString(t, "container_id", string(containerID))
}

func TestDockerRuntime_DispatchContainer(t *testing.T) {
	suites := []struct {
		container                   shared.Container
		expectedDockerConfig        *container.Config
		expectedDockerHostConfig    *container.HostConfig
		expectedDockerNetworkConfig *network.NetworkingConfig
	}{
		// Simple test with command, arguments and an image all set
		{
			container: shared.Container{
				Image:      "test",
				Args:       []string{"arg1", "arg2"},
				EntryPoint: "entrypoint",
			},
			expectedDockerConfig: &container.Config{
				Image:      "test",
				Entrypoint: []string{"entrypoint"},
				Cmd:        []string{"arg1", "arg2"},
			},
			expectedDockerHostConfig: &container.HostConfig{
				Privileged: false,
			},
			expectedDockerNetworkConfig: &network.NetworkingConfig{
				EndpointsConfig: map[string]*network.EndpointSettings{
					"network": {
						NetworkID: string(""),
					},
				},
			},
		},

		// Test the hostname can be set
		{
			container: shared.Container{
				Hostname: "hostName",
			},
			expectedDockerConfig: &container.Config{
				Hostname: "hostName",
			},
			expectedDockerHostConfig: &container.HostConfig{},
			expectedDockerNetworkConfig: &network.NetworkingConfig{
				EndpointsConfig: map[string]*network.EndpointSettings{
					"network": {
						NetworkID: string(""),
						Aliases:   []string{"hostName", "hostname"}, // We expect both original and lowercase aliases because docker DNS is case sensitive
					},
				},
			},
		},

		// Test that the workspace volume will be mounted
		{
			container: shared.Container{
				WorkingDir: "workdir",
			},
			expectedDockerConfig: &container.Config{
				WorkingDir: "workdir",
			},
			expectedDockerHostConfig: &container.HostConfig{
				Mounts: []mount.Mount{ // We expect a mount point of our working directory specified in the container
					{
						Type:   mount.TypeBind,
						Source: "",
						Target: "workdir",
					},
				},
			},
			expectedDockerNetworkConfig: &network.NetworkingConfig{
				EndpointsConfig: map[string]*network.EndpointSettings{
					"network": {
						NetworkID: string(""),
					},
				},
			},
		},

		// Test that Privileged mode can be enabled
		{
			container: shared.Container{
				Privileged: true,
			},
			expectedDockerConfig: &container.Config{},
			expectedDockerHostConfig: &container.HostConfig{
				Privileged: true,
			},
			expectedDockerNetworkConfig: &network.NetworkingConfig{
				EndpointsConfig: map[string]*network.EndpointSettings{
					"network": {
						NetworkID: string(""),
					},
				},
			},
		},

		// Test that resource limits are honoured
		{
			container: shared.Container{
				Resources: &shared.ContainerResources{
					Limits: &shared.ContainerResourcesUnits{
						CPU:    0.1,
						Memory: "500m",
					},
				},
			},
			expectedDockerConfig: &container.Config{},
			expectedDockerHostConfig: &container.HostConfig{
				Resources: container.Resources{
					Memory:   524288000,
					NanoCPUs: 100000000,
				},
			},
			expectedDockerNetworkConfig: &network.NetworkingConfig{
				EndpointsConfig: map[string]*network.EndpointSettings{
					"network": {
						NetworkID: string(""),
					},
				},
			},
		},

		// Test environment variables
		{
			container: shared.Container{
				Environment: map[string]string{
					"test_key": "test_val",
				},
			},
			expectedDockerConfig: &container.Config{
				Env: []string{"test_key=test_val"},
			},
			expectedDockerHostConfig: &container.HostConfig{},
			expectedDockerNetworkConfig: &network.NetworkingConfig{
				EndpointsConfig: map[string]*network.EndpointSettings{
					"network": {
						NetworkID: string(""),
					},
				},
			},
		},
	}

	for i, suite := range suites {
		t.Run(
			fmt.Sprintf("suites[%d]", i),
			func(t *testing.T) {
				controller := gomock.NewController(t)
				client := mock_client.NewMockCommonAPIClient(controller)
				dockerRuntime := runtime.DockerRuntime{
					Client: client,
				}

				client.
					EXPECT().
					ImagePull(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(&test.NoOpReadCloser{Reader: bytes.NewReader([]byte(""))}, nil)

				client.
					EXPECT().
					ContainerCreate(gomock.Any(), gomock.Eq(suite.expectedDockerConfig), gomock.Eq(suite.expectedDockerHostConfig), gomock.Eq(suite.expectedDockerNetworkConfig), gomock.Any()).
					Return(container.ContainerCreateCreatedBody{ID: "container_id"}, nil)

				client.
					EXPECT().
					ContainerStart(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)

				containerID, err := dockerRuntime.DispatchContainer(context.TODO(), shared.JobID(""), suite.container)

				test.ExpectError(t, nil, err)
				test.ExpectString(t, "container_id", string(containerID))
			},
		)
	}
}

func TestDockerRuntime_Initialize(t *testing.T) {
	controller := gomock.NewController(t)
	client := mock_client.NewMockCommonAPIClient(controller)
	dockerRuntime := runtime.DockerRuntime{
		Client: client,
	}

	// Test the case where creating the network works as expected
	client.
		EXPECT().
		NetworkCreate(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(types.NetworkCreateResponse{}, nil)

	err := dockerRuntime.Initialize(context.TODO(), shared.JobID(""), "")
	test.ExpectError(t, nil, err)

	// Test the case where creating the network returns an error
	mockError := errors.New("error_creating_network")
	client.
		EXPECT().
		NetworkCreate(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(types.NetworkCreateResponse{}, mockError)

	err = dockerRuntime.Initialize(context.TODO(), shared.JobID(""), "")
	test.ExpectErrorLike(t, mockError, err)
}

func TestDockerRuntime_Terminate(t *testing.T) {
	controller := gomock.NewController(t)
	client := mock_client.NewMockCommonAPIClient(controller)
	dockerRuntime := runtime.DockerRuntime{
		Client: client,
	}

	// Test the case where removing the network works as expected
	client.
		EXPECT().
		NetworkRemove(gomock.Any(), gomock.Any()).
		Return(nil)

	err := dockerRuntime.Terminate(context.TODO(), shared.JobID(""))
	test.ExpectError(t, nil, err)

	// Test the case where removing the network returns an error
	mockError := errors.New("error_removing_network")
	client.
		EXPECT().
		NetworkRemove(gomock.Any(), gomock.Any()).
		Return(mockError)

	err = dockerRuntime.Terminate(context.TODO(), shared.JobID(""))
	test.ExpectErrorLike(t, mockError, err)
}

func TestDockerRuntime_TerminateContainer(t *testing.T) {
	controller := gomock.NewController(t)
	client := mock_client.NewMockCommonAPIClient(controller)
	dockerRuntime := runtime.DockerRuntime{
		Client: client,
	}

	// Test the case where removing the container works as expected
	client.
		EXPECT().
		ContainerRemove(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil)

	err := dockerRuntime.TerminateContainer(context.TODO(), shared.ContainerID(""))
	test.ExpectError(t, nil, err)

	// Test the case where removing the container returns an error
	mockError := errors.New("error_removing_container")
	client.
		EXPECT().
		ContainerRemove(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(mockError)

	err = dockerRuntime.TerminateContainer(context.TODO(), shared.ContainerID(""))
	test.ExpectErrorLike(t, mockError, err)
}

func TestDockerRuntime_CopyLogsForContainer_Errors(t *testing.T) {
	controller := gomock.NewController(t)
	client := mock_client.NewMockCommonAPIClient(controller)
	dockerRuntime := runtime.DockerRuntime{
		Client: client,
	}

	// Check we are returned an error if the logs cannot be retrieved
	mockError := errors.New("error_getting_logs")
	client.
		EXPECT().
		ContainerLogs(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, mockError)

	err := dockerRuntime.CopyLogsForContainer(context.TODO(), shared.ContainerID(""), &test.NoOpWriteCloser{Writer: ioutil.Discard}, &test.NoOpWriteCloser{Writer: ioutil.Discard})
	test.ExpectErrorLike(t, mockError, err)

	// Check we are returned an error if we cannot read the logs
	mockError = errors.New("error_reading_log_reader")
	client.
		EXPECT().
		ContainerLogs(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&test.NoOpReadCloser{Reader: &test.ErroringReader{Error: mockError}}, nil)

	err = dockerRuntime.CopyLogsForContainer(context.TODO(), shared.ContainerID(""), &test.NoOpWriteCloser{Writer: ioutil.Discard}, &test.NoOpWriteCloser{Writer: ioutil.Discard})
	test.ExpectErrorLike(t, mockError, err)
}
