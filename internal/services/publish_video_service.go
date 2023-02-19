package services

import (
	"mini_tiktok/internal/dao/models"
	"mini_tiktok/internal/dao/mysql"
)

type PostVideoFlow struct {
	playUrl  string
	coverUrl string
	title    string
	userId   int64

	video *models.Video
}

func PostVideo(uid int64, videoPath string, coverPath string, title string) error {
	return NewPostVideoFlow(uid, videoPath, coverPath, title).Do()
}
func NewPostVideoFlow(uid int64, videoPath string, coverPath string, title string) *PostVideoFlow {
	return &PostVideoFlow{
		playUrl:  videoPath,
		coverUrl: coverPath,
		title:    title,
		userId:   uid,
	}
}

func (p *PostVideoFlow) Do() error {
	if err := p.publish(); err != nil {
		return err
	}
	return nil
}

func (p *PostVideoFlow) publish() error {
	video := &models.Video{
		UserId:   p.userId,
		Title:    p.title,
		PlayUrl:  p.playUrl,
		CoverUrl: p.coverUrl,
	}
	if err := mysql.AddVideo(video); err != nil {
		return err
	}
	return nil
}
