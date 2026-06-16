package middleware

import (
	"net/http"
	"os"
	"slices"
	"strings"

	"godance/internal/dto"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func abortUnauthorized(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, dto.Error{Error: "unauthorized", Code: "UNAUTHORIZED"})
}

func RequireAuth(c *gin.Context) {
	header := c.GetHeader("Authorization")

	if !strings.HasPrefix(header, "Bearer ") {
		abortUnauthorized(c)
		return
	}

	tokenString := strings.TrimPrefix(header, "Bearer ")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		abortUnauthorized(c)
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		sub, ok := claims["sub"].(float64)
		if !ok {
			abortUnauthorized(c)
			return
		}
		c.Set("userID", uint(sub))
		c.Set("role", claims["role"])
		c.Next()
	} else {
		abortUnauthorized(c)
	}
}

// UserID возвращает ID аутентифицированного пользователя из контекста.
// Работает только после middleware RequireAuth; иначе вернёт 0.
func UserID(c *gin.Context) uint {
	v, ok := c.Get("userID")
	if !ok {
		return 0
	}
	id, _ := v.(uint)
	return id
}

func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _ := c.Get("role")
		roleStr, _ := role.(string)
		if slices.Contains(roles, roleStr) {
			c.Next()
			return
		}
		c.AbortWithStatusJSON(http.StatusForbidden, dto.Error{Error: "forbidden", Code: "FORBIDDEN"})
	}
}
