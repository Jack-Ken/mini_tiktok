package controller

import (
	"errors"
	"fmt"
	"mini_tiktok/internal/dao/models"
	"mini_tiktok/internal/services"
	"mini_tiktok/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func Feed_Hanlder(c *gin.Context) {
	//todo
	p := NewProxyFeedVideoList(c)
	token, ok := c.GetQuery("token")
	// 无登录状态
	if !ok {
		fmt.Println("11111")
		if err := p.DoNoToken(); err != nil {
			p.FeedVideoListError(err.Error())
		}
		return
	}
	// 有登录状态
	if err := p.DoWithToken(token); err != nil {
		p.FeedVideoListError(err.Error())
	}
}

type ProxyFeedVideoList struct {
	*gin.Context
}

func NewProxyFeedVideoList(c *gin.Context) *ProxyFeedVideoList {
	return &ProxyFeedVideoList{c}
}

// 未登录的视频流处理
func (p *ProxyFeedVideoList) DoNoToken() error {
	rawTimeStamp := p.Query("latest_time")
	var latestTime time.Time
	intTime, err := strconv.ParseInt(rawTimeStamp, 10, 64)
	if err != nil {
		latestTime = time.Unix(0, intTime*1e6) //注意：前端传来的时间戳是以ms为单位的
	}

	// 调用service层接口，获取videoList
	videoList, err := services.QueryFeedVideoList(0, latestTime)
	if err != nil {
		return err
	}
	p.FeedVideoListSuccess(videoList)
	return nil
}

// 登录的视频流处理
func (p *ProxyFeedVideoList) DoWithToken(token string) error {
	mc, err := utils.ParseToken(token)
	if err != nil {
		return errors.New("token不正确")
	}
	if time.Now().Unix() > mc.ExpiresAt.Unix() {
		return errors.New("token超时")
	}
	rawTimeStamp := p.Query("latest_time")
	var latestTime time.Time
	intTime, err := strconv.ParseInt(rawTimeStamp, 10, 64)
	if err != nil {
		latestTime = time.Unix(0, intTime*1e6) //注意：前端传来的时间戳是以ms为单位的
	}
	// 调用service层接口，获取videoList
	videoList, err := services.QueryFeedVideoList(0, latestTime)
	if err != nil {
		return err
	}
	p.FeedVideoListSuccess(videoList)
	// 调用service层接口，获取videoList
	return nil

}

func (p *ProxyFeedVideoList) FeedVideoListError(msg string) {
	p.JSON(http.StatusOK, models.FeedResponse{
		Response: models.Response{
			StatusCode: 1,
			StatusMsg:  msg,
		},
	})

}

func (p *ProxyFeedVideoList) FeedVideoListSuccess(videoList *services.FeedVideoList) {
	p.JSON(http.StatusOK, models.FeedResponse{
		Response: models.Response{
			StatusCode: 0,
			StatusMsg:  "success",
		},
		VideoList: videoList.VideoList,
		NextTime:  videoList.NextTime,
	})
}
