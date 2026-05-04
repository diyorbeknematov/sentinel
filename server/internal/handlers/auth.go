package handlers

import (
	"errors"
	"fmt"
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

	if body.Username == "" || body.Password == "" || body.Email == "" {
		h.logger.Warn("username, pasword or email cannot be empty")
		errorResponse(c, http.StatusBadRequest, errors.New("username, pasword or email cannot be empty"))
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

	if body.Username == "" || body.Password == "" {
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
	fmt.Println("accessToken:", accessToken.Token, body)
	c.JSON(http.StatusOK, authResponse{
		Token: accessToken.Token,
	})
}

// @Description Forgot Password
// @Summary Forgot Password
// @Tags Auth
// @Accept json
// @Produce json
// @Param email body models.ForgotPasswordRequest true "Email"
// @Success 200 {object} SuccessResponse
// @Failure 400,404,500 {object} ErrorResponse
// @Router /sentinel/forgot-password [post]
func (h *Handler) forgotPassword(c *gin.Context) {
	var req models.ForgotPasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn(err.Error())
		errorResponse(c, http.StatusBadRequest, err)
		return
	}

	err := h.service.Authorization.ForgotPassword(req.Email)
	if err != nil {
		h.logger.Error(err.Error())
		errorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Reset password link sent to your email",
	})
}

// @Description Reset Password
// @Summary Reset Password
// @Tags Auth
// @Accept json
// @Produce json
// @Param email body models.ResetPasswordRequest true "Email"
// @Success 200 {object} SuccessResponse
// @Failure 400,404,500 {object} ErrorResponse
// @Router /sentinel/reset-password [post]
func (h *Handler) resetPassword(c *gin.Context) {
	var req models.ResetPasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn(err.Error())
		errorResponse(c, http.StatusBadRequest, err)
		return
	}

	err := h.service.Authorization.ResetPassword(req.Token, req.Password)
	if err != nil {
		h.logger.Error(err.Error())
		errorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Password successfully updated",
	})
}
