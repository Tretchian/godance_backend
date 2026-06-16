package main

import (
	"log"
	"os"
	"time"

	"godance/config"
	"godance/internal/gateway"
	feedbackrepo "godance/internal/repository/feedback"
	feedbackservice "godance/internal/service/feedback"
)

func main() {
	db, err := config.NewDB()
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}

	service := feedbackservice.NewService(
		feedbackrepo.NewRepository(db),
		gateway.NewStubPayment(),
	)

	interval := workerInterval()
	log.Printf("feedback timeout worker started, interval=%s", interval)

	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		if processed, err := service.ProcessExpired(); err != nil {
			log.Printf("process expired: %v", err)
		} else if processed > 0 {
			log.Printf("processed %d expired feedback request(s)", processed)
		}
		<-ticker.C
	}
}

func workerInterval() time.Duration {
	if v := os.Getenv("WORKER_INTERVAL"); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return time.Minute
}
