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
)

// Runtime is used an interface to the runtime config for handling containers.
// Implementors should make sure that init,terminate are called in order otherwise containers could leak.
type Runtime interface {
	// Initialize will ready the runtime environment, for example creating kube service/docker network
	Initialize(context context.Context, jobID shared.JobID, workDir string) error

	// DispatchContainer will dispatch a container to the runtime, it will return as soon as the container is dispatched
	DispatchContainer(context context.Context, jobID shared.JobID, container shared.Container) (shared.ContainerID, error)

	// WaitForContainer waits for a container to satisfy the waiting condition, this will block until it does
	WaitForContainer(context context.Context, id shared.ContainerID, condition shared.ContainerWaitCondition) error

	// CopyLogsForContainer will copy logs from the container to the writers, it will block until the container terminates
	CopyLogsForContainer(context context.Context, id shared.ContainerID, stdOut io.WriteCloser, stdErr io.WriteCloser) error

	// TerminateContainer will stop a container in the runtime config and remove it
	TerminateContainer(context context.Context, containerID shared.ContainerID) error

	// Terminate will remove any services/networks created during the init
	Terminate(context context.Context, pipeline shared.JobID) error
}

// Factory is used for creating instances of our runtime
// This is used to prevent propagating configuration down into lower levels of the code base.
type Factory interface {
	// Create creates an instance of a runtime, error if one could not be created.
	Create() (Runtime, error)
}
