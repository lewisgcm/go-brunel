package trigger

import (
	"context"
	"go-brunel/internal/pkg/shared"
)

// LocalRepository is used to denote a job that is running locally
var LocalRepository = shared.Repository{}

// Trigger provides a channel for job events, i.e anytime a job needs processing
// it will become available on the channel.
type Trigger interface {
	Await(ctx context.Context) <-chan Event
}

// Event holds the job for processing, the working directory and a channel to send back the state of the job when done.
// The state channel is used for recording the state of the job.
type Event struct {
	Job      shared.Job
	JobState chan shared.JobState
	WorkDir  string
	Context  context.Context
}
