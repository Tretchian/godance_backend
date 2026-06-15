package main

import (
	authhandler "godance/internal/handler/auth"
	competitionhandler "godance/internal/handler/competition"
	feedbackhandler "godance/internal/handler/feedback"
	judgehandler "godance/internal/handler/judge"
	registrationhandler "godance/internal/handler/registration"
	competitionrepo "godance/internal/repository/competition"
	feedbackrepo "godance/internal/repository/feedback"
	judgerepo "godance/internal/repository/judge"
	registrationrepo "godance/internal/repository/registration"
	userrepo "godance/internal/repository/user"
	authservice "godance/internal/service/auth"
	competitionservice "godance/internal/service/competition"
	feedbackservice "godance/internal/service/feedback"
	judgeservice "godance/internal/service/judge"
	registrationservice "godance/internal/service/registration"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func setupRoutes(r *gin.Engine, db *gorm.DB) {
	api := r.Group("/api/v1")

	repo := competitionrepo.NewRepository(db)
	service := competitionservice.NewService(repo)
	competitionHandler := competitionhandler.NewHandler(service)
	competitionHandler.Register(api)

	feedbackRepo := feedbackrepo.NewRepository(db)
	feedbackService := feedbackservice.NewService(feedbackRepo)
	feedbackHandler := feedbackhandler.NewHandler(feedbackService)
	feedbackHandler.Register(api)

	judgeRepo := judgerepo.NewRepository(db)
	judgeService := judgeservice.NewService(judgeRepo)
	judgeHandler := judgehandler.NewHandler(judgeService)
	judgeHandler.Register(api)

	registrationRepo := registrationrepo.NewRepository(db)
	registrationService := registrationservice.NewService(registrationRepo)
	registrationHandler := registrationhandler.NewHandler(registrationService)
	registrationHandler.Register(api)

	userRepo := userrepo.NewRepository(db)
	authService := authservice.NewService(userRepo)
	authHandler := authhandler.NewHandler(authService)
	authHandler.Register(api)
}
