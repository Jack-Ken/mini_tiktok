package main

import (
	"context"
	"fmt"
	"mini_tiktok/internal/initialize"
	"mini_tiktok/internal/router"
	validator "mini_tiktok/pkg"
	"mini_tiktok/pkg/snowflake"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func main() {
	// 1、加载配置文件
	if err := initialize.Load_Config(); err != nil {
		fmt.Printf("initialize config setting failed, err:%v \n", err)
		return
	}
	// 2、初始化日志
	if err := initialize.Init_Log(initialize.Conf.LogConfig, initialize.Conf.AppConfig.Mode); err != nil {
		fmt.Printf("initialize logger failed, err:%v \n", err)
		return
	}
	defer zap.L().Sync()
	zap.L().Debug("logger initialize success...")
	// 3、初始化MySQL连接
	if err := initialize.Init_Mysql(initialize.Conf.MySqlConfig); err != nil {
		fmt.Printf("initialize mysql link failed, err:%v \n", err)
		return
	}
	// 4、初始化Redis连接
	if err := initialize.Init_Redis(initialize.Conf.RedisConfig); err != nil {
		fmt.Printf("initialize redis link failed, err:%v \n", err)
		return
	}
	defer initialize.Close_Redis()
	//雪花算法生成ID的初始化,genID用于创建新的ID，且已经做了互斥处理

	//fmt.Println(genID.GetID())
	if err := snowflake.Init(initialize.Conf.AppConfig.StartTime, initialize.Conf.AppConfig.MachineID); err != nil {
		fmt.Printf("snowflake init failed, err:%v \n", err)
		return
	}

	// 初始化gin框架内置的校验器的翻译器
	if err := validator.InitTrans("zh"); err != nil {
		fmt.Printf("valitador trans init failed, err:%v \n", err)
		return
	}

	// 5、注册路由
	r := router.InitRouter(initialize.Conf.AppConfig.Mode)
	// 6、启动服务(优雅关节）

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", initialize.Conf.AppConfig.Port),
		Handler: r,
	}

	go func() {
		// 开启一个goroutine启动服务
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Fatal("listen: %s\n", zap.Error(err))
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
