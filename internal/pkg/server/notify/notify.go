package notify

import "go-brunel/internal/pkg/shared"

// Notify will notify the status of a job, the receiver of the notification is determined by the implementation
type Notify interface {
	Notify(id shared.JobID) error
}
