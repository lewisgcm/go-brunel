package store

import (
	"go-brunel/internal/pkg/shared"
	"time"
)

type Log struct {
	JobID   shared.JobID   `bson:"-"`
	Message string         `bson:"message"`
	LogType shared.LogType `bson:"type"`
	Stage 	string		   `bson:"stage"`
	Time    time.Time
}

type ContainerLog struct {
	ContainerID shared.ContainerID `bson:"container_id"`
	Message     string             `bson:"message"`
	LogType     shared.LogType     `bson:"type"`
	Time        time.Time
}

type LogStore interface {
	// Log should messages of a given type for a job
	Log(l Log) error

	// ContainerLog should log a message for the given containerID
	ContainerLog(l ContainerLog) error

	FilterLogByJobIDFromTime(id shared.JobID, t time.Time) ([]Log, error)

	FilterContainerLogByContainerIDFromTime(id shared.ContainerID, t time.Time) ([]ContainerLog, error)
}
