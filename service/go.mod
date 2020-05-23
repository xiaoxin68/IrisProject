module QianfengCmsProject/service

go 1.14

replace (
    QianfengCmsProject/models => ../models
    QianfengCmsProject/utils => ../utils
)

require (
	QianfengCmsProject/models v0.0.0-00010101000000-000000000000 // indirect
	github.com/go-xorm/xorm v0.7.9 // indirect
)
