package initialize

import (
	"fmt"

	"github.com/fsnotify/fsnotify"

	"github.com/spf13/viper"
)

// Conf用于保存配置信息
var Conf = new(Config)

type Config struct {
	*AppConfig   `mapstructure:"app"`
	*LogConfig   `mapstructure:"log"`
	*MySqlConfig `mapstructure:"mysql"`
	*RedisConfig `mapstructure:"redis"`
	*QiNiuConfig `mapstructure:"qiniu"`
}

type AppConfig struct {
	Name      string `mapstructure:"name"`
	Mode      string `mapstructure:"mode"`
	Version   string `mapstructure:"version"`
	StartTime string `mapstructure:"start_time"`
	MachineID int64  `mapstructure:"machine_id"`
	Port      int    `mapstructure:"port"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	FileName   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
}

type MySqlConfig struct {
	Host               string `mapstructure:"host"`
	Port               int    `mapstructure:"port"`
	User               string `mapstructure:"user"`
	Password           string `mapstructure:"password"`
	DbName             string `mapstructure:"dbname"`
	Charset            string `mapstructure:"charset"`
	MaxOpenConnections int    `mapstructure:"max_open_connections"`
	MaxIdleConnections int    `mapstructure:"max_idle_connections"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	Db       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

type QiNiuConfig struct {
	Accesskey   string `mapstructure:"accesskey"`
	Sercetkey   string `mapstructure:"sercetkey"`
	Bucket      string `mapstructure:"bucket"`
	Qiniuserver string `mapstructure:"qiniuserver"`
}

// 加载配置文件
func Load_Config() (err error) {
	// 读取配置文件
	viper.SetConfigName("config.yaml") // 配置文件名称
	viper.SetConfigType("yaml")        // 如果配置文件中没有扩展名，则需要配置此项(只适用于从远程获取配置信息，从本地查找的话会忽略此代码)
	viper.AddConfigPath("../config")   // 查找配置文件所在路径

	if err = viper.ReadInConfig(); err != nil { // 处理读取配置文件的错误
		fmt.Printf("Fatal error config file: %v \n", err)
		return
	}

	if err = viper.Unmarshal(Conf); err != nil {
		fmt.Printf("viper.Unmarshal config file failed: %v \n", err)
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("The config file has cahnged...")
		if err = viper.Unmarshal(Conf); err != nil {
			fmt.Printf("viper.Unmarshal config file failed: %v \n", err)
		}
	})
	return
}
