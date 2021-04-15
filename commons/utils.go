package commons

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"math/rand"
	"reflect"
	"strings"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytesRmndr(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

//JWT
func newJwtCustomClaims(Id int, name, email string) *JwtCustomClaims {
	return &JwtCustomClaims{
		Id: Id, Name: name, Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 30).Unix(),
		}}
}
func GetToken(Id int, name, email string) (string, error) {
	claims := newJwtCustomClaims(Id, name, email)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(SECRET))
}

//tag标签 转换工具
func GetFieldByTag(data interface{}, tagName string, filterZeroVal bool) (map[string]interface{}, error) {
	fields := make(map[string]interface{})

	value := reflect.ValueOf(data).Elem()
	typ := reflect.TypeOf(data).Elem()

	//类型为struct
	kd := value.Kind()
	if reflect.Struct != kd {
		return nil, errors.New("the param should be struct")
	}

	filedNum := value.NumField()
	for i := 0; i < filedNum; i++ {
		val := value.Field(i).Interface()

		t := value.Field(i).Type()
		//判断字段是什么类型 过滤
		switch t.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			//整型 判断是否过滤0值
			if filterZeroVal {
				vStr := fmt.Sprintf("%v", val)
				if 0 == strings.Compare(vStr, "0") {
					//跳出本次循环不加入map
					continue
				}
			}
		case reflect.String:
			if filterZeroVal {
				vStr := val.(string)
				if "" == vStr {
					continue
				}
			}
		case reflect.Interface:
			if value.Field(i).IsNil() {
				continue
			}
		case reflect.Ptr:
			if value.Field(i).IsNil() {
				continue
			}
		default:
			//不支持float类型
			continue
		}
		key := typ.Field(i).Tag.Get(tagName)
		if "" != key {
			fields[key] = val
		}
	}
	return fields, nil
}
func GetMapsByTag(data interface{}, tagName string, filterZeroVal bool) ([]map[string]interface{}, error) {
	res := make([]map[string]interface{}, 0)
	value := reflect.ValueOf(data).Elem()

	//类型为struct
	kd := value.Kind()
	if reflect.Slice != kd {
		return nil, errors.New("the param should be struct")
	}

	for i := 0; i < value.Len(); i++ {
		v := value.Index(i).Interface()
		m, err := GetFieldByTag(v, tagName, filterZeroVal)
		if err != nil {
			return nil, err
		}
		res = append(res, m)
	}
	return res, nil
}

func MapToStruct(fields map[string]interface{}, data interface{}, tagName string) error {
	value := reflect.ValueOf(data).Elem()
	typ := reflect.TypeOf(data).Elem()
	kind := value.Kind()

	if reflect.Struct != kind {
		return errors.New("the data should be struct")
	}

	fieldsNum := value.NumField()
	for i := 0; i < fieldsNum; i++ {
		key := typ.Field(i).Tag.Get(tagName)
		if v, ok := fields[key]; ok {
			value.Field(i).Set(reflect.ValueOf(v))
		}
	}
	return nil
}

//必须都是指针 Ptr
func StructToStruct(from interface{}, to interface{}, tagName string) error {
	fields, err := GetFieldByTag(from, tagName, true)
	if err != nil {
		return err
	}
	return MapToStruct(fields, to, tagName)
}
