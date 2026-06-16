package judge

import (
	"net/http"
	"strconv"

	"godance/internal/dto"
	"godance/internal/httpx"
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
		c.JSON(http.StatusBadRequest, dto.Error{Error: "invalid competition id", Code: "BAD_REQUEST"})
		return
	}
	result, err := h.service.ListForCompetition(uint(id))
	if err != nil {
		httpx.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *Handler) Assign(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error{Error: "invalid competition id", Code: "BAD_REQUEST"})
		return
	}

	var req dto.AssignJudgeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Bind(c, err)
		return
	}

	result, err := h.service.Assign(uint(id), req.JudgeID, middleware.UserID(c))
	if err != nil {
		httpx.Error(c, err)
		return
	}
	c.JSON(http.StatusCreated, result)
}

func (h *Handler) Unassign(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error{Error: "invalid competition id", Code: "BAD_REQUEST"})
		return
	}
	judgeID, err := strconv.ParseUint(c.Param("judge_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error{Error: "invalid judge id", Code: "BAD_REQUEST"})
		return
	}

	if err := h.service.Unassign(uint(id), uint(judgeID), middleware.UserID(c)); err != nil {
		httpx.Error(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) GetCatalog(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	result, err := h.service.ListCatalog(page, limit)
	if err != nil {
		httpx.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *Handler) GetProfile(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error{Error: "invalid judge id", Code: "BAD_REQUEST"})
		return
	}

	result, err := h.service.GetProfile(uint(id))
	if err != nil {
		httpx.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}
