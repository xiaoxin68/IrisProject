package service

import (
	"QianfengCmsProject/models"
	"github.com/go-xorm/xorm"
	"github.com/kataras/iris/v12"
)

//用户功能模块接口标准
type UserService interface {
	//用户列表
	GetUserList(offset int,limit int)[]*models.User
}

//用户功能服务实现结构体
type userService struct {
	Engin *xorm.Engine
}

func NewUserService(engin *xorm.Engine) UserService {
	return &userService{
		Engin: engin,
	}
}

//用户列表
func (uc *userService)GetUserList(offset int,limit int)[]*models.User  {
	var userList []*models.User
	err := uc.Engin.Where("del_flag=?",0).Limit(limit,offset).Find(&userList)
	if err != nil {
		iris.New().Logger().Error(err.Error())
		panic(err.Error())
		return nil
	}
	return userList
}