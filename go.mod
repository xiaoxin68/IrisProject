module QianfengCmsProject

go 1.14

require (
	QianfengCmsProject/controller v0.0.0-00010101000000-000000000000 // indirect
	QianfengCmsProject/datasource v0.0.0-00010101000000-000000000000 // indirect
	github.com/kataras/iris/v12 v12.1.8 // indirect
	github.com/shurcooL/sanitized_anchor_name v1.0.0 // indirect
)

replace (
	QianfengCmsProject/config => ./config
	QianfengCmsProject/controller => ./controller
	QianfengCmsProject/datasource => ./datasource
	QianfengCmsProject/models => ./models
	QianfengCmsProject/service => ./service
	QianfengCmsProject/utils => ./utils
)
