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
		judges.DELETE(
			"/:judge_id",
			middleware.RequireAuth,
			middleware.RequireRole(string(types.UserRoleOrganizer)),
			h.Unassign,
		)
	}

	catalog := rg.Group("/judges")
	{
		catalog.GET("", h.GetCatalog)
		catalog.GET("/:id", h.GetProfile)
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

func (h *Handler) Unassign(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid competition id"})
		return
	}
	judgeID, err := strconv.ParseUint(c.Param("judge_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid judge id"})
		return
	}

	if err := h.service.Unassign(uint(id), uint(judgeID), middleware.UserID(c)); err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{"error": msg})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) GetCatalog(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	result, err := h.service.ListCatalog(page, limit)
	if err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{"error": msg})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *Handler) GetProfile(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid judge id"})
		return
	}

	result, err := h.service.GetProfile(uint(id))
	if err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{"error": msg})
		return
	}
	c.JSON(http.StatusOK, result)
}

func mapError(err error) (int, string) {
	switch {
	case errors.Is(err, domain.ErrCompetitionNotFound),
		errors.Is(err, domain.ErrUserNotFound),
		errors.Is(err, domain.ErrAssignmentNotFound):
		return http.StatusNotFound, err.Error()
	case errors.Is(err, domain.ErrNotOwner):
		return http.StatusForbidden, err.Error()
	case errors.Is(err, domain.ErrAlreadyAssigned),
		errors.Is(err, domain.ErrJudgeHasRequests):
		return http.StatusConflict, err.Error()
	case errors.Is(err, domain.ErrJudgeLimitReached),
		errors.Is(err, domain.ErrCompetitionClosed),
		errors.Is(err, domain.ErrNotJudge):
		return http.StatusBadRequest, err.Error()
	default:
		return http.StatusInternalServerError, "internal error"
	}
}
