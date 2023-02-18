package models

import (
	"time"
)

type Comment struct {
	Id         int64     `json:"id"`
	UserId     int64     `json:"-"` // 用于1用户对多评论关系的id
	VideoId    int64     `json:"-"`
	User       User      `json:"user" gorm:"-"`
	Content    string    `json:"content"`
	CreateDate string    `json:"create-date" gorm:"-"`
	CreatedAt  time.Time // 创建时间

}
