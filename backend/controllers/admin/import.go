package admin

import (
	"backend/controllers"
	"backend/services"
)

type ImportController struct {
	controllers.BaseController
}

func (c *ImportController) Post() {
	file, header, err := c.GetFile("file")
	if err != nil {
		c.Fail(err)
		return
	}
	defer file.Close()
	batch, err := (&services.ImportService{}).CreateBatch(c.User().TenantId, c.User().UserId, header.Filename, file)
	if err != nil {
		c.Fail(err)
		return
	}
	c.Success(batch)
}

func (c *ImportController) Get() {
	page, size := c.PageQuery()
	if c.Ctx.Input.Query("id") != "" {
		b, err := (&services.ImportService{}).GetBatchDetail(c.User().TenantId, c.Ctx.Input.Query("id"))
		if err != nil {
			c.Fail(err)
			return
		}
		c.Success(b)
		return
	}
	list, total, err := (&services.ImportService{}).ListBatches(c.User().TenantId, size, (page-1)*size)
	if err != nil {
		c.Fail(err)
		return
	}
	c.Success(controllers.PageData(list, total, page, size))
}
