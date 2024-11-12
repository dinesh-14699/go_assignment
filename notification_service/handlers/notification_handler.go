package handlers

import (
	"encoding/json"
	"net/http"
	"notification_service/config"
	"notification_service/email"
	"notification_service/services"
)

type NotificationHandler struct {
	cfg *config.Config
}

func NewNotificationHandler(cfg *config.Config) *NotificationHandler {
	return &NotificationHandler{cfg}
}

func (h *NotificationHandler) SendNotification(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		To           string `json:"to"`
		Subject      string `json:"subject"`
		Country      string `json:"country"`
		EmailService string `json:"email_service,omitempty"` 
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var emailService email.EmailService

	if requestBody.EmailService == "mailgun" {
		emailService = &email.MailgunService{}
	} else {
		emailService = &email.SMTPService{
			From:     h.cfg.SMTPFrom,
			Password: h.cfg.SMTPPassword,
			Host:     h.cfg.SMTPHost,
			Port:     h.cfg.SMTPPort,
		}
	}

	covidData := []services.CovidData{}
	data, err := services.FetchCovidDataFromUrl(requestBody.Country)

	if err == nil {
		covidData = append(covidData, *data)
	}

	if err := emailService.SendEmail(requestBody.To, requestBody.Subject, covidData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"message": "Notification sent successfully!",
		"to":      requestBody.To,
		"subject": requestBody.Subject,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response) 
}
