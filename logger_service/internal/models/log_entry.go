package models

import "time"

type LogEntry struct {
	ID          string    `json:"id" bson:"_id,omitempty"`
	Message     string    `json:"message" bson:"message"`
	Level       string    `json:"level" bson:"level"`
	UserID      string    `json:"user_id" bson:"user_id"`
	Username    string    `json:"username" bson:"username"`
	ServiceType string    `json:"service_type" bson:"service_type"`
	Timestamp   time.Time `json:"timestamp" bson:"timestamp"`
}
