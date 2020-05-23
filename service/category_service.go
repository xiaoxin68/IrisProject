package service

import (
	"QianfengCmsProject/models"
	"github.com/go-xorm/xorm"
	"github.com/kataras/iris/v12"
)

type CategoryService interface {
	GetCategoryByShopId(shopId int64) ([]models.FoodCategory, error)
	AddCategory(category *models.FoodCategory) bool
	DeleteShop(restaurantId int) bool
}

func NewCategoryService(db *xorm.Engine) CategoryService {
	return &categorySevice{
		engine: db,
	}
}

type categorySevice struct {
	engine *xorm.Engine
}

func (cs *categorySevice) DeleteShop(restaurantId int) bool {
	shop := models.Shop{ShopId: restaurantId, Dele: 1}
	_, err := cs.engine.Where("shop_id=?", restaurantId).Cols("dele").Update(&shop)
	if err != nil {
		iris.New().Logger().Info(err.Error())
	}
	return err == nil
}

func (cs *categorySevice) AddCategory(category *models.FoodCategory) bool {
	iris.New().Logger().Info(category)

	_, err := cs.engine.Insert(category)

	if err != nil {
		iris.New().Logger().Info(err.Error())
	}
	return err == nil
}

func (cs *categorySevice) GetCategoryByShopId(shopId int64) ([]models.FoodCategory, error) {
	categories := make([]models.FoodCategory, 0)
	err := cs.engine.Where("restaurant_id = ?", shopId).Find(&categories)
	if err != nil {
		iris.New().Logger().Info(err.Error())
	}
	iris.New().Logger().Info(categories)
	return categories, err
}
