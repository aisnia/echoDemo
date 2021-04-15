package main

import (
	"flag"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"io"
	"learn_together/commons"
	"learn_together/controllers"
	"learn_together/initer"
	myjwt "learn_together/mymiddleware/jwt"
	"learn_together/service/cache"
	"learn_together/service/db"
	"learn_together/service/email"
	"net/http"
	"os"
	"time"
)

var (
	config   initer.Config
	confPath = flag.String("c", "configs/learn_together.toml", "specify the configuration file, default is configs/learn_together.toml")
	logFile  io.Writer
)

func main() {
	//加载配置
	config = initer.LoadConfig(*confPath)

	//logFile, err := os.OpenFile(config.LogPath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	//if err != nil {
	//	panic(err.Error())
	//}
	//defer logFile.Close()

	logFile = os.Stdout

	Init()
	//1. echo 初始化
	e := echo.New()

	//1.1日志
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
	InitMiddleware(e)

	e.Logger.Fatal(e.StartServer(s))

}

func Init() {
	//xorm
	db.InitXorm(&config)
	//redis
	cache.InitRedis(&config)
	//email
	email.InitEmial(&config, logFile)

}

func InitMiddleware(e *echo.Echo) {
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
	route.GET("/sendCode", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "hello world")
	})

	e.POST("/sendCode", controllers.SendEmailCode)
	e.POST("/register", controllers.Register)
	//e.GET("/sendCode", controllers.SendEmailCode)

	//6、auth 权限校验 用的auth_control的服务

	//auth.AuthInit(&config)
	//e.Use(casbinmiddleware.MiddlewareWithConfig(casbinmiddleware.Config{
	//	Enforcer: auth.Enforcer,
	//}))

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
