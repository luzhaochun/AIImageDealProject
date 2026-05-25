package services

import (
	"backend/models"
	"backend/pkg/apperr"

	"github.com/beego/beego/v2/client/orm"
	"github.com/google/uuid"
)

type EditorService struct{}

func (s *EditorService) SaveVersion(tenantId, userId, imageId, layerData string) (*models.ImageVersion, error) {
	o := orm.NewOrm()
	var img models.Image
	if err := o.QueryTable(new(models.Image)).Filter("tenant_id", tenantId).Filter("id", imageId).Filter("assigned_to", userId).One(&img); err != nil {
		return nil, apperr.ErrNotFound
	}
	if img.Status != models.StatusInProgress && img.Status != models.StatusAssigned && img.Status != models.StatusRejected {
		return nil, apperr.ErrInvalidTransition
	}
	if img.Status == models.StatusAssigned {
		img.Status = models.StatusInProgress
		_, _ = o.Update(&img, "Status")
	}

	var maxVer models.ImageVersion
	_ = o.QueryTable(new(models.ImageVersion)).Filter("image_id", imageId).OrderBy("-version_no").One(&maxVer)
	next := maxVer.VersionNo + 1
	if next < 1 {
		next = 1
	}

	_, _ = o.QueryTable(new(models.ImageVersion)).Filter("image_id", imageId).Update(orm.Params{"is_current": 0})

	ver := &models.ImageVersion{
		Id:        uuid.New().String(),
		ImageId:   imageId,
		VersionNo: next,
		LayerData: layerData,
		IsCurrent: 1,
		CreatedBy: userId,
	}
	if _, err := o.Insert(ver); err != nil {
		return nil, err
	}
	return ver, nil
}

func (s *EditorService) SubmitForReview(tenantId, userId, imageId string) (*models.Image, error) {
	o := orm.NewOrm()
	var img models.Image
	if err := o.QueryTable(new(models.Image)).Filter("tenant_id", tenantId).Filter("id", imageId).Filter("assigned_to", userId).One(&img); err != nil {
		return nil, apperr.ErrNotFound
	}
	if err := MustTransition(img.Status, models.StatusPending1stReview); err != nil {
		return nil, err
	}
	count, _ := o.QueryTable(new(models.ImageVersion)).Filter("image_id", imageId).Filter("is_current", 1).Count()
	if count == 0 {
		return nil, apperr.New("VALIDATION_ERROR", "请先保存至少一个版本", apperr.ErrValidation)
	}
	img.Status = models.StatusPending1stReview
	_, err := o.Update(&img, "Status")
	return &img, err
}

func (s *EditorService) GetCurrentVersion(imageId string) (*models.ImageVersion, error) {
	o := orm.NewOrm()
	var ver models.ImageVersion
	err := o.QueryTable(new(models.ImageVersion)).Filter("image_id", imageId).Filter("is_current", 1).One(&ver)
	return &ver, err
}

func (s *EditorService) ListVersions(imageId string) ([]models.ImageVersion, error) {
	o := orm.NewOrm()
	var list []models.ImageVersion
	_, err := o.QueryTable(new(models.ImageVersion)).Filter("image_id", imageId).OrderBy("-version_no").All(&list)
	return list, err
}
