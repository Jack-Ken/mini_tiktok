package controller

import (
	"mini_tiktok/internal/dao/models"
	service "mini_tiktok/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func PublishListHandler(c *gin.Context) {
	// 参数准备
	uid := c.GetInt64("userid")
	p := NewProxyPublishList(c)
	if err := p.QueryPublishVideoList(uid); err != nil {
		p.PublishListError(err.Error())
	}
}

type ProxyPublishList struct {
	*gin.Context
}

func NewProxyPublishList(c *gin.Context) *ProxyPublishList {
	return &ProxyPublishList{c}
}

func (p *ProxyPublishList) QueryPublishVideoList(uid int64) error {
	// 调用service执行服务
	videoList, err := service.PublishListService(uid)
	if err != nil {
		return err
	}
	p.PublishListSuccess(videoList)
	return nil
}

func (p *ProxyPublishList) PublishListError(msg string) error {
	p.JSON(http.StatusOK, models.PublishListResponse{
		Response: models.Response{
			StatusCode: 1,
			StatusMsg:  msg,
		},
	})
	return nil
}

func (p *ProxyPublishList) PublishListSuccess(videoList *[]*models.Video) error {
	p.JSON(http.StatusOK, models.PublishListResponse{
		Response: models.Response{
			StatusCode: 0,
			StatusMsg:  "success",
		},
		VideoList: *videoList,
	})
	return nil
}
