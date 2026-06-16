package video

import (
	"errors"
	"net/http"
	"strconv"

	"godance/internal/domain"
	"godance/internal/dto"
	"godance/internal/middleware"
	"godance/internal/service/video"
	types "godance/internal/type"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *video.Service
}

func NewHandler(service *video.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Register(rg *gin.RouterGroup) {
	videos := rg.Group("/videos")
	{
		videos.POST("/upload-url",
			middleware.RequireAuth,
			middleware.RequireRole(string(types.UserRoleParticipant)),
			h.CreateUploadURL,
		)
		videos.POST("/:id/confirm",
			middleware.RequireAuth,
			middleware.RequireRole(string(types.UserRoleParticipant)),
			h.Confirm,
		)
		videos.GET("/:id", middleware.RequireAuth, h.GetView)
	}
}

func (h *Handler) CreateUploadURL(c *gin.Context) {
	var req dto.UploadUrlRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	result, err := h.service.CreateUploadURL(middleware.UserID(c), req)
	if err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{"error": msg})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *Handler) Confirm(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid video id"})
		return
	}

	var req dto.VideoConfirmRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.service.ConfirmUpload(middleware.UserID(c), uint(id), req)
	if err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{"error": msg})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *Handler) GetView(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid video id"})
		return
	}

	result, err := h.service.GetViewURL(middleware.UserID(c), uint(id))
	if err != nil {
		status, msg := mapError(err)
		c.JSON(status, gin.H{"error": msg})
		return
	}
	c.JSON(http.StatusOK, result)
}

func mapError(err error) (int, string) {
	switch {
	case errors.Is(err, domain.ErrVideoNotFound),
		errors.Is(err, domain.ErrFeedbackRequestNotFound):
		return http.StatusNotFound, err.Error()
	case errors.Is(err, domain.ErrNotVideoOwner),
		errors.Is(err, domain.ErrNotRequestParticipant),
		errors.Is(err, domain.ErrRequestNotAwaitingVideo):
		return http.StatusForbidden, err.Error()
	case errors.Is(err, domain.ErrVideoGone):
		return http.StatusGone, err.Error()
	case errors.Is(err, domain.ErrVideoObjectMissing),
		errors.Is(err, domain.ErrInvalidStatusTransition):
		return http.StatusBadRequest, err.Error()
	default:
		return http.StatusInternalServerError, "internal error"
	}
}
