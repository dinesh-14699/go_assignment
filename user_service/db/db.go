package db

import (
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "log"
    "user_service/models" 
)

var DB *gorm.DB

func InitDB(dsn string) error {
    var err error
    DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        return err 
    }

    log.Println("Database connection established successfully")

    if err := DB.AutoMigrate(&models.User{}); err != nil {
        return err 
    }

    log.Println("Database migration completed successfully")
    return nil
}
