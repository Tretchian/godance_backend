package payment

import (
	"net/http"

	"godance/internal/dto"
	"godance/internal/httpx"
	"godance/internal/service/feedback"
	"godance/internal/service/registration"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	feedback     *feedback.Service
	registration *registration.Service
}

func NewHandler(feedbackService *feedback.Service, registrationService *registration.Service) *Handler {
	return &Handler{feedback: feedbackService, registration: registrationService}
}

func (h *Handler) Register(rg *gin.RouterGroup) {
	webhooks := rg.Group("/payments/webhooks")
	{
		webhooks.POST("/feedback", h.FeedbackWebhook)
		webhooks.POST("/registration", h.RegistrationWebhook)
	}
}

// FeedbackWebhook принимает события платёжного шлюза по запросам ОС
// (held / captured / released / failed) и драйвит статусы запроса.
func (h *Handler) FeedbackWebhook(c *gin.Context) {
	var payload dto.PaymentWebhook
	if err := c.ShouldBindJSON(&payload); err != nil {
		httpx.Bind(c, err)
		return
	}

	if err := h.feedback.HandleFeedbackWebhook(payload); err != nil {
		httpx.Error(c, err)
		return
	}
	c.Status(http.StatusOK)
}

// RegistrationWebhook подтверждает оплату взноса участника.
func (h *Handler) RegistrationWebhook(c *gin.Context) {
	var payload dto.PaymentWebhook
	if err := c.ShouldBindJSON(&payload); err != nil {
		httpx.Bind(c, err)
		return
	}

	if err := h.registration.HandleRegistrationWebhook(payload); err != nil {
		httpx.Error(c, err)
		return
	}
	c.Status(http.StatusOK)
}
