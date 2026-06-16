package dto

type RegistrationItem struct {
	ID            uint   `json:"id"`
	ParticipantID uint   `json:"participant_id"`
	FullName      string `json:"full_name"`
	PaymentStatus string `json:"payment_status"`
	RegisteredAt  string `json:"registered_at"`
}

type RegistrationListResponse struct {
	Data       []RegistrationItem `json:"data"`
	Pagination Pagination         `json:"pagination"`
}

type RegistrationCreated struct {
	RegistrationID uint   `json:"registration_id"`
	PaymentURL     string `json:"payment_url"`
}
