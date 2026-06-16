package competition

import (
	"errors"
	"net/http"
	"strconv"

	"godance/internal/domain"
	"godance/internal/dto"
	"godance/internal/middleware"
	"godance/internal/service/competition"
	types "godance/internal/type"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *competition.Service
}

func NewHandler(service *competition.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Register(rg *gin.RouterGroup) {
	competitions := rg.Group("/competitions")
	{
		competitions.GET("/", h.GetList)
		competitions.GET("/:id", h.GetCompetition)
		competitions.GET(
			"/:id/summary",
			middleware.RequireAuth,
			middleware.RequireRole(string(types.UserRoleOrganizer)),
			h.GetSummary,
		)
		competitions.POST(
			"",
			middleware.RequireAuth,
			middleware.RequireRole(string(types.UserRoleOrganizer)),
			h.CreateCompetition,
		)
		competitions.PATCH(
			"/:id",
			middleware.RequireAuth,
			middleware.RequireRole(string(types.UserRoleOrganizer)),
			h.PatchCompetition,
		)
	}
}

func (h *Handler) GetList(c *gin.Context) {
	status := c.Query("status")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	result, err := h.service.GetPage(status, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *Handler) GetCompetition(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid competition id"})
		return
	}

	result, err := h.service.GetCompetition(uint(id))
	if err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{"error": msg})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *Handler) GetSummary(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid competition id"})
		return
	}

	result, err := h.service.Summary(uint(id), middleware.UserID(c))
	if err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{"error": msg})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *Handler) CreateCompetition(c *gin.Context) {
	var req dto.CreateCompetitionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.service.CreateCompetition(middleware.UserID(c), req)
	if err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{"error": msg})
		return
	}
	c.JSON(http.StatusCreated, result)
}

func (h *Handler) PatchCompetition(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid competition id"})
		return
	}

	var req dto.UpdateCompetitionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.service.PatchCompetition(uint(id), middleware.UserID(c), req)
	if err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{"error": msg})
		return
	}
	c.JSON(http.StatusOK, result)
}

func mapError(err error) (int, string) {
	switch {
	case errors.Is(err, domain.ErrCompetitionNotFound):
		return http.StatusNotFound, err.Error()
	case errors.Is(err, domain.ErrNotOwner):
		return http.StatusForbidden, err.Error()
	default:
		return http.StatusInternalServerError, "internal error"
	}
}
