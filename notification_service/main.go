package main

import (
	"fmt"
	"log"
	"net/http"
	"notification_service/config"
	"notification_service/handlers"
	"notification_service/middleware"
	"github.com/dinesh-14699/go_assignment/common_utils/pubsubservice"
	"notification_service/subscribers"

	"github.com/dinesh-14699/go_assignment/common_utils/logger"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	logger.InitLogger("http://localhost:8084/logs", "notifiaction-service")
    logger.Log.Info("Application has started")

	if err := godotenv.Load(); err != nil {
		logger.Log.Error("No .env file found, relying on environment variables.")
		log.Println("No .env file found, relying on environment variables.")
	}

	cfg := config.LoadConfig()

	handler := handlers.NewNotificationHandler(cfg)
    subscriber := subscribers.NewSubscribers(cfg)

	projectID := "go-lang-440709"
	location := "../gcp.json"

	err := pubsubservice.InitializePubSubClient(projectID, location)
	if err != nil {
		logger.Log.Fatalf("Failed to initialize Pub/Sub client: %v", err)
	} else {
		logger.Log.Info("initialize Pub/Sub client")
	}

	r := chi.NewRouter()
	r.Use(middleware.TokenValidationMiddleware) 

	r.Post("/send-notification", handler.SendNotification)

    go startSubscriber(subscriber)

	logger.Log.Println("Starting Notification Service on port:", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		logger.Log.Fatal("Failed to start server:", err)
	}
}

func startSubscriber(subscriber *subscribers.Subscribers) {
	subscriptionID := "MySub"
	updateChan := make(chan string)

	go pubsubservice.ReceiveMessages(subscriptionID, updateChan)

	fmt.Println("Listening for messages...")

	for msg := range updateChan {
		subscriber.ReceiveMessages(msg)
	}
}

