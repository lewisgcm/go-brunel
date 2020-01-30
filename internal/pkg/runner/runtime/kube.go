/*
 * Author: Lewis Maitland
 *
 * Copyright (c) 2019 Lewis Maitland
 */

package runtime

import (
	"context"
	"encoding/json"
	"fmt"
	"go-brunel/internal/pkg/shared"
	"go-brunel/internal/pkg/shared/util"
	"io"
	"k8s.io/client-go/rest"
	"strings"
	"time"

	"github.com/docker/go-units"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
)

const (
	selector           = "subdomain"
	watchContainerName = "watcher"
	pollInterval       = time.Second
)

type KubeRuntime struct {
	Client          rest.Interface
	VolumeClaimName string
	Namespace       string
	WorkDir         string
}

// safeJobID will return a kubernetes safe namespace name, Kubernetes doesnt like it when they start with a number :'(
func safeJobID(id shared.JobID) string {
	return fmt.Sprintf("job-%s", id)
}

func (pipeline *KubeRuntime) Initialize(context context.Context, jobID shared.JobID, _ string) error {
	err := pipeline.
		Client.
		Post().
		Context(context).
		Namespace(pipeline.Namespace).
		Resource(string(corev1.ResourceServices)).
		Body(&corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name: safeJobID(jobID),
			},
			Spec: corev1.ServiceSpec{
				ClusterIP: "None",
				Selector: map[string]string{
					selector: safeJobID(jobID),
				},
				Ports: []corev1.ServicePort{
					{
						Port: 1234,
						Name: "port",
					},
				},
			},
		}).
		Do().
		Error()

	return errors.Wrap(err, "error initializing runner")
}

func (pipeline *KubeRuntime) Terminate(context context.Context, jobID shared.JobID) error {
	err := pipeline.
		Client.
		Delete().
		Context(context).
		Namespace(pipeline.Namespace).
		Resource(string(corev1.ResourceServices)).
		Name(safeJobID(jobID)).
		Do().
		Error()

	return errors.Wrap(err, "failed to remove kubernetes service")
}

func (pipeline *KubeRuntime) WaitForContainer(context context.Context, id shared.ContainerID, condition shared.ContainerWaitCondition) error {
	for {
		select {
		case <-context.Done():
			return errors.New("context cancelled waiting for container")
		default:
			var pod corev1.Pod
			err := pipeline.
				Client.
				Get().
				Resource(string(corev1.ResourcePods)).
				Namespace(pipeline.Namespace).
				Name(string(id)).
				Do().
				Into(&pod)

			if err != nil {
				return errors.Wrap(err, "failed to get pod")
			}

			if pod.Status.Phase == corev1.PodFailed || pod.Status.Phase == "CrashLoopBackOff" {
				return errors.New("pod could not be scheduled: " + string(pod.Status.Phase))
			}

			// This method is now pretty horrible, we basically need to check for the pod container being
			// ready that is not our watcher container.
			for _, status := range pod.Status.ContainerStatuses {
				if status.Name != watchContainerName {
					if status.State.Terminated != nil {
						if status.State.Terminated.ExitCode != 0 {
							return fmt.Errorf("container exited with non zero exit status: %d", status.State.Terminated.ExitCode)
						} else if condition.State == shared.ContainerWaitRunning {
							return errors.New("container completed whilst waiting for it to be ready")
						}
						return nil
					}

					if status.State.Running != nil && (condition.State&shared.ContainerWaitRunning) != 0 {
						return nil
					}
				}
			}

			// Check for any image waiting errors etc
			for _, status := range pod.Status.ContainerStatuses {
				if status.State.Waiting != nil && status.State.Waiting.Message != "" {
					return errors.New("failure waiting for pod: " + status.State.Waiting.Message)
				}
			}
			time.Sleep(pollInterval)
		}
	}
}

type dockerLogFormat struct {
	Log    string
	Stream string
	Time   time.Time
}

// writeContainerLogsSince will fetch the logs for our 'watcher' container, those logs are the raw docker log output. i.e json
// encoded logging information, we then decode the logs line by line and write them to the correct writer. We then return the number of logs
// written to the writers or an error
func (pipeline *KubeRuntime) writeContainerLogsSince(
	ctx context.Context,
	id shared.ContainerID,
	stdOut io.WriteCloser,
	stdErr io.WriteCloser,
	since time.Time,
) (int64, time.Time, error) {
	logReq := pipeline.
		Client.
		Get().
		Context(ctx).
		Namespace(pipeline.Namespace).
		Name(string(id)).
		Resource(string(corev1.ResourcePods)).
		SubResource("log").
		VersionedParams(
			&corev1.PodLogOptions{
				Container: watchContainerName,
				SinceTime: &metav1.Time{
					Time: since,
				},
			},
			scheme.ParameterCodec,
		)

	stream, err := logReq.Stream()
	if err != nil {
		return 0, time.Time{}, errors.Wrap(err, "error executing remote command")
	}
	var counter int64

	newSince := since
	logWriter := &util.LoggerWriter{
		Recorder: func(log string) error {

			// The log is raw JSON, we need to decode it
			var parsedLog dockerLogFormat
			if err := json.NewDecoder(strings.NewReader(log)).Decode(&parsedLog); err != nil {
				return errors.Wrap(err, "error decoding log")
			}

			// It seems as though the since kubernetes API variable isnt quite accurate so we do an
			// extra check here to make sure that the logs didnt occur before our last check
			if parsedLog.Time.UnixNano() <= since.UnixNano() {
				return nil
			}
			counter++

			// We set our new since to the last log, this way its easy for our caller to get a correct since check
			newSince = parsedLog.Time
			switch parsedLog.Stream {
			case "stderr":
				_, err = stdErr.Write([]byte(parsedLog.Log))
				if err != nil {
					return errors.Wrap(err, "error writing log to stderr")
				}
			case "stdout":
				_, err = stdOut.Write([]byte(parsedLog.Log))
				if err != nil {
					return errors.Wrap(err, "error writing log to stdout")
				}
			}
			return nil
		},
	}
	_, err = io.Copy(logWriter, stream)
	if err != nil {
		return counter, time.Time{}, errors.Wrap(err, "error reading log stream")
	}

	if err := stream.Close(); err != nil {
		return counter, time.Time{}, errors.Wrap(err, "error closing stream")
	}

	if err := logWriter.Close(); err != nil {
		return counter, time.Time{}, errors.Wrap(err, "error closing log writer")
	}

	return counter, newSince, nil
}

// TODO this is diabolical by all accords, and i would very much like to burn all of this with fire. However kubernetes
// does not have a way to get separate stderr/stdout logs from the API, so we need to do all the hooky stuff to make it happen.
// If you look at the DispatchContainer method, we create two containers the one the user wants to deploy and a 'watcher'. The watcher
// container essentially dumps the raw docker logs for the user container to stdout. This allows us to use the kubernetes logs endpoint to get the
// raw logging stream, decode it and write it to the proper channels.
func (pipeline *KubeRuntime) CopyLogsForContainer(ctx context.Context, id shared.ContainerID, stdOut io.WriteCloser, stdErr io.WriteCloser) error {
	var sinceCheck = time.Time{}
	var wroteLines int64

	for {
		select {
		case <-ctx.Done():
			return errors.New("context has been cancelled")
		default:
			var pod corev1.Pod
			err := pipeline.
				Client.
				Get().
				Context(ctx).
				Resource(string(corev1.ResourcePods)).
				Namespace(pipeline.Namespace).
				Name(string(id)).Do().Into(&pod)
			if err != nil {
				return errors.Wrap(err, "error getting pod for logs")
			}

			wroteLines, sinceCheck, err = pipeline.writeContainerLogsSince(ctx, id, stdOut, stdErr, sinceCheck)
			if err != nil {
				return errors.Wrap(err, "error getting logs")
			}

			// We want to exit as soon as our non 'watcher' container is done
			for _, status := range pod.Status.ContainerStatuses {
				if status.Name != watchContainerName {
					if status.State.Terminated != nil && wroteLines == 0 {
						return nil
					}
				}
			}
			time.Sleep(pollInterval)
		}
	}
}

func (pipeline *KubeRuntime) DispatchContainer(context context.Context, jobID shared.JobID, container shared.Container) (shared.ContainerID, error) {
	var zero int64
	var containerID = shared.ContainerID("brunel-container-" + strings.Replace(uuid.New().String(), "-", "", -1))

	// Check if the service exists, if not fail our config has not been initialized correctly
	if err := pipeline.
		Client.
		Get().
		Namespace(pipeline.Namespace).
		Resource(string(corev1.ResourceServices)).
		Name(safeJobID(jobID)).
		Do().
		Error(); err != nil {
		return shared.EmptyContainerID, errors.Wrap(err, "error getting service, may not exist")
	}

	// Create our command for the container
	var command []string
	var mounts []corev1.VolumeMount

	if container.WorkingDir != "" {
		mounts = append(
			mounts,
			corev1.VolumeMount{
				Name:      "workspace",
				MountPath: container.WorkingDir,
				SubPath:   string(jobID),
			},
		)
	}

	if container.EntryPoint != "" {
		command = []string{container.EntryPoint}
	}

	// Create our container config
	var env []corev1.EnvVar
	if container.Environment != nil && len(container.Environment) > 0 {
		for key, value := range container.Environment {
			env = append(env, corev1.EnvVar{Name: key, Value: value})
		}
	}

	// Used for handling resource requests/limits
	resources := corev1.ResourceRequirements{}
	if container.Resources != nil {
		// Set the kube resource limits
		if container.Resources.Limits != nil {
			// NOTE we use the docker function for converting the string to bytes, the reason for this is
			// it gives us consistency
			memory, _ := units.RAMInBytes(container.Resources.Limits.Memory)

			resources.Limits = corev1.ResourceList{
				corev1.ResourceCPU: *resource.NewScaledQuantity(
					int64(container.Resources.Limits.CPU*1e9),
					resource.Nano,
				),
				corev1.ResourceMemory: *resource.NewQuantity(memory, resource.DecimalSI),
			}
		}

		// Set the kube resource request
		if container.Resources.Requests != nil {
			// NOTE we use the docker function for converting the string to bytes, the reason for this is
			// it gives us consistency
			memory, _ := units.RAMInBytes(container.Resources.Requests.Memory)

			resources.Requests = corev1.ResourceList{
				corev1.ResourceCPU: *resource.NewScaledQuantity(
					int64(container.Resources.Requests.CPU*1e9),
					resource.Nano,
				),
				corev1.ResourceMemory: *resource.NewQuantity(memory, resource.DecimalSI),
			}
		}
	}

	// Create our pod based off of our spec
	err := pipeline.Client.
		Post().
		Context(context).
		Namespace(pipeline.Namespace).
		Resource(string(corev1.ResourcePods)).
		Body(&corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name: string(containerID),
				Labels: map[string]string{
					selector: safeJobID(jobID),
				},
			},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{
						Name:         string(containerID),
						Env:          env,
						Image:        container.Image,
						Command:      command,
						Args:         container.Args,
						WorkingDir:   container.WorkingDir,
						Stdin:        true,
						VolumeMounts: mounts,
						SecurityContext: &corev1.SecurityContext{
							Privileged: &container.Privileged,
						},
						Resources: resources,
					},
					// This is more of our stderr/stdout ugliness, here we execute a tail on the POD_ID container logs
					// We can then get the raw logs and parse into stderr/stdout as required
					{
						Name:    watchContainerName,
						Image:   "busybox",
						Command: []string{"sh", "-c", "--"},
						// Tail the entire file, retry and be quiet about it :)
						Args: []string{"tail -F -q -n +1 /var/log/pods/*$POD_ID/*" + string(containerID) + "/0.log"},
						Env: []corev1.EnvVar{
							{
								Name: "POD_ID",
								ValueFrom: &corev1.EnvVarSource{
									FieldRef: &corev1.ObjectFieldSelector{
										FieldPath: "metadata.uid",
									},
								},
							},
						},
						VolumeMounts: []corev1.VolumeMount{
							{
								Name:      "pod-logs",
								MountPath: "/var/log/pods",
								ReadOnly:  true,
							},
							{
								Name:      "container-logs",
								MountPath: "/var/lib/docker/containers",
								ReadOnly:  true,
							},
						},
					},
				},
				Hostname:  container.Hostname,
				Subdomain: safeJobID(jobID),
				DNSConfig: &corev1.PodDNSConfig{
					Searches: []string{
						fmt.Sprintf("%s.%s.svc.cluster.local", safeJobID(jobID), pipeline.Namespace),
					},
				},
				RestartPolicy: corev1.RestartPolicyNever,
				Volumes: []corev1.Volume{
					{
						Name: "workspace",
						VolumeSource: corev1.VolumeSource{
							PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
								ClaimName: pipeline.VolumeClaimName,
								ReadOnly:  false,
							},
						},
					},

					// Mount the logging directory for k8s so that we can access the raw logs
					// This way we can decode them and actually get proper stdout/stderr separation
					// Its pretty ugly but until its supported in k8s, we do it this way
					{
						Name: "pod-logs",
						VolumeSource: corev1.VolumeSource{
							HostPath: &corev1.HostPathVolumeSource{
								Path: "/var/log/pods",
							},
						},
					},
					{
						Name: "container-logs",
						VolumeSource: corev1.VolumeSource{
							HostPath: &corev1.HostPathVolumeSource{
								Path: "/var/lib/docker/containers",
							},
						},
					},
				},
				TerminationGracePeriodSeconds: &zero,
			},
		}).
		Do().
		Error()

	return containerID, errors.Wrap(err, "error creating pod")
}

func (pipeline *KubeRuntime) TerminateContainer(context context.Context, containerID shared.ContainerID) error {
	err := pipeline.
		Client.
		Delete().
		Context(context).
		Namespace(pipeline.Namespace).
		Resource(string(corev1.ResourcePods)).
		Name(string(containerID)).
		Do().
		Error()

	return errors.Wrap(err, "error terminating container")
}
