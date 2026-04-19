package handlers

import (
	"errors"
	"net/http"

	"github.com/diyorbek/sentinel/internal/models"
	"github.com/gin-gonic/gin"
)

type authResponse struct {
	Token string `json:"token"`
}

// @Description Register User
// @Summary Register User
// @Tags Auth
// @Accept json
// @Produce json
// @Param signup body models.Register true "Register"
// @Success 200 {object} authResponse
// @Failure 400,404,500 {object} ErrorResponse
// @Router /sentinel/register [post]
func (h *Handler) register(c *gin.Context) {
	var body models.Register

	if err := c.ShouldBindJSON(&body); err != nil {
		h.logger.Warn(err.Error())
		errorResponse(c, http.StatusBadRequest, err)
		return
	}

	if body.UserName == "" || body.Password == "" {
		h.logger.Warn("username or pasword cannot be empty")
		errorResponse(c, http.StatusBadRequest, errors.New("username or pasword cannot be empty"))
		return
	}

	accessToken, _, err := h.service.Authorization.Register(body)
	if err != nil {
		h.logger.Error(err.Error())
		errorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, authResponse{
		Token: accessToken.Token,
	})
}

// @Description Login User
// @Summary Login User
// @Tags Auth
// @Accept json
// @Produce json
// @Param login body models.Login true "Login"
// @Success 200 {object} authResponse
// @Failure 400,404,500 {object} ErrorResponse
// @Router /sentinel/login [post]
func (h *Handler) login(c *gin.Context) {
	var body models.Login

	if err := c.ShouldBindJSON(&body); err != nil {
		h.logger.Warn(err.Error())
		errorResponse(c, http.StatusBadRequest, err)
		return
	}

	if body.UserName == "" || body.Password == "" {
		h.logger.Warn("error: Username and password are required")
		errorResponse(c, http.StatusBadRequest, errors.New("error: Username and password are required"))
		return
	}

	accessToken, _, err := h.service.Authorization.Login(body)
	if err != nil {
		h.logger.Error(err.Error())
		errorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, authResponse{
		Token: accessToken.Token,
	})
}
