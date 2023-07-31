package client

import (
	"context"
	"mini_tiktok/internal/dao/models"
	authService "mini_tiktok/internal/rpc/rpcGen/auth"
	service "mini_tiktok/internal/services"
	"mini_tiktok/utils"
)

type AuthService struct {
	*authService.UnimplementedAuthServiceServer
}

func (auth *AuthService) Register(ctx context.Context, req *authService.RegisterRequest) (*authService.RegisterResponse, error) {
	r := models.RegisterRequest{
		Username: req.Username,
		Password: req.Password,
	}

	// 2、业务处理（services定义，此处只做调用处理）
	userInfo, err := service.Register_service(&r)
	if err != nil {
		return &authService.RegisterResponse{
			StatusCode: -1,
			StatusMsg:  err.Error(),
		}, err
	}
	// 3、获取token
	token, err := utils.GenToken(userInfo.Name, userInfo.Id)
	if err != nil {
		return &authService.RegisterResponse{
			StatusCode: -1,
			StatusMsg:  err.Error(),
		}, err
	}
	return &authService.RegisterResponse{
		StatusCode: 0,
		StatusMsg:  "success",
		UserId:     uint32(userInfo.Id),
		Token:      token,
	}, nil
}

func (auth *AuthService) Login(ctx context.Context, req *authService.LoginRequest) (*authService.LoginResponse, error) {
	r := models.LoginRequest{
		Username: req.Username,
		Password: req.Password,
	}
	// 2、业务处理（services定义，此处只做调用处理）
	loginResp, err := service.Login_Service(&r)
	if err != nil {
		return &authService.LoginResponse{
			StatusCode: -1,
			StatusMsg:  err.Error(),
		}, err
	}
	// 3、获取token值
	token, err := utils.GenToken(loginResp.Username, loginResp.UserId)
	if err != nil {
		return &authService.LoginResponse{
			StatusCode: -1,
			StatusMsg:  err.Error(),
		}, err
	}
	return &authService.LoginResponse{
		StatusCode: 0,
		StatusMsg:  "success",
		UserId:     uint32(loginResp.Id),
		Token:      token,
	}, nil
}
func (auth *AuthService) QueryInfo(ctx context.Context, req *authService.QueryInfoRequest) (*authService.QueryInfoResponse, error) {
	userResp, err := service.Info_Service(req.UserId)
	if err != nil {
		return &authService.QueryInfoResponse{
			StatusCode: -1,
			StatusMsg:  err.Error(),
		}, err
	}
	return &authService.QueryInfoResponse{
		StatusCode: -1,
		StatusMsg:  err.Error(),
		User: &authService.User{
			Id:            userResp.Id,
			Name:          userResp.Name,
			FollowCount:   userResp.FollowCount,
			FollowerCount: userResp.FollowerCount,
			IsFollow:      userResp.IsFollow,
		},
	}, nil
}
