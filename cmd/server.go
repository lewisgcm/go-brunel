/*
 * Author: Lewis Maitland
 *
 * Copyright (c) 2019 Lewis Maitland
 */

package main

import (
	"fmt"
	"go-brunel/internal/pkg/server"
	"go-brunel/internal/pkg/server/endpoint/api/container"
	"go-brunel/internal/pkg/server/endpoint/api/environment"
	"go-brunel/internal/pkg/server/endpoint/api/hook"
	"go-brunel/internal/pkg/server/endpoint/api/job"
	"go-brunel/internal/pkg/server/endpoint/api/repository"
	"go-brunel/internal/pkg/server/endpoint/api/user"
	"go-brunel/internal/pkg/server/endpoint/remote"
	"net/http"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/spf13/viper"
)

func FileServer(router *chi.Mux) {
	root := "./web"
	fs := http.FileServer(http.Dir(root))

	router.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		if _, err := os.Stat(root + r.RequestURI); os.IsNotExist(err) {
			http.StripPrefix(r.RequestURI, fs).ServeHTTP(w, r)
		} else {
			fs.ServeHTTP(w, r)
		}
	})
}

func loadServerConfig() (error, server.Config) {
	serverConfig := server.Config{}

	conf := viper.New()

	conf.AutomaticEnv()
	conf.SetEnvPrefix("BRUNEL")
	conf.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	conf.SetConfigName("brunel")
	conf.AddConfigPath("./")
	err := conf.ReadInConfig()

	conf.SetDefault("listen", ":8085")

	if err != nil {
		switch err.(type) {
		default:
			panic(fmt.Errorf("fatal error loading config file: %s \n", err))
		case viper.ConfigFileNotFoundError:
			log.Warning("no config file found. Using defaults and environment variables")
		}
	}

	return conf.Unmarshal(&serverConfig), serverConfig
}

func main() {
	err, serverConfig := loadServerConfig()
	if err != nil {
		log.Fatal(err)
	}

	err = serverConfig.Validate()
	if err != nil {
		log.Fatal(err)
	}

	jobStore, err := serverConfig.GetJobStore()
	if err != nil {
		log.Fatal(err)
	}

	userStore, err := serverConfig.GetUserStore()
	if err != nil {
		log.Fatal(err)
	}

	repositoryStore, err := serverConfig.GetRepositoryStore()
	if err != nil {
		log.Fatal(err)
	}

	environmentStore, err := serverConfig.GetEnvironmentStore()
	if err != nil {
		log.Fatal(err)
	}

	logStore, err := serverConfig.GetLogStore()
	if err != nil {
		log.Fatal(err)
	}

	containerStore, err := serverConfig.GetContainerStore()
	if err != nil {
		log.Fatal(err)
	}

	stageStore, err := serverConfig.GetStageStore()
	if err != nil {
		log.Fatal(err)
	}

	notifier, err := serverConfig.GetNotifier()
	if err != nil {
		log.Fatal(err)
	}

	oauths, err := serverConfig.GetOAuthProviders()
	if err != nil {
		log.Fatal(err)
	}

	err = remote.Server(
		jobStore,
		logStore,
		containerStore,
		repositoryStore,
		environmentStore,
		stageStore,
		notifier,
		*serverConfig.Remote.Credentials,
		serverConfig.Remote.Listen,
	)
	if err != nil {
		log.Fatal(err)
	}

	jwtSerializer := serverConfig.GetJWTSerializer()

	router := chi.NewRouter()
	router.Use(
		middleware.DefaultCompress,
		middleware.RedirectSlashes,
		// security.Middleware("keymatch_model.conf", "routes.csv", jwtSerializer),
		middleware.Recoverer,
	)
	router.Mount("/api/hook", hook.Routes(serverConfig.WebHook, jobStore, repositoryStore, notifier))
	router.Mount("/api/environment", environment.Routes(environmentStore))
	router.Mount("/api/repository", repository.Routes(repositoryStore, jobStore))
	router.Mount("/api/job", job.Routes(jobStore, logStore, stageStore, containerStore, repositoryStore, jwtSerializer))
	router.Mount("/api/container", container.Routes(logStore, containerStore, jwtSerializer))
	router.Mount("/api/user", user.Routes(userStore, oauths, jwtSerializer))
	FileServer(router)

	walkFunc := func(method string, route string, handler http.Handler, middleware ...func(http.Handler) http.Handler) error {
		log.Info("registering route: ", method, " ", route)
		return nil
	}
	if err := chi.Walk(router, walkFunc); err != nil {
		log.Fatal(err.Error())
	}

	log.Infof("listening for http connections on: '%s'", serverConfig.Listen)
	log.Fatal(http.ListenAndServe(serverConfig.Listen, router))
}
