package router

// 包含 API v1 路由
// @tag.name 用户
// @tag.description 用户相关接口
// @tag.name 社区
// @tag.description 社区相关接口
// @tag.name 帖子
// @tag.description 帖子相关接口
// @tag.name 投票
// @tag.description 投票相关接口
import (
	"net/http"
	"time"

	"web-app/logger"
	"web-app/middlewares"

	"github.com/gin-gonic/gin"
	// swagger
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-contrib/pprof"

	"web-app/controller"
)

func Setup(mode string) *gin.Engine {
	if mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode) // 开发模式
	}

	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true)) // 放在这里就是全网站限速

	// Swagger 文档路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.LoadHTMLFiles("templates/index.html")
	r.Static("/static", "./static")

	// 访问首页
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	v1 := r.Group("/api/v1")

	// 注册
	v1.POST("/signup", controller.SignUpHandler)
	// 登录
	v1.POST("/login", controller.LoginHandler)

	v1.GET("/posts", controller.GetPostListHandler)      // 帖子列表（分页）
	v1.GET("/posts2", controller.GetPostListHandler2) // 帖子列表（分页）
	v1.GET("/community", controller.CommunityHandler)           // 社区列表
	v1.GET("/community/:id", controller.CommunityDetailHandler) // 社区详情
	v1.GET("/post/:id", controller.GetPostDetailHandler) // 帖子详情
	
	// v1.Use(middlewares.JWTAuthMiddleware())
	v1.Use(middlewares.JWTAuthMiddleware(), middlewares.RateLimitMiddleware(2*time.Second, 1)) // 需要登录认证之后才能访问的接口
	// 下面这些需要认证
	// api 限速
	{
		v1.POST("/post", controller.CreatePostHandler)       // 发帖
		v1.POST("/vote", controller.PostVoteController) // 点赞踩)
	}

	pprof.Register(r) // 注册性能分析相关的路由
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "404",
		})
	})

	return r
}
