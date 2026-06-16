package main

import (
	"godance/config"
	"godance/internal/domain"
	"godance/internal/models"
	"godance/internal/storage"

	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	db, err := config.NewDB()

	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}

	db.AutoMigrate(
		&models.User{},
		&models.Profile{},
		&models.Competition{},
		&models.Registration{},
		&models.Video{},
		&models.FeedbackRequest{},
		&models.FeedbackResponse{},
		&models.FeedbackRating{},
		&models.Notification{},
		&models.Payment{},
		&models.AuditLog{},
	)

	var store domain.VideoStorage
	if os.Getenv("MINIO_ENDPOINT") == "" {
		log.Println("MINIO_ENDPOINT not set: using stub video storage")
		store = storage.NewStub()
	} else {
		minioStore, err := storage.NewMinIO()
		if err != nil {
			log.Fatal("failed to init MinIO storage:", err)
		}
		store = minioStore
	}

	setupRoutes(r, db, store)

	r.Run(":8080")
}
