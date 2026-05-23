package main

import (
	"godance/config"
	"godance/internal/models"

	"log"

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

	setupRoutes(r, db)

	r.Run(":8080")
}
