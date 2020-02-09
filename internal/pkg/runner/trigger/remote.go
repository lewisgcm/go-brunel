/*
 * Author: Lewis Maitland
 *
 * Copyright (c) 2019 Lewis Maitland
 */

package trigger

import (
	"context"
	"github.com/pkg/errors"
	"go-brunel/internal/pkg/runner/remote"
	"go-brunel/internal/pkg/shared"
	"log"
	"time"
)

type RemoteTrigger struct {
	Remote      remote.Remote
	BaseWorkDir string
}

func (trigger *RemoteTrigger) getCancellationChannel(ctx context.Context, id shared.JobID) chan bool {
	cancelled := make(chan bool)
	go func() {
	loop:
		for {
			select {
			case <-ctx.Done():
				close(cancelled)
				break loop
			default:
				isCancelled, err := trigger.
					Remote.
					HasBeenCancelled(id)

				if err != nil {
					close(cancelled)
					break loop
				}

				if isCancelled {
					cancelled <- isCancelled
					break loop
				}
				time.Sleep(time.Second)
			}
		}
	}()
	return cancelled
}

func (trigger *RemoteTrigger) Await(ctx context.Context) <-chan Event {
	channel := make(chan Event, 1)

	go func() {
	loop:
		for {
			select {
			case <-ctx.Done():
				break loop

			default:
				job, err := trigger.
					Remote.
					GetNextAvailableJob()

				if err != nil {
					log.Println(errors.Wrap(err, "error waiting for next job"))
					break loop
				}

				if job == nil {
					time.Sleep(time.Second)
					continue loop
				}

				jobCtx, jobCancel := context.WithCancel(ctx)
				stateChannel := make(chan shared.JobState)
				channel <- Event{
					Job:      *job,
					JobState: stateChannel,
					WorkDir:  trigger.BaseWorkDir + "/" + string(job.ID) + "/",
					Context:  jobCtx,
				}

				cancelCtx, cancelCancel := context.WithCancel(ctx)
				cancelledChan := trigger.getCancellationChannel(cancelCtx, job.ID)

			jobLoop:
				for {
					select {
					case result := <-stateChannel:
						if err = trigger.
							Remote.
							SetJobState(job.ID, result); err != nil {
							log.Println(errors.Wrap(err, "error setting job state"))
						}
						cancelCancel()
						break jobLoop
					case isCancelled, ok := <-cancelledChan:
						if (ok && isCancelled) || !ok {
							<-stateChannel
							jobCancel()
							cancelCancel()
							break jobLoop
						}
					default:
						time.Sleep(time.Millisecond * 200)
					}
				}
			}
		}
		close(channel)
	}()

	return channel
}
