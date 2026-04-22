package workspace

import (
	"errors"
	workspaceDto "glaze/dto/workspace"
	"glaze/logger"
	"glaze/response"
	"glaze/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Handler struct {
	Service
}

func NewHandler(s Service) *Handler {
	return &Handler{s}
}

func (h *Handler) GetAllWorkspaces(c *gin.Context) {
	user, err := utils.ExtractUser(c)
	if err != nil {
		response.Unauthorized(c, errors.New("unauthorized"))
		return
	}

	workspaces, err := h.Service.GetAllWorkspaces(c, user.ID)
	if err != nil {
		logger.Logger.Error("error getting workspaces", zap.Error(err))
		response.InternalError(c, errors.New("error getting workspaces"))
		return
	}

	response.OK(c, workspaces)
}

func (h *Handler) CreateWorkspace(c *gin.Context) {
	user, err := utils.ExtractUser(c)
	if err != nil {
		response.Unauthorized(c, errors.New("unauthorized"))
		return
	}

	var req workspaceDto.CreateWorkspaceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Logger.Error("Failed to bind json ", zap.Error(err))
		response.BadRequest(c, errors.New("invalid request body"))
		return
	}

	workspace, err := h.Service.CreateWorkspace(c, user.ID, req.Name)
	if err != nil {
		logger.Logger.Error("Failed to create workspace", zap.Error(err))
		response.InternalError(c, err)
		return
	}

	response.OK(c, workspace)
}

func (h *Handler) GetWorkspace(c *gin.Context) {
	user, err := utils.ExtractUser(c)
	if err != nil {
		response.Unauthorized(c, errors.New("unauthorized"))
		return
	}

	var req workspaceDto.GetWorkspaceByIDReq

	if err := c.ShouldBindUri(&req); err != nil {
		response.BadRequest(c, errors.New("invalid request"))
		return
	}

	id, err := uuid.Parse(req.ID)
	if err != nil {
		response.BadRequest(c, errors.New("invalid request"))
		return
	}
	workspace, err := h.Service.GetWorkspace(c, user.ID, id)
	if err != nil {
		logger.Logger.Error("Failed to get workspace", zap.Error(err))
		response.InternalError(c, err)
		return
	}

	response.OK(c, workspace)
}
