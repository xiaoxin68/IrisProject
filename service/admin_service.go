package service

import (
	"QianfengCmsProject/models"
	"github.com/go-xorm/xorm"
)

//定义AdminService接口
type AdminService interface {
	//通过管理员用户名+密码 获取管理员实体 如果查询到，返回管理员实体，并返回true
	//否则 返回 nil ，false
	GetByAdminNameAndPassword(username, password string) (models.Admin, bool)
	GetAdminCount() (int64, error)
	SaveAvatarImg(adminId int64, fileName string) bool
}

//在我们实际的开发过程中，我们往往将数据提供服务模块设计成接口，这样设计的目的是接口定义和具体
//的功能编程实现了分离，有助于我们在不同的实现方案之间进行切换，成本非常小
func NewAdminService(db *xorm.Engine) AdminService {
	return &adminSevice{
		engine: db,
	}
}

//实现类
type adminSevice struct {
	engine *xorm.Engine
}

//保存头像信息
func (ac *adminSevice) SaveAvatarImg(adminId int64, fileName string) bool {
	admin := models.Admin{Avatar: fileName}
	_, err := ac.engine.ID(adminId).Cols("avatar").Update(&admin)
	return err != nil
}

func (ac *adminSevice) GetAdminCount() (int64, error) {
	count, err := ac.engine.Count(new(models.Admin))

	if err != nil {
		panic(err.Error())
		return 0, err
	}
	return count, nil
}

/**
 * 通过用户名和密码查询管理员
 */
func (ac *adminSevice) GetByAdminNameAndPassword(username, password string) (models.Admin, bool) {
	var admin models.Admin

	ac.engine.Where("admin_name = ? and pwd = ? ", username, password).Get(&admin)

	//fmt.Println(admin,"............",admin.AdminId != 0)

	return admin, admin.AdminId != 0
}
