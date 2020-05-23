package controller

import (
	"QianfengCmsProject/models"
	"QianfengCmsProject/service"
	"QianfengCmsProject/utils"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"strconv"
)

type CategoryController struct {
	Ctx     iris.Context
	Service service.CategoryService
}

/**
 * 添加食品种类实体
 */
type CategoryEntity struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	RestaurantId string `json:"restaurant_id"`
}

func (cc *CategoryController) BeforeActivation(a mvc.BeforeActivation) {
	a.Handle("GET", "/getcategory/{shopId}", "GetCategoryByShopId")

	//删除商铺记录
	a.Handle("DELETE", "/restaurant/{restaurant_id}", "DeleteRestaurant")
}

//返回对应商铺的食品种类
func (cc *CategoryController) GetCategoryByShopId() mvc.Result {
	shopIdStr := cc.Ctx.Params().Get("shopId")
	if shopIdStr == "" {
		iris.New().Logger().Info(shopIdStr)
		return mvc.Response{
			Object: map[string]interface{}{
				"status":  utils.RECODE_FAIL,
				"message": utils.Recode2Text(utils.RESPMSG_ERROR_CATEGORIES),
			},
		}
	}

	shopId, err := strconv.Atoi(shopIdStr)
	if err != nil {
		iris.New().Logger().Info(shopId)
		return mvc.Response{
			Object: map[string]interface{}{
				"status":  utils.RECODE_FAIL,
				"message": utils.Recode2Text(utils.RESPMSG_ERROR_CATEGORIES),
			},
		}
	}

	///调用服务实体功能类查询商铺对应的食品种类列表
	categories, err := cc.Service.GetCategoryByShopId(int64(shopId))
	if err != nil {
		iris.New().Logger().Info(err.Error())
		return mvc.Response{
			Object: map[string]interface{}{
				"status":  utils.RECODE_FAIL,
				"message": utils.Recode2Text(utils.RESPMSG_ERROR_CATEGORIES),
			},
		}
	}

	//返回对应店铺的食品种类类型
	return mvc.Response{
		Object: map[string]interface{}{
			"status":        utils.RECODE_OK,
			"category_list": &categories,
		},
	}
}

func (cc *CategoryController) PostAddcategory() mvc.Result {

	var categoryEntity *CategoryEntity
	cc.Ctx.ReadJSON(&categoryEntity)

	if categoryEntity.Name == "" {
		return mvc.Response{
			Object: map[string]interface{}{
				"status":  utils.RECODE_FAIL,
				"message": utils.Recode2Text(utils.RESPMSG_ERROR_CATEGORYADD),
			},
		}
	}

	iris.New().Logger().Info(categoryEntity)
	restaurant_id, _ := strconv.Atoi(categoryEntity.RestaurantId)
	//构造要添加的数据记录
	foodCategory := &models.FoodCategory{
		CategoryName:     categoryEntity.Name,
		CategoryDesc:     categoryEntity.Description,
		Level:            1,
		ParentCategoryId: 0,
		RestaurantId:     int64(restaurant_id),
	}

	saveSucc := cc.Service.AddCategory(foodCategory)
	if !saveSucc {
		return mvc.Response{
			Object: map[string]interface{}{
				"status":  utils.RECODE_FAIL,
				"message": utils.Recode2Text(utils.RESPMSG_ERROR_CATEGORYADD),
			},
		}
	}

	//成功
	return mvc.Response{
		Object: map[string]interface{}{
			"status":  utils.RECODE_OK,
			"message": utils.Recode2Text(utils.RESPMSG_SUCCESS_CATEGORYADD),
		},
	}
}

func (cc *CategoryController) DeleteRestaurant() mvc.Result {
	restaurant_id := cc.Ctx.Params().Get(("restaurant_id"))

	shopId, err := strconv.Atoi(restaurant_id)

	if err != nil {
		return mvc.Response{
			Object: map[string]interface{}{
				"status":  utils.RECODE_FAIL,
				"type":    utils.RESPMSG_HASNOACCESS,
				"message": utils.Recode2Text(utils.RESPMSG_HASNOACCESS),
			},
		}
	}

	delete := cc.Service.DeleteShop(shopId)
	if !delete {
		return mvc.Response{
			Object: map[string]interface{}{
				"status":  utils.RECODE_FAIL,
				"type":    utils.RESPMSG_HASNOACCESS,
				"message": utils.Recode2Text(utils.RESPMSG_HASNOACCESS),
			},
		}
	} else {
		return mvc.Response{
			Object: map[string]interface{}{
				"status":  utils.RECODE_OK,
				"type":    utils.RESPMSG_SUCCESS_DELETESHOP,
				"message": utils.Recode2Text(utils.RESPMSG_SUCCESS_DELETESHOP),
			},
		}
	}
}
