package handlers

import (
	"errors"
	"net/http"

	"github.com/diyorbek/sentinel/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type createAccountResponse struct {
	Id uuid.UUID `json:"id"`
}

// @Description Create Account
// @Summary Create Account
// @Tags account
// @Accept json
// @Produce json
// @Param create body models.CreateAccount true "Create Account"
// @Success 201 {object} createAccountResponse
// @Failure 400,500 {object} ErrorResponse
// @Router /sentinel/accounts [post]
// @Security BearerAuth
func (h *Handler) CreateAccount(ctx *gin.Context) {
	var body models.CreateAccount

	if err := ctx.ShouldBindJSON(&body); err != nil {
		h.logger.Warn(err.Error())
		errorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	id, err := h.service.Account.CreateAccount(body)
	if err != nil {
		h.logger.Error(err.Error())
		errorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	h.logger.Info("account created successfully")

	ctx.JSON(http.StatusCreated, createAccountResponse{
		Id: id,
	})
}

// @Description Get Me
// @Summary Get current account
// @Tags account
// @Accept json
// @Produce json
// @Success 200 {object} models.Account
// @Failure 401,500 {object} ErrorResponse
// @Router /sentinel/me [get]
// @Security BearerAuth
func (h *Handler) GetMe(ctx *gin.Context) {
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

	account, err := h.service.Account.GetAccountByID(id)
	if err != nil {
		h.logger.Error(err.Error())
		errorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, account)
}

// @Description Get Account by ID
// @Summary Get Account
// @Tags account
// @Accept json
// @Produce json
// @Param id path string true "account id"
// @Success 200 {object} models.Account
// @Failure 400,404,500 {object} ErrorResponse
// @Router /sentinel/accounts/{id} [get]
// @Security BearerAuth
func (h *Handler) GetAccount(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		errorResponse(ctx, http.StatusBadRequest, errors.New("id is empty"))
		return
	}

	account, err := h.service.Account.GetAccountByID(id)
	if err != nil {
		h.logger.Error(err.Error())
		errorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, account)
}

// @Description Update Account
// @Summary Update Account
// @Tags account
// @Accept json
// @Produce json
// @Param id path string true "account id"
// @Param update body models.UpdateAccount true "update account"
// @Success 200 {object} SuccessResponse
// @Failure 400,500 {object} ErrorResponse
// @Router /sentinel/accounts/{id} [put]
// @Security BearerAuth
func (h *Handler) UpdateAccount(ctx *gin.Context) {
	var req models.UpdateAccount

	idParam := ctx.Param("id")
	if idParam == "" {
		errorResponse(ctx, http.StatusBadRequest, errors.New("id is empty"))
		return
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		errorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	req.Id = idParam

	if err := h.service.Account.UpdateAccount(req); err != nil {
		h.logger.Error(err.Error())
		errorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	ctx.JSON(http.StatusOK, SuccessResponse{
		Message: "Account updated successfully",
	})
}

// @Description Delete Account
// @Summary Delete Account
// @Tags account
// @Accept json
// @Produce json
// @Param id path string true "account id"
// @Success 200 {object} SuccessResponse
// @Failure 400,500 {object} ErrorResponse
// @Router /sentinel/accounts/{id} [delete]
// @Security BearerAuth
func (h *Handler) DeleteAccount(ctx *gin.Context) {
	idParam := ctx.Param("id")
	if idParam == "" {
		errorResponse(ctx, http.StatusBadRequest, errors.New("id is empty"))
		return
	}

	if err := h.service.Account.DeleteAccount(idParam); err != nil {
		h.logger.Error(err.Error())
		errorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, SuccessResponse{
		Message: "Account deleted successfully",
	})
}