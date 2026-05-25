package services

import (
	"backend/models"
	"backend/pkg/apperr"

	"github.com/beego/beego/v2/client/orm"
)

type ArchiveService struct{}

func (s *ArchiveService) Complete(tenantId, imageId string) (*models.Image, error) {
	o := orm.NewOrm()
	var img models.Image
	if err := o.QueryTable(new(models.Image)).Filter("tenant_id", tenantId).Filter("id", imageId).One(&img); err != nil {
		return nil, apperr.ErrNotFound
	}
	from := img.Status
	if from != models.Status2ndReviewPassed && from != models.Status1stReviewPassed {
		return nil, apperr.ErrInvalidTransition
	}
	if err := MustTransition(from, models.StatusCompleted); err != nil {
		// allow 1st passed -> completed when no 2nd review
		if from == models.Status1stReviewPassed {
			img.Status = models.StatusCompleted
			_, err := o.Update(&img, "Status")
			return &img, err
		}
		return nil, err
	}
	img.Status = models.StatusCompleted
	_, err := o.Update(&img, "Status")
	return &img, err
}

func (s *ArchiveService) ListArchive(tenantId, category string, limit, offset int) ([]models.Image, int64, error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(models.Image)).Filter("tenant_id", tenantId).Filter("status__in", models.Status2ndReviewPassed, models.StatusCompleted)
	if category != "" {
		qs = qs.Filter("category", category)
	}
	total, _ := qs.Count()
	var list []models.Image
	_, err := qs.OrderBy("-updated_at").Limit(limit, offset).All(&list)
	return list, total, err
}
