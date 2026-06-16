package main

import (
	"testing"

	"godance/internal/storage"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// TestSetupRoutesNoPanic проверяет, что все группы роутов регистрируются без
// конфликтов (например, несовпадающих имён wildcard на общем префиксе
// /competitions/:id/...). Register-методы хендлеров не обращаются к БД,
// поэтому пустого *gorm.DB достаточно.
func TestSetupRoutesNoPanic(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	setupRoutes(r, &gorm.DB{}, storage.NewStub())
}
