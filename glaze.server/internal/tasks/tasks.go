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
	Branch       string `json:"branch"`
	CommitSHA    string `json:"commit_sha"`
}

func NewBuildTask(deploymentID, repoFullName, branch, commitSHA string) (*asynq.Task, error) {
	payload, err := json.Marshal(BuildPayload{
		DeploymentID: deploymentID,
		RepoFullName: repoFullName,
		Branch:       branch,
		CommitSHA:    commitSHA,
	})
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(TypeBuildDeployment, payload), nil
}
