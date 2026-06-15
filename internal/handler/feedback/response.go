package feedback

import (
	"net/http"
	"strconv"

	"godance/internal/dto"
	"godance/internal/middleware"

	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateResponse(c *gin.Context) {
	var req dto.CreateFeedbackResponse
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.service.CreateResponse(middleware.UserID(c), req)
	if err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{"error": msg})
		return
	}
	c.JSON(http.StatusCreated, result)
}

func (h *Handler) GetResponse(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid response id"})
		return
	}

	result, err := h.service.GetResponse(uint(id), middleware.UserID(c))
	if err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{"error": msg})
		return
	}
	c.JSON(http.StatusOK, result)
}
