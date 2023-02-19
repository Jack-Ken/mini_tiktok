package mysql

import (
	db "mini_tiktok/internal/initialize"

	"gorm.io/gorm"
)

func GetVideoFavorState(userId int64, videoId int64) (bool, error) {
	var count int64
	err := db.SqlSession.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("select * from user_favor_videos where user_id = ? and video_id = ?", userId, videoId).Count(&count).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return false, err
	}
	if count == 0 {
		return false, nil
	}
	return true, nil
}
