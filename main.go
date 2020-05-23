package main

import (
	"QianfengCmsProject/config"
	"QianfengCmsProject/controller"
	"QianfengCmsProject/datasource"
	"QianfengCmsProject/models"
	"QianfengCmsProject/service"
	"QianfengCmsProject/utils"
	"encoding/json"
	iris "github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	app := newApp()

	//应用App设置
	configation(app)

	//路由设置
	mvcHandle(app)

	config := config.InitConfig()
	addr := "localhost:" + config.Port
	app.Run(
		iris.Addr(addr),                               //在端口8080进行监听
		iris.WithoutServerError(iris.ErrServerClosed), //无服务错误提示
		iris.WithOptimizations,                        //对json数据序列化更快的配置
	)
}

//构建App
func newApp() *iris.Application {
	app := iris.New()

	//设置日志级别  开发阶段为debug
	app.Logger().SetLevel("debug")

	//注册静态资源
	app.HandleDir("/static", "./static")
	app.HandleDir("/manage/static", "./static")
	//访问：http://localhost:9000/img/a.jpg  ==>  将会把/img/a.jpg映射成为/static/img/a.jpg
	app.HandleDir("/img", "./static/img")

	//注册视图文件
	app.RegisterView(iris.HTML("./static", ".html"))

	app.Get("/", func(context context.Context) {
		context.View("index.html")
	})

	return app
}

/**
 * 项目设置
 */
func configation(app *iris.Application) {

	//配置 字符编码
	app.Configure(iris.WithConfiguration(iris.Configuration{
		Charset: "UTF-8",
	}))

	//错误配置
	//未发现错误
	app.OnErrorCode(iris.StatusNotFound, func(context context.Context) {
		context.JSON(iris.Map{
			"errmsg": iris.StatusNotFound,
			"msg":    " not found ",
			"data":   iris.Map{},
		})
	})

	app.OnErrorCode(iris.StatusInternalServerError, func(context context.Context) {
		context.JSON(iris.Map{
			"errmsg": iris.StatusInternalServerError,
			"msg":    " interal error ",
			"data":   iris.Map{},
		})
	})
}

//MVC 架构模式处理
func mvcHandle(app *iris.Application) {
	//启用session
	sessManager := sessions.New(sessions.Config{
		Cookie:  "sessioncookie",
		Expires: 24 * time.Hour,
	})

	//获取redis实例
	redis := datasource.NewRedis()
	//设置session的同步位置为redis
	sessManager.UseDatabase(redis)

	//实例化mysql数据库引擎
	engine := datasource.NewMysqlEngine()

	//管理员模块功能
	adminService := service.NewAdminService(engine)

	admin := mvc.New(app.Party("/admin")) //设置路由组
	admin.Register(
		adminService,
		sessManager.Start,
	)
	//通过mvc的Handle方法进行控制器的指定
	admin.Handle(new(controller.AdminController))

	//统计功能模块
	statisService := service.NewStatisService(engine)
	statis := mvc.New(app.Party("/statis/{model}/{date}/"))
	statis.Register(
		statisService,
		sessManager.Start,
	)
	statis.Handle(new(controller.StatisController))

	//用户模块
	userService := service.NewUserService(engine)
	user := mvc.New(app.Party("/v1/users"))
	user.Register(
		userService,
		sessManager.Start,
	)
	user.Handle(new(controller.UserController))

	//商铺模块
	shopService := service.NewShopService(engine)
	shop := mvc.New(app.Party("/shopping/restaurants/"))
	shop.Register(
		shopService,
		sessManager.Start,
	)
	shop.Handle(new(controller.ShopController))

	//类别模块
	categoryService := service.NewCategoryService(engine)
	category := mvc.New(app.Party("/shopping/"))
	category.Register(
		categoryService,
	)
	category.Handle(new(controller.CategoryController)) //控制器

	//food模块
	//食品模块
	foodService := service.NewFoodService(engine)
	foodMvc := mvc.New(app.Party("/shopping/v2/foods/"))
	foodMvc.Register(
		foodService,
	)
	foodMvc.Handle(new(controller.FoodController)) //控制器

	//地址pois检索
	app.Get("/v1/pois", func(context context.Context) {
		path := context.Request().URL.String()
		iris.New().Logger().Info(path)

		rs, err := http.Get("https://elm.cangdu.org" + path)
		if err != nil {
			context.JSON(iris.Map{
				"status":  utils.RECODE_FAIL,
				"type":    utils.RESPMSG_ERROR_SEARCHADDRESS,
				"message": utils.Recode2Text(utils.RESPMSG_ERROR_SEARCHADDRESS),
			})
			return
		}

		//请求成功
		body, err := ioutil.ReadAll(rs.Body)
		var searchList []*models.PoiSearch
		//安马歇尔 马歇尔
		json.Unmarshal(body, &searchList)
		context.JSON(&searchList)
	})

	//文件上传
	//请求url：/admin/update/avatar/1
	//Content-Disposition: form-data; name="file"; filename="bg.jpeg"
	//Content-Type: image/jpeg
	app.Post("/admin/upload/avatar/{adminId}", func(context context.Context) {
		adminId := context.Params().Get("adminId")
		iris.New().Logger().Info(adminId)
		//获取上传的文件信息
		file, info, err := context.FormFile("file")
		if err != nil {
			iris.New().Logger().Info(err.Error())
			context.JSON(iris.Map{
				"status":  utils.RECODE_FAIL,
				"type":    utils.RESPMSG_ERROR_PICTUREADD,
				"failure": utils.Recode2Text(utils.RESPMSG_ERROR_PICTUREADD),
			})
			return
		}
		defer file.Close()
		//获取文件名
		fname := info.Filename
		//创建copy的文件
		out, err := os.OpenFile("./uploads/"+fname, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			iris.New().Logger().Info("文件路径：" + err.Error())
			context.JSON(iris.Map{
				"status":  utils.RECODE_FAIL,
				"type":    utils.RESPMSG_ERROR_PICTUREADD,
				"failure": utils.Recode2Text(utils.RESPMSG_ERROR_PICTUREADD),
			})
			return
		}
		defer out.Close()
		_, err = io.Copy(out, file) //复制文件到指定文件中
		if err != nil {
			context.JSON(iris.Map{
				"status":  utils.RECODE_FAIL,
				"type":    utils.RESPMSG_ERROR_PICTUREADD,
				"failure": utils.Recode2Text(utils.RESPMSG_ERROR_PICTUREADD),
			})
			return
		}
		//根据adminId更新它的avatar信息
		intAdminId, _ := strconv.Atoi(adminId)
		adminService.SaveAvatarImg(int64(intAdminId), fname)
		context.JSON(iris.Map{
			"status":     utils.RECODE_OK,
			"image_path": fname,
		})
	})

	//图片上传
	app.Post("/v1/adding/{model}", func(context context.Context) {
		model := context.Params().Get("model")
		iris.New().Logger().Info(model)

		file, info, err := context.FormFile("file")
		if err != nil {
			iris.New().Logger().Info(err.Error())
			context.JSON(iris.Map{
				"status":  utils.RECODE_FAIL,
				"type":    utils.RESPMSG_ERROR_PICTUREADD,
				"faliure": utils.Recode2Text(utils.RESPMSG_ERROR_PICTUREADD),
			})
			return
		}
		defer file.Close()
		fname := info.Filename

		//判断上传目录是否存在
		isExist, err := utils.PathExists("./uploads/" + model)

		if !isExist {
			err = os.Mkdir("./uploads/"+model, 0777)
			if err != nil {
				context.JSON(iris.Map{
					"status":  utils.RECODE_FAIL,
					"type":    utils.RESPMSG_ERROR_PICTUREADD,
					"failure": utils.Recode2Text(utils.RESPMSG_ERROR_PICTUREADD),
				})
				return
			}
		}

		out, err := os.OpenFile("./uploads/"+model+"/"+fname, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			iris.New().Logger().Info(err.Error())
			context.JSON(iris.Map{
				"status":  utils.RECODE_FAIL,
				"type":    utils.RESPMSG_ERROR_PICTUREADD,
				"failure": utils.Recode2Text(utils.RESPMSG_ERROR_PICTUREADD),
			})
			return
		}

		iris.New().Logger().Info("文件路径：" + out.Name())
		defer out.Close()

		_, err = io.Copy(out, file)
		if err != nil {
			context.JSON(iris.Map{
				"status":  utils.RECODE_FAIL,
				"type":    utils.RESPMSG_ERROR_PICTUREADD,
				"failure": utils.Recode2Text(utils.RESPMSG_ERROR_PICTUREADD),
			})
			return
		}

		//上传成功
		context.JSON(iris.Map{
			"status":     utils.RECODE_OK,
			"image_path": fname,
		})
	})
}
