package webhooks

import (
	"context"
	"errors"
	"fmt"
	"glaze/config"
	webhookDto "glaze/dto/webhook"
	"glaze/internal/tasks"
	"glaze/logger"
	"glaze/models"
	"strings"
	"time"

	"github.com/hibiken/asynq"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type service struct {
	time.Duration
	AsynqClient *asynq.Client
	DB          *gorm.DB
}

func NewService(asynqClient *asynq.Client) Service {
	return &service{
		time.Duration(20) * time.Second,
		asynqClient,
		config.DB,
	}
}

func (s *service) createDeployment(c context.Context) (string, error) {

	return "", nil
}

func (s *service) Github(c context.Context, eventType string) (string, error) {

	if eventType == "ping" {
		return "pong", nil
	}

	if eventType == "push" {
		logger.Logger.Info("Received github push event")

		return "push event", nil
	}
	return "", nil
}

func (s *service) GithubPush(c context.Context, payload webhookDto.PushPayload) (string, error) {
	var projects []models.Project
	err := s.DB.Joins("JOIN repositories ON repositories.id = projects.repository_id").
		Where("repositories.external_id = ?", fmt.Sprintf("%d", payload.Repository.ID)).
		Find(&projects).Error
	if err != nil {
		logger.Logger.Error("Failed to get repositories from database", zap.Error(err))
		return "", err
	}
	if len(projects) == 0 {
		logger.Logger.Error("Failed to get repositories from database")
		return "", errors.New("no repositories found")
	}

	branch := strings.TrimPrefix(payload.Ref, "refs/heads/")
	count := 0

	for _, project := range projects {
		if project.DeployBranch != branch {
			continue
		}

		deployment := &models.Deployment{
			ProjectID:  project.ID,
			CommitHash: payload.After,
			CommitMsg:  payload.HeadCommit.Message,
			AuthorName: payload.HeadCommit.Author.Name,
			Branch:     branch,
			Status:     models.StatusQueued,
		}
		s.DB.Create(deployment)

		task, _ := tasks.NewBuildTask(deployment.ID.String(), project.RepoFullName)
		s.AsynqClient.Enqueue(task)

		count++
	}

	return fmt.Sprintf("Triggered %d deployments", count), nil
}
