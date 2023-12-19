package app

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/alloyzeus/go-azfl/errors"
)

// https://docs.aws.amazon.com/AmazonECS/latest/developerguide/task-metadata-endpoint-v3.html
type ecsTaskMetadata struct {
	Cluster  string `json:"Cluster"`
	TaskARN  string `json:"TaskARN"`
	Revision string `json:"Revision"`
}

func getECSTaskID() (taskID string, taskRev string, err error) {
	ecsContainerMetadataURI := os.Getenv("ECS_CONTAINER_METADATA_URI")
	if ecsContainerMetadataURI == "" {
		return "", "", nil
	}

	client := &http.Client{}
	resp, err := client.Get(ecsContainerMetadataURI + "/task")
	if err != nil {
		return "", "", errors.Wrap("task metadata fetch", err)
	}

	var taskMetadata ecsTaskMetadata
	err = json.NewDecoder(resp.Body).Decode(&taskMetadata)
	if err != nil {
		return "", "", errors.Wrap("task metadata loading", err)
	}

	parts := strings.Split(taskMetadata.TaskARN, "/")
	if len(parts) == 2 {
		return parts[1], taskMetadata.Revision, nil
	}
	return taskMetadata.TaskARN, taskMetadata.Revision, nil
}
