package store

import (
	"go-brunel/internal/pkg/shared"
	"time"
)

type Stage struct {
	ID        shared.StageID    `bson:"id"`
	JobID     shared.JobID      `bson:"job_id"`
	State     shared.StageState `bson:"state"`
	StartedAt *time.Time        `bson:"started_at,omitempty"`
	StoppedAt *time.Time        `bson:"stopped_at,omitempty"`
}

type StageStore interface {
	AddOrUpdate(stage Stage) error

	Get(jobID shared.JobID) ([]Stage, error)
}
