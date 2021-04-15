package logic

import (
	"context"
	"learn_together/commons"
	"learn_together/dao"
	"learn_together/models/bean"
	"learn_together/models/vo"
	"learn_together/service/cache"
)

type LoginService interface {
	Register(registerReq *vo.RegisterReq) (int64, error)
	Login(login *bean.Login) (bool, error)
}

type LoginServiceImpl struct {
}

func (*LoginServiceImpl) Register(registerReq *vo.RegisterReq) (int64, error) {
	//判断验证码
	value, err := cache.Rdb.Get(context.Background(), registerReq.Token).Result()
	if err != nil || value == "" {
		return -1, err
	}
	loginDao := &dao.MysqlLoginDao{}
	login := &bean.Login{}
	_ = commons.StructToStruct(registerReq, login, "json")
	return loginDao.AddLogin(login)
}
func (*LoginServiceImpl) Login(login *bean.Login) (bool, error) {
	loginDao := &dao.MysqlLoginDao{}
	return loginDao.Login(login)
}

func (*LoginServiceImpl) CheckPassword(id int, password string) (bool, error) {
	loginDao := &dao.MysqlLoginDao{}
	login := &bean.Login{
		Id: id,
	}
	has, err := loginDao.SelectById(login)
	if !has || err != nil {
		return has, err
	}
	if password != login.Password {
		return false, nil
	}
	return true, nil
}

func (l *LoginServiceImpl) UpdatePassword(id int, password string, newPassword string) (bool, error) {
	has, err := l.CheckPassword(id, password)
	if !has || err != nil {
		return has, err
	}
	loginDao := &dao.MysqlLoginDao{}
	login := &bean.Login{
		Id:       id,
		Password: newPassword,
	}
	_, err = loginDao.UpdatePassword(login)
	if err != nil {
		return false, err
	}
	return true, nil
}


