package feedback

import (
	"net/http"
	"strconv"

	"godance/internal/dto"
	"godance/internal/httpx"
	"godance/internal/middleware"

	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateRequest(c *gin.Context) {
	var req dto.CreateFeedbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Bind(c, err)
		return
	}

	result, err := h.service.CreateRequest(middleware.UserID(c), req)
	if err != nil {
		httpx.Error(c, err)
		return
	}
	c.JSON(http.StatusCreated, result)
}

func (h *Handler) GetRequest(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error{Error: "invalid request id", Code: "BAD_REQUEST"})
		return
	}

	result, err := h.service.GetRequest(uint(id), middleware.UserID(c))
	if err != nil {
		httpx.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *Handler) PatchRequest(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error{Error: "invalid request id", Code: "BAD_REQUEST"})
		return
	}

	var req dto.FeedbackRequestAction
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Bind(c, err)
		return
	}

	// Поддерживается единственное действие — confirm (валидируется биндингом).
	result, err := h.service.ConfirmRequest(uint(id), middleware.UserID(c))
	if err != nil {
		httpx.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}
