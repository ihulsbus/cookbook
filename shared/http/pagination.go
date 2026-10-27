package http

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	m "github.com/ihulsbus/cookbook/shared/models"
)

func ParsePagination(ctx *gin.Context) (m.PaginationRequest, error) {
	page, err := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if err != nil {
		return m.PaginationRequest{}, errors.New("page must be a valid integer")
	}
	limit, err := strconv.Atoi(ctx.DefaultQuery("limit", "25"))
	if err != nil {
		return m.PaginationRequest{}, errors.New("limit must be a valid integer")
	}

	return m.NormalizePagination(page, limit), nil
}
