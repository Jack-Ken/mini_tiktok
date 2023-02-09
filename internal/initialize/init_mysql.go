package initialize

import (
	"fmt"

	"gorm.io/gorm/schema"

	"go.uber.org/zap"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var SqlSession *gorm.DB

func Init_Mysql(cfg *MySqlConfig) (err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DbName,
		cfg.Charset,
	)
	SqlSession, err = gorm.Open(mysql.New(mysql.Config{DSN: dsn}), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		zap.L().Error("Connect DB failed", zap.Error(err))
		return
	}
	sqlDb, err := SqlSession.DB()

	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDb.SetMaxIdleConns(cfg.MaxIdleConnections)

	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDb.SetMaxOpenConns(cfg.MaxOpenConnections)

	return
}
