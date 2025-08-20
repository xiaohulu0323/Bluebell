package router

import (
	"net/http"

	"web-app/logger"

	"github.com/gin-gonic/gin"

	"web-app/controller"
)

func Setup() *gin.Engine {
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	r.POST("/signup", controller.SignUpHandler)

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Welcome to Bluebell!")
	})

	return r
}
