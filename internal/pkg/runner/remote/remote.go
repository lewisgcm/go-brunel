package remote

import "go-brunel/internal/pkg/shared"

// Remote is an interface that defines all expected communication between a runner and server
type Remote interface {
	// GetNextAvailableJob should be atomic and only ever return a single job to a unique runner.
	// Returns nil if no job is found, client is expected to retry afterwards
	GetNextAvailableJob() (*shared.Job, error)

	// SetJobState set the state of the job with given id
	SetJobState(id shared.JobID, state shared.JobState) error

	// Check if the job is cancelled
	HasBeenCancelled(id shared.JobID) (bool, error)

	// Log should store messages of a given type for a job
	Log(id shared.JobID, message string, logType shared.LogType, stage string) error

	// Add stage will add a new stage to the pipeline output
	AddStage(name string) (shared.StageID, error)

	// AddContainer should log a container against a given JobID. The containerID is the ID returned by the
	// runtime config the container is running in (e.g docker, kube etc).
	AddContainer(id shared.JobID, containerID shared.ContainerID, meta shared.ContainerMeta, container shared.Container, state shared.ContainerState) error

	// SetContainerState should set the state of a container
	SetContainerState(id shared.ContainerID, state shared.ContainerState) error

	// ContainerLog should log a message for the given containerID
	ContainerLog(id shared.ContainerID, message string, logType shared.LogType) error

	SearchForValue(searchPath []string, name string) (string, error)

	SearchForSecret(searchPath []string, name string) (string, error)
}
