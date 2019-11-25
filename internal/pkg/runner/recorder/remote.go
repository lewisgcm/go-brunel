/*
 * Author: Lewis Maitland
 *
 * Copyright (c) 2019 Lewis Maitland
 */

package recorder

import (
	"github.com/pkg/errors"
	"go-brunel/internal/pkg/runner/remote"
	"go-brunel/internal/pkg/shared"
)

type RemoteRecorder struct {
	Remote remote.Remote
}

func (recorder *RemoteRecorder) RecordLog(jobID shared.JobID, log string, logType shared.LogType, stage string) error {
	return errors.Wrap(
		recorder.Remote.Log(jobID, log, logType, stage),
		"error recording log",
	)
}

func (recorder *RemoteRecorder) RecordContainer(jobID shared.JobID, containerID shared.ContainerID, meta shared.ContainerMeta, container shared.Container, state shared.ContainerState) error {
	return errors.Wrap(
		recorder.Remote.AddContainer(jobID, containerID, meta, container, state),
		"error recording container",
	)
}

func (recorder *RemoteRecorder) RecordContainerState(containerID shared.ContainerID, state shared.ContainerState) error {
	return errors.Wrap(
		recorder.Remote.SetContainerState(containerID, state),
		"error updating container status",
	)
}

func (recorder *RemoteRecorder) RecordContainerLog(containerID shared.ContainerID, log string, logType shared.LogType) error {
	return errors.Wrap(
		recorder.Remote.ContainerLog(containerID, log, logType),
		"error recording container log",
	)
}
