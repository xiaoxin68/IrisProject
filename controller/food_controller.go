package controller

import (
	"QianfengCmsProject/models"
	"QianfengCmsProject/service"
	"QianfengCmsProject/utils"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"strconv"
)

/**
 * 食品模块控制器
 */
type FoodController struct {
	Ctx     iris.Context
	Service service.FoodService
}

type AddFoodEntity struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Activity    string `json:"activity"`
	ImagePath   string `json:"image_path"`
	CategoryId  string `json:"category_id"`
}

/**
 * url：foods/count?
 * type：Get
 * desc：获取所有的食品记录总数
 */
func (fc *FoodController) GetCount() mvc.Result {
	iris.New().Logger().Info(" 食品记录总数 ")
	result, err := fc.Service.GetFoodCount()

	if err != nil {
		return mvc.Response{
			Object: map[string]interface{}{
				"status": utils.RESPMSG_FAIL,
				"count":  0,
			},
		}
	}

	//查询成功
	return mvc.Response{
		Object: map[string]interface{}{
			"status": utils.RESPMSG_OK,
			"count":  result,
		},
	}
}

func (fc *FoodController) Get() mvc.Result {
	offset, err := strconv.Atoi(fc.Ctx.URLParam("offset"))
	limit, err := strconv.Atoi(fc.Ctx.URLParam("limit"))
	if err != nil {
		offset = 0
		limit = 20
	}
	list, err := fc.Service.GetFoodList(offset, limit)
	if err != nil {
		return mvc.Response{
			Object: map[string]interface{}{
				"status":  utils.RECODE_FAIL,
				"type":    utils.RESPMSG_ERROR_FOODLIST,
				"message": utils.Recode2Text(utils.RESPMSG_ERROR_FOODLIST),
			},
		}
	}

	return mvc.Response{
		Object: &list,
	}
}

func (fc *FoodController) PostAddfood() mvc.Result {
	var foodEntity AddFoodEntity

	err := fc.Ctx.ReadJSON(&foodEntity)

	if err != nil {
		return mvc.Response{
			Object: map[string]interface{}{
				"status":  utils.RECODE_FAIL,
				"type":    utils.RESPMSG_ERROR_FOODADD,
				"message": utils.Recode2Text(utils.RESPMSG_ERROR_FOODADD),
			},
		}
	}

	iris.New().Logger().Info(foodEntity)
	category_id, _ := strconv.Atoi(foodEntity.CategoryId)
	
	food := &models.Food{
		Name:        foodEntity.Name,
		Description: foodEntity.Description,
		ImagePath:   foodEntity.ImagePath,
		Activity:    foodEntity.Activity,
		CategoryId:  int64(category_id),
		DelFlag:     0,
		Rating:      0,
	}

	isSuccess := fc.Service.SaveFood(food)
	if !isSuccess {
		return mvc.Response{
			Object: map[string]interface{}{
				"status":  utils.RECODE_FAIL,
				"type":    utils.RESPMSG_ERROR_FOODADD,
				"message": utils.Recode2Text(utils.RESPMSG_ERROR_FOODADD),
			},
		}
	}

	return mvc.Response{
		Object: map[string]interface{}{
			"status":  utils.RECODE_OK,
			"success": utils.Recode2Text(utils.RESPMSG_SUCCESS_FOODADD),
		},
	}
}
