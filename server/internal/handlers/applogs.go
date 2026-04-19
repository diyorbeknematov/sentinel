package handlers

import (
	"errors"
	"net/http"

	"github.com/diyorbek/sentinel/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type appLogListResponse struct {
	Data  []models.Log `json:"data"`
	Total int          `json:"total"`
	Page  int          `json:"page"`
	Limit int          `json:"limit"`
}

// @Description Get applog by id
// @Summary Get applog by id
// @Tags logs
// @Accept json 
// @Produce json 
// @Param id path string true "app log id"
// @Success 200 {object} models.Log
// @Failure 400,401,404,500 {object} ErrorResponse
// @Router /sentinal/api/applogs/{id} [get]
func (h *Handler) GetAppLogByID(ctx *gin.Context) {
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

	appLog, err := h.service.AppLog.GetLogByID(id)
	if err != nil {
		h.logger.Error(err.Error())
		errorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, appLog)
}

// @Description Get app logs by filter 
// @Summary Get app logs by filter 
// @Tags logs 
// Accept json
// @Produce json 
// @Param agent_id query string false "agent id"
// @Param user_id query string false "user id"
// @Param type query string false "type"
// @Param level query string false "level"
// @Param from query string false "filter from recorded_at (RFC3339 format)"
// @Param to query string false "Filter to recorded_at (RFC3339 format)"
// @Param limit query int false "Limit" default(10)
// @Param page query int false "Page" default(1)
// @Sucess 200 {object} appLogListResponse
// @Failure 400,401,404,500 {object} ErrorResponse
// @Router /sentinal/api/applogs [get]
func (h *Handler) GetListAppLog(ctx *gin.Context) {
	var filter models.FilterAppLog

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

	appLogs, total, err := h.service.AppLog.ListLogs(filter)
	if err != nil {
		h.logger.Error(err.Error())
		errorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, appLogListResponse{
		Data:  appLogs,
		Total: total,
		Limit: filter.Limit,
		Page:  page,
	})
}
