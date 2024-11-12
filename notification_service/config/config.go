package config

import (
	"os"
)

type Config struct {
	SMTPFrom     string
	SMTPPassword string
	SMTPHost     string
	SMTPPort     string
	Port         string
}

func LoadConfig() *Config {
	return &Config{
		SMTPFrom:     os.Getenv("SMTP_FROM"),
		SMTPPassword: os.Getenv("SMTP_PASSWORD"),
		SMTPHost:     os.Getenv("SMTP_HOST"),
		SMTPPort:     os.Getenv("SMTP_PORT"),
		Port:         os.Getenv("PORT"),
	}
}
