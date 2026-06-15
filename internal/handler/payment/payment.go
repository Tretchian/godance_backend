package payment

import (
	"net/http"

	"godance/internal/dto"
	"godance/internal/service/feedback"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	feedback *feedback.Service
}

func NewHandler(feedbackService *feedback.Service) *Handler {
	return &Handler{feedback: feedbackService}
}

func (h *Handler) Register(rg *gin.RouterGroup) {
	webhooks := rg.Group("/payments/webhooks")
	{
		webhooks.POST("/feedback", h.FeedbackWebhook)
	}
}

// FeedbackWebhook принимает события платёжного шлюза по запросам ОС
// (held / captured / released / failed) и драйвит статусы запроса.
func (h *Handler) FeedbackWebhook(c *gin.Context) {
	var payload dto.PaymentWebhook
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.feedback.HandleFeedbackWebhook(payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}
