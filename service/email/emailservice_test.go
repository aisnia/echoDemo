package email

import (
	"fmt"
	jemail "github.com/jordan-wright/email"
	"net/smtp"
	"testing"
)

func Test(t *testing.T) {
	//config := initer.LoadConfig("../../configs/learn_together.toml")
	//InitEmial(&config, os.Stdout)
	e := jemail.NewEmail()
	e.From = "1350017101@qq.com"
	eNum := "1935457604@qq.com"
	//eNum := "1350017101@qq.com"
	code := "???"
	e.To = []string{eNum}
	e.Subject = "验证码邮件"
	e.HTML = []byte(`
		【一起学】您的验证码：<u>` + code + `<u>,您正在进行身份验证，打死都不要告诉别人`)
	//Send(e)
	err :=e.Send("smtp.qq.com:25",smtp.PlainAuth("","1350017101@qq.com","kvrabzhbeirpfijc","smtp.qq.com"))
	if err != nil{
		fmt.Println(err)
	}
}
