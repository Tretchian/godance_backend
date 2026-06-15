package types

type RegistrationPaymentStatus string

const (
	RegistrationPaymentStatusPending RegistrationPaymentStatus = "pending"
	RegistrationPaymentStatusPaid    RegistrationPaymentStatus = "paid"
)
