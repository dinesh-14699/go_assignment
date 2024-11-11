package main

import (
	"covid_handler/cache"
	"covid_handler/handlers"
	// "covid_handler/middleware"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

func main() {
    router := chi.NewRouter()
    cache.InitializeCache("54.221.19.22:6379", "", 0)

    // router.Use(middleware.Logger)
    // router.Use(middleware.TokenValidationMiddleware) 

    router.Get("/covid/{country}", handlers.GetCovidData)
    
    router.Get("/covid-report-table", handlers.GenerateCovidReportTable)
	router.Get("/covid-report-graph", handlers.GenerateCovidReportGraph)
	router.Get("/covid-report-download", handlers.DownloadCovidData)
	router.Get("/covid-report-trend", handlers.GenerateCovidTrendGraph)


    logrus.Info("Starting server on port 8082...")
    if err := http.ListenAndServe(":8082", router); err != nil {
        log.Fatal("Server failed to start:", err)
    }
}
