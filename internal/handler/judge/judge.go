package judge

import (
	"errors"
	"net/http"
	"strconv"

	"godance/internal/domain"
	"godance/internal/dto"
	"godance/internal/middleware"
	"godance/internal/service/judge"
	types "godance/internal/type"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *judge.Service
}

func NewHandler(service *judge.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Register(rg *gin.RouterGroup) {
	judges := rg.Group("/competitions/:id/judges")
	{
		judges.GET("", h.GetList)
		judges.POST(
			"",
			middleware.RequireAuth,
			middleware.RequireRole(string(types.UserRoleOrganizer)),
			h.Assign,
		)
	}
}

func (h *Handler) GetList(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid competition id"})
		return
	}
	result, err := h.service.ListForCompetition(uint(id))
	if err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{"error": msg})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *Handler) Assign(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid competition id"})
		return
	}
	userID := middleware.UserID(c)

	var req dto.AssignJudgeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.service.Assign(uint(id), req.JudgeID, userID)
	if err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{"error": msg})
		return
	}
	c.JSON(http.StatusCreated, result)
}

func mapError(err error) (int, string) {
	switch {
	case errors.Is(err, domain.ErrCompetitionNotFound),
		errors.Is(err, domain.ErrUserNotFound):
		return http.StatusNotFound, err.Error()
	case errors.Is(err, domain.ErrNotOwner):
		return http.StatusForbidden, err.Error()
	case errors.Is(err, domain.ErrAlreadyAssigned):
		return http.StatusConflict, err.Error()
	case errors.Is(err, domain.ErrJudgeLimitReached),
		errors.Is(err, domain.ErrCompetitionClosed),
		errors.Is(err, domain.ErrNotJudge):
		return http.StatusBadRequest, err.Error()
	default:
		return http.StatusInternalServerError, "internal error"
	}
}
