package runner

import (
	dockerclient "github.com/docker/docker/client"
	"github.com/pkg/errors"
	"go-brunel/internal/pkg/runner/environment"
	"go-brunel/internal/pkg/runner/pipeline"
	"go-brunel/internal/pkg/runner/recorder"
	"go-brunel/internal/pkg/runner/remote"
	"go-brunel/internal/pkg/runner/runtime"
	"go-brunel/internal/pkg/runner/trigger"
	"go-brunel/internal/pkg/runner/vcs"
	"go-brunel/internal/pkg/shared"
	credentials "go-brunel/internal/pkg/shared/remote"
)

type Config struct {
	// WorkingDirectory for the job, if nil it will default to the current working directory.
	// Jobs are cloned in this directory in the format WorkingDirectory + jobID, in kubernetes it should be the root of the
	// volumeClaimName volume mount location.
	WorkingDirectory string `mapstructure:"working-directory"`

	Runtime    shared.RuntimeType
	Kubernetes *shared.KubernetesConfig

	Remote *struct {
		Endpoint    string
		Credentials *credentials.Credentials
	}
}

func (config *Config) Valid() error {

	if config.Remote != nil && config.Remote.Credentials == nil {
		return errors.New("no remote.credentials have been supplied for connecting to remote server")
	} else if config.Runtime == shared.RuntimeTypeKubernetes && config.Kubernetes == nil {
		return errors.New("kubernetes configuration should be provided when using kubernetes as a runtime")
	} else if config.WorkingDirectory == "" {
		return errors.New("working-directory should not be empty")
	}

	return nil
}

func (config *Config) Trigger() (trigger.Trigger, error) {
	if config.Remote != nil {
		r, e := config.remote()
		if e != nil {
			return nil, e
		}
		return &trigger.RemoteTrigger{
			Remote:      r,
			BaseWorkDir: config.WorkingDirectory,
		}, nil
	}
	return &trigger.LocalTrigger{WorkDir: config.WorkingDirectory}, nil
}

func (config *Config) JobHandler() (*pipeline.JobHandler, error) {
	jobRecorder, err := config.recorder()
	if err != nil {
		return nil, err
	}

	runtimeFactory, err := config.runtimeFactory()
	if err != nil {
		return nil, err
	}

	environmentFactory, err := config.environment()
	if err != nil {
		return nil, err
	}

	return &pipeline.JobHandler{
		RuntimeFactory: runtimeFactory,
		Recorder:       jobRecorder,
		WorkSpace: &pipeline.LocalWorkSpace{
			VCS:                &vcs.GitVCS{},
			EnvironmentFactory: environmentFactory,
			Recorder:           jobRecorder,
		},
	}, nil
}

func (config *Config) runtimeFactory() (runtime.Factory, error) {
	switch config.Runtime {
	case shared.RuntimeTypeKubernetes:
		client, err := config.Kubernetes.GetKubernetesRESTClient()
		if err != nil {
			return nil, err
		}
		return &runtime.KubeRuntimeFactory{
			Client:          client,
			Namespace:       config.Kubernetes.Namespace,
			VolumeClaimName: config.Kubernetes.VolumeClaimName,
		}, nil
	}

	client, err := dockerclient.NewEnvClient()
	if err != nil {
		return nil, errors.Wrap(err, "error creating docker client")
	}
	return &runtime.DockerRuntimeFactory{
		Client: client,
	}, nil
}

func (config *Config) remote() (remote.Remote, error) {
	if config.Remote != nil {
		return remote.NewRPCClient(*config.Remote.Credentials, config.Remote.Endpoint)
	}
	return nil, nil
}

func (config *Config) recorder() (recorder.Recorder, error) {
	if config.Remote != nil {
		r, e := config.remote()
		if e != nil {
			return nil, e
		}
		return &recorder.RemoteRecorder{
			Remote: r,
		}, nil
	}
	return &recorder.LocalRecorder{}, nil
}

func (config *Config) environment() (environment.Factory, error) {
	if config.Remote != nil {
		r, e := config.remote()
		if e != nil {
			return nil, e
		}
		return &environment.RemoteEnvironmentFactory{
			Remote: r,
		}, nil
	}
	return &environment.LocalEnvironmentFactory{}, nil
}
