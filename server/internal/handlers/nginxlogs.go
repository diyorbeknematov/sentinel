package handlers

import (
	"errors"
	"net/http"

	"github.com/diyorbek/sentinel/internal/models"
	"github.com/gin-gonic/gin"
)

type nginxListResponse struct {
	Data  []models.NginxLogResponse `json:"data"`
	Total int                       `json:"total"`
	Page  int                       `json:"page"`
	Limit int                       `json:"limit"`
}

// @Description Get nginx log by id
// @Summary Get nginx log by id
// @Tags logs
// @Accept json
// @Produce json
// @Param id path string ture "nginx log id"
// @Success 200 {object} models.NginxLog
// @Failure 400,401,404,500 {object} ErrorResponse
// @Router /sentinel/nginxlogs/{id} [get]
// @Security BearerAuth
func (h *Handler) GetNginxLogByID(ctx *gin.Context) {
	// get id from path param
	id := ctx.Param("id")
	if id == "" {
		h.logger.Warn("empty param value")
		errorResponse(ctx, http.StatusBadRequest, errors.New("param value is emtpy"))
		return
	}

	nginxLog, err := h.service.NginxLog.GetNginxLogByID(id)
	if err != nil {
		h.logger.Error(err.Error())
		errorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, nginxLog)
}

// @Description get nginxlogs by filter
// @Summary get nginxlogs by filter
// @Tags logs
// @Accept json
// @Produce json
// @Param agent_id query string false "agent id"
// @Pram method query string false "method"
// @Param status query string false "status"
// @Param from query string false "filter from recorded_at (RFC3339 format)"
// @Param to query string false "Filter to recorded_at (RFC3339 format)"
// @Param limit query int false "Limit" default(10)
// @Param page query int false "Page" default(1)
// @Sucess 200 {object} nginxListResponse
// @Failure 400,401,404,500 {object} ErrorResponse
// @Router /sentinel/nginxlogs [get]
// @Security BearerAuth
func (h *Handler) ListNginxLogs(ctx *gin.Context) {
	var filter models.FilterNginxLog

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

	nginxLogs, total, err := h.service.NginxLog.ListNginxLogs(filter)
	if err != nil {
		h.logger.Error(err.Error())
		errorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, nginxListResponse{
		Data:  nginxLogs,
		Total: total,
		Limit: filter.Limit,
		Page:  page,
	})
}
