package models

// user表的结构体定义

type User struct {
	Id            int64      `json:"id" gorm:"id,omitempty"`
	Name          string     `json:"name" gorm:"name,omitempty"`
	FollowCount   int64      `json:"follow_count" gorm:"follow_count,omitempty"`
	FollowerCount int64      `json:"follower_count" gorm:"follower_count,omitempty"`
	IsFollow      bool       `json:"is_follow" gorm:"is_follow,omitempty"`
	User          *Login     `json:"-"`                                     //用户与账号密码之间的一对一
	Videos        []*Video   `json:"-"`                                     //用户与投稿视频的一对多
	Follows       []*User    `json:"-" gorm:"many2many:user_relations;"`    //用户之间的多对多
	FavorVideos   []*Video   `json:"-" gorm:"many2many:user_favor_videos;"` //用户与点赞视频之间的多对多
	Comments      []*Comment `json:"-"`
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

//type UserLoginRegisterResponse struct {
//	UserId int64  `json:"user_id"` // ,omitempty
//	Token  string `json:"token"`
//}

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
	User
}
