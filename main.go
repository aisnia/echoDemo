package main

import (
	"flag"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"io"
	"learn_together/commons"
	"learn_together/init"
	"learn_together/mymiddleware/casbinmiddleware"
	myjwt "learn_together/mymiddleware/jwt"
	"learn_together/service/auth"
	"learn_together/service/db"
	"net/http"
	"os"
	"time"
)

var (
	config   init.Config
	confPath = flag.String("c", "configs/learn_together.toml", "specify the configuration file, default is configs/learn_together.toml")
	logFile  io.Writer
)

func main() {
	//加载配置
	config = init.LoadConfig(*confPath)
	//1. echo 初始化
	db.InitXorm(&config)
	e := echo.New()

	//1.1日志
	//logFile, err := os.OpenFile(config.LogPath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	//if err != nil {
	//	panic(err.Error())
	//}
	//defer logFile.Close()

	logFile = os.Stdout

	e.Logger.SetOutput(logFile)
	//日志格式
	e.Logger.SetHeader("${time_rfc3339} ${level}")

	//1.2 校验器
	e.Validator = &DefaultValidator{validator: validator.New()}
	//1.3 Server
	s := &http.Server{
		Addr:         config.Server.Addr,
		ReadTimeout:  time.Duration(config.Server.ReadTimeout) * time.Minute,
		WriteTimeout: time.Duration(config.Server.WriteTimeout) * time.Minute,
	}
	//1.4 是否有横幅
	e.HideBanner = false
	//1.5 中间件
	Init(e)

	e.Logger.Fatal(e.StartServer(s))

}
func Init(e *echo.Echo) {
	//1、恢复
	e.Use(middleware.Recover())
	//2、请求ID
	e.Use(middleware.RequestID())
	//3、跨域
	//e.Use(middleware.CORS())
	//4、日志
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "id=${id},status=${status},time=${time_rfc3339_nano},remote_ip=${remote_ip},host=${host},method=${method},uri=${uri},latency=${latency},query=${query},form=${form}\n",
		Output: logFile,
	}))

	//5.jwt 判断是否登录
	jwtMiddleware := myjwt.JWTWithConfig(myjwt.JWTConfig{
		SigningKey:  []byte(commons.SECRET),
		Claims:      &JwtCustomClaims{},
		TokenLookup: "form:token",
		SuccessHandler: func(c echo.Context) {
			
		},
	})

	route := e.Group("/learn", jwtMiddleware)
	route.GET("/hello", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "hello world")
	})
	//6、auth 权限校验 用的auth_control的服务

	auth.AuthInit(&config)
	e.Use(casbinmiddleware.MiddlewareWithConfig(casbinmiddleware.Config{
		Enforcer: auth.Enforcer,
	}))

}

type DefaultValidator struct {
	validator *validator.Validate
}

func (cv *DefaultValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

type JwtCustomClaims struct {
	Id       int64  `json:"id"`
	UserName string `json:"username"`
	jwt.StandardClaims
}
