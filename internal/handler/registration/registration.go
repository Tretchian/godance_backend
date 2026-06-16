package registration

import (
	"net/http"
	"strconv"

	"godance/internal/dto"
	"godance/internal/httpx"
	"godance/internal/middleware"
	"godance/internal/service/registration"
	types "godance/internal/type"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *registration.Service
}

func NewHandler(service *registration.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Register(rg *gin.RouterGroup) {
	registrations := rg.Group("/competitions/:id/registrations")
	{
		registrations.GET(
			"",
			middleware.RequireAuth,
			middleware.RequireRole(string(types.UserRoleOrganizer)),
			h.GetList,
		)
		registrations.POST(
			"",
			middleware.RequireAuth,
			middleware.RequireRole(string(types.UserRoleParticipant)),
			h.Create,
		)
	}
}

func (h *Handler) GetList(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error{Error: "invalid competition id", Code: "BAD_REQUEST"})
		return
	}

	page := clampPage(c.DefaultQuery("page", "1"))
	limit := clampLimit(c.DefaultQuery("limit", "20"))

	result, err := h.service.ListForCompetition(uint(id), middleware.UserID(c), page, limit)
	if err != nil {
		httpx.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *Handler) Create(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error{Error: "invalid competition id", Code: "BAD_REQUEST"})
		return
	}

	result, err := h.service.Register(uint(id), middleware.UserID(c))
	if err != nil {
		httpx.Error(c, err)
		return
	}
	c.JSON(http.StatusCreated, result)
}

func clampPage(raw string) int {
	page, err := strconv.Atoi(raw)
	if err != nil || page < 1 {
		return 1
	}
	return page
}

func clampLimit(raw string) int {
	limit, err := strconv.Atoi(raw)
	if err != nil || limit < 1 || limit > 100 {
		return 20
	}
	return limit
}
