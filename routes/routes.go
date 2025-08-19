
package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"web-app/logger"
)

func Setup()*gin.Engine{
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	r.GET("/", func(c *gin.Context){
		c.String(http.StatusOK, "Welcome to Bluebell!")
	})

	return r
}