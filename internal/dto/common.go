package dto

type Pagination struct {
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
	Total int64 `json:"total"`
}

type Error struct {
	Error string `json:"error"`
	Code  string `json:"code"`
}

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ValidationError struct {
	Error   string       `json:"error"`
	Code    string       `json:"code"`
	Details []FieldError `json:"details"`
}
