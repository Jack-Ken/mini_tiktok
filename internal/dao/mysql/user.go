package mysql

import (
	"errors"
	"mini_tiktok/internal/dao/models"
	db "mini_tiktok/internal/initialize"
	"mini_tiktok/utils"

	"gorm.io/gorm"
)

// dao层主要的操作是将对数据库的操作封装，用于被services层调用

// register
// 检查用户名是否已经存在

func CheckUserExist(username string) (err error) {
	var user models.User
	if !errors.Is(db.SqlSession.Where("username=?", username).First(&user).Error, gorm.ErrRecordNotFound) {
		return errors.New("用户已存在")
	}
	//if !errors.Is(global.App.DY_DB.Model(&model.User{}).Where("username = ?", user.Username).First(&u).Error, gorm.ErrRecordNotFound) {
	//	return errors.New("this username is registered already"), user
	//}
	return
}

// 向数据库中插入注册的新用户数据

func InserUser(user *models.User) (err error) {
	if err = db.SqlSession.Create(user).Error; err != nil {
		return err
	}
	return
}

//login
//
func Login(r *models.LoginRequest) (u *models.User, err error) {
	var user models.User
	if errors.Is(db.SqlSession.Where("username = ?", r.Username).First(&user).Error, gorm.ErrRecordNotFound) {
		// 用户名错误
		return nil, errors.New("用户不存在")
	}
	if user.Password != utils.EncryptPassword(r.Password) {
		return nil, errors.New("密码错误")
	}
	return &user, nil
}

//userinfo

func UserInfo(id int64, name string) (u *models.User, err error) {
	var user models.User
	if errors.Is(db.SqlSession.Where("username = ? AND id = ?", name, id).First(&user).Error, gorm.ErrRecordNotFound) {
		// 用户名错误
		return nil, errors.New("用户不存在")
	}
	return &user, nil
}
