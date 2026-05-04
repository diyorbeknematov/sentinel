package handlers

import (
	"errors"
	"net/http"

	"github.com/diyorbek/sentinel/internal/models"
	"github.com/gin-gonic/gin"
)

type statsResponse struct {
	Data models.DashboardStats `json:"data"`
}

// @Description Get dashboard stats
// @Summary Get dashboard statistics
// @Tags dashboard
// @Accept json
// @Produce json
// @Param agent_id query string false "agent id"
// @Param period query string false "time period: 1h, 24h, 7d" default(1h)
// @Success 200 {object} statsResponse
// @Failure 400,401,500 {object} ErrorResponse
// @Router /sentinel/stats [get]
// @Security ApiKeyAuth
func (h *Handler) GetDashboardStats(ctx *gin.Context) {
	var filter models.StatsFilter

	if err := ctx.ShouldBindQuery(&filter); err != nil {
		h.logger.Warn("invalid query params", "error", err)
		errorResponse(ctx, http.StatusBadRequest, errors.New("invalid query parameters"))
		return
	}

	// Default period to 1h if not provided
	if filter.Period == "" {
		filter.Period = "1h"
	}

	stats, err := h.service.Stats.GetDashboardStats(filter)
	if err != nil {
		h.logger.Error("failed to get dashboard stats", "error", err)
		errorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, statsResponse{Data: stats})
}

// @Description Get log volume stats
// @Summary Get log volume statistics for an agent
// @Tags dashboard
// @Accept json
// @Produce json
// @Param agent_id query string false "agent id"
// @Success 200 {array} models.LogVolumeStats
// @Failure 400,401,500 {object} ErrorResponse
// @Router /sentinel/logvolume [get]
// @Security ApiKeyAuth
func (h *Handler) GetLogVolumeStats(ctx *gin.Context) {
	agentID := ctx.Query("agent_id")

	stats, err := h.service.Stats.GetLogVolumeStats(agentID)
	if err != nil {
		h.logger.Error("failed to get log volume stats", "error", err)
		errorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, stats)
}