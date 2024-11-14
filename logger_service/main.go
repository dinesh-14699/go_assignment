package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"logger_service/config"
	"logger_service/internal/database"
	"logger_service/internal/repositories"
	"logger_service/internal/services"
	"logger_service/api"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	cfg := config.LoadConfig()

	db := database.Connect(cfg.MongoURI)
	repo := repositories.NewLoggerRepository(db) 
	loggerService := services.NewLoggerService(repo)

	router := chi.NewRouter()
	api.RegisterRoutes(router, loggerService) 

	log.Println("Logger service started on :8084")
	log.Fatal(http.ListenAndServe(":8084", router)) 
}
