package handlers

import (
	"errors"
	"net/http"

	"github.com/diyorbek/sentinel/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type metricListResponse struct {
	Data  []models.Metric `json:"data"`
	Total int             `json:"total"`
	Page  int             `json:"page"`
	Limit int             `json:"limit"`
}

// @Description Get metric by id 
// @Summary Get metric by id 
// @Tags logs 
// @Accept json 
// @Produce json
// @Param id path string true "metric id"
// @Success 200 {object} models.Metric
// @Failure 400,401,404,500 {object} ErrorResponse
// @Router /sentinel/api/metrics/{id} [get]
func (h *Handler) GetMetricsByID(ctx *gin.Context) {
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

	metric, err := h.service.Metric.GetMetricsByID(id)
	if err != nil {
		h.logger.Error(err.Error())
		errorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, metric)
}

// @Description Get metrics by filter 
// @Summary Get metric by filter 
// @Tags logs
// @Accept json 
// @Produce json 
// @Param agent_id query string false "agent id"
// @Param from query string false "filter from recorded_at (RFC3339 format)"
// @Param to query string false "Filter to recorded_at (RFC3339 format)"
// @Param limit query int false "Limit" default(10)
// @Param page query int false "Page" default(1)
// @Sucess 200 {object} metricListResponse
// @Failure 400,401,404,500 {object} ErrorResponse
// @Router /sentinel/api/metrics [get]
func (h *Handler) ListMetrics(ctx *gin.Context) {
	var filter models.FilterMetrics

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

	metrics, total, err := h.service.Metric.ListMetrics(filter)
	if err != nil {
		h.logger.Error(err.Error())
		errorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, metricListResponse{
		Data: metrics,
		Total: total,
		Limit: filter.Limit,
		Page: page,
	})
}
