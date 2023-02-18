package models

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}

// response中有关User和Video的信息表

//type UserInfo struct {
//	ID            int64  `json:"id"`   // 用户ID
//	UserName      string `json:"name"` // 用户名
//	FollowCount   int64  `json:"follow_count"`
//	FollowerCount int64  `json:"follower_count"`
//	IsFollow      bool   `json:"is_follow"`
//}

// response结构中使用的信息
//required int64 id = 1; // 视频唯一标识
//required User author = 2; // 视频作者信息
//required string play_url = 3; // 视频播放地址
//required string cover_url = 4; // 视频封面地址
//required int64 favorite_count = 5; // 视频的点赞总数
//required int64 comment_count = 6; // 视频的评论总数
//required bool is_favorite = 7; // true-已点赞，false-未点赞
//required string title = 8; // 视频标题

//type VideoInfo struct {
//	ID            int64    `gorm:"primarykey"`                                  // 视频ID
//	UserID        int64    `json:"author_id,omitempty"`                         // 发布作者
//	User          UserInfo `json:"author" gorm:"foreignKey:UserID"`             // user信息
//	PlayUrl       string   `json:"play_url,omitempty" gorm:"default:testName"`  // 视频地址
//	CoverUrl      string   `json:"cover_url,omitempty" gorm:"default:testName"` // 封面地址
//	FavoriteCount int64    `json:"favorite_count" gorm:"default:0"`             // 点赞数量
//	CommentCount  int64    `json:"comment_count" gorm:"default:0"`              // 评论数量
//	IsFavorite    bool     `json:"is_favorite" gorm:"default:false"`            // 是否点赞
//	Title         string   `json:"title, omitempty" gorm:"comment:视频说明"`        // 投稿时添加的文字
//}
