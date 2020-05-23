module QianfengCmsProject/datasource

go 1.14

replace (
	QianfengCmsProject/config => ../config
	QianfengCmsProject/models => ../models
	QianfengCmsProject/utils => ../utils
)

require (
	QianfengCmsProject/config v0.0.0-00010101000000-000000000000
	QianfengCmsProject/models v0.0.0-00010101000000-000000000000
	github.com/go-sql-driver/mysql v1.5.0
	github.com/go-xorm/xorm v0.7.9
	github.com/kataras/iris/v12 v12.1.8 // indirect
)
