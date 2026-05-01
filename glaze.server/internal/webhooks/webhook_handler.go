package webhooks

import (
	"errors"
	webhookDto "glaze/dto/webhook"
	"glaze/response"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Service
}

func NewHandler(s Service) *Handler {
	return &Handler{s}
}

func (h *Handler) Github(c *gin.Context) {
	eventType := c.GetHeader("X-GitHub-Event")

	var payload webhookDto.PushPayload
	if eventType == "push" {
		if err := c.ShouldBindJSON(&payload); err != nil {
			response.BadRequest(c, errors.New("invalid payload"))
			return
		}
		_, _ = h.Service.GithubPush(c, payload)
	}

	if eventType == "ping" {
		c.JSON(200, gin.H{"message": "pong"})
	}
	//github, err := h.Service.GithubPush(c, eventType, )
	//if err != nil {
	//	c.JSON(500, gin.H{"error": err.Error()})
	//	return
	//}

	c.JSON(201, gin.H{"message": eventType})
}
