package main

import (
	"fmt"
	"log"
	"net/http"
	"notification_service/config"
	"notification_service/handlers"
	"notification_service/pubsubservice"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on environment variables.")
	}

	cfg := config.LoadConfig()

	handler := handlers.NewNotificationHandler(cfg)

	projectID := "go-lang-440709"
	location := "../gcp.json"

	err := pubsubservice.InitializePubSubClient(projectID, location)
	if err != nil {
		log.Fatalf("Failed to initialize Pub/Sub client: %v", err)
	} else {
		logrus.Info("initialize Pub/Sub client")
	}

	r := chi.NewRouter()
	r.Post("/send-notification", handler.SendNotification)

    go startSubscriber()

	log.Println("Starting Notification Service on port:", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func startSubscriber() {
	subscriptionID := "MySub"
	updateChan := make(chan string)

	go pubsubservice.ReceiveMessages(subscriptionID, updateChan)

	fmt.Println("Listening for messages...")

	for msg := range updateChan {
		fmt.Printf("Processed message: %s\n", msg)
	}
}

