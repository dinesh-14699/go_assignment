package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type CovidData struct {
	Country      string `json:"country"`
	Region       string `json:"region"`
	Cases        int    `json:"cases"`
	Deaths       int    `json:"deaths"`
	Recovered    int    `json:"recovered"`
	ActiveCases  int    `json:"active"`
	Critical     int    `json:"critical"`
	LastUpdated  int64  `json:"updated"`
	LastUpdatedFormatted string `json:"lastUpdatedFormatted"`
}


func FetchCovidDataFromUrl(region string) (*CovidData, error) {
	url := fmt.Sprintf("https://disease.sh/v3/covid-19/countries/%s", region)
  
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching data: %v", err)
	}
	defer resp.Body.Close()

	var data CovidData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	data.LastUpdatedFormatted = time.Unix(data.LastUpdated/1000, 0).Format(time.RFC3339)


	return &data, nil
}
