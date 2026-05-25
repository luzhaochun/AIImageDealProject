package review

import (
	"encoding/json"

	"backend/controllers"
	"backend/models"
	"backend/services"
)

type Controller struct {
	controllers.BaseController
}

func (c *Controller) Queue() {
	page, size := c.PageQuery()
	round, _ := c.GetInt("round", 1)
	status := models.StatusPending1stReview
	if round == 2 {
		status = models.StatusPending2ndReview
	}
	list, total, err := (&services.ImageService{}).ListByStatus(c.User().TenantId, status, size, (page-1)*size)
	if err != nil {
		c.Fail(err)
		return
	}
	c.Success(controllers.PageData(list, total, page, size))
}

func (c *Controller) Submit() {
	id := c.Ctx.Input.Param(":id")
	var req services.SubmitReviewReq
	_ = json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	img, err := (&services.ReviewService{}).Submit(c.User().TenantId, c.User().UserId, id, req)
	if err != nil {
		c.Fail(err)
		return
	}
	c.Success(img)
}

func (c *Controller) ClaimSecond() {
	img, err := (&services.ImageService{}).ClaimSecondReview(c.User().TenantId, c.User().UserId)
	if err != nil {
		c.Fail(err)
		return
	}
	c.Success(img)
}

func (c *Controller) Stats() {
	round, _ := c.GetInt("round", 1)
	st, err := (&services.ReviewService{}).Stats(c.User().TenantId, round)
	if err != nil {
		c.Fail(err)
		return
	}
	c.Success(st)
}

func (c *Controller) SecondPool() {
	page, size := c.PageQuery()
	list, total, err := (&services.ImageService{}).ListByStatus(
		c.User().TenantId, models.Status1stReviewPassed, size, (page-1)*size,
	)
	if err != nil {
		c.Fail(err)
		return
	}
	c.Success(controllers.PageData(list, total, page, size))
}

func (c *Controller) ListImageReviews() {
	id := c.Ctx.Input.Param(":id")
	list, err := (&services.ReviewService{}).ListReviews(id)
	if err != nil {
		c.Fail(err)
		return
	}
	c.Success(list)
}

type ArchiveController struct {
	controllers.BaseController
}

func (c *ArchiveController) List() {
	page, size := c.PageQuery()
	list, total, err := (&services.ArchiveService{}).ListArchive(c.User().TenantId, c.GetString("category"), size, (page-1)*size)
	if err != nil {
		c.Fail(err)
		return
	}
	c.Success(controllers.PageData(list, total, page, size))
}

func (c *ArchiveController) Complete() {
	id := c.Ctx.Input.Param(":id")
	img, err := (&services.ArchiveService{}).Complete(c.User().TenantId, id)
	if err != nil {
		c.Fail(err)
		return
	}
	c.Success(img)
}

func (c *ArchiveController) Export() {
	var req services.ExportRequest
	_ = json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	job, err := (&services.ArchiveService{}).CreateExport(c.User().TenantId, req)
	if err != nil {
		c.Fail(err)
		return
	}
	c.Success(job)
}

func (c *ArchiveController) ExportList() {
	list, err := (&services.ArchiveService{}).ListExports(c.User().TenantId)
	if err != nil {
		c.Fail(err)
		return
	}
	c.Success(list)
}
