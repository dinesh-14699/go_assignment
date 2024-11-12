package email

import "notification_service/services"

type MailgunService struct {
}

func (m *MailgunService) SendEmail(to string, subject string, covidData []services.CovidData) error {
	return nil
}
