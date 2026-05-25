package models

import "github.com/beego/beego/v2/client/orm"

func init() {
	orm.RegisterModel(
		new(Tenant),
		new(User),
		new(ImportBatch),
		new(Image),
		new(ImageVersion),
		new(ReviewRecord),
	)
}
