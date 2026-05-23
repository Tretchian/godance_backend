package types

type FeedbackRequestStatus string

const (
	FeedbackRequestStatusAwaitingPayment      FeedbackRequestStatus = "awaiting_payment"
	FeedbackRequestStatusAwaitingVideo        FeedbackRequestStatus = "awaiting_video"
	FeedbackRequestStatusAwaitingConfirmation FeedbackRequestStatus = "awaiting_confirmation"
	FeedbackRequestStatusPending              FeedbackRequestStatus = "pending"
	FeedbackRequestStatusCompleted            FeedbackRequestStatus = "completed"
	FeedbackRequestStatusRefunded             FeedbackRequestStatus = "refunded"
)
