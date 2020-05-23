package controller

import (
	"QianfengCmsProject/service"
	"QianfengCmsProject/utils"
	"fmt"
	iris "github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
	"strings"
)

type StatisController struct {
	//iris框架自动为每个请求都绑定上下文对象：可作为接受参数
	Ctx iris.Context

	//admin功能实体：引入Service接口
	Service service.StatisService

	//session对象：存储session信息
	Session *sessions.Session
}

var (
	ADMINMODULE = "ADMIN_"
	USERMODULE  = "USER_"
	ORDERMODULE = "ORDER_"
)

/* /statis/admin/2019-03-10/count*/
func (sc *StatisController) GetCount() mvc.Result {
	path := sc.Ctx.Path()

	var pathSlice []string
	if path != "" {
		pathSlice = strings.Split(path, "/")
	}

	//不符合请求格式
	if len(pathSlice) != 5 {
		return mvc.Response{
			Object: map[string]interface{}{
				"status": utils.RECODE_FAIL,
				"count":  0,
			},
		}
	}

	//将最前面的去掉
	pathSlice = pathSlice[1:]
	model := pathSlice[1]
	date := pathSlice[2]
	var result int64
	switch model {
	case "user":
		fmt.Println("GetCount--->user")
	case "order":
		fmt.Println("order--->user")
	case "admin":
		adminStatis := sc.Session.Get(ADMINMODULE + date)
		if adminStatis != nil {
			adminStatis = adminStatis.(float64)
			return mvc.Response{
				Object: map[string]interface{}{
					"status": utils.RECODE_OK,
					"count":  adminStatis,
				},
			}
		} else {
			result = sc.Service.GetAdminDailyCount(date)
			sc.Session.Set(ADMINMODULE, result)
		}
	}

	return mvc.Response{
		Object: map[string]interface{}{
			"status": utils.RECODE_OK,
			"count":  result,
		},
	}
}

