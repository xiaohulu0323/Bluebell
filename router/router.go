package router

import (
	"net/http"

	"web-app/logger"
	"web-app/middlewares"

	"github.com/gin-gonic/gin"

	"web-app/controller"
)

func Setup(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	v1 := r.Group("/api/v1")

	// 注册
	v1.POST("/signup", controller.SignUpHandler)
	// 登录
	v1.POST("/login", controller.LoginHandler)

	v1.Use(middlewares.JWTAuthMiddleware()) // 需要登录认证之后才能访问的接口
	
	{
		v1.GET("/community", controller.CommunityHandler)           // 社区列表
		v1.GET("/community/:id", controller.CommunityDetailHandler) // 社区详情

		v1.POST("/post", controller.CreatePostHandler)       // 发帖
		v1.GET("/post/:id", controller.GetPostDetailHandler) // 帖子详情
	}

	
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "404",
		})
	})

	return r
}



