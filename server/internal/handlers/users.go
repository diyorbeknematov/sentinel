package handlers

import (
	"errors"
	"net/http"

	"github.com/diyorbek/sentinel/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type createUserResponse struct {
	Id uuid.UUID `json:"id"`
}

// @Description Create User
// @Summary Create User
// @Tags user
// @Accept json
// @Produce json
// @Param create body models.CreateUser true "Create User"
// @Success 201 {object} createUserResponse 
// @Failure 400,404,500 {object} ErrorResponse
// @Router /sentinal/api/users [post]
func (h *Handler) CreateUser(ctx *gin.Context) {
	var body models.CreateUser
	if err := ctx.ShouldBindJSON(&body); err != nil {
		h.logger.Warn(err.Error())
		errorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	id, err := h.service.User.CreateUser(body)
	if err != nil {
		h.logger.Error(err.Error())
		errorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	h.logger.Info("created user sucessfully (create user)")
	ctx.JSON(http.StatusCreated, createUserResponse{
		Id: id,
	})
}

// @Description Get User 
// @Summary Get User 
// @Tags user 
// @Accept json 
// @Produce json 
// @Param id path string ture "user id"
// @Success 200 {object} models.User
// @Failure 400,401,404,500 {object} ErrorResponse
// @Router /sentinal/api/users/{id} [get] 
func (h *Handler) GetUser(ctx *gin.Context) {
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

	user, err := h.service.User.GetUser(id)
	if err != nil {
		h.logger.Error(err.Error())
		errorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	h.logger.Info("query is done sucessfully (get user)")
	ctx.JSON(http.StatusOK, user)
}

// @Description Update User 
// @Summary Update User
// @Tags user 
// @Accept json 
// @Produce json 
// @Param id path string true "user id"
// @Param update body models.UpdateUser true "update user"
// @Success 200 {object} SuccessResponse
// @Failure 400,401,404,500 {object} ErrorResponse
// @Router /sentinal/api/users/{id} [put]
func (h *Handler) UpdateUser(ctx *gin.Context) {
	var req models.UpdateUser

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

	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.logger.Warn(err.Error())
		errorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	req.Id = id

	if err := h.service.User.UpdateUser(req); err != nil {
		h.logger.Error(err.Error())
		errorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	h.logger.Info("user updated succussfully (update user)")
	ctx.JSON(http.StatusOK, SuccessResponse{
		Message: "user updatated successly",
	})
}

// @Description Update User Role
// @Summary Update User Role
// @Tags user 
// @Accept json 
// @Produce json 
// @Param id path string true "user id"
// @Param update body models.UpdateRole true "update user role"
// @Success 200 {object} SuccessResponse
// @Failure 400,401,404,500 {object} ErrorResponse
// @Router /sentinal/api/users/{id}/role [put]
func (h *Handler) UpdateUserRole(ctx *gin.Context) {
	var body models.UpdateRole

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

	if err := ctx.ShouldBindJSON(&body); err != nil {
		h.logger.Warn(err.Error())
		errorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	body.Id = id
	err = h.service.User.UpdateUserRole(body)
	if err != nil {
		h.logger.Error(err.Error())
		errorResponse(ctx, http.StatusInternalServerError, err)
		return
	}
}

// @Description Delete User
// @Summary Delete User 
// @Tags user 
// @Accept json 
// @Produce json 
// @Param id path string ture "user id"
// @Success 200 {object} SuccessResponse
// @Failure 400,401,404,500 {object} ErrorResponse
// @Router /sentinal/api/users/{id} [delete]
func (h *Handler) DeleteUser(ctx *gin.Context) {
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

	err = h.service.User.DeleteUser(id)
	if err != nil {
		h.logger.Error(err.Error())
		errorResponse(ctx, http.StatusInternalServerError, err)
	}

	ctx.JSON(http.StatusOK, SuccessResponse{
		Message: "User deleted successfully",
	})
}
