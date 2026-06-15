package feedback

import (
	"errors"
	"net/http"

	"godance/internal/domain"
	"godance/internal/middleware"
	"godance/internal/service/feedback"
	types "godance/internal/type"

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

	requests := g.Group("/requests")
	requests.POST("",
		middleware.RequireAuth,
		middleware.RequireRole(string(types.UserRoleParticipant)),
		h.CreateRequest,
	)
	requests.GET("/:id", middleware.RequireAuth, h.GetRequest)
	requests.PATCH("/:id",
		middleware.RequireAuth,
		middleware.RequireRole(string(types.UserRoleParticipant)),
		h.PatchRequest,
	)

	responses := g.Group("/responses")
	responses.POST("",
		middleware.RequireAuth,
		middleware.RequireRole(string(types.UserRoleJudge)),
		h.CreateResponse,
	)
	responses.GET("/:id", middleware.RequireAuth, h.GetResponse)

	ratings := g.Group("/ratings")
	ratings.POST("",
		middleware.RequireAuth,
		middleware.RequireRole(string(types.UserRoleParticipant)),
		h.CreateRating,
	)
}

func mapError(err error) (int, string) {
	switch {
	case errors.Is(err, domain.ErrCompetitionNotFound),
		errors.Is(err, domain.ErrFeedbackRequestNotFound),
		errors.Is(err, domain.ErrFeedbackResponseNotFound):
		return http.StatusNotFound, err.Error()
	case errors.Is(err, domain.ErrNotRequestParticipant),
		errors.Is(err, domain.ErrNotRequestJudge):
		return http.StatusForbidden, err.Error()
	case errors.Is(err, domain.ErrResponseAlreadyExists),
		errors.Is(err, domain.ErrRatingAlreadyExists):
		return http.StatusConflict, err.Error()
	case errors.Is(err, domain.ErrJudgeNotAssigned),
		errors.Is(err, domain.ErrInvalidStatusTransition):
		return http.StatusBadRequest, err.Error()
	default:
		return http.StatusInternalServerError, "internal error"
	}
}
