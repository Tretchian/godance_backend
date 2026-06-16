package auth

import (
	"godance/internal/dto"
	"godance/internal/service/auth"
	"net/http"

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

func (h *Handler) Refresh(c *gin.Context) {
	var req dto.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.service.Refresh(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *Handler) CreateUser(c *gin.Context) {
	var req dto.RegistrationRequest
	c.ShouldBindJSON(&req)

	result, err := h.service.CreateUser(req)

	if err != nil {
		c.JSON(http.StatusBadRequest, err)
	}
	c.JSON(http.StatusOK, result)
}

func (h *Handler) Login(c *gin.Context) {
	var req dto.LoginRequest
	c.ShouldBindJSON(&req)

	result, err := h.service.Login(req)

	if err != nil {
		c.JSON(http.StatusBadRequest, err)
	}
	c.JSON(http.StatusOK, result)
}
