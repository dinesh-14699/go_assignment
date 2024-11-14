package repositories

import (
	"context"
	"logger_service/internal/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type LoggerRepository interface {
	InsertLog(logEntry models.LogEntry) error
	FindLogs(serviceType, level, username, userID, message string, startDate, endDate time.Time) ([]models.LogEntry, error)
}

type loggerRepository struct {
	collection *mongo.Collection
}

func NewLoggerRepository(db *mongo.Database) LoggerRepository {
	return &loggerRepository{
		collection: db.Collection("logs"),
	}
}

func (r *loggerRepository) InsertLog(logEntry models.LogEntry) error {
	_, err := r.collection.InsertOne(context.TODO(), logEntry)
	return err
}

func (r *loggerRepository) FindLogs(serviceType, level, username, userID, message string, startDate, endDate time.Time) ([]models.LogEntry, error) {
	filter := bson.M{}
	if serviceType != "" {
		filter["service_type"] = serviceType
	}
	if level != "" {
		filter["level"] = level
	}
	if username != "" {
		filter["username"] = username
	}
	if userID != "" {
		filter["user_id"] = userID
	}
	if message != "" {
		filter["message"] = bson.M{"$regex": message, "$options": "i"}
	}
	if !startDate.IsZero() || !endDate.IsZero() {
		dateFilter := bson.M{}
		if !startDate.IsZero() {
			dateFilter["$gte"] = startDate
		}
		if !endDate.IsZero() {
			dateFilter["$lte"] = endDate
		}
		filter["timestamp"] = dateFilter
	}

	cursor, err := r.collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var logs []models.LogEntry
	if err = cursor.All(context.Background(), &logs); err != nil {
		return nil, err
	}
	return logs, nil
}
