package services

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type LogData struct {
	ServiceType string `json:"service_type"`
	Level       string `json:"level"`
	Message     string `json:"message"`
	UserId      string `json:"user_id"`
	Username    string `json:"username"`
	Timestamp   string `json:"timestamp"`
}

var LogServiceURL = "http://localhost:8084/logs"

func SendLog(serviceType, level, message, user_id, username string) {
	logData := LogData{
		ServiceType: serviceType,
		Level:       level,
		UserId:      user_id,
        Username:    username,
		Message:     message,
		Timestamp:   time.Now().Format(time.RFC3339),
	}

	logJSON, _ := json.Marshal(logData)
	req, err := http.NewRequest("POST", LogServiceURL, bytes.NewBuffer(logJSON))
	if err != nil {
		logrus.Error("Error creating log request:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logrus.Error("Error sending log to logger_service:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logrus.Error("Logger service returned non-OK status:", resp.StatusCode)
	}
}

func SetLogServiceURL(url string) {
	LogServiceURL = url
}
