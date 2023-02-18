package models

import (
	"time"
)

type Message struct {
	Id         int64     `json:"id,omitempty"`
	ToUserId   int64     `json:"to_user_id"`
	FromUserId int64     `json:"from_user_id"`
	Content    string    `json:"content"`
	CreateTime string    `json:"create_time"`
	CreatedAt  time.Time `json:"-"` // 创建时间
}
