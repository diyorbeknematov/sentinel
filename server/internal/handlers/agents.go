package handlers

import (
	"errors"
	"net/http"

	"github.com/diyorbek/sentinel/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type createAgentResponse struct {
	Id     uuid.UUID `json:"id"`
	Messge string    `json:"message"`
}

type agentListResponse struct {
	Data  []models.Agent `json:"data"`
	Total int            `json:"total"`
	Page  int            `json:"page"`
	Limit int            `json:"limit"`
}

// @Description Create Agent
// @Summary Create Agent
// @Tags agents
// @Accept json
// @Produce json 
// @Param create body models.CreateAgent true "Create Agent"
// @Success 201 {object} createAgentResponse
// @Failure 400,404,500 {object} ErrorResponse
// @Router /sentinel/api/agents [post]
func (h *Handler) CreateAgent(ctx *gin.Context) {
	var body models.CreateAgent

	if err := ctx.ShouldBindJSON(&body); err != nil {
		h.logger.Warn(err.Error())
		errorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	id, err := h.service.Agent.CreateAgent(body)
	if err != nil {
		h.logger.Error(err.Error())
		errorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	h.logger.Info("Agent created successfully (create agent)")
	ctx.JSON(http.StatusCreated, createAgentResponse{
		Id:     id,
		Messge: "agent created successfully",
	})
}

// @Description Get Agent by id 
// @Summary Get Agent by id 
// @Tags agents
// @Accept json 
// @Produce json 
// @Param id path string true "agent id"
// @Success 200 {object} models.Agent
// @Failure 400,401,404,500 {object} ErrorResponse
// @Router /sentinel/api/agents/{id} [get]
func (h *Handler) GetAgentByID(ctx *gin.Context) {
	paramValue := ctx.Param("id")
	if paramValue == "" {
		h.logger.Warn("empty param value")
		errorResponse(ctx, http.StatusBadRequest, errors.New("param value is emtpy"))
		return
	}

	id, err := uuid.Parse(paramValue)
	if err != nil {
		h.logger.Error(err.Error())
		errorResponse(ctx, http.StatusBadRequest, errors.New("agent id is wrong"))
		return
	}

	agent, err := h.service.Agent.GetAgentByID(id)
	if err != nil {
		h.logger.Error(err.Error())
		errorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, agent)
}

// @Description Retrieve a paginated list of agents 
// @Summary Get list of agents
// @Tags agents
// @Accept json
// @Produce json
// @Param name query string false "agent name"
// @Param from query string false "Filter from last_seen (RFC3339 format)"
// @Param to query string false "Filter to last_seen (RFC3339 format)"
// @Param limit query int false "Limit" default(10)
// @Param page query int false "Page" default(1)
// @Success 200 {object} agentListResponse
// @Failure 400,401,404,500 {object} ErrorResponse
// @Router /sentinel/api/agents [get]
func (h *Handler) ListAgents(ctx *gin.Context) {
	var filter models.FilterAgent

	if err := ctx.ShouldBindQuery(&filter); err != nil {
		h.logger.Warn(err.Error())
		errorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	page, err := getPageQuery(ctx)
	if err != nil {
		h.logger.Warn(err.Error())
		errorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	filter.Offset = (page - 1) * filter.Limit

	agents, total, err := h.service.Agent.ListAgents(filter)
	if err != nil {
		h.logger.Error(err.Error())
		errorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, agentListResponse{
		Data:  agents,
		Total: total,
		Limit: filter.Limit,
		Page:  page,
	})
}

// @Description Update Last Seen
// @Summary Update Last Seen
// @Tags agents
// @Accept json
// @Produce json 
// @Param id path string true "agent id"
// @Param update body models.UpdateLastSeen true "update last seen"
// @Success 200 {object} SuccessResponse
// @Failure 400,401,404,500 {object} ErrorResponse
// @Router /sentinel/api/agents/{id}/lastseen [put]
func (h *Handler) UpdateLastSeen(ctx *gin.Context) {
	var paramValue = ctx.Param("id")
	if paramValue == "" {
		h.logger.Warn("empty param value")
		errorResponse(ctx, http.StatusBadRequest, errors.New("param value is empty"))
		return
	}

	id, err := uuid.Parse(paramValue)
	if err != nil {
		h.logger.Error(err.Error())
		errorResponse(ctx, http.StatusBadRequest, errors.New("user id is wrong"))
		return
	}

	var body models.UpdateLastSeen
	if err := ctx.ShouldBindJSON(&body); err != nil {
		h.logger.Warn(err.Error())
		errorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	body.Id = id
	err = h.service.Agent.UpdateLastSeen(body)
	if err != nil {
		h.logger.Error(err.Error())
		errorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, SuccessResponse{
		Message: "updated lastseen successfully",
	})
}

// @Description Delete Agent
// @Summary Delete Agent
// @Tags agents 
// @Accept json 
// @Produce json 
// @Param id path string true "agent id"
// @Success 200 {object} SuccessResponse
// @Failure 400,401,404,500 {object} ErrorResponse
// @Router /sentinel/api/agents/{id} [delete]
func (h *Handler) DeleteAgent(ctx *gin.Context) {
	var paramValue = ctx.Param("id")
	if paramValue == "" {
		h.logger.Warn("empty param value")
		errorResponse(ctx, http.StatusBadRequest, errors.New("param value is empty"))
		return
	}

	id, err := uuid.Parse(paramValue)
	if err != nil {
		h.logger.Error(err.Error())
		errorResponse(ctx, http.StatusBadRequest, errors.New("user id is wrong"))
		return
	}

	err = h.service.Agent.DeleteAgent(id)
	if err != nil {
		h.logger.Error(err.Error())
		errorResponse(ctx, http.StatusInsufficientStorage, err)
		return
	}

	ctx.JSON(http.StatusOK, SuccessResponse{
		Message: "agent deleted successflly",
	})
}