package dto

type UploadUrlRequest struct {
	RequestID   uint   `json:"request_id" binding:"required"`
	Filename    string `json:"filename" binding:"required"`
	ContentType string `json:"content_type" binding:"required,oneof=video/mp4 video/quicktime"`
	FileSize    int64  `json:"file_size" binding:"required"`
}

type UploadUrlResponse struct {
	VideoID   uint   `json:"video_id"`
	UploadURL string `json:"upload_url"`
	ExpiresAt string `json:"expires_at"`
}

type VideoConfirmRequest struct {
	RequestID uint    `json:"request_id" binding:"required"`
	Comment   *string `json:"comment" binding:"omitempty,max=1000"`
}

type VideoConfirmResponse struct {
	VideoID       uint   `json:"video_id"`
	Status        string `json:"status"`
	ExpiresAt     string `json:"expires_at"`
	RequestStatus string `json:"request_status"`
}

type VideoViewResponse struct {
	VideoID      uint   `json:"video_id"`
	ViewURL      string `json:"view_url"`
	URLExpiresAt string `json:"url_expires_at"`
	VideoLength  string `json:"video_length"`
}
