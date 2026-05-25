package services

import (
	"backend/models"
	"backend/pkg/apperr"

	"github.com/beego/beego/v2/client/orm"
	"github.com/google/uuid"
)

type ReviewService struct{}

type SubmitReviewReq struct {
	Round   int    `json:"round"`
	Result  string `json:"result"`
	Comment string `json:"comment"`
}

func (s *ReviewService) Submit(tenantId, reviewerId, imageId string, req SubmitReviewReq) (*models.Image, error) {
	o := orm.NewOrm()
	var img models.Image
	if err := o.QueryTable(new(models.Image)).Filter("tenant_id", tenantId).Filter("id", imageId).One(&img); err != nil {
		return nil, apperr.ErrNotFound
	}

	var ver models.ImageVersion
	_ = o.QueryTable(new(models.ImageVersion)).Filter("image_id", imageId).Filter("is_current", 1).One(&ver)

	rec := &models.ReviewRecord{
		Id:         uuid.New().String(),
		ImageId:    imageId,
		VersionId:  ver.Id,
		Round:      req.Round,
		Result:     req.Result,
		Comment:    req.Comment,
		ReviewerId: reviewerId,
	}
	if _, err := o.Insert(rec); err != nil {
		return nil, err
	}

	var to string
	if req.Round == 1 {
		if img.Status != models.StatusPending1stReview {
			return nil, apperr.ErrInvalidTransition
		}
		if req.Result == models.ReviewPass {
			to = models.Status1stReviewPassed
		} else {
			to = models.StatusRejected
		}
	} else if req.Round == 2 {
		if img.Status != models.StatusPending2ndReview {
			return nil, apperr.ErrInvalidTransition
		}
		if req.Result == models.ReviewPass {
			to = models.Status2ndReviewPassed
		} else {
			to = models.StatusRejected
		}
	} else {
		return nil, apperr.ErrValidation
	}

	if err := MustTransition(img.Status, to); err != nil {
		return nil, err
	}
	img.Status = to
	_, err := o.Update(&img, "Status")
	return &img, err
}

func (s *ReviewService) ListReviews(imageId string) ([]models.ReviewRecord, error) {
	o := orm.NewOrm()
	var list []models.ReviewRecord
	_, err := o.QueryTable(new(models.ReviewRecord)).Filter("image_id", imageId).OrderBy("-created_at").All(&list)
	return list, err
}
