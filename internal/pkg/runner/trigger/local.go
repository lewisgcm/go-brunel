/*
 * Author: Lewis Maitland
 *
 * Copyright (c) 2019 Lewis Maitland
 */

package trigger

import (
	"context"
	"fmt"
	"go-brunel/internal/pkg/shared"
	"time"
)

type LocalTrigger struct {
	WorkDir string
}

func (trigger *LocalTrigger) Await(ctx context.Context) <-chan Event {
	channel := make(chan Event)

	go func() {
		stateChan := make(chan shared.JobState)
		channel <- Event{
			Job: shared.Job{
				Repository: LocalRepository,
				ID:         shared.JobID(fmt.Sprintf("local-build-%s", time.Now().Format("2006-01-02-15-04-05"))),
			},
			WorkDir:  trigger.WorkDir,
			JobState: stateChan,
			Context:  ctx,
		}
		<-stateChan
		close(channel)
	}()

	return channel
}
