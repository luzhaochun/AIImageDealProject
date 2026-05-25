package filters

import (
	"strings"

	"backend/pkg/ctxutil"
	"backend/services"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
)

func JWTAuth(ctx *context.Context) {
	path := ctx.Input.URL()
	if strings.HasPrefix(path, "/api/v1/auth/login") || strings.HasPrefix(path, "/api/v1/health") || strings.HasPrefix(path, "/static/") {
		return
	}
	auth := ctx.Input.Header("Authorization")
	if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
		ctx.Output.SetStatus(401)
		_ = ctx.Output.JSON(map[string]string{"code": "UNAUTHORIZED", "message": "未登录"}, false, false)
		return
	}
	token := strings.TrimPrefix(auth, "Bearer ")
	claims, err := services.ParseToken(token)
	if err != nil {
		ctx.Output.SetStatus(401)
		_ = ctx.Output.JSON(map[string]string{"code": "UNAUTHORIZED", "message": "token 无效"}, false, false)
		return
	}
	ctx.Input.SetData(ctxutil.BeegoUserKey, &ctxutil.UserContext{
		UserId:   claims.UserId,
		TenantId: claims.TenantId,
		Role:     claims.Role,
		Email:    claims.Email,
	})
}

func RequireRole(roles ...string) beego.FilterFunc {
	return func(ctx *context.Context) {
		uc := GetUser(ctx)
		if uc == nil {
			ctx.Output.SetStatus(401)
			return
		}
		for _, r := range roles {
			if uc.Role == r {
				return
			}
		}
		ctx.Output.SetStatus(403)
		_ = ctx.Output.JSON(map[string]string{"code": "FORBIDDEN", "message": "无权限"}, false, false)
	}
}

func GetUser(ctx *context.Context) *ctxutil.UserContext {
	v := ctx.Input.GetData(ctxutil.BeegoUserKey)
	if v == nil {
		return nil
	}
	uc, _ := v.(*ctxutil.UserContext)
	return uc
}

func init() {
	beego.InsertFilter("/api/v1/*", beego.BeforeRouter, JWTAuth)
}
