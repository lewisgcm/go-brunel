/*
 * Author: Lewis Maitland
 *
 * Copyright (c) 2019 Lewis Maitland
 */

package shared

import (
	"context"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/pkg/errors"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type RuntimeType string

type PersistenceType string

type EnvironmentProviderType string

type NotificationType string

// The constants defined here are all common configuration options used in brunel.
// At some point the plan would be to split these out into configuration structs for use
// in any methods that need them.
const (
	ConfigFile                              = "config-file"
	WorkingDirectory                        = "working-directory"
	EnvironmentFile                         = "env-file"
	PersistenceTypeMongo   PersistenceType  = "mongo"
	RuntimeTypeKubernetes  RuntimeType      = "kubernetes"
	NotificationTypeGitLab NotificationType = "gitlab"
)

type MongoConfig struct {
	Uri string
	DB  string
}

type KubernetesConfig struct {
	// ConfigFile is the location of the kubernetes configuration file to use when connecting to kubernetes
	ConfigFile string `mapstructure:"config-file"`

	// namespace is the kubernetes namespace to use for running jobs
	Namespace string

	// volumeClaimName is the name of the volume claim within kubernetes that we will use for cloning and building jobs
	VolumeClaimName string `mapstructure:"volume-claim-name"`
}

type GitLabConfig struct {
	URL    string
	Secret string
}

func (config *KubernetesConfig) GetKubernetesClient() (kubernetes.Interface, error) {
	kubernetesConfig, err := clientcmd.BuildConfigFromFlags("", config.ConfigFile)
	if err != nil {
		return nil, errors.Wrap(err, "error creating kubernetes client configuration")
	}

	kubernetesClient, err := kubernetes.NewForConfig(kubernetesConfig)
	if err != nil {
		return nil, errors.Wrap(err, "error creating kubernetes client")
	}
	return kubernetesClient, nil
}

func (config *KubernetesConfig) GetKubernetesRESTClient() (rest.Interface, error) {
	client, err := config.GetKubernetesClient()
	if err != nil {
		return nil, err
	}
	return client.CoreV1().RESTClient(), nil
}

func (config *MongoConfig) GetMongoDatabase() (*mongo.Database, error) {
	mongoClient, err := mongo.NewClient(config.Uri)
	if err != nil {
		return nil, errors.Wrap(err, "error creating mongodb client")
	}

	err = mongoClient.Connect(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, "error connecting to mongodb")
	}

	return mongoClient.Database(config.DB), nil
}
