package handlers

import (
	"encoding/json"
	"net/http"
	"notification_service/config"
	"notification_service/email"
	"notification_service/services"

	"github.com/dinesh-14699/go_assignment/common_utils/logger"
)

type NotificationHandler struct {
	cfg *config.Config
}

func NewNotificationHandler(cfg *config.Config) *NotificationHandler {
	return &NotificationHandler{cfg}
}

func (h *NotificationHandler) SendNotification(w http.ResponseWriter, r *http.Request) {
	logger.Log.Info("Received request to send notification")

	var requestBody struct {
		To           string `json:"to"`
		Subject      string `json:"subject"`
		Country      string `json:"country"`
		EmailService string `json:"email_service,omitempty"` 
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		logger.Log.Errorf("Error decoding request payload: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	logger.Log.Infof("Parsed request body: To=%s, Subject=%s, Country=%s, EmailService=%s", requestBody.To, requestBody.Subject, requestBody.Country, requestBody.EmailService)

	var emailService email.EmailService

	if requestBody.EmailService == "mailgun" {
		emailService = &email.MailgunService{}
		logger.Log.Info("Using Mailgun email service")
	} else {
		emailService = &email.SMTPService{
			From:     h.cfg.SMTPFrom,
			Password: h.cfg.SMTPPassword,
			Host:     h.cfg.SMTPHost,
			Port:     h.cfg.SMTPPort,
		}
		logger.Log.Info("Using SMTP email service")
	}

	covidDataNew := []services.CovidData{}
	data, err := services.FetchCovidDataFromUrl(requestBody.Country)

	if err == nil {
		covidDataNew = append(covidDataNew, *data)
	}

	if err := emailService.SendEmail(requestBody.To, requestBody.Subject, covidDataNew); err != nil {
		logger.Log.Errorf("Error sending email: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logger.Log.Infof("Email sent to %s with subject %s", requestBody.To, requestBody.Subject)

	response := map[string]string{
		"message": "Notification sent successfully!",
		"to":      requestBody.To,
		"subject": requestBody.Subject,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response) 
	logger.Log.Info("Notification process completed successfully")
}

