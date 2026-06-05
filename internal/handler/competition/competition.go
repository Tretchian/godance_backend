package competition

import (
	"net/http"
	"strconv"

	"godance/internal/dto"
	"godance/internal/middleware"
	"godance/internal/models"
	"godance/internal/service/competition"
	types "godance/internal/type"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
		competitions.GET(
			"/:id/summary",
			middleware.RequireAuth,
			middleware.RequireRole(string(types.UserRoleOrganizer)),
			h.GetCompetition,
		)
		competitions.POST(
			"",
			middleware.RequireAuth,
			middleware.RequireRole("organizer"),
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

func (h *Handler) CreateCompetition(c *gin.Context) {
	user := c.MustGet("userID").(models.User)
	var req dto.CreateCompetitionRequest
	c.ShouldBindJSON(&req)

	result, err := h.service.CreateCompetition(user, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}
	c.JSON(http.StatusOK, result)
}

func (h *Handler) PatchCompetition(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)

	var req dto.UpdateCompetitionRequest
	c.ShouldBindJSON(&req)

	result, err := h.service.PatchCompetition(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, result)
}

func (h *Handler) GetCompetition(c *gin.Context) {
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
