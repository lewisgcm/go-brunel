package remote

import "go-brunel/internal/pkg/shared"

type SetJobStateRequest struct {
	Id    shared.JobID
	State shared.JobState
}

type LogRequest struct {
	Id      shared.JobID
	Message string
	StageID shared.StageID
	LogType shared.LogType
}

type SetStageStateRequest struct {
	Id    shared.StageID
	JobID shared.JobID
	State shared.StageState
}

type AddContainerRequest struct {
	Id          shared.JobID
	ContainerID shared.ContainerID
	Meta        shared.ContainerMeta
	Container   shared.Container
	State       shared.ContainerState
}

type SetContainerStateRequest struct {
	Id    shared.ContainerID
	State shared.ContainerState
}

type ContainerLogRequest struct {
	Id      shared.ContainerID
	Message string
	LogType shared.LogType
}

type SearchForXRequest struct {
	SearchPath []string
	Name       string
}

type GetNextAvailableJobResponse struct {
	Job *shared.Job
}

type Empty struct {
}
