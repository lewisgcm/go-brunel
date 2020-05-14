package store

import (
	"errors"
	"go-brunel/internal/pkg/shared"
	"strings"
	"time"
)

type JobListPage struct {
	Count int64 `bson:"job_count"`
	Jobs  []Job `bson:"jobs"`
}

type Job struct {
	ID            shared.JobID          `bson:"-"`
	RepositoryID  RepositoryID          `bson:"-"`
	EnvironmentID *shared.EnvironmentID `bson:"-"`
	Commit        shared.Commit
	State         shared.JobState
	StartedBy     string     `bson:"started_by"`
	StoppedBy     *string    `bson:"stopped_by"`
	CreatedAt     time.Time  `bson:"created_at"`
	StartedAt     *time.Time `bson:"started_at"`
	StoppedAt     *time.Time `bson:"stopped_at"`
}

func (job *Job) Clean() {
	job.Commit.Branch = strings.TrimSpace(job.Commit.Branch)
	job.Commit.Revision = strings.TrimSpace(job.Commit.Revision)
}

func (job *Job) IsValid() error {
	if len(job.Commit.Branch) == 0 || len(job.Commit.Revision) == 0 {
		return errors.New("branch and revision are required")
	}

	return nil
}

type JobDetail struct {
	Job    `bson:",inline"`
	Stages []struct {
		Name      string
		Container []Container
	}
}

type JobStore interface {
	Next() (*Job, error)

	Get(id shared.JobID) (*Job, error)

	Add(job Job) (*Job, error)

	UpdateStoppedAtByID(id shared.JobID, t time.Time) error

	UpdateStateByID(id shared.JobID, s shared.JobState) error

	CancelByID(id shared.JobID, userID string) error

	FilterByRepositoryID(
		repositoryID RepositoryID,
		filter string,
		pageIndex int64,
		pageSize int64,
		sortColumn string,
		sortOrder int,
	) (JobListPage, error)

	Delete(id shared.JobID) error
}
