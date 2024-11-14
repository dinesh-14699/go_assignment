package api

import (
	"encoding/json"
	"net/http"
	"time"

	"logger_service/internal/services"

	"github.com/go-chi/chi/v5"
)

type LogRequest struct {
	Message     string `json:"message"`
	Level       string `json:"level"`
	UserID      string `json:"user_id"`
	Username    string `json:"username"`
	ServiceType string `json:"service_type"`
}

func RegisterRoutes(router *chi.Mux, loggerService services.LoggerService) {
	router.Post("/logs", func(w http.ResponseWriter, r *http.Request) {
		var logReq LogRequest
		if err := json.NewDecoder(r.Body).Decode(&logReq); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		err := loggerService.Log(logReq.Message, logReq.Level, logReq.UserID, logReq.Username, logReq.ServiceType)
		if err != nil {
			http.Error(w, "Failed to log message", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK) // Respond with 200 OK
	})

	router.Get("/logs", func(w http.ResponseWriter, r *http.Request) {
		handleGetLogs(w, r, loggerService)
	})
}


func handleGetLogs(w http.ResponseWriter, r *http.Request, loggerService services.LoggerService) {
	query := r.URL.Query()

	serviceType := query.Get("service_type")
	level := query.Get("level")
	username := query.Get("username")
	userID := query.Get("user_id")
	message := query.Get("message")

	var startDate, endDate time.Time
	var err error
	if query.Get("start_date") != "" {
		startDate, err = time.Parse("2006-01-02", query.Get("start_date"))
		if err != nil {
			http.Error(w, "Invalid start_date format", http.StatusBadRequest)
			return
		}
	}
	if query.Get("end_date") != "" {
		endDate, err = time.Parse("2006-01-02", query.Get("end_date"))
		if err != nil {
			http.Error(w, "Invalid end_date format", http.StatusBadRequest)
			return
		}
	}

	logs, err := loggerService.SearchLogs(serviceType, level, username, userID, message, startDate, endDate)
	if err != nil {
		http.Error(w, "Error retrieving logs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
}
