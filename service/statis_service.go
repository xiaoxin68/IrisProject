package service

import (
	"QianfengCmsProject/models"
	"fmt"
	"github.com/go-xorm/xorm"
	"time"
)

//统计功能模块接口标准
type StatisService interface {
	GetAdminDailyCount(date string) int64
}

//统计功能服务实现结构体
type statisService struct {
	Engin *xorm.Engine
}

func NewStatisService(engin *xorm.Engine) StatisService {
	return &statisService{
		Engin: engin,
	}
}

func (ss *statisService) GetAdminDailyCount(date string) int64 {
	if date == "NaN-NaN-NaN" { //当日增长数据请求
		date = time.Now().Format("2006-01-02")
	}

	//查询如期data格式解析
	startDate, err := time.Parse("2006-01-02 15:04:05", date+" 00:00:00")
	if err != nil {
		return 0
	}

	endDate := startDate.AddDate(0, 0, 1)
	result, err := ss.Engin.Where("create_time between ? and ? and status = 0",
		startDate.Format("2006-01-02 15:04:05"),
		endDate.Format("2006-01-02 15:04:05")).Count(models.Admin{})
	if err != nil {
		return 0
	}
	fmt.Println(result)
	return result
}
