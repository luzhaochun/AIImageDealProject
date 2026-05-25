package services

import (
	"context"
	"fmt"

	"backend/models"
	"backend/pkg/apperr"
	"backend/pkg/idgen"
	"backend/pkg/redisx"

	"github.com/beego/beego/v2/client/orm"
	"github.com/google/uuid"
)

type ImageService struct{}

type CategoryCount struct {
	Name           string `json:"name"`
	AvailableCount int64  `json:"available_count"`
}

func (s *ImageService) ListCategories(tenantId string) ([]CategoryCount, error) {
	o := orm.NewOrm()
	var rows []orm.Params
	_, err := o.Raw(`
		SELECT category AS name, COUNT(*) AS cnt
		FROM images
		WHERE tenant_id = ? AND status = ?
		GROUP BY category
	`, tenantId, models.StatusPendingAssign).Values(&rows)
	if err != nil {
		return nil, err
	}
	out := make([]CategoryCount, 0, len(rows))
	for _, r := range rows {
		name := fmt.Sprint(r["name"])
		cnt := int64(0)
		if v, ok := r["cnt"]; ok {
			switch n := v.(type) {
			case int64:
				cnt = n
			case int32:
				cnt = int64(n)
			case float64:
				cnt = int64(n)
			}
		}
		out = append(out, CategoryCount{Name: name, AvailableCount: cnt})
	}
	return out, nil
}

func (s *ImageService) Claim(tenantId, userId, category, tenantSlug string) (*models.Image, error) {
	o := orm.NewOrm()
	var img models.Image
	err := o.DoTx(func(ctx context.Context, txOrm orm.TxOrmer) error {
		var id string
		if err := txOrm.Raw(`
			SELECT id FROM images
			WHERE tenant_id = ? AND status = ? AND category = ?
			ORDER BY created_at ASC LIMIT 1
			FOR UPDATE SKIP LOCKED
		`, tenantId, models.StatusPendingAssign, category).QueryRow(&id); err != nil {
			return apperr.New("NO_IMAGE_AVAILABLE", "该类目下暂无可用图片", apperr.ErrNoImageAvailable)
		}
		if id == "" {
			return apperr.New("NO_IMAGE_AVAILABLE", "该类目下暂无可用图片", apperr.ErrNoImageAvailable)
		}
		if err := txOrm.QueryTable(new(models.Image)).Filter("id", id).One(&img); err != nil {
			return err
		}
		if err := MustTransition(img.Status, models.StatusAssigned); err != nil {
			return err
		}
		img.Status = models.StatusAssigned
		img.AssignedTo = userId
		img.GlobalNo = idgen.GlobalNo(tenantSlug)
		_, err := txOrm.Update(&img, "Status", "AssignedTo", "GlobalNo")
		return err
	})
	if err != nil {
		return nil, err
	}
	return &img, nil
}

func (s *ImageService) Discard(tenantId, userId, imageId, reason string) error {
	n, err := redisx.IncrDiscardCount(context.Background(), tenantId, userId)
	if err != nil {
		return err
	}
	if n > 10 {
		return apperr.New("DISCARD_LIMIT_EXCEEDED", "今日换图次数已达上限", apperr.ErrDiscardLimit)
	}
	o := orm.NewOrm()
	var img models.Image
	if err := o.QueryTable(new(models.Image)).Filter("tenant_id", tenantId).Filter("id", imageId).Filter("assigned_to", userId).One(&img); err != nil {
		return apperr.ErrNotFound
	}
	if err := MustTransition(img.Status, models.StatusDiscarded); err != nil {
		return err
	}
	img.Status = models.StatusDiscarded
	img.DiscardReason = reason
	_, err = o.Update(&img, "Status", "DiscardReason")
	return err
}

func (s *ImageService) StartEditing(tenantId, userId, imageId string) (*models.Image, error) {
	o := orm.NewOrm()
	var img models.Image
	if err := o.QueryTable(new(models.Image)).Filter("tenant_id", tenantId).Filter("id", imageId).Filter("assigned_to", userId).One(&img); err != nil {
		return nil, apperr.ErrNotFound
	}
	from := img.Status
	if from == models.StatusRejected {
		from = models.StatusRejected
	}
	if img.Status == models.StatusAssigned {
		if err := MustTransition(img.Status, models.StatusInProgress); err != nil {
			return nil, err
		}
		img.Status = models.StatusInProgress
		_, err := o.Update(&img, "Status")
		return &img, err
	}
	if img.Status == models.StatusRejected {
		if err := MustTransition(img.Status, models.StatusInProgress); err != nil {
			return nil, err
		}
		img.Status = models.StatusInProgress
		_, err := o.Update(&img, "Status")
		return &img, err
	}
	if img.Status == models.StatusInProgress {
		return &img, nil
	}
	return nil, apperr.ErrInvalidTransition
}

func (s *ImageService) ListTasks(tenantId, userId, statusFilter string, limit, offset int) ([]models.Image, int64, error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(models.Image)).Filter("tenant_id", tenantId).Filter("assigned_to", userId)
	if statusFilter != "" {
		qs = qs.Filter("status", statusFilter)
	} else {
		qs = qs.Filter("status__in", models.StatusAssigned, models.StatusInProgress, models.StatusRejected, models.StatusPending1stReview)
	}
	total, _ := qs.Count()
	var list []models.Image
	_, err := qs.OrderBy("-updated_at").Limit(limit, offset).All(&list)
	return list, total, err
}

func (s *ImageService) ListByStatus(tenantId, status string, limit, offset int) ([]models.Image, int64, error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(models.Image)).Filter("tenant_id", tenantId).Filter("status", status)
	total, _ := qs.Count()
	var list []models.Image
	_, err := qs.OrderBy("created_at").Limit(limit, offset).All(&list)
	return list, total, err
}

type ImageDetail struct {
	models.Image
	LastReviewComment string `json:"last_review_comment,omitempty"`
	AssigneeName      string `json:"assignee_name,omitempty"`
	VersionNo         int    `json:"version_no,omitempty"`
}

func (s *ImageService) GetDetail(tenantId, imageId string) (*ImageDetail, error) {
	img, err := s.Get(tenantId, imageId)
	if err != nil {
		return nil, err
	}
	d := &ImageDetail{Image: *img}
	o := orm.NewOrm()
	if img.AssignedTo != "" {
		var u models.User
		if o.QueryTable(new(models.User)).Filter("id", img.AssignedTo).One(&u) == nil {
			d.AssigneeName = u.DisplayName
		}
	}
	var ver models.ImageVersion
	if o.QueryTable(new(models.ImageVersion)).Filter("image_id", imageId).Filter("is_current", 1).One(&ver) == nil {
		d.VersionNo = ver.VersionNo
	}
	if img.Status == models.StatusRejected {
		var rec models.ReviewRecord
		if o.QueryTable(new(models.ReviewRecord)).Filter("image_id", imageId).Filter("result", models.ReviewReject).OrderBy("-created_at").One(&rec) == nil {
			d.LastReviewComment = rec.Comment
		}
	}
	return d, nil
}

func (s *ImageService) Get(tenantId, imageId string) (*models.Image, error) {
	o := orm.NewOrm()
	var img models.Image
	err := o.QueryTable(new(models.Image)).Filter("tenant_id", tenantId).Filter("id", imageId).One(&img)
	return &img, err
}

func (s *ImageService) GetTenantSlug(tenantId string) (string, error) {
	o := orm.NewOrm()
	var t models.Tenant
	if err := o.QueryTable(new(models.Tenant)).Filter("id", tenantId).One(&t); err != nil {
		return "IMG", err
	}
	return t.Slug, nil
}

func (s *ImageService) ClaimSecondReview(tenantId, reviewerId string) (*models.Image, error) {
	o := orm.NewOrm()
	var img models.Image
	err := o.DoTx(func(ctx context.Context, txOrm orm.TxOrmer) error {
		var id string
		if err := txOrm.Raw(`
			SELECT id FROM images
			WHERE tenant_id = ? AND status = ?
			ORDER BY updated_at ASC LIMIT 1
			FOR UPDATE SKIP LOCKED
		`, tenantId, models.Status1stReviewPassed).QueryRow(&id); err != nil || id == "" {
			return apperr.New("NO_IMAGE_AVAILABLE", "暂无待二审图片", apperr.ErrNoImageAvailable)
		}
		if err := txOrm.QueryTable(new(models.Image)).Filter("id", id).One(&img); err != nil {
			return err
		}
		if err := MustTransition(img.Status, models.StatusPending2ndReview); err != nil {
			return err
		}
		img.Status = models.StatusPending2ndReview
		_, err := txOrm.Update(&img, "Status")
		return err
	})
	if err != nil {
		return nil, err
	}
	return &img, nil
}

func NewImageID() string { return uuid.New().String() }
