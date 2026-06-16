package auth

import (
	"net/http"

	"godance/internal/dto"
	"godance/internal/httpx"
	"godance/internal/service/auth"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *auth.Service
}

func NewHandler(service *auth.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Register(rg *gin.RouterGroup) {
	auth := rg.Group("/auth")
	{
		auth.POST("/register", h.CreateUser)
		auth.POST("/login", h.Login)
		auth.POST("/refresh", h.Refresh)
	}
}

func (h *Handler) CreateUser(c *gin.Context) {
	var req dto.RegistrationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Bind(c, err)
		return
	}

	result, err := h.service.CreateUser(req)
	if err != nil {
		httpx.Error(c, err)
		return
	}
	c.JSON(http.StatusCreated, result)
}

func (h *Handler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Bind(c, err)
		return
	}

	result, err := h.service.Login(req)
	if err != nil {
		httpx.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *Handler) Refresh(c *gin.Context) {
	var req dto.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Bind(c, err)
		return
	}

	result, err := h.service.Refresh(req.RefreshToken)
	if err != nil {
		httpx.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}
