/*
 * Author: Lewis Maitland
 *
 * Copyright (c) 2019 Lewis Maitland
 */

package runtime

import (
	dockerclient "github.com/docker/docker/client"
	"k8s.io/client-go/rest"
)

type KubeRuntimeFactory struct {
	Client          rest.Interface
	Namespace       string
	VolumeClaimName string
}

func (factory *KubeRuntimeFactory) Create() (Runtime, error) {
	return &KubeRuntime{
		Client:          factory.Client,
		Namespace:       factory.Namespace,
		VolumeClaimName: factory.VolumeClaimName,
	}, nil
}

type DockerRuntimeFactory struct {
	Client dockerclient.CommonAPIClient
}

func (factory *DockerRuntimeFactory) Create() (Runtime, error) {
	return &DockerRuntime{
		Client: factory.Client,
	}, nil
}
