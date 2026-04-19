package handlers

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SuccessResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Message string `json:"error"`
}

func errorResponse(c *gin.Context, status int, err error) {
	c.JSON(status, ErrorResponse{
		Message: err.Error(),
	})
}

func getPageQuery(ctx *gin.Context) (int, error) {
	pageStr := ctx.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		return 1, errors.New("invalid page parameter")
	}
	return page, nil
}
