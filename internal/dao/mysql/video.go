package mysql

import (
	"errors"
	"mini_tiktok/internal/dao/models"
	db "mini_tiktok/internal/initialize"
	"time"
)

// 获取视频
func QueryVideoListByLimtAndTime(limit int, latestTime time.Time, videoList *[]*models.Video) error {
	if videoList == nil {
		return errors.New("空指针错误")
	}
	return db.SqlSession.Model(&models.Video{}).Where("created_at<?", latestTime).Order("created_at DESC").
		Limit(limit).Find(videoList).Error
}

func QueryVideoListByUserId(uid int64, videoList *[]*models.Video) error {
	if videoList == nil {
		return errors.New("空指针错误")
	}
	return db.SqlSession.Where("user_id=?", uid).
		Select([]string{"id", "user_id", "play_url", "cover_url", "favorite_count", "comment_count", "is_favorite", "title"}).
		Find(videoList).Error

}

func AddVideo(video *models.Video) error {
	if video == nil {
		return errors.New("AddVideo video 空指针")
	}
	return db.SqlSession.Create(video).Error
}
