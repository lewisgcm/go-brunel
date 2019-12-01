package store

import (
	"go-brunel/internal/pkg/shared"
	"time"
)

type Container struct {
	ID          string                `bson:"-"`
	JobID       shared.JobID          `bson:"-"`
	ContainerID shared.ContainerID    `bson:"container_id"`
	State       shared.ContainerState `bson:"state"`
	Spec        shared.Container      `bson:"spec"`
	Meta        shared.ContainerMeta  `bson:"meta"`
	CreatedAt   time.Time             `bson:"created_at"`
	StartedAt   *time.Time            `bson:"started_at"`
	StoppedAt   *time.Time            `bson:"stopped_at"`
}

type ContainerStore interface {
	Add(container Container) error

	UpdateStateByContainerID(id shared.ContainerID, state shared.ContainerState) error

	UpdateStoppedAtByContainerID(id shared.ContainerID, time time.Time) error

	UpdateStartedAtByContainerID(id shared.ContainerID, t time.Time) error

	FilterByJobID(jobID shared.JobID) ([]Container, error)
}
