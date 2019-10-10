// +build integration !unit

package runtime

import (
	"context"
	"go-brunel/internal/pkg/runner/runtime"
	"go-brunel/internal/pkg/shared"
	"io"
	v12 "k8s.io/api/core/v1"
	"k8s.io/client-go/rest"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	dockerclient "github.com/docker/docker/client"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
)

type testSuite struct {
	runtimes []runtimeTestEnvironment
}

type mockBufferWriteCloser struct {
	io.Writer
}

func (mwc *mockBufferWriteCloser) Close() error {
	return nil
}

func setup(t *testing.T) testSuite {
	if testing.Short() {
		t.Skip("skipping integration tests in short mode.")
	}

	var runtimes []runtimeTestEnvironment

	// Passing in empty arguments will create a docker runtime
	dockerClient, err := dockerclient.NewEnvClient()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	dockerRuntime := runtime.DockerRuntime{
		Client: dockerClient,
	}
	runtimes = append(runtimes, &dockerRuntimeTestEnvironment{
		dockerRuntime: &dockerRuntime,
		client:        dockerClient,
	})

	// Create the kube runtime
	kubeConfig := shared.KubernetesConfig{
		ConfigFile:      "/Users/lewis/.kube/config",
		Namespace:       "brunel",
		VolumeClaimName: "brunel-workspace-volume-claim",
	}
	restClient, err := kubeConfig.GetKubernetesRESTClient()
	if err != nil {
		t.Error(err, "error configuring kubernetes runtime config")
		t.FailNow()
	}
	kubeRuntime := runtime.KubeRuntime{
		Client:          restClient,
		Namespace:       kubeConfig.Namespace,
		VolumeClaimName: kubeConfig.VolumeClaimName,
	}
	runtimes = append(runtimes, &kubeRuntimeTestEnvironment{
		kubeRuntime: &kubeRuntime,
		namespace:   "brunel",
		client:      restClient,
	})

	return testSuite{
		runtimes: runtimes,
	}
}

type runtimeTestEnvironment interface {
	runtime() runtime.Runtime
	containerExists(id shared.ContainerID, t *testing.T) bool
	jobRuntimeExists(id shared.JobID, t *testing.T) bool
}

type kubeRuntimeTestEnvironment struct {
	kubeRuntime runtime.Runtime
	client      rest.Interface
	namespace   string
}

func (env *kubeRuntimeTestEnvironment) runtime() runtime.Runtime {
	return env.kubeRuntime
}

func (env *kubeRuntimeTestEnvironment) containerExists(id shared.ContainerID, t *testing.T) bool {
	var pods v12.PodList
	err := env.
		client.
		Get().
		Resource("pods").
		Namespace(env.namespace).
		VersionedParams(&v1.ListOptions{FieldSelector: "metadata.name=" + string(id)}, v1.ParameterCodec).
		Do().
		Into(&pods)
	if err != nil {
		t.Error(errors.Wrap(err, "error getting kube pod"))
	}
	return len(pods.Items) > 0
}

func (env *kubeRuntimeTestEnvironment) jobRuntimeExists(id shared.JobID, t *testing.T) bool {
	var services v12.ServiceList
	err := env.
		client.
		Get().
		Resource("services").
		Namespace(env.namespace).
		VersionedParams(&v1.ListOptions{FieldSelector: "metadata.name=job-" + string(id)}, v1.ParameterCodec).
		Do().
		Into(&services)
	if err != nil {
		t.Error(errors.Wrap(err, "error getting kube service"))
	}
	return len(services.Items) > 0
}

type dockerRuntimeTestEnvironment struct {
	dockerRuntime runtime.Runtime
	client        *client.Client
}

func (env *dockerRuntimeTestEnvironment) runtime() runtime.Runtime {
	return env.dockerRuntime
}

func (env *dockerRuntimeTestEnvironment) containerExists(id shared.ContainerID, t *testing.T) bool {
	searchFilters := filters.NewArgs()
	searchFilters.Add("id", string(id))

	containers, err := env.client.ContainerList(context.Background(), types.ContainerListOptions{Filters: searchFilters})
	if err != nil {
		t.Error(errors.Wrap(err, "error getting docker container"))
	}
	return len(containers) > 0
}

func (env *dockerRuntimeTestEnvironment) jobRuntimeExists(id shared.JobID, t *testing.T) bool {
	searchFilters := filters.NewArgs()
	searchFilters.Add("name", string(id))
	networks, err := env.client.NetworkList(context.Background(), types.NetworkListOptions{Filters: searchFilters})
	if err != nil {
		t.Error(errors.Wrap(err, "error getting kube service"))
	}
	return len(networks) > 0
}
