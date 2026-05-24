package middleware

import (
	"net/http"
	"os"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func RequireAuth(c *gin.Context) {
	header := c.GetHeader("Authorization")

	if !strings.HasPrefix(header, "Bearer ") {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	tokenString := strings.TrimPrefix(header, "Bearer ")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {

		c.Set("userID", claims["sub"])
		c.Set("role", claims["role"])
		c.Next()
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
}

func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.MustGet("role").(string)
		if slices.Contains(roles, role) {
			c.Next()
			return
		}
		c.AbortWithStatus(http.StatusForbidden)
	}
}
