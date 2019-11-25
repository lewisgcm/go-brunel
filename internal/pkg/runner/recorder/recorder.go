/*
 * Author: Lewis Maitland
 *
 * Copyright (c) 2019 Lewis Maitland
 */

package recorder

import (
	"go-brunel/internal/pkg/shared"
)

// recorder is responsible for recording information about the state of a job and its container
// the data should be considered as read only, as some implementations (text) only write this to stdout
type Recorder interface {
	RecordLog(jobID shared.JobID, log string, logType shared.LogType, stage string) error
	RecordContainer(jobID shared.JobID, containerID shared.ContainerID, meta shared.ContainerMeta, container shared.Container, state shared.ContainerState) error
	RecordContainerState(containerID shared.ContainerID, state shared.ContainerState) error
	RecordContainerLog(containerID shared.ContainerID, log string, logType shared.LogType) error
}
