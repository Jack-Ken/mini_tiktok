package router

import (
	"mini_tiktok/internal/controller"
	"mini_tiktok/internal/initialize"
	"mini_tiktok/utils"

	"github.com/gin-gonic/gin"
)

func InitRouter(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(initialize.GinLogger(), initialize.GinRecovery(true))
	base := r.Group("/douyin")

	base.GET("/feed", controller.Feed_Hanlder) // 视频流接口

	userGroup := base.Group("/user")
	{
		userGroup.GET("/", utils.JWTAuthMiddleware(), controller.QueryInfo_Hanlder) // 用户信息
		userGroup.POST("/register", controller.Register_Hanlder)                    // 用户注册接口
		userGroup.POST("/login", controller.Login_Hanlder)                          // 用户登录接口
	}

	publishGroup := base.Group("/publish")
	{
		publishGroup.POST("/action", utils.JWTAuthMiddleware(), controller.PublishVideosHandler) // 视频投稿
		publishGroup.GET("/list", utils.JWTAuthMiddleware(), controller.PublishListHandler)      // 发布列表
	}

	favoriteGroup := base.Group("/favorite")
	{
		favoriteGroup.POST("/action") // 赞操作
		favoriteGroup.GET("/list")    // 喜欢列表
	}

	commentGroup := base.Group("/comment")
	{
		commentGroup.POST("/action") // 评论操作
		commentGroup.GET("/list")    // 视频评论列表
	}

	relationGroup := base.Group("/relation")
	{
		relationGroup.POST("/action")       // 关注操作
		relationGroup.GET("/follow/list")   // 用户关注列表
		relationGroup.GET("/follower/list") // 用户粉丝列表
		relationGroup.GET("/friend/list")   // 用户好友列表
	}

	return r
}
