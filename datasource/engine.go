package datasource

import (
	"QianfengCmsProject/config"
	"QianfengCmsProject/models"
	_ "github.com/go-sql-driver/mysql" //不能忘记导入
	"github.com/go-xorm/xorm"
	"github.com/kataras/iris/v12"
)

func NewMysqlEngine() *xorm.Engine {
	//engine, err := xorm.NewEngine("mysql", "root:123456@/iris?charset=utf8")
	initConfig := config.InitConfig()
	if initConfig == nil{
		return nil
	}
	database := initConfig.DataBase

	dataSourceName := database.User + ":" + database.Pwd + "@tcp(" + database.Host + ")/" + database.Database + "?charset=utf8"

	engine,err := xorm.NewEngine(database.Drive,dataSourceName)

	iris.New().Logger().Info(database)

	err = engine.Sync2(
		/*new(models.Permission),*/
		/*new(models.City),*/
		//new(models.AdminPermission),
		new(models.Admin),
		new(models.User),
		new(models.Food),
		new(models.Shop),
		new(models.FoodCategory),
		//new(models.UserOrder),
	)

	if err != nil {
		panic(err.Error())
	}

	//设置显示sql语句
	engine.ShowSQL(true)
	engine.SetMaxOpenConns(10)

	return engine
}
