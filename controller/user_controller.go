package controller

import (
	"QianfengCmsProject/service"
	"QianfengCmsProject/utils"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
	"strconv"
)

//每一页最大的内容
const MaxLimit = 50

/**
 * 用户控制器结构体：用来实现处理用户模块的接口的请求，并返回给客户端
 */
type UserController struct {
	//iris框架自动为每个请求都绑定上下文对象：可作为接受参数
	Ctx iris.Context

	//admin功能实体：引入Service接口
	Service service.UserService

	//session对象：存储session信息
	Session *sessions.Session
}

/**
 * 获取用户总数
 * 请求类型：Get
 * 请求Url：/v1/users/list?offset=0&limit=20
 */
func (uc *UserController) GetList() mvc.Result {
	offsetStr := uc.Ctx.FormValue("offset")
	limitStr := uc.Ctx.FormValue("limit")

	var offset int
	var limit int

	//判断offset和limit两个变量任意一个都不能为""
	if offsetStr == "" || limitStr == "" {
		return mvc.Response{
			Object: map[string]interface{}{
				"status":  utils.RECODE_FAIL,
				"type":    utils.RESPMSG_ERROR_USERLIST,
				"message": utils.Recode2Text(utils.RESPMSG_ERROR_USERLIST),
			},
		}
	}

	//将offset和limit转化成int64类型
	offset, err := strconv.Atoi(offsetStr)
	limit, err = strconv.Atoi(limitStr)
	if err != nil {
		return mvc.Response{
			Object: map[string]interface{}{
				"status":  utils.RECODE_FAIL,
				"type":    utils.RESPMSG_ERROR_USERLIST,
				"message": utils.Recode2Text(utils.RESPMSG_ERROR_USERLIST),
			},
		}
	}

	//做页数的检出
	if offset <= 0 {
		offset = 0
	}

	//做最大的限制
	if limit > MaxLimit {
		limit = MaxLimit
	}

	//调用service中提供的方法，分页查询
	userList := uc.Service.GetUserList(offset, limit)

	if len(userList) == 0 {
		return mvc.Response{
			Object: map[string]interface{}{
				"status":  utils.RECODE_FAIL,
				"type":    utils.RESPMSG_ERROR_USERLIST,
				"message": utils.Recode2Text(utils.RESPMSG_ERROR_USERLIST),
			},
		}
	}

	//将查询到的用户数据进行转换成前端需要的内容
	var respList []interface{}
	for _, user := range userList {
		respList = append(respList, user.UserToRespDesc())
	}

	//返回用户列表
	return mvc.Response{
		Object: &respList,
	}
}
