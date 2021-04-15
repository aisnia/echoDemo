package bean

import "time"

//登录类
type Login struct {
	Id         int       `json:"id"`
	Name       string    `json:"name" xorm:"name" `
	Email      string    `json:"email" xorm:"email" `
	Password   string    `json:"password" xorm:"password"`
	CreateTime time.Time `json:"create_time" xorm:"created"`
	UpdateTime time.Time `json:"update_time" xorm:"created updated"`
	DeleteAt   time.Time `json:"delete_at" xorm:"deleted"`
}

//用户类
type User struct {
	Id       int       `json:"id"`
	UId      int       `json:"u_id" xorm:"u_id"`
	Sex      int       `json:"sex" xorm:"sex"`
	Age      int       `json:"age" xorm:"age"`
	Birthday time.Time `json:"birthday" xorm:"birthday"`
	Phone    string    `json:"phone" xorm:"phone"`
	Address  string    `json:"address" xorm:"address"`
	School   string    `json:"school" xorm:"school" `
	State    int       `json:"state" xorm:"state"`
}

//
