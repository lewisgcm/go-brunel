/*
 * Author: Lewis Maitland
 *
 * Copyright (c) 2019 Lewis Maitland
 */

package main

import (
	"context"
	"fmt"
	"go-brunel/internal/pkg/runner"
	"log"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

func loadRunnerConfig() (error, runner.Config) {
	runnerConfig := runner.Config{}

	conf := viper.New()

	conf.AutomaticEnv()
	conf.SetEnvPrefix("BRUNEL")
	conf.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	conf.SetConfigName("runner")
	conf.AddConfigPath("./")
	err := conf.ReadInConfig()

	if err != nil {
		switch err.(type) {
		default:
			panic(fmt.Errorf("fatal error loading config file: %s \n", err))
		case viper.ConfigFileNotFoundError:
			log.Println("no config file found. Using defaults and environment variables")
		}
	}

	return conf.Unmarshal(&runnerConfig), runnerConfig
}

func main() {
	err, runnerConfig := loadRunnerConfig()
	if err != nil {
		log.Fatal(err)
	}

	err = runnerConfig.Valid()
	if err != nil {
		log.Fatal(errors.Wrap(err, "supplied configuration or command flags are invalid"))
	}

	log.Println("creating job trigger")
	jobTrigger, err := runnerConfig.Trigger()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("configuring job handler")
	jobHandler, err := runnerConfig.JobHandler()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("waiting for jobs...")
	for event := range jobTrigger.Await(context.Background()) {
		jobHandler.Handle(event)
	}
}
