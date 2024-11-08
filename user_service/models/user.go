package models

import "time"

type User struct {
    ID         uint      `gorm:"primaryKey"`
    Username   string    `gorm:"not null"`
    Email      string    `gorm:"uniqueIndex;not null"`
    Password   string    `gorm:"not null"`
	Role       string    `gorm:"not null"`
    Country    string `json:"country"`
	State      string `json:"state"`
	City       string `json:"city"`
	PostalCode string `json:"postal_code"`
    CreatedAt time.Time
}
