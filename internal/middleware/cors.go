package middleware

import (
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORS возвращает CORS-middleware с белым списком origin-ов плюс превью-деплои
// *.vercel.app. Credentials не используются: фронт шлёт Bearer-токен в заголовке
// Authorization, а не cookie.
func CORS(allowedOrigins []string) gin.HandlerFunc {
	allowed := make(map[string]bool, len(allowedOrigins))
	for _, o := range allowedOrigins {
		if o = strings.TrimSpace(o); o != "" {
			allowed[o] = true
		}
	}

	return cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
		AllowOriginFunc: func(origin string) bool {
			if allowed[origin] {
				return true
			}
			return strings.HasSuffix(origin, ".vercel.app")
		},
	})
}
