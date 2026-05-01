package webhooks

import (
	"context"
	webhookDto "glaze/dto/webhook"
)

type Service interface {
	Github(c context.Context, eventType string) (string, error)
	GithubPush(c context.Context, payload webhookDto.PushPayload) (string, error)
	createDeployment(c context.Context) (string, error)
}
