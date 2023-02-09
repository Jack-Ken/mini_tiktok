package services

import (
	"mini_tiktok/internal/dao/models"
	"mini_tiktok/internal/dao/mysql"
	"mini_tiktok/pkg/snowflake"
	"mini_tiktok/utils"
)

// service层存放业务逻辑代码

func Register_service(r *models.RegisterRequest) (u *models.User, err error) {
	// 1、判断用户存不存在
	if err = mysql.CheckUserExist(r.Username); err != nil {
		return nil, err
	}
	// 2、生成UID
	newID := snowflake.G.GetID()
	// 构造一个User实例
	user := &models.User{
		Username: r.Username,
		Password: utils.EncryptPassword(r.Password), // 密码加密
		ID:       newID,
	}
	//3、保存进数据库
	if err := mysql.InserUser(user); err != nil {
		return nil, err
	}
	return user, nil
}

func Login_Service(r *models.LoginRequest) (u *models.User, err error) {
	var user *models.User
	user, err = mysql.Login(r)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func Info_Service(id int64, name string) (u *models.User, err error) {
	var user *models.User
	user, err = mysql.UserInfo(id, name)
	if err != nil {
		return nil, err
	}
	return user, nil
}
