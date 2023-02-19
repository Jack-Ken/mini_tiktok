package controller

import (
	"mini_tiktok/internal/dao/models"
	service "mini_tiktok/internal/services"
	"mini_tiktok/pkg"
	"mini_tiktok/utils"
	"net/http"

	"github.com/go-playground/validator/v10"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

// login handler

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

	// 2、业务处理（services定义，此处只做调用处理）
	loginResp, err := service.Login_Service(&r)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.RegisterResponse{
			Response: models.Response{
				StatusCode: -1,
				StatusMsg:  err.Error(),
			},
		})
		return
	}
	// 3、获取token值
	token, err := utils.GenToken(loginResp.Username, loginResp.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.RegisterResponse{
			Response: models.Response{
				StatusCode: -1,
				StatusMsg:  err.Error(),
			},
		})
		return
	}
	// 4、返回响应
	c.JSON(http.StatusOK, models.RegisterResponse{
		Response: models.Response{
			StatusCode: 0,
			StatusMsg:  "success",
		},
		UserId: loginResp.UserId,
		Token:  token,
	})

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

	// 对请求参数进行参数校验
	// 2、业务处理（services定义，此处只做调用处理）
	userInfo, err := service.Register_service(&r)
	if err != nil {
		zap.L().Error("login failed", zap.String("username", r.Username), zap.String("password", r.Password), zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.RegisterResponse{
			Response: models.Response{
				StatusCode: -1,
				StatusMsg:  err.Error(),
			},
		})
		return
	}
	// 3、获取token
	token, err := utils.GenToken(userInfo.Name, userInfo.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.RegisterResponse{
			Response: models.Response{
				StatusCode: -1,
				StatusMsg:  err.Error(),
			},
		})
		return
	}
	// 4、返回响应
	c.JSON(http.StatusOK, models.RegisterResponse{
		Response: models.Response{
			StatusCode: 0,
			StatusMsg:  "success",
		},
		UserId: userInfo.Id,
		Token:  token,
	})
}

func QueryInfo_Hanlder(c *gin.Context) {
	//todo
	//username := c.GetString("username")
	userid := c.GetInt64("userid")

	userResp, err := service.Info_Service(userid)
	if err != nil {
		zap.L().Error("get user information failed", zap.String("username", userResp.Name), zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.RegisterResponse{
			Response: models.Response{
				StatusCode: -1,
				StatusMsg:  err.Error(),
			},
		})
		return
	}

	//user := models.UserInfo{
	//	ID:            resp.ID,
	//	UserName:      resp.Username,
	//	FollowCount:   resp.FollowCount,
	//	FollowerCount: resp.FollowerCount,
	//	IsFollow:      resp.IsFollow,
	//}

	c.JSON(http.StatusOK, models.UserInfoResponse{
		Response: models.Response{
			StatusCode: 0,
			StatusMsg:  "success",
		},
		User: *userResp,
	})

}
