package feedback

import (
	"github.com/gin-gonic/gin"
	"godance/internal/dto"
)

func (h *Handler) CreateRating(c *gin.Context) {
	var req dto.CreateFeedbackRating
	c.ShouldBindJSON(&req)
}
