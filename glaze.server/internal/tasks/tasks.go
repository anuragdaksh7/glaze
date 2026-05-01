package tasks

import (
	"encoding/json"

	"github.com/hibiken/asynq"
)

const (
	TypeBuildDeployment = "deployment:build"
)

type BuildPayload struct {
	DeploymentID string `json:"deployment_id"`
	RepoFullName string `json:"repo_full_name"`
}

func NewBuildTask(deploymentID, repoFullName string) (*asynq.Task, error) {
	payload, err := json.Marshal(BuildPayload{
		DeploymentID: deploymentID,
		RepoFullName: repoFullName,
	})
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(TypeBuildDeployment, payload), nil
}
