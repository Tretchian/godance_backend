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

// PaymentWebhookStatus — статусы, приходящие от платёжного шлюза в вебхуках.
type PaymentWebhookStatus string

const (
	PaymentWebhookStatusSucceeded     PaymentWebhookStatus = "succeeded"
	PaymentWebhookStatusFailed        PaymentWebhookStatus = "failed"
	PaymentWebhookStatusHeld          PaymentWebhookStatus = "held"
	PaymentWebhookStatusCaptured      PaymentWebhookStatus = "captured"
	PaymentWebhookStatusReleased      PaymentWebhookStatus = "released"
	PaymentWebhookStatusPaidOut       PaymentWebhookStatus = "paid_out"
	PaymentWebhookStatusCaptureFailed PaymentWebhookStatus = "capture_failed"
)
