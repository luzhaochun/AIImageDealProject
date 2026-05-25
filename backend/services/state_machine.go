package services

import (
	"backend/models"
	"backend/pkg/apperr"
)

var transitions = map[string][]string{
	models.StatusPendingStorage:   {models.StatusPendingAssign},
	models.StatusPendingAssign:    {models.StatusAssigned},
	models.StatusAssigned:         {models.StatusDiscarded, models.StatusInProgress},
	models.StatusInProgress:       {models.StatusPending1stReview},
	models.StatusPending1stReview: {models.Status1stReviewPassed, models.StatusRejected},
	models.StatusRejected:         {models.StatusInProgress},
	models.Status1stReviewPassed:  {models.StatusPending2ndReview, models.StatusCompleted},
	models.StatusPending2ndReview: {models.Status2ndReviewPassed, models.StatusRejected},
	models.Status2ndReviewPassed:  {models.StatusCompleted},
}

func CanTransition(from, to string) bool {
	allowed, ok := transitions[from]
	if !ok {
		return false
	}
	for _, s := range allowed {
		if s == to {
			return true
		}
	}
	return false
}

func MustTransition(from, to string) error {
	if !CanTransition(from, to) {
		return apperr.New("INVALID_TRANSITION", "当前状态不允许该操作: "+from+" -> "+to, apperr.ErrInvalidTransition)
	}
	return nil
}
