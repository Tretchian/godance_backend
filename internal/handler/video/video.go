package video

import (
	"net/http"
	"strconv"

	"godance/internal/dto"
	"godance/internal/httpx"
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
		httpx.Bind(c, err)
		return
	}

	result, err := h.service.CreateUploadURL(middleware.UserID(c), req)
	if err != nil {
		httpx.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *Handler) Confirm(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error{Error: "invalid video id", Code: "BAD_REQUEST"})
		return
	}

	var req dto.VideoConfirmRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Bind(c, err)
		return
	}

	result, err := h.service.ConfirmUpload(middleware.UserID(c), uint(id), req)
	if err != nil {
		httpx.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *Handler) GetView(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error{Error: "invalid video id", Code: "BAD_REQUEST"})
		return
	}

	result, err := h.service.GetViewURL(middleware.UserID(c), uint(id))
	if err != nil {
		httpx.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}
