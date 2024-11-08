package config

import (
    "log"
    "os"

    "github.com/joho/godotenv"
)

var (
    DSN       string 
    JWTSecret string
)

func LoadConfig() {
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found, reading environment variables")
    }

    DSN = os.Getenv("DATABASE_DSN")
    if DSN == "" {
        log.Fatal("DATABASE_DSN is required but not set")
    }

    JWTSecret = os.Getenv("JWT_SECRET")
    if JWTSecret == "" {
        log.Fatal("JWT_SECRET is required but not set")
    }
}
