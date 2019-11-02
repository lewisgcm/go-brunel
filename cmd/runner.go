/*
 * Author: Lewis Maitland
 *
 * Copyright (c) 2019 Lewis Maitland
 */

package main

import (
	"context"
	"flag"
	"go-brunel/internal/pkg/runner"
	"go-brunel/internal/pkg/shared"
	"log"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func main() {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	flag.String(shared.ConfigFile, "", "configuration file for the runner")
	flag.String(shared.WorkingDirectory, dir + "/", "the working directory for pipelines")
	flag.String(shared.EnvironmentFile, "", "an environment file to load environment variables from")

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	viper.AutomaticEnv()
	err = viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		log.Fatal(err)
	}

	configFile := viper.GetString(shared.ConfigFile)
	if configFile != "" {
		viper.SetConfigFile(configFile)
		if err = viper.ReadInConfig(); err != nil {
			log.Fatal(errors.Wrap(err, "error reading configuration"))
		}
	}

	var runnerConfig runner.Config
	err = viper.Unmarshal(&runnerConfig)
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
