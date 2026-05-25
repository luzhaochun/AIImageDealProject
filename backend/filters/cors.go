package filters

import (
	"net/http"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
)

// CORS 开发环境跨域（骨架）
func CORS(ctx *context.Context) {
	ctx.Output.Header("Access-Control-Allow-Origin", "*")
	ctx.Output.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
	ctx.Output.Header("Access-Control-Allow-Headers", "Origin,Content-Type,Authorization")
	if ctx.Input.Method() == http.MethodOptions {
		ctx.Output.SetStatus(http.StatusNoContent)
	}
}

func init() {
	beego.InsertFilter("*", beego.BeforeRouter, CORS)
}
