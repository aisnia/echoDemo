package controllers

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/hashicorp/go-uuid"
	jemail "github.com/jordan-wright/email"
	"github.com/labstack/echo/v4"
	"learn_together/commons"
	"learn_together/logic"
	"learn_together/models/bean"
	"learn_together/models/vo"
	"learn_together/service/cache"
	"learn_together/service/email"
	"net/http"
	"time"
)

var c = context.Background()

//发送邮箱的验证码
func SendEmailCode(ctx echo.Context) error {

	eNum := ctx.FormValue("email")
	fmt.Println(eNum)
	//1.生成随机的字符串
	code := commons.RandStringBytesRmndr(4)
	//2. 生成对应的 id
	id, err := uuid.GenerateUUID()
	if err != nil {
		return err
	}
	//验证码存在5分钟
	cache.Rdb.Set(c, id, code, 5*time.Minute)

	//发送对应的邮件
	e := jemail.NewEmail()
	e.From = "ailengy369@qq.com"
	e.To = []string{eNum}
	e.Subject = "验证码邮件"
	e.HTML = []byte(`
		【一起学】您的验证码：  <a>` + code + `</a>,您正在进行身份验证，打死都不要告诉别人`)
	email.Send(e)
	res := map[string]interface{}{
		"token": id,
	}
	return ctx.JSON(http.StatusOK, commons.Resp{Code: 0, Data: res})
}

//注册
func Register(ctx echo.Context) error {
	login := &vo.RegisterReq{}
	if err := ctx.Bind(login); err != nil {
		return ctx.JSON(http.StatusOK, commons.Resp{Code: 1, Message: "Bind Error"})
	}

	if err := ctx.Validate(login); err != nil {
		return ctx.JSON(http.StatusOK, commons.Resp{Code: 1, Message: "Validate Error"})
	}

	loginService := &logic.LoginServiceImpl{}
	_, err := loginService.Register(login)
	if err != nil {
		return ctx.JSON(http.StatusOK, commons.Resp{Code: 1, Message: "Register error"})
	}
	return ctx.JSON(http.StatusOK, commons.Resp{Code: 0, Message: "Ok"})
}

//登录
func Login(ctx echo.Context) error {
	login := &bean.Login{}
	if err := ctx.Bind(login); err != nil {
		return ctx.JSON(http.StatusOK, commons.Resp{Code: 1, Message: "Bind Error"})
	}
	//验证参数有没有
	if login.Password == "" || login.Email == "" {
		return ctx.JSON(http.StatusOK, commons.Resp{Code: 1, Message: "Param Error"})
	}
	//登录
	loginService := &logic.LoginServiceImpl{}
	exit, err := loginService.Login(login)
	if err != nil || !exit {
		return ctx.JSON(http.StatusOK, commons.Resp{Code: 1, Message: "用户或密码错误"})
	}

	//生成JWT
	token, err := commons.GetToken(login.Id, login.Name, login.Email)
	if err != nil {
		return ctx.JSON(http.StatusOK, commons.Resp{Code: 1, Message: "生成Token失败"})
	}

	data := map[string]interface{}{
		"token": token,
	}
	return ctx.JSON(http.StatusOK, commons.Resp{Code: 0, Data: data})
}

//密码检查
func CheckPassword(ctx echo.Context) error {
	token := ctx.Get("user").(*jwt.Token)
	claim := token.Claims.(*commons.JwtCustomClaims)
	password := ctx.FormValue("password")

	loginService := &logic.LoginServiceImpl{}
	has, err := loginService.CheckPassword(claim.Id, password)
	if !has || err != nil {
		return ctx.JSON(http.StatusOK, commons.Resp{Code: -1, Message: "原密码错误"})
	}
	return ctx.JSON(http.StatusOK, commons.Resp{Code: 0})
}

//修改密码
func UpdatePassword(ctx echo.Context) error {
	token := ctx.Get("user").(*jwt.Token)
	claim := token.Claims.(*commons.JwtCustomClaims)
	password := ctx.FormValue("password")
	newPassword := ctx.FormValue("newPassword")

	loginService := &logic.LoginServiceImpl{}
	has, err := loginService.UpdatePassword(claim.Id, password, newPassword)
	if !has || err != nil {
		return ctx.JSON(http.StatusOK, commons.Resp{Code: -1, Message: "原密码错误"})
	}

	return ctx.JSON(http.StatusOK, commons.Resp{Code: 0})
}

//修改个人资料
func UpdateUserInfo(ctx echo.Context) error {
	user := &bean.User{}
	if err := ctx.Bind(user); err != nil {
		return ctx.JSON(http.StatusOK, commons.Resp{Code: -1, Message: "参数错误"})
	}

	if err := ctx.Validate(user); err != nil {
		return ctx.JSON(http.StatusOK, commons.Resp{Code: -1, Message: "参数格式不对"})
	}

	return nil
}
