package commons

import (
	"github.com/dgrijalva/jwt-go"
)

//返回值
type Resp struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

const SECRET = "learn_together"

type JwtCustomClaims struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	jwt.StandardClaims
}
