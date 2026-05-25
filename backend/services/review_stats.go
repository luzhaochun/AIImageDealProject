package services

import (
	"time"

	"backend/models"

	"github.com/beego/beego/v2/client/orm"
)

type ReviewStats struct {
	QueueCount    int64   `json:"queue_count"`
	TodayReviewed int64   `json:"today_reviewed"`
	PassRate      float64 `json:"pass_rate"`
}

func (s *ReviewService) Stats(tenantId string, round int) (*ReviewStats, error) {
	o := orm.NewOrm()
	status := models.StatusPending1stReview
	if round == 2 {
		status = models.StatusPending2ndReview
	}
	queue, _ := o.QueryTable(new(models.Image)).Filter("tenant_id", tenantId).Filter("status", status).Count()

	start := time.Now().Truncate(24 * time.Hour)
	var today int64
	_ = o.Raw(`
		SELECT COUNT(*) FROM review_records r
		INNER JOIN images i ON i.id = r.image_id
		WHERE i.tenant_id = ? AND r.round = ? AND r.created_at >= ?
	`, tenantId, round, start).QueryRow(&today)

	var passed, total int64
	_ = o.Raw(`
		SELECT
			SUM(CASE WHEN r.result = ? THEN 1 ELSE 0 END),
			COUNT(*)
		FROM review_records r
		INNER JOIN images i ON i.id = r.image_id
		WHERE i.tenant_id = ? AND r.round = ? AND r.created_at >= ?
	`, models.ReviewPass, tenantId, round, start).QueryRow(&passed, &total)

	rate := 0.0
	if total > 0 {
		rate = float64(passed) / float64(total) * 100
	}
	return &ReviewStats{QueueCount: queue, TodayReviewed: today, PassRate: rate}, nil
}
