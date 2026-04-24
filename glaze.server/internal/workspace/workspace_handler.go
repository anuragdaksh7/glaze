package workspace

import (
	"errors"
	integrationDto "glaze/dto/integration"
	workspaceDto "glaze/dto/workspace"
	"glaze/logger"
	"glaze/models"
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

func (h *Handler) UpdateWorkspace(c *gin.Context) {
	user, err := utils.ExtractUser(c)
	if err != nil {
		response.Unauthorized(c, errors.New("unauthorized"))
		return
	}

	var reqUri workspaceDto.GetWorkspaceByIDReq
	if err := c.ShouldBindUri(&reqUri); err != nil {
		response.BadRequest(c, errors.New("invalid request params"))
		return
	}
	workspaceID, err := uuid.Parse(reqUri.ID)
	if err != nil {
		response.BadRequest(c, errors.New("invalid workspace id"))
		return
	}

	var reqBody workspaceDto.UpdateWorkspaceRequest
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		response.BadRequest(c, errors.New("invalid request body"))
		return
	}

	res, err := h.Service.UpdateWorkspace(c, user.ID, workspaceID, &reqBody)
	if err != nil {
		logger.Logger.Error("Failed to update workspace", zap.Error(err))
		if err.Error() == "action not allowed: must be owner or admin" || err.Error() == "unauthorized or workspace not found" {
			response.Unauthorized(c, err)
			return
		}
		response.InternalError(c, err)
		return
	}

	response.OK(c, res)
}

func (h *Handler) DeleteWorkspace(c *gin.Context) {
	user, err := utils.ExtractUser(c)
	if err != nil {
		response.Unauthorized(c, errors.New("unauthorized"))
		return
	}

	var reqUri workspaceDto.GetWorkspaceByIDReq
	if err := c.ShouldBindUri(&reqUri); err != nil {
		response.BadRequest(c, errors.New("invalid request params"))
		return
	}
	workspaceID, err := uuid.Parse(reqUri.ID)
	if err != nil {
		response.BadRequest(c, errors.New("invalid workspace id"))
		return
	}

	err = h.Service.DeleteWorkspace(c, user.ID, workspaceID)
	if err != nil {
		logger.Logger.Error("Failed to delete workspace", zap.Error(err))
		if err.Error() == "action not allowed: must be owner" || err.Error() == "unauthorized or workspace not found" {
			response.Unauthorized(c, err)
			return
		}
		response.InternalError(c, err)
		return
	}

	response.OK(c, map[string]string{"message": "workspace deleted successfully"})
}

func (h *Handler) ListWorkspaceMembers(c *gin.Context) {
	user, err := utils.ExtractUser(c)
	if err != nil {
		response.Unauthorized(c, errors.New("unauthorized"))
		return
	}

	var reqUri workspaceDto.GetWorkspaceByIDReq
	if err := c.ShouldBindUri(&reqUri); err != nil {
		response.BadRequest(c, errors.New("invalid request params"))
		return
	}
	workspaceID, err := uuid.Parse(reqUri.ID)
	if err != nil {
		response.BadRequest(c, errors.New("invalid workspace id"))
		return
	}

	members, err := h.Service.ListWorkspaceMembers(c, user.ID, workspaceID)
	if err != nil {
		logger.Logger.Error("Failed to list workspace members", zap.Error(err))
		response.InternalError(c, err)
		return
	}

	response.OK(c, members)
}

func (h *Handler) UpdateWorkspaceMemberRole(c *gin.Context) {
	user, err := utils.ExtractUser(c)
	if err != nil {
		response.Unauthorized(c, errors.New("unauthorized"))
		return
	}

	var reqUri workspaceDto.WorkspaceMemberReq
	if err := c.ShouldBindUri(&reqUri); err != nil {
		response.BadRequest(c, errors.New("invalid request params"))
		return
	}
	workspaceID, err := uuid.Parse(reqUri.WorkspaceID)
	if err != nil {
		response.BadRequest(c, errors.New("invalid workspace id"))
		return
	}
	targetUserID, err := uuid.Parse(reqUri.UserID)
	if err != nil {
		response.BadRequest(c, errors.New("invalid user id"))
		return
	}

	var reqBody workspaceDto.UpdateWorkspaceMemberRoleRequest
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		response.BadRequest(c, errors.New("invalid request body"))
		return
	}

	// Assuming you import models in the handler, or you just cast to models.WorkspaceRole
	err = h.Service.UpdateWorkspaceMemberRole(c, user.ID, workspaceID, targetUserID, models.WorkspaceRole(reqBody.Role))
	if err != nil {
		logger.Logger.Error("Failed to update workspace member role", zap.Error(err))
		response.InternalError(c, err)
		return
	}

	response.OK(c, map[string]string{"message": "member role updated successfully"})
}

func (h *Handler) RemoveWorkspaceMember(c *gin.Context) {
	user, err := utils.ExtractUser(c)
	if err != nil {
		response.Unauthorized(c, errors.New("unauthorized"))
		return
	}

	var reqUri workspaceDto.WorkspaceMemberReq
	if err := c.ShouldBindUri(&reqUri); err != nil {
		response.BadRequest(c, errors.New("invalid request params"))
		return
	}
	workspaceID, err := uuid.Parse(reqUri.WorkspaceID)
	if err != nil {
		response.BadRequest(c, errors.New("invalid workspace id"))
		return
	}
	targetUserID, err := uuid.Parse(reqUri.UserID)
	if err != nil {
		response.BadRequest(c, errors.New("invalid user id"))
		return
	}

	err = h.Service.RemoveWorkspaceMember(c, user.ID, workspaceID, targetUserID)
	if err != nil {
		logger.Logger.Error("Failed to remove workspace member", zap.Error(err))
		response.InternalError(c, err)
		return
	}

	response.OK(c, map[string]string{"message": "member removed successfully"})
}

func (h *Handler) GetIntegrations(c *gin.Context) {
	user, err := utils.ExtractUser(c)
	if err != nil {
		response.Unauthorized(c, errors.New("unauthorized"))
		return
	}

	var reqUri workspaceDto.GetWorkspaceByIDReq
	if err := c.ShouldBindUri(&reqUri); err != nil {
		response.BadRequest(c, errors.New("invalid request params"))
		return
	}
	workspaceID, err := uuid.Parse(reqUri.ID)
	if err != nil {
		response.BadRequest(c, errors.New("invalid workspace id"))
		return
	}

	res, err := h.Service.ListIntegrations(c, user.ID, workspaceID)
	if err != nil {
		logger.Logger.Error("Failed to list integrations", zap.Error(err))
		response.InternalError(c, err)
		return
	}

	response.OK(c, res)
}

func (h *Handler) ConnectGithub(c *gin.Context) {
	user, err := utils.ExtractUser(c)
	if err != nil {
		response.Unauthorized(c, errors.New("unauthorized"))
		return
	}

	var reqUri workspaceDto.GetWorkspaceByIDReq
	if err := c.ShouldBindUri(&reqUri); err != nil {
		response.BadRequest(c, errors.New("invalid request params"))
		return
	}
	workspaceID, err := uuid.Parse(reqUri.ID)
	if err != nil {
		response.BadRequest(c, errors.New("invalid workspace id"))
		return
	}

	url, err := h.Service.ConnectGithub(c, user.ID, workspaceID)
	if err != nil {
		logger.Logger.Error("Failed to connect github", zap.Error(err))
		response.InternalError(c, err)
		return
	}

	// redirect
	response.OK(c, map[string]string{"url": url})
}

func (h *Handler) GithubCallback(c *gin.Context) {
	user, err := utils.ExtractUser(c)
	if err != nil {
		response.Unauthorized(c, errors.New("unauthorized"))
		return
	}

	var reqParam integrationDto.GithubCallbackReq
	if err := c.ShouldBind(&reqParam); err != nil {
		response.BadRequest(c, errors.New("invalid request params"))
		return
	}

	res, err := h.Service.GithubCallback(c, user.ID, reqParam.Code, reqParam.State)
	if err != nil {
		logger.Logger.Error("Failed to github callback", zap.Error(err))
		response.InternalError(c, err)
		return
	}

	response.OK(c, res)
}

func (h *Handler) DeleteIntegration(c *gin.Context) {}
