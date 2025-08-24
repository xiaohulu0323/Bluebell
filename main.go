package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"web-app/controller"
	"web-app/dao/mysql"
	"web-app/dao/redis"
	"web-app/logger"
	"web-app/pkg/snowflake"
	"web-app/router"
	"web-app/settings"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func main() {
	// if len(os.Args) < 2 {
	// 	fmt.Println("need config file.eg : bluebell config.yaml")
	// 	return
	// }

	// 1. 加载配置
	var err error
	if len(os.Args) >= 2 {
		// 有参数，使用指定的配置文件
		err = settings.Init(os.Args[1])
	} else {
		// 没参数，使用默认配置文件
		err = settings.Init("")
	}
	
	if err != nil {
		fmt.Printf("settings.Init() failed, err: %v \n", err)
		return
	}

	// 2. 初始化日志
	if err := logger.Init(settings.Conf.LogConfig, settings.Conf.Mode); err != nil {
		fmt.Printf("logger.Init() failed, err: %v \n", err)
		return
	}
	defer zap.L().Sync() // 确保日志在程序退出前被写入
	zap.L().Debug("logger init success...")

	// 3. 初始化MySQL连接
	if err := mysql.Init(settings.Conf.MySQLConfig); err != nil {
		fmt.Printf("mysql.Init() failed, err: %v \n", err)
		return
	}
	defer mysql.Close() // 确保在程序退出时关闭数据库连接

	// 4. 初始化redis连接
	if err := redis.Init(settings.Conf.RedisConfig); err != nil {
		fmt.Printf("redis.Init() failed, err: %v \n", err)
		return
	}
	defer redis.Close() // 确保在程序退出时关闭redis连接

	if err := snowflake.Init(settings.Conf.StartTime, settings.Conf.MachineID); err != nil {
		fmt.Printf("snowflake.Init() failed, err: %v \n", err)
		return
	}

	// 初始化gin框架内置的校验器使用的翻译器
	if err := controller.InitTrans("zh"); err != nil {
		fmt.Printf("controller.InitTrans() failed, err: %v \n", err)
		return
	}

	// 5. 注册路由
	r := router.Setup(settings.Conf.Mode)

	// 6. 启动服务（优雅关机）
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", viper.GetInt("port")),
		Handler: r,
	}

	go func() {
		// 开启一个goroutine启动服务
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 等待中断信号来优雅地关闭服务器，为关闭服务器操作设置一个5秒的超时
	quit := make(chan os.Signal, 1) // 创建一个接收信号的通道
	// kill 默认会发送 syscall.SIGTERM 信号
	// kill -2 发送 syscall.SIGINT 信号，我们常用的Ctrl+C就是触发系统SIGINT信号
	// kill -9 发送 syscall.SIGKILL 信号，但是不能被捕获，所以不需要添加它
	// signal.Notify把收到的 syscall.SIGINT或syscall.SIGTERM 信号转发给quit
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 此处不会阻塞
	<-quit                                               // 阻塞在此，当接收到上述两种信号时才会往下执行
	zap.L().Info("Shutdown Server ...")
	// 创建一个5秒超时的context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// 5秒内优雅关闭服务（将未处理完的请求处理完再关闭服务），超过5秒就超时退出
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("Server Shutdown: ", zap.Error(err))
	}

	zap.L().Info("Server exiting")
}
