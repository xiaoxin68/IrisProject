package service

import (
	"QianfengCmsProject/models"
	"github.com/go-xorm/xorm"
	"github.com/kataras/iris/v12"
)

type FoodService interface {
	GetFoodCount() (int64, error)
	GetFoodList(offset int, limit int) ([]models.Food, error)
	SaveFood(food *models.Food) bool

}

func NewFoodService(db *xorm.Engine) FoodService {
	return &foodService{
		engine: db,
	}
}

type foodService struct {
	engine *xorm.Engine
}

func (fs *foodService) SaveFood(food *models.Food) bool {
	iris.New().Logger().Info(food)

	_, err := fs.engine.Insert(food)

	if err != nil {
		iris.New().Logger().Info(err.Error())
	}
	return err == nil
}

func (fs *foodService) GetFoodList(offset int, limit int) ([]models.Food, error) {
	foodList := make([]models.Food,0)
	err := fs.engine.Where(" del_flag = 0 ").Limit(limit,offset).Find(&foodList)
	return foodList,err
}

/**
 * 获取食品总记录数
 */
func (fs *foodService) GetFoodCount() (int64, error) {
	count, err := fs.engine.Where(" del_flag = 0 ").Count(new(models.Food))
	return count, err
}


