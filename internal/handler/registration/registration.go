package registration

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup) {
	registrations := rg.Group("/registrations")
	{
		registrations.GET("/")
	}
}
