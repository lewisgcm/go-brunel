/*
 * Author: Lewis Maitland
 *
 * Copyright (c) 2019 Lewis Maitland
 */

package main

import (
	"flag"
	"go-brunel/internal/pkg/server"
	"go-brunel/internal/pkg/server/endpoint/api/container"
	"go-brunel/internal/pkg/server/endpoint/api/hook"
	"go-brunel/internal/pkg/server/endpoint/api/job"
	"go-brunel/internal/pkg/server/endpoint/api/repository"
	"go-brunel/internal/pkg/server/endpoint/api/user"
	"go-brunel/internal/pkg/server/endpoint/remote"
	"go-brunel/internal/pkg/shared"
	"net/http"

	log "github.com/Sirupsen/logrus"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func main() {
	flag.String(shared.ConfigFile, "", "configuration file for the server")
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	viper.AutomaticEnv()
	err := viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		log.Fatal(err)
	}

	configFile := viper.GetString(shared.ConfigFile)
	if configFile != "" {
		viper.SetConfigFile(configFile)
		if err = viper.ReadInConfig(); err != nil {
			log.Fatal(errors.Wrap(err, "error reading configuration"))
		}
	} else {
		log.Fatal("no configuration file has been provided")
	}

	var serverConfig server.Config
	err = viper.Unmarshal(&serverConfig)
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
	router.Mount("/api/hook", hook.Routes(jobStore, repositoryStore, notifier))
	router.Mount("/api/repository", repository.Routes(repositoryStore, jobStore))
	router.Mount("/api/job", job.Routes(jobStore, logStore, containerStore, jwtSerializer))
	router.Mount("/api/container", container.Routes(logStore, jwtSerializer))
	router.Mount("/api/user", user.Routes(userStore, oauths, jwtSerializer))

	walkFunc := func(method string, route string, handler http.Handler, middleware ...func(http.Handler) http.Handler) error {
		log.Info("registering route: ", method, " ", route)
		return nil
	}
	if err := chi.Walk(router, walkFunc); err != nil {
		log.Fatal(err.Error())
	}

	log.Fatal(http.ListenAndServe(":8085", router))
}
