package models

import (
	"time"

	"gorm.io/gorm"
)

// user表的结构体定义

type User struct {
	ID            int64          `gorm:"primarykey"` // 主键ID
	CreatedAt     time.Time      // 创建时间
	UpdatedAt     time.Time      // 更新时间
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`                            // 删除时间
	FollowCount   int64          `json:"follow_count,omitempty" gorm:"default:0"`   // 关注数
	FollowerCount int64          `json:"follower_count,omitempty" gorm:"default:0"` // 粉丝数
	IsFollow      bool           `json:"is_follow,omitempty" gorm:"default:false"`  // 当前用户是否关注
	Username      string         `json:"username" gorm:"comment:username" `         // 登录账号
	Password      string         `json:"password" gorm:"comment:password"`          // 登录密码
}

// 定义请求的参数结构体

type RegisterRequest struct {
	Username string `json:"username" gorm:"not null; comment:username for register;" form:"username" binding:"required"` // 用户名称
	Password string `json:"password" gorm:"not null; comment:password for register;" form:"password" binding:"required"` // 用户密码
}

type LoginRequest struct {
	Username string `json:"username" gorm:"not null; comment:username for register;" form:"username" binding:"required"` // 用户名称
	Password string `json:"password" gorm:"not null; comment:password for register;" form:"password" binding:"required"` // 用户密码
}

type UserInfoRequest struct {
	UserId int64  `json:"user_id" form:"user_id" binding:"required"`
	Token  string `json:"token" form:"token" binding:"required"`
}

// 定义响应参数结构体

type RegisterResponse struct {
	Response
	UserId int64  `json:"user_id"` // ,omitempty
	Token  string `json:"token"`
}

type LoginResponse struct {
	Response
	UserId int64  `json:"user_id"` // ,omitempty
	Token  string `json:"token"`
}

type UserInfoResponse struct {
	Response
	UserInfo
}

type UserInfo struct {
	ID            int64  `json:"id"`
	UserName      string `json:"name"`
	FollowCount   int64  `json:"follow_count"`
	FollowerCount int64  `json:"follower_count"`
	IsFollow      bool   `json:"is_follow"`
}
