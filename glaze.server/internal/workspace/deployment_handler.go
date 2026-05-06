package workspace

import (
	"errors"
	deploymentDto "glaze/dto/deployment"
	"glaze/logger"
	"glaze/response"
	"glaze/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (h *Handler) ListDeployments(c *gin.Context) {
	user, err := utils.ExtractUser(c)
	if err != nil {
		response.Unauthorized(c, errors.New("unauthorized"))
		return
	}

	var req deploymentDto.DeploymentsByProjectReq
	if err := c.ShouldBindUri(&req); err != nil {
		response.BadRequest(c, errors.New("invalid request params"))
		return
	}

	workspaceID, err := uuid.Parse(req.WorkspaceID)
	if err != nil {
		response.BadRequest(c, errors.New("invalid workspace id"))
		return
	}
	projectID, err := uuid.Parse(req.ProjectID)
	if err != nil {
		response.BadRequest(c, errors.New("invalid project id"))
		return
	}

	res, err := h.Service.ListDeployments(c, user.ID, workspaceID, projectID)
	if err != nil {
		logger.Logger.Error("Failed to list deployments", zap.Error(err))
		response.InternalError(c, err)
		return
	}

	response.OK(c, res)
}

func (h *Handler) GetDeployment(c *gin.Context) {
	user, err := utils.ExtractUser(c)
	if err != nil {
		response.Unauthorized(c, errors.New("unauthorized"))
		return
	}

	var req deploymentDto.DeploymentByIDReq
	if err := c.ShouldBindUri(&req); err != nil {
		response.BadRequest(c, errors.New("invalid request params"))
		return
	}

	workspaceID, err := uuid.Parse(req.WorkspaceID)
	if err != nil {
		response.BadRequest(c, errors.New("invalid workspace id"))
		return
	}
	projectID, err := uuid.Parse(req.ProjectID)
	if err != nil {
		response.BadRequest(c, errors.New("invalid project id"))
		return
	}
	deploymentID, err := uuid.Parse(req.DeploymentID)
	if err != nil {
		response.BadRequest(c, errors.New("invalid deployment id"))
		return
	}

	res, err := h.Service.GetDeployment(c, user.ID, workspaceID, projectID, deploymentID)
	if err != nil {
		logger.Logger.Error("Failed to get deployment", zap.Error(err))
		response.InternalError(c, err)
		return
	}

	response.OK(c, res)
}

func (h *Handler) TriggerDeployment(c *gin.Context) {
	user, err := utils.ExtractUser(c)
	if err != nil {
		response.Unauthorized(c, errors.New("unauthorized"))
		return
	}

	var req deploymentDto.DeploymentsByProjectReq
	if err := c.ShouldBindUri(&req); err != nil {
		response.BadRequest(c, errors.New("invalid request params"))
		return
	}

	workspaceID, err := uuid.Parse(req.WorkspaceID)
	if err != nil {
		response.BadRequest(c, errors.New("invalid workspace id"))
		return
	}
	projectID, err := uuid.Parse(req.ProjectID)
	if err != nil {
		response.BadRequest(c, errors.New("invalid project id"))
		return
	}

	var body deploymentDto.CreateDeploymentRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		response.BadRequest(c, errors.New("invalid request body"))
		return
	}

	res, err := h.Service.TriggerDeployment(c, user.ID, workspaceID, projectID, &body)
	if err != nil {
		logger.Logger.Error("Failed to trigger deployment", zap.Error(err))
		response.InternalError(c, err)
		return
	}

	response.Created(c, res)
}
