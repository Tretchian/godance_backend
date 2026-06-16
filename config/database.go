package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

const (
	dbConnectAttempts = 10
	dbConnectBackoff  = 3 * time.Second
)

func NewDB() (*gorm.DB, error) {
	godotenv.Load()

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSLMODE"),
	)

	// Ретраим подключение, чтобы переждать старт БД в docker-compose.
	var err error
	for attempt := 1; attempt <= dbConnectAttempts; attempt++ {
		var db *gorm.DB
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{TranslateError: true})
		if err == nil {
			if sqlDB, e := db.DB(); e != nil {
				err = e
			} else if err = sqlDB.Ping(); err == nil {
				DB = db
				return db, nil
			}
		}
		log.Printf("db connect attempt %d/%d failed: %v", attempt, dbConnectAttempts, err)
		time.Sleep(dbConnectBackoff)
	}
	return nil, err
}
