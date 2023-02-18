package models

type Login struct {
	Id       int64 `gorm:"primary_key"`
	UserId   int64
	Username string `gorm:""`
	Password string `gorm:"size:200;notnull"`
}
