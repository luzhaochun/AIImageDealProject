package controllers

import (
	"context"

	"backend/pkg/redisx"

	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
)

type HealthController struct {
	BaseController
}

func (c *HealthController) Get() {
	mysqlOK := false
	var one int
	if err := orm.NewOrm().Raw("SELECT 1").QueryRow(&one); err == nil && one == 1 {
		mysqlOK = true
	}
	redisOK := redisx.Ping(context.Background()) == nil
	addr, _ := beego.AppConfig.String("redis_addr")
	c.Success(map[string]interface{}{
		"status":  "ok",
		"service": "imagedeal-api",
		"mysql":   mysqlOK,
		"redis":   redisOK,
		"redis_addr": addr,
	})
}
