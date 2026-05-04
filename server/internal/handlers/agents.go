package handlers

import (
	"errors"
	"net/http"

	"github.com/diyorbek/sentinel/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type createAgentResponse struct {
	Id           uuid.UUID `json:"id"`
	KafkaBrokers []string  `json:"kafka_brokers"`
	KafkaTopic   string    `json:"kafka_topic"`
	Message      string    `json:"message"`
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
// @Router /sentinel/agents [post]
// @Security BearerAuth
func (h *Handler) CreateAgent(ctx *gin.Context) {
	var body models.CreateAgent

	if err := ctx.ShouldBindJSON(&body); err != nil {
		h.logger.Warn(err.Error())
		errorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	userIDRaw, exists := ctx.Get("account_id")
	if !exists {
		h.logger.Warn("unauthorized: account_id not found in context")
		errorResponse(ctx, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	id, ok := userIDRaw.(string)
	if !ok {
		errorResponse(ctx, http.StatusBadRequest, errors.New("invalid account_id type"))
		return
	}

	body.Id = id

	resp, err := h.service.Agent.CreateAgent(body)
	if err != nil {
		h.logger.Error(err.Error())
		errorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	h.logger.Info("Agent created successfully (create agent)")
	ctx.JSON(http.StatusCreated, createAgentResponse{
		Id:           resp.ID,
		KafkaBrokers: resp.KafkaBrokers,
		KafkaTopic:   resp.KafkaTopic,
		Message:      "agent created successfully",
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
// @Router /sentinel/agents/{id} [get]
// @Security BearerAuth
func (h *Handler) GetAgentByID(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		h.logger.Warn("empty param value")
		errorResponse(ctx, http.StatusBadRequest, errors.New("param value is emtpy"))
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
// @Param account_id query string false "account id"
// @Param name query string false "agent name"
// @Param from query string false "Filter from last_seen (RFC3339 format)"
// @Param to query string false "Filter to last_seen (RFC3339 format)"
// @Param limit query int false "Limit" default(10)
// @Param page query int false "Page" default(1)
// @Success 200 {object} agentListResponse
// @Failure 400,401,404,500 {object} ErrorResponse
// @Router /sentinel/agents [get]
// @Security BearerAuth
func (h *Handler) ListAgents(ctx *gin.Context) {
	var filter models.FilterAgent

	if err := ctx.ShouldBindQuery(&filter); err != nil {
		h.logger.Warn(err.Error())
		errorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	userIDRaw, exists := ctx.Get("account_id")
	if !exists {
		h.logger.Warn("unauthorized: account_id not found in context")
		errorResponse(ctx, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	id, ok := userIDRaw.(string)
	if !ok {
		errorResponse(ctx, http.StatusBadRequest, errors.New("invalid account_id type"))
		return
	}

	filter.AccountID = id

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
// @Param update body models.HeartbeatRequest true "update last seen"
// @Success 200 {object} SuccessResponse
// @Failure 400,401,404,500 {object} ErrorResponse
// @Router /sentinel/agents/heartbeat [put]
// @Security BearerAuth
func (h *Handler) Heartbeat(ctx *gin.Context) {
	var body models.HeartbeatRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		h.logger.Warn(err.Error())
		errorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	err := h.service.Agent.UpdateLastSeen(body.AgentID)
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
// @Router /sentinel/agents/{id} [delete]
// @Security BearerAuth
func (h *Handler) DeleteAgent(ctx *gin.Context) {
	var id = ctx.Param("id")
	if id == "" {
		h.logger.Warn("empty param value")
		errorResponse(ctx, http.StatusBadRequest, errors.New("param value is empty"))
		return
	}

	err := h.service.Agent.DeleteAgent(id)
	if err != nil {
		h.logger.Error(err.Error())
		errorResponse(ctx, http.StatusInsufficientStorage, err)
		return
	}

	ctx.JSON(http.StatusOK, SuccessResponse{
		Message: "agent deleted successflly",
	})
}
