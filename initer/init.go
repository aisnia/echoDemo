package initer

import (
	"github.com/spf13/viper"
)

func LoadConfig(confPath string) (config Config) {
	viper.SetConfigFile(confPath)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			panic("config is no found")
		} else {
			// Config file was found but another error was produced
			panic("config is error")
		}
	}

	viper.Unmarshal(&config)
	return config
	//fmt.Println(config)
}

type Config struct {
	LogPath    string       //日志路径
	HideBanner bool         //是否隐藏Banner
	Server     ServerConfig //echo服务器的配置
	Auth       AuthConfig   //Auth 服务
	Mysql      MysqlConfig
	Redis      RedisConfig
	Email      EmailConfig  //邮件服务
}
type ServerConfig struct {
	Addr         string //监听的socket地址端口
	ReadTimeout  int
	WriteTimeout int
}
type AuthConfig struct {
	AuthAddr        string //服务地址
	AuthServiceName string
	ModelPath       string
	PolicyPath      string //初始化策略
	IsLoadPolicy    bool
	Driver          string
	Connection      string
	DbSpecified     bool
}
type RedisConfig struct {
	Addr string //监听的socket地址端口
	PWD  string
	DB   int
}
type MysqlConfig struct {
	Connection string
}

type EmailConfig struct {
	Addr     string // smtp.qq.com:25  服务器地址
	Username string // 用来发送邮件的用户
	Password string //对应的授权码
	Host     string //域名 用来安全认证
	PoolSize int    //Email连接池的个数
}
