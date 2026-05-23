package feedback

import (
	"godance/internal/service/feedback"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *feedback.Service
}

func NewHandler(service *feedback.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Register(r *gin.RouterGroup) {
	g := r.Group("/feedback")

	ratings := g.Group("/ratings")
	ratings.POST("", h.CreateRating)

	requests := g.Group("/requests")
	requests.POST("", h.CreateRequest)
	requests.GET("/:id", h.GetRequest)
	requests.PATCH("/:id", h.PatchRequest)

	responses := g.Group("/responses")
	responses.POST("", h.CreateResponse)
	responses.GET("/:id", h.GetResponse)
}
