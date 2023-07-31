package controller

import (
	"context"
	"log"
	"mini_tiktok/internal/dao/models"
	authService "mini_tiktok/internal/rpc/rpcGen/auth"
	"mini_tiktok/pkg"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/go-playground/validator/v10"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

// login handler

var AuthClient authService.AuthServiceClient

const auth_address = "127.0.0.1:8889"

func init() {
	// 连接服务器
	conn, err := grpc.Dial(auth_address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("net.Connect err: %v", err)
	}
	defer conn.Close()

	// 建立gRPC连接
	AuthClient = authService.NewAuthServiceClient(conn)
}

func Login_Hanlder(c *gin.Context) {
	//todo
	//1、获取参数和参数校验
	var r models.LoginRequest
	if err := c.ShouldBind(&r); err != nil {
		// 请求参数有误，直接返回失败响应
		zap.L().Error("login request ShouldBindJson error", zap.Error(err))
		// 判断err是不是validator.ValidationErrors类型
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			c.JSON(http.StatusBadRequest, models.RegisterResponse{
				Response: models.Response{
					StatusCode: -1,
					StatusMsg:  err.Error(),
				},
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": pkg.RemoveTopStruct(errs.Translate(pkg.Trans))})

		//
		return
	}

	// 内部调用grpc服务
	resp, err := AuthClient.Login(context.Background(), &authService.LoginRequest{
		Username: r.Username,
		Password: r.Password,
	})
	if err != nil {
		zap.L().Error("user login failed", zap.Error(err))
	}

	c.JSON(http.StatusOK, resp)
}

// register handler

func Register_Hanlder(c *gin.Context) {
	//todo
	//1、获取参数和参数校验
	var r models.RegisterRequest
	if err := c.ShouldBind(&r); err != nil {
		// 请求参数有误，直接返回失败响应
		zap.L().Error("register request ShouldBindJson error", zap.Error(err))
		// 判断err是不是validator.ValidationErrors类型
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			c.JSON(http.StatusBadRequest, models.RegisterResponse{
				Response: models.Response{
					StatusCode: -1,
					StatusMsg:  err.Error(),
				},
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": pkg.RemoveTopStruct(errs.Translate(pkg.Trans))})

		//
		return
	}
	// 内部调用grpc服务
	resp, err := AuthClient.Register(context.Background(), &authService.RegisterRequest{
		Username: r.Username,
		Password: r.Password,
	})

	if err != nil {
		zap.L().Error("user register failed", zap.Error(err))
	}
	c.JSON(http.StatusOK, resp)
}

func QueryInfo_Hanlder(c *gin.Context) {
	//todo
	//username := c.GetString("username")
	userid := c.GetInt64("userid")

	resp, err := AuthClient.QueryInfo(context.Background(), &authService.QueryInfoRequest{
		UserId: userid,
	})

	if err != nil {
		zap.L().Error("get user information failed", zap.String("username", resp.User.Name), zap.Error(err))
	}

	c.JSON(http.StatusOK, resp)

}
