/*
 * Author: Lewis Maitland
 *
 * Copyright (c) 2019 Lewis Maitland
 */

package recorder

import (
	"bytes"
	"encoding/json"
	"go-brunel/internal/pkg/shared"
	"log"
)

type LocalRecorder struct {
}

func (recorder LocalRecorder) RecordLog(jobID shared.JobID, message string, logType shared.LogType, stageID shared.StageID) error {
	log.Println(message)
	return nil
}

func (recorder LocalRecorder) RecordStageState(jobID shared.JobID, id shared.StageID, state shared.StageState) error {
	switch state {
	case shared.StageStateRunning:
		log.Println("running stage ", id)
		break
	case shared.StageStateError:
		log.Println("error executing stage ", id)
		break
	case shared.StageStateSuccess:
		log.Println("completed stage ", id)
		break
	}

	return nil
}

func (recorder LocalRecorder) RecordContainer(jobID shared.JobID, containerID shared.ContainerID, meta shared.ContainerMeta, container shared.Container, state shared.ContainerState) error {
	log.Printf("starting container in stage [%s] with id [%s] and spec: \n", meta.StageID, containerID)
	log.Println("---")
	b := bytes.Buffer{}
	if e := json.NewEncoder(&b).Encode(container); e != nil {
		return e
	}
	log.Print(b.String())
	log.Println("---")
	return nil
}

func (recorder LocalRecorder) RecordContainerState(containerID shared.ContainerID, state shared.ContainerState) error {
	t := "starting"
	switch state {
	case shared.ContainerStateRunning:
		t = "running"
	case shared.ContainerStateStopped:
		t = "stopped"
	}
	log.Printf("container with id %s is now %s", containerID, t)
	return nil
}

func (recorder LocalRecorder) RecordContainerLog(containerID shared.ContainerID, message string, logType shared.LogType) error {
	log.Printf("%s: %s", containerID, message)
	return nil
}
