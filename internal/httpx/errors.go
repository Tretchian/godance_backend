// Package httpx централизует преобразование доменных ошибок в HTTP-ответы
// единого формата ({error, code} / ValidationError) согласно OpenAPI.
package httpx

import (
	"errors"
	"net/http"
	"reflect"
	"strings"

	"godance/internal/domain"
	"godance/internal/dto"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Error пишет доменную ошибку в ответ с подходящими статусом и кодом.
func Error(c *gin.Context, err error) {
	status, code := classify(err)
	msg := err.Error()
	if status == http.StatusInternalServerError {
		msg = "internal error"
	}
	c.JSON(status, dto.Error{Error: msg, Code: code})
}

// Bind пишет ошибку разбора/валидации тела запроса: 422 ValidationError для
// ошибок валидатора (с полями), иначе 400 BAD_REQUEST (например, битый JSON).
func Bind(c *gin.Context, err error) {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		details := make([]dto.FieldError, 0, len(ve))
		for _, fe := range ve {
			details = append(details, dto.FieldError{
				Field:   fe.Field(),
				Message: validationMessage(fe),
			})
		}
		c.JSON(http.StatusUnprocessableEntity, dto.ValidationError{
			Error:   "Ошибка валидации",
			Code:    "VALIDATION_ERROR",
			Details: details,
		})
		return
	}
	c.JSON(http.StatusBadRequest, dto.Error{Error: err.Error(), Code: "BAD_REQUEST"})
}

func classify(err error) (int, string) {
	switch {
	case errors.Is(err, domain.ErrCompetitionNotFound),
		errors.Is(err, domain.ErrUserNotFound),
		errors.Is(err, domain.ErrFeedbackRequestNotFound),
		errors.Is(err, domain.ErrFeedbackResponseNotFound),
		errors.Is(err, domain.ErrVideoNotFound),
		errors.Is(err, domain.ErrAssignmentNotFound),
		errors.Is(err, domain.ErrRegistrationNotFound):
		return http.StatusNotFound, "NOT_FOUND"

	case errors.Is(err, domain.ErrNotOwner),
		errors.Is(err, domain.ErrNotRequestParticipant),
		errors.Is(err, domain.ErrNotRequestJudge),
		errors.Is(err, domain.ErrNotVideoOwner),
		errors.Is(err, domain.ErrRequestNotAwaitingVideo):
		return http.StatusForbidden, "FORBIDDEN"

	case errors.Is(err, domain.ErrAlreadyAssigned),
		errors.Is(err, domain.ErrAlreadyRegistered),
		errors.Is(err, domain.ErrResponseAlreadyExists),
		errors.Is(err, domain.ErrRatingAlreadyExists),
		errors.Is(err, domain.ErrJudgeHasRequests),
		errors.Is(err, domain.ErrEmailTaken):
		return http.StatusConflict, "CONFLICT"

	case errors.Is(err, domain.ErrJudgeLimitReached),
		errors.Is(err, domain.ErrCompetitionClosed),
		errors.Is(err, domain.ErrNotJudge),
		errors.Is(err, domain.ErrRegistrationClosed),
		errors.Is(err, domain.ErrParticipantLimitReached),
		errors.Is(err, domain.ErrJudgeNotAssigned),
		errors.Is(err, domain.ErrInvalidStatusTransition),
		errors.Is(err, domain.ErrVideoObjectMissing):
		return http.StatusBadRequest, "BAD_REQUEST"

	case errors.Is(err, domain.ErrInvalidToken),
		errors.Is(err, domain.ErrInvalidCredentials):
		return http.StatusUnauthorized, "UNAUTHORIZED"

	case errors.Is(err, domain.ErrVideoGone):
		return http.StatusGone, "GONE"

	default:
		return http.StatusInternalServerError, "INTERNAL_ERROR"
	}
}

func validationMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "обязательное поле"
	case "email":
		return "некорректный email"
	case "min":
		return "значение меньше минимального (" + fe.Param() + ")"
	case "max":
		return "значение больше максимального (" + fe.Param() + ")"
	case "oneof":
		return "недопустимое значение (ожидается одно из: " + fe.Param() + ")"
	default:
		return "некорректное значение"
	}
}

// RegisterValidatorTagName заставляет валидатор использовать имена полей из
// json-тегов, чтобы поле в ValidationError совпадало с телом запроса.
func RegisterValidatorTagName() {
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		return
	}
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}
