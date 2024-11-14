package services

import (
	"logger_service/internal/models"
	"time"
)

type LoggerService interface {
	Log(message, level, userID, username, serviceType string) error
	SearchLogs(serviceType, level, username, userID, message string, startDate, endDate time.Time) ([]models.LogEntry, error)
}