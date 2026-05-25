package ctxutil

type UserContext struct {
	UserId   string
	TenantId string
	Role     string
	Email    string
}

type ctxKey struct{}

// Set on beego context via Data map
const BeegoUserKey = "user_ctx"
