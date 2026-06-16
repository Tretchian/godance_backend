package types

type VideoStatus string

const (
	VideoStatusUploading VideoStatus = "uploading"
	VideoStatusActive    VideoStatus = "active"
	VideoStatusDeleted   VideoStatus = "deleted"
)
