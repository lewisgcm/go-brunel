/*
 * Author: Lewis Maitland
 *
 * Copyright (c) 2019 Lewis Maitland
 */

package runtime

import (
	"context"
	"go-brunel/internal/pkg/shared"
	"io"
	"io/ioutil"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	dockercontainer "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	dockerclient "github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/go-units"
	"github.com/pkg/errors"
)

type DockerRuntime struct {
	Client  dockerclient.CommonAPIClient
	WorkDir string
}

func (pipeline *DockerRuntime) Initialize(ctx context.Context, jobID shared.JobID, workDir string) error {
	pipeline.WorkDir = workDir

	_, err := pipeline.
		Client.
		NetworkCreate(
			ctx,
			string(jobID),
			types.NetworkCreate{},
		)
	return errors.Wrap(err, "failed to initialize runner")
}

func (pipeline *DockerRuntime) Terminate(ctx context.Context, jobID shared.JobID) error {
	return errors.Wrap(
		pipeline.
			Client.
			NetworkRemove(ctx, string(jobID)),
		"error terminating pipeline",
	)
}

func (pipeline *DockerRuntime) WaitForContainer(ctx context.Context, id shared.ContainerID, condition shared.ContainerWaitCondition) error {
	for {
		select {
		case <-ctx.Done():
			return errors.New("context cancelled waiting for container")
		default:
			status, err := pipeline.Client.ContainerInspect(ctx, string(id))
			if err != nil {
				return errors.Wrap(err, "error inspecting container")
			}

			// If the container is exited, it has been stopped
			if status.State.Status == "exited" {
				if (condition.State&shared.ContainerWaitStopped) != 0 && status.State.ExitCode != 0 {
					return errors.New("container has exited with non zero exit status")
				} else if condition.State == shared.ContainerWaitRunning {
					return errors.New("container has exited whilst waiting for it to be running")
				}
				return nil
			} else if status.State.Status == "running" && (condition.State&shared.ContainerWaitRunning) != 0 {
				return nil
			}
			time.Sleep(time.Second)
		}
	}
}

func (pipeline *DockerRuntime) CopyLogsForContainer(ctx context.Context, id shared.ContainerID, stdOut io.WriteCloser, stdErr io.WriteCloser) error {
	reader, err := pipeline.Client.ContainerLogs(ctx, string(id), types.ContainerLogsOptions{
		ShowStderr: true,
		ShowStdout: true,
		Follow:     true,
		Timestamps: false,
	})

	if err != nil {
		return errors.Wrap(err, "error getting container logs")
	}
	_, err = stdcopy.StdCopy(stdOut, stdErr, reader)

	return errors.Wrap(err, "error copying logs")
}

func (pipeline *DockerRuntime) DispatchContainer(ctx context.Context, jobID shared.JobID, container shared.Container) (shared.ContainerID, error) {
	var mounts []mount.Mount
	if container.WorkingDir != "" {
		mounts = append(mounts, mount.Mount{
			Type:   mount.TypeBind,
			Source: pipeline.WorkDir,
			Target: container.WorkingDir,
		})
	}

	// Pull our image
	reader, err := pipeline.Client.ImagePull(ctx, container.Image, types.ImagePullOptions{})
	if err != nil {
		return shared.EmptyContainerID, errors.Wrap(err, "error pulling container image")
	}

	// Wait for the pull to complete
	_, err = io.Copy(ioutil.Discard, reader)
	if err != nil {
		return shared.EmptyContainerID, errors.Wrap(err, "error copying container pull request output")
	}

	var aliases []string
	if container.Hostname != "" {
		aliases = append(aliases, container.Hostname)
		aliases = append(aliases, strings.ToLower(container.Hostname))
	}

	var envVariables []string
	if container.Environment != nil {
		for key, value := range container.Environment {
			envVariables = append(envVariables, key+"="+value)
		}
	}

	resources := dockercontainer.Resources{}
	if container.Resources != nil && container.Resources.Limits != nil {
		// We need to convert out CPU count into dockers NanoCPUS, this is the same as dockers --cpus flag
		resources.NanoCPUs = int64(container.Resources.Limits.CPU * 1e9)

		// We also need to convert out memory string limit into something docker can use
		memory, _ := units.RAMInBytes(container.Resources.Limits.Memory)
		resources.Memory = memory
	}

	var entrypoint []string
	if container.EntryPoint != "" {
		entrypoint = []string{container.EntryPoint}
	}

	// Create a container from the image
	resp, err := pipeline.Client.ContainerCreate(ctx, &dockercontainer.Config{
		Image:      container.Image,
		Entrypoint: entrypoint,
		Cmd:        container.Args,
		Hostname:   container.Hostname,
		Env:        envVariables,
		WorkingDir: container.WorkingDir,
	}, &dockercontainer.HostConfig{
		Mounts:     mounts,
		Privileged: container.Privileged,
		Resources:  resources,
	}, &network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{
			"network": {
				NetworkID: string(jobID),
				Aliases:   aliases,
			},
		},
	}, "")
	if err != nil {
		return shared.EmptyContainerID, errors.Wrap(err, "error creating container")
	}

	if err := pipeline.Client.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return shared.ContainerID(resp.ID), errors.Wrap(err, "error starting container")
	}

	return shared.ContainerID(resp.ID), nil
}

func (pipeline *DockerRuntime) TerminateContainer(ctx context.Context, containerID shared.ContainerID) error {
	return errors.Wrap(
		pipeline.
			Client.
			ContainerRemove(
				ctx,
				string(containerID),
				types.ContainerRemoveOptions{
					Force:         true,
					RemoveVolumes: true,
				},
			),
		"error terminating container",
	)
}
