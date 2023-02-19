package models

import (
	"time"

	"gorm.io/gorm"
)

// video表的结构体定义

//type Video struct {
//	ID        int64          `gorm:"primarykey"` // 主键ID
//	CreatedAt time.Time      // 创建时间
//	UpdatedAt time.Time      // 更新时间
//	DeletedAt gorm.DeletedAt `json:"-" gorm:"index" `                     // 删除时间
//	AuthorID  int64          `json:"author_id" gorm:"comment:author_id" ` // 视频作者的id
//	PlayUrl   string         `json:"play_url" gorm:"comment:play_url"`    // 播放地址
//	CoverUrl  string         `json:"cover_url" gorm:"comment:cover_url"`  // 封面地址
//	Title     string         `json:"title" gorm:"comment:title"`          // 视频标题
//}

type Video struct {
	Id            int64          `json:"id,omitempty" gorm:"primarykey"` // 视频ID
	UserId        int64          `json:"-" gorm:"author_id"`
	Author        User           `json:"author,omitempty" gorm:"-"` //这里应该是作者对视频的一对多的关系，而不是视频对作者，故gorm不能存他，但json需要返回它
	PlayUrl       string         `json:"play_url,omitempty" gorm:"play_url"`
	CoverUrl      string         `json:"cover_url,omitempty" gorm:"cover_url"`
	FavoriteCount int64          `json:"favorite_count,omitempty"`
	CommentCount  int64          `json:"comment_count,omitempty"`
	IsFavorite    bool           `json:"is_favorite,omitempty"`
	Title         string         `json:"title,omitempty" gorm:"title"`
	Users         []*User        `json:"-" gorm:"many2many:user_favor_videos;"`
	Comments      []*Comment     `json:"-"`
	CreatedAt     time.Time      `json:"-"`               // 创建时间
	UpdatedAt     time.Time      `json:"-"`               // 更新时间
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index" ` // 删除时间
}

// request参数接结构体

type FeedRequest struct {
	LatestTime int64  `json:"latest_time" form:"latest_time" binding:"required"`
	Token      string `json:"token" form:"token" binding:"required"`
}

type PublishListRequest struct {
	UserId int64  `json:"user_id" form:"user_id" binding:"required"`
	Token  string `json:"token" form:"token" binding:"required"`
}

// response 参数结构体

type FeedResponse struct {
	Response
	VideoList []*Video `json:"video_list,omitempty"`
	NextTime  int64    `json:"next_time,omitempty"`
}

type PublishListResponse struct {
	Response
	VideoList []*Video `json:"video_list,omitempty"`
}
