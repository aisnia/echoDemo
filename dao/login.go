package dao

import (
	"learn_together/models/bean"
	"learn_together/service/db"
)

type LoginDao interface {
	AddLogin(login *bean.Login) (int64, error)
	Login(login *bean.Login) (bool, error)
	SelectById(login *bean.Login) (bool, error)
}

type MysqlLoginDao struct {
}

func (*MysqlLoginDao) AddLogin(login *bean.Login) (int64, error) {
	return db.Engine.InsertOne(login)
}

func (*MysqlLoginDao) Login(login *bean.Login) (bool, error) {
	return db.Engine.Where("email = ? and password = ?", login.Email, login.Password).Get(login)
}

func (*MysqlLoginDao) SelectById(login *bean.Login) (bool, error) {
	return db.Engine.ID(login.Id).Get(login)
}

func (*MysqlLoginDao) UpdatePassword(login *bean.Login) (int64,error) {
	return db.Engine.ID(login.Id).Cols("password").Update(login)
}
