package handlers

import (
	"errors"
	"net/http"

	"github.com/diyorbek/sentinel/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type alertListResponse struct {
	Data  []models.Alert `json:"data"`
	Total int            `json:"total"`
	Page  int            `json:"page"`
	Limit int            `json:"limit"`
}

// @Description Get Alert by id 
// @Summary Get Alert by id 
// @Tags alerts
// @Accept json 
// @Produce json 
// @Param id path string true "alert id"
// @Sucess 200 {object} models.Alert
// @Failure 400,401,404,500 {object} ErrorResponse
// @Router /sentinal/api/alerts/{id} [get]
func (h *Handler) GetAlertByID(ctx *gin.Context) {
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

	alert, err := h.service.Alert.GetAlertByID(id)
	if err != nil {
		h.logger.Error(err.Error())
		errorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, alert)
}

// @Description Get alerts by filter 
// @Summary Get alerts by filter 
// @Tags alerts 
// @Accept json 
// @Produce json 
// @Param agent_id query string false "agent id"
// @Pram type query string false "type"
// @Param severity query string false "severity"
// @Param is_read query string false "is read"
// @Param from query string false "filter from created_at (RFC3339 format)"
// @Param to query string false "Filter to created_at (RFC3339 format)"
// @Param limit query int false "Limit" default(10)
// @Param page query int false "Page" default(1)
// @Sucess 200 {object} nginxListResponse
// @Failure 400,401,404,500 {object} ErrorResponse
// @Router /sentinal/api/alerts [get]
func (h *Handler) ListAlerts(ctx *gin.Context) {
	var filter models.FilterAlert

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

	alerts, total, err := h.service.Alert.ListAlerts(filter)
	if err != nil {
		h.logger.Error(err.Error())
		errorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, alertListResponse{
		Data:  alerts,
		Total: total,
		Limit: filter.Limit,
		Page:  page,
	})
}

// @Description Mark alert read
// @Summary Mark alert read 
// @Tags alerts 
// @Accept json 
// @Produce json 
// @Param id path string true "alert id"
// @Param mark_read body models.MarkAlertRead true "mark alert reaad"
// @Success 200 {object} SuccessResponse
// @Failure 400,401,404,500 {object} ErrorResponse
// @Router /sentinal/api/alerts/{id}/markread [put]
func (h *Handler) MarkAlertRead(ctx *gin.Context) {
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

	var body models.MarkAlertRead
	if err = ctx.ShouldBindJSON(&body); err != nil {
		h.logger.Warn(err.Error())
		errorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	body.Id = id

	err = h.service.Alert.MarkAlertRead(body)
	if err != nil {
		h.logger.Error(err.Error())
		errorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, SuccessResponse{
		Message: "updated laret read successfully",
	})
}
