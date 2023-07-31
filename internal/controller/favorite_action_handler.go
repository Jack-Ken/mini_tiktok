package controller

import (
	"mini_tiktok/internal/dao/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func FavoriteActionHandler(c *gin.Context) {
	//uid := c.GetInt64("userid")
	actionType := c.PostForm("action_type")
	p := NewProxyFavoriteAction(c)
	if actionType == "" {
		p.FavoriteActionError("操作为空!!!")
	}
	action, _ := strconv.Atoi(actionType)
	//strconv.ParseInt(actionType, 10, 32)

	// 判断点赞操作类型：取消或者点赞
	if action == 1 { // 点赞

	}
	// 取消点赞

	// 更新数据库，根据video_id修改video表中的favorite_count和在user_favor_videos表中添加数据
}

type ProxyFavoriteAction struct {
	*gin.Context
}

func NewProxyFavoriteAction(c *gin.Context) *ProxyFavoriteAction {
	return &ProxyFavoriteAction{c}
}

func (p *ProxyFavoriteAction) FavoriteActionError(msg string) {
	p.JSON(http.StatusBadRequest, models.Response{
		StatusCode: -1,
		StatusMsg:  msg,
	})
}

func (p *ProxyFavoriteAction) FavoriteActionSuccess() {
	p.JSON(http.StatusOK, models.Response{
		StatusCode: 0,
		StatusMsg:  "success",
	})
}
