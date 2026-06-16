package feedback

import (
	"net/http"

	"godance/internal/dto"
	"godance/internal/httpx"
	"godance/internal/middleware"

	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateRating(c *gin.Context) {
	var req dto.CreateFeedbackRating
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Bind(c, err)
		return
	}

	result, err := h.service.CreateRating(middleware.UserID(c), req)
	if err != nil {
		httpx.Error(c, err)
		return
	}
	c.JSON(http.StatusCreated, result)
}
