package main

import (
	// "covid_handler/cache"
	"covid_handler/handlers"
	"covid_handler/pubsubservice"

	// "covid_handler/middleware"
	"log"
	"net/http"
	"covid_handler/logger"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

func main() {
    router := chi.NewRouter()
    // cache.InitializeCache("100.29.10.135:6379", "", 0)

    // router.Use(middleware.Logger)
    // router.Use(middleware.TokenValidationMiddleware) 

	projectID := "go-lang-440709"
	location := "../gcp.json"
	logger.InitLogger("http://localhost:8084/logs", "covid-data-service")
    logger.Log.Info("Application has started")
    

	err := pubsubservice.InitializePubSubClient(projectID, location)
	if err != nil {
		log.Fatalf("Failed to initialize Pub/Sub client: %v", err)
	} else {
		logrus.Info("initialize Pub/Sub client")
	}

    router.Get("/covid/{country}", handlers.GetCovidData)
    router.Get("/covid-report-table", handlers.GenerateCovidReportTable)
	router.Get("/covid-report-graph", handlers.GenerateCovidReportGraph)
	router.Get("/covid-report-download", handlers.DownloadCovidData)
	router.Get("/covid-report-trend", handlers.GenerateCovidTrendGraph)
	router.Post("/send-covid-notification", handlers.FetchCovidDataAndPublish)


    logrus.Info("Starting server on port 8082...")
    if err := http.ListenAndServe(":8082", router); err != nil {
        log.Fatal("Server failed to start:", err)
    }
}
