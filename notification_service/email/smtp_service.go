package email

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"notification_service/services"
	"path/filepath"
)

type SMTPService struct {
	From     string
	Password string
	Host     string
	Port     string
}

func (s *SMTPService) SendEmail(to string, subject string, covidData []services.CovidData) error {
	
	templateDir, err := filepath.Abs("templates")
    if err != nil {
        return fmt.Errorf("failed to get template directory: %w", err)
    }

    t, err := template.ParseFiles(filepath.Join(templateDir, "daily_covid_data_email.html"))
    if err != nil {
        return fmt.Errorf("failed to parse template: %w", err)
    }

    emailContent := struct {
        To        string
        CovidData []services.CovidData
    }{
        To:        to,
        CovidData: covidData,
    }

    var bodyBuffer bytes.Buffer
    if err := t.Execute(&bodyBuffer, emailContent); err != nil {
        return fmt.Errorf("failed to execute template: %w", err)
    }

    msg := "MIME-Version: 1.0\r\n" +
        "Content-Type: text/html; charset=\"UTF-8\"\r\n" +
        "From: " + s.From + "\r\n" +
        "To: " + to + "\r\n" +
        "Subject: " + subject + "\r\n\r\n" +
        bodyBuffer.String()

		auth := smtp.PlainAuth("", s.From, s.Password, s.Host)

		err = smtp.SendMail(s.Host+":"+s.Port, auth, s.From, []string{to}, []byte(msg))

		if err != nil {
			return fmt.Errorf("failed to send email: %w", err)
		}
	
		return nil
}
