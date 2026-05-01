package task

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
	CommitHash   string `json:"commit_hash"`
}

func NewBuildTask(deploymentID, repoFullName, commitHash string) (*asynq.Task, error) {
	payload, err := json.Marshal(BuildPayload{
		DeploymentID: deploymentID,
		RepoFullName: repoFullName,
		CommitHash:   commitHash,
	})
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(TypeBuildDeployment, payload), nil
}
