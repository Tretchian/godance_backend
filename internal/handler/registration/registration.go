package registration

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup) {
	users := rg.Group("/registrations")
	{
		users.GET("/")
	}
}
