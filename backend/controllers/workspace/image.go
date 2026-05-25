package workspace

import (
	"encoding/json"

	"backend/controllers"
	"backend/models"
	"backend/services"
)

type ImageController struct {
	controllers.BaseController
}

func (c *ImageController) Claim() {
	var body struct {
		Category string `json:"category"`
	}
	_ = json.Unmarshal(c.Ctx.Input.RequestBody, &body)
	slug, _ := (&services.ImageService{}).GetTenantSlug(c.User().TenantId)
	img, err := (&services.ImageService{}).Claim(c.User().TenantId, c.User().UserId, body.Category, slug)
	if err != nil {
		c.Fail(err)
		return
	}
	c.Success(img)
}

func (c *ImageController) List() {
	page, size := c.PageQuery()
	lib := c.GetString("library", "tasks")
	svc := &services.ImageService{}
	var list []models.Image
	var total int64
	var err error
	switch lib {
	case "tasks":
		list, total, err = svc.ListTasks(c.User().TenantId, c.User().UserId, c.GetString("status"), size, (page-1)*size)
	case "pending_1st_review", "first_review":
		list, total, err = svc.ListByStatus(c.User().TenantId, models.StatusPending1stReview, size, (page-1)*size)
	default:
		list, total, err = svc.ListTasks(c.User().TenantId, c.User().UserId, "", size, (page-1)*size)
	}
	if err != nil {
		c.Fail(err)
		return
	}
	c.Success(controllers.PageData(list, total, page, size))
}

func (c *ImageController) Discard() {
	id := c.Ctx.Input.Param(":id")
	var body struct {
		Reason string `json:"reason"`
	}
	_ = json.Unmarshal(c.Ctx.Input.RequestBody, &body)
	if err := (&services.ImageService{}).Discard(c.User().TenantId, c.User().UserId, id, body.Reason); err != nil {
		c.Fail(err)
		return
	}
	c.Success(map[string]string{"ok": "true"})
}

func (c *ImageController) StartEditing() {
	id := c.Ctx.Input.Param(":id")
	img, err := (&services.ImageService{}).StartEditing(c.User().TenantId, c.User().UserId, id)
	if err != nil {
		c.Fail(err)
		return
	}
	c.Success(img)
}

func (c *ImageController) Get() {
	id := c.Ctx.Input.Param(":id")
	img, err := (&services.ImageService{}).GetDetail(c.User().TenantId, id)
	if err != nil {
		c.Fail(err)
		return
	}
	c.Success(img)
}

type EditorController struct {
	controllers.BaseController
}

func (c *EditorController) SaveVersion() {
	id := c.Ctx.Input.Param(":id")
	var body struct {
		LayerData interface{} `json:"layer_data"`
	}
	_ = json.Unmarshal(c.Ctx.Input.RequestBody, &body)
	b, _ := json.Marshal(body.LayerData)
	ver, err := (&services.EditorService{}).SaveVersion(c.User().TenantId, c.User().UserId, id, string(b))
	if err != nil {
		c.Fail(err)
		return
	}
	c.Success(ver)
}

func (c *EditorController) Submit() {
	id := c.Ctx.Input.Param(":id")
	img, err := (&services.EditorService{}).SubmitForReview(c.User().TenantId, c.User().UserId, id)
	if err != nil {
		c.Fail(err)
		return
	}
	c.Success(img)
}

func (c *EditorController) Versions() {
	id := c.Ctx.Input.Param(":id")
	list, err := (&services.EditorService{}).ListVersions(id)
	if err != nil {
		c.Fail(err)
		return
	}
	c.Success(list)
}

type CategoriesController struct {
	controllers.BaseController
}

func (c *CategoriesController) List() {
	list, err := (&services.ImageService{}).ListCategories(c.User().TenantId)
	if err != nil {
		c.Fail(err)
		return
	}
	c.Success(list)
}
