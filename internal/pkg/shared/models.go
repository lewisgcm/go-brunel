/*
 * Author: Lewis Maitland
 *
 * Copyright (c) 2019 Lewis Maitland
 */

package shared

// ContainerID is a type for a containers id, for example a docker container id (string)
type ContainerID string

// ContainerState denotes the state of a container
type ContainerState uint8

// ContainerWaitState are the states that we can wait for a container to be in
type ContainerWaitState uint8

// JobID is the ID of a job, it should be a string an be friendly for both kube and docker
type JobID string

// JobState denotes the state of a job for example whether it is running, failed, success
type JobState uint8

// LogType is used to denote the type of log message, this is used for recording job and container logs
type LogType uint8

// Exported constants here are for various purposes. The LogTypeStd* are for use in the recorders to specify the type of log line being recorded.
// JobState is used to denote what state a job is in.
const (
	LogTypeStdOut LogType = 1
	LogTypeStdErr LogType = 2

	JobStateWaiting    JobState = 0
	JobStateProcessing JobState = 1
	JobStateFailed     JobState = 2
	JobStateSuccess    JobState = 3
	JobStateCancelled  JobState = 4

	ContainerStateStarting ContainerState = 0
	ContainerStateRunning  ContainerState = 1
	ContainerStateStopped  ContainerState = 2

	ContainerWaitRunning ContainerWaitState = 1 << 0
	ContainerWaitStopped ContainerWaitState = 1 << 1

	// EmptyContainerID denotes an empty container ID, used in error returns
	EmptyContainerID ContainerID = ""
)

// ContainerWaitCondition are used as conditions when waiting for a container
type ContainerWaitCondition struct {
	// Wait for the container state before we are done waiting
	State ContainerWaitState
}

type ContainerResources struct {
	Limits   *ContainerResourcesUnits
	Requests *ContainerResourcesUnits
}

type ContainerResourcesUnits struct {
	CPU    float32 `yaml:"CPU"` // This is the same as dockers --cpus
	Memory string
}

// Container is used for defining a container for dispatch as part of the pipeline.
// Add commands to be run in the container shell.
type Container struct {
	Image       string
	Environment map[string]string
	Hostname    string
	EntryPoint  string
	Args        []string
	WorkingDir  string `yaml:"working_dir"`
	Privileged  bool
	Resources   *ContainerResources
}

// ContainerMeta is used for handling additional container meta data such as the containers stage or if it is a service.
type ContainerMeta struct {
	Stage   string
	Service bool
}

// Stage defines a runnable stage that can be restricted to specific environments with
// sidecar services and steps
type Stage struct {
	Environments []string
	Services     []Container
	Steps        []Container
}

// Spec is used for defining the pipeline
type Spec struct {
	Version     string
	Description string
	Maintainers []string
	Stages      map[string]Stage
}

// Commit is used to denote a commit for a job
type Commit struct {
	Branch   string
	Revision string
}

// Job is used to denote a job that should be processed
type Job struct {
	ID         JobID
	Repository Repository
	Commit     Commit
	State      JobState
}

// Repository is used to denote a single VCS repository known to the system
type Repository struct {
	Project string
	Name    string
	URI     string
}
