package email

import "notification_service/services"

type EmailService interface {
	SendEmail(to string, subject string, covidData []services.CovidData) error
}
