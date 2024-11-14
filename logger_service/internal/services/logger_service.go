package services

import (
	"logger_service/internal/models"
	"logger_service/internal/repositories"
	"time"
)



type loggerService struct {
	repo repositories.LoggerRepository
}

func NewLoggerService(repo repositories.LoggerRepository) LoggerService {
	return &loggerService{repo}
}

func (s *loggerService) Log(message, level, userID, username, serviceType string) error {
	logEntry := models.LogEntry{
		Message:     message,
		Level:       level,
		UserID:      userID,
		Username:    username,
		ServiceType: serviceType,
		Timestamp:   time.Now(),
	}

	return s.repo.InsertLog(logEntry)
}

func (s *loggerService) SearchLogs(serviceType, level, username, userID, message string, startDate, endDate time.Time) ([]models.LogEntry, error) {
	return s.repo.FindLogs(serviceType, level, username, userID, message, startDate, endDate)
}