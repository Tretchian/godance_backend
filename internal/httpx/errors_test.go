package httpx

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"godance/internal/domain"

	"github.com/gin-gonic/gin"
)

func init() { gin.SetMode(gin.TestMode) }

func TestErrorMapsDomainError(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	Error(c, domain.ErrCompetitionNotFound)

	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusNotFound)
	}
	var body map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if body["code"] != "NOT_FOUND" {
		t.Fatalf("code = %q, want NOT_FOUND", body["code"])
	}
	if body["error"] == "" {
		t.Fatal("error message is empty")
	}
}

func TestErrorHidesInternalDetails(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	Error(c, errString("db exploded: secret dsn"))

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 500", w.Code)
	}
	var body map[string]string
	_ = json.Unmarshal(w.Body.Bytes(), &body)
	if body["code"] != "INTERNAL_ERROR" {
		t.Fatalf("code = %q, want INTERNAL_ERROR", body["code"])
	}
	if body["error"] != "internal error" {
		t.Fatalf("error = %q, want generic message (no leak)", body["error"])
	}
}

type errString string

func (e errString) Error() string { return string(e) }
