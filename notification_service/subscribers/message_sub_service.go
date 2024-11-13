package subscribers

import (
	"encoding/json"
	"fmt"
	"notification_service/config"
	"notification_service/email"
	"notification_service/services"
)

type Subscribers struct {
	cfg *config.Config
}

func NewSubscribers(cfg *config.Config) *Subscribers {
	return &Subscribers{cfg}
}

func (s *Subscribers) ReceiveMessages(message string) {
	fmt.Printf("Processed message: %s\n", message)

	var genericPayload map[string]interface{}
	err := json.Unmarshal([]byte(message), &genericPayload)
	if err != nil {
		fmt.Printf("Error unmarshaling message: %v\n", err)
		return
	}

	fmt.Printf("Unmarshaled data: %+v\n", genericPayload)

	to, toExists := genericPayload["to"].(string)
	subject, subjectExists := genericPayload["subject"].(string)
	country, countryExists := genericPayload["country"].(string)

	if !toExists || !subjectExists || !countryExists {
		fmt.Println("Required fields 'to' or 'subject''or' country are missing in the payload")
		return
	}

	emailService := &email.SMTPService{
		From:     s.cfg.SMTPFrom,
		Password: s.cfg.SMTPPassword,
		Host:     s.cfg.SMTPHost,
		Port:     s.cfg.SMTPPort,
	}

	covidDataNew := []services.CovidData{}
	data, err := services.FetchCovidDataFromUrl(country)

	if err == nil {
		covidDataNew = append(covidDataNew, *data)
	}

	

	if err := emailService.SendEmail(to, subject, covidDataNew); err != nil {
		fmt.Printf("Error sending email: %v\n", err)
		return
	}

	fmt.Println("Email sent successfully to:", to)
}
