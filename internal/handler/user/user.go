package user

import (
	"godance/internal/service/user"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *user.Service
}

func NewHandler(service *user.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Register(rg *gin.RouterGroup) {
}
