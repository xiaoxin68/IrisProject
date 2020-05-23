package service

import (
	"QianfengCmsProject/models"
	"github.com/go-xorm/xorm"
)

type ShopService interface {
	GetShopList(offset int, limit int) []models.Shop
}

func NewShopService(db *xorm.Engine) ShopService {
	return &shopSevice{
		engine: db,
	}
}

type shopSevice struct {
	engine *xorm.Engine
}

func (ss *shopSevice) GetShopList(offset int, limit int) []models.Shop {
	shopList := make([]models.Shop,0)
	ss.engine.Where(" dele = 0").Limit(limit,offset).Find(&shopList)
	return shopList
}
