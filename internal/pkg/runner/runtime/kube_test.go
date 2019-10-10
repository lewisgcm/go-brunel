package runtime_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"go-brunel/internal/pkg/runner/runtime"
	"go-brunel/internal/pkg/shared"
	"go-brunel/test"
	"io"
	"io/ioutil"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest/fake"
	"net/http"
	"strings"
	"testing"
	"time"
)

func fakeRESTClient() fake.RESTClient {
	return fake.RESTClient{
		GroupVersion:         schema.GroupVersion{Group: "", Version: "v1"},
		NegotiatedSerializer: serializer.DirectCodecFactory{CodecFactory: scheme.Codecs},
	}
}

func TestKubeRuntime_Initialize(t *testing.T) {
	fakeRest := fakeRESTClient()
	kubeRuntime := runtime.KubeRuntime{
		Client:    &fakeRest,
		Namespace: "test",
	}

	// Test happy path
	fakeRest.Resp = &http.Response{
		StatusCode: 200,
		Body: &test.NoOpReadCloser{
			Reader: bytes.NewReader([]byte("")),
		},
	}
	err := kubeRuntime.Initialize(context.TODO(), shared.JobID("id"), "")
	test.ExpectError(t, nil, err)

	// Test expected error on kube client error
	expectedError := errors.New("kube_init_error")
	fakeRest.Err = expectedError
	err = kubeRuntime.TerminateContainer(context.TODO(), shared.ContainerID("id"))
	test.ExpectErrorLike(t, expectedError, err)
}

func TestKubeRuntime_Terminate(t *testing.T) {
	fakeRest := fakeRESTClient()
	kubeRuntime := runtime.KubeRuntime{
		Client:    &fakeRest,
		Namespace: "test",
	}

	// Test happy path
	fakeRest.Resp = &http.Response{
		StatusCode: 200,
		Body: &test.NoOpReadCloser{
			Reader: bytes.NewReader([]byte("")),
		},
	}
	err := kubeRuntime.Initialize(context.TODO(), shared.JobID("id"), "")
	test.ExpectError(t, nil, err)

	// Test expected error on kube client error
	expectedError := errors.New("kube_terminate_error")
	fakeRest.Err = expectedError
	err = kubeRuntime.Terminate(context.TODO(), shared.JobID("id"))
	test.ExpectErrorLike(t, expectedError, err)
}

func TestKubeRuntime_TerminateContainer(t *testing.T) {
	fakeRest := fakeRESTClient()
	kubeRuntime := runtime.KubeRuntime{
		Client:    &fakeRest,
		Namespace: "test",
	}

	// Test happy path
	fakeRest.Resp = &http.Response{
		StatusCode: 200,
		Body: &test.NoOpReadCloser{
			Reader: bytes.NewReader([]byte("")),
		},
	}
	err := kubeRuntime.TerminateContainer(context.TODO(), shared.ContainerID("id"))
	test.ExpectError(t, nil, err)

	// Test expected error on kube client error
	expectedError := errors.New("kube_terminate_error")
	fakeRest.Err = expectedError
	err = kubeRuntime.TerminateContainer(context.TODO(), shared.ContainerID("id"))
	test.ExpectErrorLike(t, expectedError, err)
}

func TestKubeRuntime_WaitForContainer(t *testing.T) {
	suites := []struct {
		waitCondition shared.ContainerWaitCondition
		expectedError error
		respBody      v1.Pod
		respError     error
	}{
		// Test error getting pod
		{
			waitCondition: shared.ContainerWaitCondition{
				State: shared.ContainerWaitRunning | shared.ContainerWaitStopped,
			},
			expectedError: errors.New("pod_get_errors"),
			respBody:      v1.Pod{},
			respError:     errors.New("pod_get_errors"),
		},
		// Test schedule error, failed
		{
			waitCondition: shared.ContainerWaitCondition{
				State: shared.ContainerWaitRunning | shared.ContainerWaitStopped,
			},
			expectedError: errors.New("pod could not be scheduled: Failed"),
			respBody: v1.Pod{
				Status: v1.PodStatus{
					Phase:             v1.PodFailed,
					ContainerStatuses: []v1.ContainerStatus{},
				},
			},
		},
		// Test wait for running, container stopped
		{
			waitCondition: shared.ContainerWaitCondition{
				State: shared.ContainerWaitRunning,
			},
			expectedError: errors.New("container completed whilst waiting for it to be ready"),
			respBody: v1.Pod{
				Status: v1.PodStatus{
					ContainerStatuses: []v1.ContainerStatus{
						{
							State: v1.ContainerState{
								Terminated: &v1.ContainerStateTerminated{
									ExitCode: 0,
								},
							},
						},
					},
				},
			},
		},
		// Test wait for stopped, container stopped non-zero exit
		{
			waitCondition: shared.ContainerWaitCondition{
				State: shared.ContainerWaitStopped,
			},
			expectedError: errors.New("container exited with non zero exit status: -1"),
			respBody: v1.Pod{
				Status: v1.PodStatus{
					ContainerStatuses: []v1.ContainerStatus{
						{
							State: v1.ContainerState{
								Terminated: &v1.ContainerStateTerminated{
									ExitCode: -1,
								},
							},
						},
					},
				},
			},
		},
		// Test wait for running, container running
		{
			waitCondition: shared.ContainerWaitCondition{
				State: shared.ContainerWaitRunning,
			},
			respBody: v1.Pod{
				Status: v1.PodStatus{
					ContainerStatuses: []v1.ContainerStatus{
						{
							State: v1.ContainerState{
								Running: &v1.ContainerStateRunning{},
							},
						},
					},
				},
			},
		},
		// Test wait for running or stopped, container stopped
		{
			waitCondition: shared.ContainerWaitCondition{
				State: shared.ContainerWaitRunning | shared.ContainerWaitStopped,
			},
			respBody: v1.Pod{
				Status: v1.PodStatus{
					ContainerStatuses: []v1.ContainerStatus{
						{
							State: v1.ContainerState{
								Terminated: &v1.ContainerStateTerminated{
									ExitCode: 0,
								},
							},
						},
					},
				},
			},
		},
		// Test wait for running or stopped, container running
		{
			waitCondition: shared.ContainerWaitCondition{
				State: shared.ContainerWaitRunning | shared.ContainerWaitStopped,
			},
			respBody: v1.Pod{
				Status: v1.PodStatus{
					ContainerStatuses: []v1.ContainerStatus{
						{
							State: v1.ContainerState{
								Running: &v1.ContainerStateRunning{},
							},
						},
					},
				},
			},
		},
		// Test wait for stopped, container stopped non-zero exit
		{
			waitCondition: shared.ContainerWaitCondition{
				State: shared.ContainerWaitStopped | shared.ContainerWaitRunning,
			},
			expectedError: errors.New("container exited with non zero exit status: -1"),
			respBody: v1.Pod{
				Status: v1.PodStatus{
					ContainerStatuses: []v1.ContainerStatus{
						{
							State: v1.ContainerState{
								Terminated: &v1.ContainerStateTerminated{
									ExitCode: -1,
								},
							},
						},
					},
				},
			},
		},
		// Test wait for stopped, container waiting error
		{
			waitCondition: shared.ContainerWaitCondition{
				State: shared.ContainerWaitStopped | shared.ContainerWaitRunning,
			},
			expectedError: errors.New("failure waiting for pod: error_waiting"),
			respBody: v1.Pod{
				Status: v1.PodStatus{
					ContainerStatuses: []v1.ContainerStatus{
						{
							State: v1.ContainerState{
								Waiting: &v1.ContainerStateWaiting{
									Message: "error_waiting",
								},
							},
						},
					},
				},
			},
		},
	}

	for i, suite := range suites {
		t.Run(
			fmt.Sprintf("suites[%d]", i),
			func(t *testing.T) {
				fakeRest := fakeRESTClient()
				kubeRuntime := runtime.KubeRuntime{
					Client:    &fakeRest,
					Namespace: "test",
				}

				podBytes, err := json.Marshal(suite.respBody)
				if err != nil {
					t.Fatal(err)
				}

				fakeRest.Err = suite.respError
				fakeRest.Resp = &http.Response{
					StatusCode: 200,
					Body: &test.NoOpReadCloser{
						Reader: bytes.NewReader(podBytes),
					},
				}

				err = kubeRuntime.WaitForContainer(context.TODO(), shared.ContainerID("id"), suite.waitCondition)
				test.ExpectErrorLike(t, suite.expectedError, err)
			},
		)
	}
}

func TestKubeRuntime_WaitForContainer_ContextTimeout(t *testing.T) {
	fakeRest := fakeRESTClient()
	kubeRuntime := runtime.KubeRuntime{
		Client:    &fakeRest,
		Namespace: "test",
	}

	podBytes, err := json.Marshal(v1.Pod{
		Status: v1.PodStatus{
			Phase: "Running",
			ContainerStatuses: []v1.ContainerStatus{
				{
					State: v1.ContainerState{
						Running: &v1.ContainerStateRunning{},
					},
				},
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	fakeRest.Client = fake.CreateHTTPClient(
		func(req *http.Request) (response *http.Response, e error) {
			return &http.Response{
				StatusCode: 200,
				Body: &test.NoOpReadCloser{
					Reader: bytes.NewReader(podBytes),
				},
			}, nil
		},
	)

	timeoutCtx, cancel := context.WithTimeout(context.Background(), time.Second)
	err = kubeRuntime.WaitForContainer(
		timeoutCtx,
		shared.ContainerID("id"),
		shared.ContainerWaitCondition{
			State: shared.ContainerWaitStopped,
		},
	)
	test.ExpectErrorLike(t, errors.New("context cancelled waiting for container"), err)
	cancel()
}

func TestKubeRuntime_DispatchContainer(t *testing.T) {
	suites := []struct {
		container       shared.Container
		servicesRespErr error
		podRespErr      error
		assert          func(t *testing.T, err error, pod v1.Pod)
	}{
		// Test error on services check (we check if the job service exists before creating container)
		{
			container:       shared.Container{},
			servicesRespErr: errors.New("bad_error"),
			podRespErr:      nil,
			assert: func(t *testing.T, err error, pod v1.Pod) {
				test.ExpectErrorLike(t, errors.New("bad_error"), err)
			},
		},
		// Test error on pod create
		{
			container:       shared.Container{},
			servicesRespErr: nil,
			podRespErr:      errors.New("bad_create_error"),
			assert: func(t *testing.T, err error, pod v1.Pod) {
				test.ExpectErrorLike(t, errors.New("bad_create_error"), err)
			},
		},
		// Test simple creation with entry point and args, verify image name
		{
			container: shared.Container{
				Image:      "someimage",
				EntryPoint: "my_entry",
				Args:       []string{"my_arg"},
			},
			servicesRespErr: nil,
			podRespErr:      nil,
			assert: func(t *testing.T, err error, pod v1.Pod) {
				test.ExpectString(t, "someimage", pod.Spec.Containers[0].Image)
				test.ExpectString(t, "my_entry", pod.Spec.Containers[0].Command[0])
				test.ExpectString(t, "my_arg", pod.Spec.Containers[0].Args[0])
			},
		},
		// Test environment variables
		{
			container: shared.Container{
				Environment: map[string]string{
					"test_key": "test_val",
				},
			},
			servicesRespErr: nil,
			podRespErr:      nil,
			assert: func(t *testing.T, err error, pod v1.Pod) {
				test.ExpectString(t, "test_key", pod.Spec.Containers[0].Env[0].Name)
				test.ExpectString(t, "test_val", pod.Spec.Containers[0].Env[0].Value)
			},
		},
		// Test workspace mounting
		{
			container: shared.Container{
				WorkingDir: "/workdir",
			},
			servicesRespErr: nil,
			podRespErr:      nil,
			assert: func(t *testing.T, err error, pod v1.Pod) {
				test.ExpectString(t, "/workdir", pod.Spec.Containers[0].VolumeMounts[0].MountPath)
			},
		},
		// Test resource limits
		{
			container: shared.Container{
				Resources: &shared.ContainerResources{
					Requests: &shared.ContainerResourcesUnits{
						CPU:    0.15,
						Memory: "150m",
					},
					Limits: &shared.ContainerResourcesUnits{
						CPU:    0.2,
						Memory: "200m",
					},
				},
			},
			servicesRespErr: nil,
			podRespErr:      nil,
			assert: func(t *testing.T, err error, pod v1.Pod) {
				test.ExpectString(t, "150m", pod.Spec.Containers[0].Resources.Requests.Cpu().String())
				test.ExpectString(t, "157286400", pod.Spec.Containers[0].Resources.Requests.Memory().String())

				test.ExpectString(t, "200m", pod.Spec.Containers[0].Resources.Limits.Cpu().String())
				test.ExpectString(t, "209715200", pod.Spec.Containers[0].Resources.Limits.Memory().String())
			},
		},
	}

	for i, suite := range suites {
		t.Run(
			fmt.Sprintf("suites[%d]", i),
			func(t *testing.T) {
				fakeRest := fakeRESTClient()
				kubeRuntime := runtime.KubeRuntime{
					Client:    &fakeRest,
					Namespace: "test",
				}

				var podRequest v1.Pod
				fakeRest.Client = fake.CreateHTTPClient(
					func(req *http.Request) (response *http.Response, e error) {
						// Mock the service request
						if strings.Contains(req.URL.Path, string(v1.ResourceServices)) {
							if suite.servicesRespErr != nil {
								return nil, suite.servicesRespErr
							}
							return &http.Response{
								StatusCode: 200,
								Body: &test.NoOpReadCloser{
									Reader: bytes.NewReader([]byte("")),
								},
							}, nil
						} else {
							err := json.NewDecoder(req.Body).Decode(&podRequest)
							if err != nil {
								return nil, err
							}

							if suite.podRespErr != nil {
								return nil, suite.podRespErr
							}

							return &http.Response{
								StatusCode: 200,
								Body: &test.NoOpReadCloser{
									Reader: bytes.NewReader([]byte("")),
								},
							}, nil
						}
					},
				)

				timeoutCtx, cancel := context.WithTimeout(context.Background(), time.Second)
				_, err := kubeRuntime.DispatchContainer(
					timeoutCtx,
					shared.JobID("id"),
					suite.container,
				)
				suite.assert(t, err, podRequest)
				cancel()
			},
		)
	}
}

func TestKubeRuntime_CopyLogsForContainer(t *testing.T) {
	suites := []struct {
		logResponseError  error
		logResponseReader io.Reader
		podResponse       v1.Pod
		podResponseError  error
		assert            func(t *testing.T, err error, stdOut string, stdErr string)
	}{
		// Test we error out if we get an error getting container logs
		{
			logResponseReader: bytes.NewReader([]byte("")),
			podResponse:       v1.Pod{},
			logResponseError:  errors.New("error_get_logs"),
			assert: func(t *testing.T, err error, stdOut string, stdErr string) {
				test.ExpectErrorLike(t, errors.New("error_get_logs"), err)
			},
		},
		// Test we error out if there is an error getting the pod
		{
			logResponseReader: bytes.NewReader([]byte("")),
			podResponse:       v1.Pod{},
			podResponseError:  errors.New("error_get_pods"),
			assert: func(t *testing.T, err error, stdOut string, stdErr string) {
				test.ExpectErrorLike(t, errors.New("error_get_pods"), err)
			},
		},
		// Test reading badly formatted log
		{
			logResponseReader: bytes.NewReader([]byte("sdsd{}\nasdasd\n")),
			podResponse: v1.Pod{
				Status: v1.PodStatus{
					ContainerStatuses: []v1.ContainerStatus{
						{
							State: v1.ContainerState{},
						},
					},
				},
			},
			assert: func(t *testing.T, err error, stdOut string, stdErr string) {
				test.ExpectErrorLike(t, errors.New("error decoding log"), err)
			},
		},
		// Test reading badly formatted log on writer close
		{
			logResponseReader: bytes.NewReader([]byte("asdasd")),
			podResponse: v1.Pod{
				Status: v1.PodStatus{
					ContainerStatuses: []v1.ContainerStatus{
						{
							State: v1.ContainerState{},
						},
					},
				},
			},
			assert: func(t *testing.T, err error, stdOut string, stdErr string) {
				test.ExpectErrorLike(t, errors.New("error decoding log"), err)
			},
		},
		// Test simple stdout and stderr logs
		{
			logResponseReader: bytes.NewReader([]byte(`{ "log": "stdout", "stream" : "stdout", "time":"2019-05-02T19:42:15.922378Z" }
{ "log": "stderr", "stream" : "stderr", "time":"2019-05-02T19:42:16.922378Z"  }`)),
			podResponse: v1.Pod{
				Status: v1.PodStatus{
					ContainerStatuses: []v1.ContainerStatus{
						{
							State: v1.ContainerState{
								Terminated: &v1.ContainerStateTerminated{},
							},
						},
					},
				},
			},
			assert: func(t *testing.T, err error, stdOut string, stdErr string) {
				test.ExpectError(t, nil, err)
				test.ExpectString(t, "stdout", stdOut)
				test.ExpectString(t, "stderr", stdErr)
			},
		},
	}

	for i, suite := range suites {
		t.Run(
			fmt.Sprintf("suites[%d]", i),
			func(t *testing.T) {
				fakeRest := fakeRESTClient()
				kubeRuntime := runtime.KubeRuntime{
					Client:    &fakeRest,
					Namespace: "test",
				}

				fakeRest.Client = fake.CreateHTTPClient(
					func(req *http.Request) (response *http.Response, e error) {
						// Mock the service request
						if strings.Contains(req.URL.Path, string("log")) {
							if suite.logResponseError != nil {
								return nil, suite.logResponseError
							}
							return &http.Response{
								StatusCode: 200,
								Body: &test.NoOpReadCloser{
									Reader: suite.logResponseReader,
								},
							}, nil
						} else {
							if suite.podResponseError != nil {
								return nil, suite.podResponseError
							}

							podBytes, err := json.Marshal(suite.podResponse)
							if err != nil {
								return nil, err
							}
							return &http.Response{
								StatusCode: 200,
								Body: &test.NoOpReadCloser{
									Reader: bytes.NewReader(podBytes),
								},
							}, nil
						}
					},
				)

				stdErrBuffer := bytes.NewBuffer([]byte(""))
				stdOutBuffer := bytes.NewBuffer([]byte(""))
				copyErr := kubeRuntime.CopyLogsForContainer(
					context.TODO(),
					shared.ContainerID("id"),
					&test.NoOpWriteCloser{Writer: stdOutBuffer},
					&test.NoOpWriteCloser{Writer: stdErrBuffer},
				)
				stdOutString, err := ioutil.ReadAll(stdOutBuffer)
				if err != nil {
					t.Fatal(err)
				}
				stdErrString, err := ioutil.ReadAll(stdErrBuffer)
				if err != nil {
					t.Fatal(err)
				}
				suite.assert(t, copyErr, string(stdOutString), string(stdErrString))
			},
		)
	}
}
