module QianfengCmsProject/controller

go 1.14

require (
	QianfengCmsProject/models v0.0.0-00010101000000-000000000000 // indirect
	QianfengCmsProject/service v0.0.0-00010101000000-000000000000 // indirect
	QianfengCmsProject/utils v0.0.0-00010101000000-000000000000 // indirect
    github.com/kataras/iris/v12 v12.1.8 // indirect
)

replace (
	QianfengCmsProject/models => ../models
	QianfengCmsProject/service => ../service
	QianfengCmsProject/utils => ../utils
)
