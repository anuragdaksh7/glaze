package integrationDto

import (
	"time"

	"github.com/google/uuid"
)

type IntegrationResponse struct {
	ID          uuid.UUID `json:"id"`
	Provider    string    `json:"provider"`
	IsActive    bool      `json:"isActive"`
	AccountName string    `json:"accountName"` // e.g. "anurag-daksh"
	CreatedAt   time.Time `json:"createdAt"`
}

type GithubCallbackReq struct {
	Code  string `form:"code" json:"code" binding:"required"`
	State string `form:"state" json:"state" binding:"required"`
}
