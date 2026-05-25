package services

import (
	"backend/models"
	"backend/pkg/apperr"

	"github.com/beego/beego/v2/client/orm"
)

type BatchDetail struct {
	models.ImportBatch
	PendingStorage int64 `json:"pending_storage"`
	PendingAssign    int64 `json:"pending_assign"`
	Downloaded       int64 `json:"downloaded"`
}

func (s *ImportService) GetBatchDetail(tenantId, batchId string) (*BatchDetail, error) {
	b, err := s.GetBatch(tenantId, batchId)
	if err != nil {
		return nil, apperr.ErrNotFound
	}
	o := orm.NewOrm()
	count := func(status string) int64 {
		n, _ := o.QueryTable(new(models.Image)).
			Filter("import_batch_id", batchId).
			Filter("status", status).Count()
		return n
	}
	pendingStorage := count(models.StatusPendingStorage)
	pendingAssign := count(models.StatusPendingAssign)
	downloaded, _ := o.QueryTable(new(models.Image)).
		Filter("import_batch_id", batchId).
		Filter("status__in", models.StatusPendingAssign, models.StatusAssigned, models.StatusDiscarded,
			models.StatusInProgress, models.StatusPending1stReview, models.StatusRejected,
			models.Status1stReviewPassed, models.StatusPending2ndReview, models.Status2ndReviewPassed, models.StatusCompleted).
		Count()
	return &BatchDetail{
		ImportBatch:    *b,
		PendingStorage: pendingStorage,
		PendingAssign:  pendingAssign,
		Downloaded:     downloaded,
	}, nil
}
