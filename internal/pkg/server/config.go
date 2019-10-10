package server

import (
	"fmt"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/gitlab"
	"github.com/pkg/errors"
	"go-brunel/internal/pkg/server/notify"
	"go-brunel/internal/pkg/server/security"
	"go-brunel/internal/pkg/server/store"
	"go-brunel/internal/pkg/server/store/mongo"
	"go-brunel/internal/pkg/shared"
	"go-brunel/internal/pkg/shared/remote"
)

type RemoteConfiguration struct {
	Listen      string
	Credentials *remote.Credentials
}

type JwtConfiguration struct {
	Secret      string
	DefaultRole string `mapstructure:"default-role"`
}

type OAuthConfiguration struct {
	Key    string
	Secret string
}

type Config struct {
	Persistence shared.PersistenceType
	Mongo       *shared.MongoConfig

	EnvironmentProvider shared.EnvironmentProviderType `mapstructure:"environment-provider"`
	Kubernetes          *shared.KubernetesConfig

	ServerName string `mapstructure:"server-name"`

	Notification shared.NotificationType
	GitLab       *shared.GitLabConfig

	Remote RemoteConfiguration

	Jwt JwtConfiguration

	OAuth map[string]OAuthConfiguration
}

func (config *Config) Validate() error {
	if config.Persistence == shared.PersistenceTypeMongo && config.Mongo == nil {
		return errors.New("invalid configuration, expecting mongoDB configuration when using mongoDB persistence")
	}
	if config.Jwt.Secret == "" || config.Jwt.DefaultRole == "" {
		return errors.New("jwt.secret and jwt.default-role must not be empty")
	}
	if config.ServerName == "" {
		return errors.New("server-name cannot be empty, this is the hostname users use to access brunel and is used for oauth")
	}
	for k, v := range config.OAuth {
		if v.Secret == "" || v.Key == "" {
			return fmt.Errorf("oauth.%s key or secret must not be empty", k)
		}
	}
	return nil
}

func (config *Config) GetNotifier() (notify.Notify, error) {
	//if config.Notification == shared.NotificationTypeGitLab {
	//	if config.GitLab == nil {
	//		return nil, errors.New("gitlab configuration must be specified when using gitlab build notifications")
	//	}
	//	repository, err := config.GetRepository()
	//	if err != nil {
	//		return nil, errors.Wrap(err, "error getting repository")
	//	}
	//	return &notify.GitLabNotify{
	//		URL:        config.GitLab.URL,
	//		Secret:     config.GitLab.Secret,
	//		Repository: repository,
	//	}, nil
	//}
	return &notify.TextNotify{}, nil
}

func (config *Config) GetJWTSerializer() security.TokenSerializer {
	return security.NewTokenSerializer(
		[]byte(config.Jwt.Secret),
		security.UserRole(config.Jwt.DefaultRole),
	)
}

func (config *Config) GetOAuthProviders() ([]goth.Provider, error) {
	var providers []goth.Provider
	for k, v := range config.OAuth {
		switch k {
		case "gitlab":
			providers = append(providers, gitlab.New(v.Key, v.Secret, config.ServerName+"/api/user/callback?provider=gitlab"))
		default:
			return providers, fmt.Errorf("uknown provider '%s'", k)
		}
	}
	return providers, nil
}

func (config *Config) GetJobRepository() (store.JobStore, error) {
	switch config.Persistence {
	case shared.PersistenceTypeMongo:
		if config.Mongo == nil {
			return nil, errors.New("no mongo configuration detected")
		}
		database, err := config.Mongo.GetMongoDatabase()
		if err != nil {
			return nil, err
		}
		return &mongo.JobStore{
			Database: database,
		}, nil
	default:
		return nil, errors.New("no persistence configuration detected")
	}
}

func (config *Config) GetUserRepository() (store.UserStore, error) {
	switch config.Persistence {
	case shared.PersistenceTypeMongo:
		if config.Mongo == nil {
			return nil, errors.New("no mongo configuration detected")
		}
		database, err := config.Mongo.GetMongoDatabase()
		if err != nil {
			return nil, err
		}
		return &mongo.UserStore{
			Database: database,
		}, nil
	default:
		return nil, errors.New("no persistence configuration detected")
	}
}

func (config *Config) GetRepositoryRepository() (store.RepositoryStore, error) {
	switch config.Persistence {
	case shared.PersistenceTypeMongo:
		if config.Mongo == nil {
			return nil, errors.New("no mongo configuration detected")
		}
		database, err := config.Mongo.GetMongoDatabase()
		if err != nil {
			return nil, err
		}
		return &mongo.RepositoryStore{
			Database: database,
		}, nil
	default:
		return nil, errors.New("no persistence configuration detected")
	}
}

func (config *Config) GetLogRepository() (store.LogStore, error) {
	switch config.Persistence {
	case shared.PersistenceTypeMongo:
		if config.Mongo == nil {
			return nil, errors.New("no mongo configuration detected")
		}
		database, err := config.Mongo.GetMongoDatabase()
		if err != nil {
			return nil, err
		}
		return &mongo.LogStore{
			Database: database,
		}, nil
	default:
		return nil, errors.New("no persistence configuration detected")
	}
}

func (config *Config) GetContainerRepository() (store.ContainerStore, error) {
	switch config.Persistence {
	case shared.PersistenceTypeMongo:
		if config.Mongo == nil {
			return nil, errors.New("no mongo configuration detected")
		}
		database, err := config.Mongo.GetMongoDatabase()
		if err != nil {
			return nil, err
		}
		return &mongo.ContainerStore{
			Database: database,
		}, nil
	default:
		return nil, errors.New("no persistence configuration detected")
	}
}
