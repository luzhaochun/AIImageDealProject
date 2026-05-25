package auth

import (
	"encoding/json"

	"backend/controllers"
	"backend/services"
)

type Controller struct {
	controllers.BaseController
}

func (c *Controller) Login() {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	_ = json.Unmarshal(c.Ctx.Input.RequestBody, &body)
	email, password := body.Email, body.Password
	res, err := (&services.AuthService{}).Login(email, password)
	if err != nil {
		c.Fail(err)
		return
	}
	c.Success(res)
}

func (c *Controller) Me() {
	u, err := (&services.AuthService{}).GetUserProfile(c.User().UserId, c.User().TenantId)
	if err != nil {
		c.Fail(err)
		return
	}
	c.Success(u)
}
