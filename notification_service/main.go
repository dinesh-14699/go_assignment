package main

import (
	"log"
	"net/http"
	"notification_service/config"
	"notification_service/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on environment variables.")
	}
	
	cfg := config.LoadConfig()

	handler := handlers.NewNotificationHandler(cfg)

	r := chi.NewRouter()
	r.Post("/send-notification", handler.SendNotification)

	log.Println("Starting Notification Service on port:", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
