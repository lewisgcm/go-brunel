// +build integration !unit

// runtime Integration tests check some 'expected' behavior of our interfaces against each runtime config.
package runtime

import (
	"context"
	"fmt"
	"go-brunel/internal/pkg/shared"
	"io/ioutil"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

// TestInitTerminate tests a simple initialization and termination of a runtime config
func TestInitTerminate(t *testing.T) {
	suite := setup(t)

	for _, r := range suite.runtimes {
		id := shared.JobID(uuid.New().String())
		if err := r.runtime().Initialize(context.TODO(), id, "./"); err != nil {
			t.Error(err)
		}

		if !r.jobRuntimeExists(id, t) {
			t.Error("expecting runtime config to have been created")
		}

		if err := r.runtime().Terminate(context.TODO(), id); err != nil {
			t.Error(err)
		}

		if r.jobRuntimeExists(id, t) {
			t.Error("expecting runtime config to have been destroyed")
		}
	}
}

func TestFaultyDispatch(t *testing.T) {
	suite := setup(t)

	for _, r := range suite.runtimes {
		id := shared.JobID(uuid.New().String())
		containerID, err := r.runtime().DispatchContainer(
			context.TODO(),
			id,
			shared.Container{
				Image:      "ubuntu:18.04",
				EntryPoint: "sh",
				Args:       []string{"-c", "--", "echo 'stdout'; echo 'stderr' > /dev/stderr"},
			},
		)
		if err == nil {
			t.Error("expecting error when dispatching container without initialized config")
		}

		if r.containerExists(containerID, t) {
			t.Error("expecting there to be no container created")
		}
	}
}

// TestSimpleContainerFailure dispatches a simple container with a faulty command that will fail
// The test follows these steps:
// 1. runtime.Initialize and check docker network exists
// 2. runtime.DispatchContainer a container, check it exists
// 3. runtime.WaitForContainer wait for the container to stop, this should return an error
// 4. runtime.TerminateContainer the container, check it has been removed
// 5. runtime.Terminate the runtime we created, delete the service
func TestSimpleContainerFailure(t *testing.T) {
	suite := setup(t)

	for _, r := range suite.runtimes {
		id := shared.JobID(uuid.New().String())

		if err := r.runtime().Initialize(context.TODO(), id, "./"); err != nil {
			t.Error(err)
		}

		if !r.jobRuntimeExists(id, t) {
			t.Error("expecting runtime config to have been created")
		}

		// This container SHOULD fail, the cat is trying to open a non existent file
		// Although this command should work, its the wait that will return an error
		containerID, err := r.runtime().DispatchContainer(
			context.TODO(),
			id,
			shared.Container{
				Image:      "ubuntu:18.04",
				EntryPoint: "sh",
				Args:       []string{"-c", "--", "cat asdasdasdsd"},
			},
		)
		if err != nil {
			t.Error(err)
		}

		if !r.containerExists(containerID, t) {
			t.Error("expecting container to have been created")
		}

		// Wait for the container to be stopped
		if err := r.runtime().WaitForContainer(context.TODO(), containerID, shared.ContainerWaitCondition{State: shared.ContainerWaitStopped}); err == nil {
			t.Error("expecting error when waiting for faulty container")
		}

		// Wait for it to be running, its already stopped we expect an error here
		if err := r.runtime().WaitForContainer(context.TODO(), containerID, shared.ContainerWaitCondition{State: shared.ContainerWaitRunning}); err == nil {
			t.Error("expecting error when waiting for faulty container")
		}

		// Wait for either stopped or running
		if err := r.runtime().WaitForContainer(context.TODO(), containerID, shared.ContainerWaitCondition{State: shared.ContainerWaitStopped | shared.ContainerWaitRunning}); err == nil {
			t.Error("expecting error when waiting for faulty container")
		}

		// Get our stdout and stderr and make sure they match the values we echoed in in the container
		// Args section. We use our mock writer, that is basically just a string writer with a noop close
		var stdOut strings.Builder
		var stdErr strings.Builder
		if err := r.runtime().CopyLogsForContainer(
			context.TODO(),
			containerID,
			&mockBufferWriteCloser{Writer: &stdOut},
			&mockBufferWriteCloser{Writer: &stdErr},
		); err != nil {
			t.Error(err)
		}

		// We only need to check our stderr in this case, this could break in the future if cats message ever changes
		// But for now it lets us know we do have logs for a borked container
		if stdErr.String() != "cat: asdasdasdsd: No such file or directory\n" {
			t.Error(fmt.Sprintf("'%s' != 'cat: asdasdasdsd: No such file or directory\n'", stdErr.String()))
		}

		if err := r.runtime().TerminateContainer(context.TODO(), containerID); err != nil {
			t.Error(err)
		}

		if r.containerExists(containerID, t) {
			t.Error("expecting container to have been destroyed")
		}

		if err := r.runtime().Terminate(context.TODO(), id); err != nil {
			t.Error(err)
		}

		if r.jobRuntimeExists(id, t) {
			t.Error("expecting runtime config to have been destroyed")
		}
	}
}

// TestSimpleContainerWaitMultiple tests a simple container dispatch and wait, but the wait could be either stopped OR running
// The test follows these steps:
// 1. runtime.Initialize
// 2. runtime.DispatchContainer a container
// 3. runtime.WaitForContainer wait for the container to stop
// 4. runtime.TerminateContainer the container
// 5. runtime.Terminate the runtime we created
func TestContainerWait(t *testing.T) {
	suite := setup(t)

	for _, r := range suite.runtimes {
		id := shared.JobID(uuid.New().String())

		if err := r.runtime().Initialize(context.TODO(), id, "./"); err != nil {
			t.Error(err)
		}

		containerID, err := r.runtime().DispatchContainer(
			context.TODO(),
			id,
			shared.Container{
				Image:      "ubuntu:18.04",
				EntryPoint: "sh",
				Args:       []string{"-c", "--", "echo 'stdout'; echo 'stderr' > /dev/stderr"},
			},
		)
		if err != nil {
			t.Error(err)
		}

		// Wait for the container to be stopped
		if err := r.runtime().WaitForContainer(context.TODO(), containerID, shared.ContainerWaitCondition{State: shared.ContainerWaitStopped}); err != nil {
			t.Error("error waiting for container", err)
		}

		// Wait for it to be running, its already stopped we expect an error here
		if err := r.runtime().WaitForContainer(context.TODO(), containerID, shared.ContainerWaitCondition{State: shared.ContainerWaitRunning}); err == nil {
			t.Error("error should be returned when attempting to wait for 'running' state on a stopped container")
		}

		// Wait for either stopped or running
		if err := r.runtime().WaitForContainer(context.TODO(), containerID, shared.ContainerWaitCondition{State: shared.ContainerWaitStopped | shared.ContainerWaitRunning}); err != nil {
			t.Error("error returned when waiting for either running or stopped", err)
		}

		// Get our stdout and stderr and make sure they match the values we echoed in in the container
		// Args section. We use our mock writer, that is basically just a string writer with a noop close
		var stdOut strings.Builder
		var stdErr strings.Builder
		if err := r.runtime().CopyLogsForContainer(
			context.TODO(),
			containerID,
			&mockBufferWriteCloser{Writer: &stdOut},
			&mockBufferWriteCloser{Writer: &stdErr},
		); err != nil {
			t.Error(err)
		}

		// Check our standard output and error, is 'stdout' + newline and 'stderr' + newline respectively
		if stdOut.String() != "stdout\n" {
			t.Error(stdOut.String(), "!=", "stdout")
		}
		if stdErr.String() != "stderr\n" {
			t.Error(stdErr.String(), "!=", "stderr")
		}

		if err := r.runtime().TerminateContainer(context.TODO(), containerID); err != nil {
			t.Error(err)
		}

		if err := r.runtime().Terminate(context.TODO(), id); err != nil {
			t.Error(err)
		}
	}
}

func TestTimeout(t *testing.T) {
	suite := setup(t)

	for _, r := range suite.runtimes {
		id := shared.JobID(uuid.New().String())

		if err := r.runtime().Initialize(context.TODO(), id, "./"); err != nil {
			t.Error(err)
		}

		containerID, err := r.runtime().DispatchContainer(
			context.TODO(),
			id,
			shared.Container{
				Image:      "ubuntu:18.04",
				EntryPoint: "sh",
				Args:       []string{"-c", "--", "while true; do echo 'hello'; sleep 1; done"},
			},
		)
		if err != nil {
			t.Error(err)
		}

		// Wait for the container to be running, so our tests later are valid
		if err := r.runtime().WaitForContainer(context.TODO(), containerID, shared.ContainerWaitCondition{State: shared.ContainerWaitRunning}); err != nil {
			t.Error("error waiting for container to be running", err)
		}

		// Now wait for the impossible, i.e our container to be stopped
		ctx, _ := context.WithTimeout(context.Background(), time.Second)
		_ = r.runtime().WaitForContainer(ctx, containerID, shared.ContainerWaitCondition{State: shared.ContainerWaitStopped})
		if ctx.Err() != context.DeadlineExceeded {
			t.Error("expecting timeout")
		}

		ctx, _ = context.WithTimeout(context.Background(), time.Second)
		_ = r.runtime().CopyLogsForContainer(ctx, containerID, &mockBufferWriteCloser{Writer: ioutil.Discard}, &mockBufferWriteCloser{Writer: ioutil.Discard})
		if ctx.Err() != context.DeadlineExceeded {
			t.Error("expecting timeout")
		}

		if err := r.runtime().TerminateContainer(context.TODO(), containerID); err != nil {
			t.Error(err)
		}

		if err := r.runtime().Terminate(context.TODO(), id); err != nil {
			t.Error(err)
		}
	}
}
