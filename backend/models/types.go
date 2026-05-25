package models

// Image status constants
const (
	StatusPendingStorage   = "pending_storage"
	StatusPendingAssign    = "pending_assign"
	StatusAssigned         = "assigned"
	StatusDiscarded        = "discarded"
	StatusInProgress       = "in_progress"
	StatusPending1stReview = "pending_1st_review"
	StatusRejected         = "rejected"
	Status1stReviewPassed  = "1st_review_passed"
	StatusPending2ndReview = "pending_2nd_review"
	Status2ndReviewPassed  = "2nd_review_passed"
	StatusCompleted        = "completed"
)

const (
	RoleAdmin    = "admin"
	RoleUser     = "user"
	RoleReviewer = "reviewer"
)

const (
	BatchProcessing = "processing"
	BatchCompleted  = "completed"
	BatchFailed     = "failed"
)

const (
	ReviewPass   = "pass"
	ReviewReject = "reject"
)
