package notify

import (
	"go-brunel/internal/pkg/shared"

	log "github.com/Sirupsen/logrus"
)

// TextNotify
type TextNotify struct {
}

func (notify *TextNotify) Notify(id shared.JobID) error {
	log.Info("job with id ", id, " has changed state ....")
	return nil
}
