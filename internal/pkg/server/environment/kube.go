package environment

import (
	"fmt"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type KubeEnvironmentProviderFactory struct {
	Namespace string
	Client    kubernetes.Interface
}

func (factory *KubeEnvironmentProviderFactory) Create(searchPath []string) kubeEnvironmentProvider {
	return kubeEnvironmentProvider{
		namespace:   factory.Namespace,
		client:      factory.Client,
		searchPaths: searchPath,
	}
}

type kubeEnvironmentProvider struct {
	namespace   string
	client      kubernetes.Interface
	searchPaths []string
}

func (kube *kubeEnvironmentProvider) GetSecret(name string) (string, error) {
	var foundSecret string
	for _, path := range kube.searchPaths {
		secrets, err := kube.
			client.
			CoreV1().
			Secrets(kube.namespace).
			List(v1.ListOptions{FieldSelector: "metadata.name=" + path})

		if err != nil {
			return "", errors.Wrap(err, "error retrieving secrets")
		}

		// If our search path returned a secret, attempt to get the value we want from it it
		if len(secrets.Items) == 1 && secrets.Items[0].Data != nil {
			if secret, ok := secrets.Items[0].Data[name]; ok {
				foundSecret = string(secret)
				break
			}
		}
	}

	if foundSecret == "" {
		return "", fmt.Errorf("failed to find secret: '%s'", name)
	}
	return foundSecret, nil
}

func (kube *kubeEnvironmentProvider) GetValue(name string) (string, error) {
	var foundConfig string
	for _, path := range kube.searchPaths {
		configs, err := kube.
			client.
			CoreV1().
			ConfigMaps(kube.namespace).
			List(v1.ListOptions{FieldSelector: "metadata.name=" + path})

		if err != nil {
			return "", errors.Wrap(err, "error retrieving secrets")
		}

		// If our search path returned a secret, attempt to get the value we want from it it
		if len(configs.Items) == 1 && configs.Items[0].Data != nil {
			if configValue, ok := configs.Items[0].Data[name]; ok {
				foundConfig = string(configValue)
				break
			}
		}
	}

	if foundConfig == "" {
		return "", fmt.Errorf("failed to find secret: '%s'", name)
	}
	return foundConfig, nil
}
