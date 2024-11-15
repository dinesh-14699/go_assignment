package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dinesh-14699/go_assignment/common_utils/cache"
	"github.com/dinesh-14699/go_assignment/common_utils/logger"
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

type CovidTimeSeriesEntry struct {
    Date      time.Time
    Cases     int
    Deaths    int
    Recovered int
}


func FetchCovidData(region string) (*CovidData, error) {

	cachedData, err := cache.GetValue(region)
	if err == nil && cachedData != "" {
		var cachedCovidData CovidData
		err := json.Unmarshal([]byte(cachedData), &cachedCovidData)
		if err == nil {
			return &cachedCovidData, nil
		}
	}

	data, err := FetchCovidDataFromUrl(region)
	if err != nil {
		return nil, err
	}

	cacheData, err := json.Marshal(data)
	if err != nil {
		logger.Log.Errorf("error marshalling data to cache: %v", err)
		return nil, fmt.Errorf("error marshalling data to cache: %v", err)
	}

	err = cache.SetValue(region, string(cacheData), 300) 
	if err != nil {
		logger.Log.Errorf("error setting cache: %v", err)
		return nil, fmt.Errorf("error setting cache: %v", err)
	}

	return data, nil
}


func FetchCovidDataFromUrl(region string) (*CovidData, error) {
	url := fmt.Sprintf("https://disease.sh/v3/covid-19/countries/%s", region)
  
	resp, err := http.Get(url)
	if err != nil {
		logger.Log.Errorf("error fetching data: %v", err)
		return nil, fmt.Errorf("error fetching data: %v", err)
	}
	defer resp.Body.Close()

	var data CovidData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		logger.Log.Errorf("error decoding response: %v", err)
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	data.LastUpdatedFormatted = time.Unix(data.LastUpdated/1000, 0).Format(time.RFC3339)


	return &data, nil
}

func FetchCovidTimeSeriesData(country string) ([]CovidTimeSeriesEntry, error) {
    url := fmt.Sprintf("https://disease.sh/v3/covid-19/historical/%s?lastdays=30", country) // Example API endpoint
    resp, err := http.Get(url)
    if err != nil {
		logger.Log.Errorf("failed to fetch data: %v", err)
        return nil, fmt.Errorf("failed to fetch data: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
		logger.Log.Errorf("error decoding response: %v", err)
        return nil, fmt.Errorf("error from API: status code %d", resp.StatusCode)
    }

    var responseData struct {
        Timeline struct {
            Cases     map[string]int `json:"cases"`
            Deaths    map[string]int `json:"deaths"`
            Recovered map[string]int `json:"recovered"`
        } `json:"timeline"`
    }
    if err := json.NewDecoder(resp.Body).Decode(&responseData); err != nil {
		logger.Log.Errorf("failed to decode API response: %v", err)
        return nil, fmt.Errorf("failed to decode API response: %w", err)
    }

    var timeSeriesData []CovidTimeSeriesEntry
    for dateStr, cases := range responseData.Timeline.Cases {
        parsedDate, err := time.Parse("1/2/06", dateStr) 
        if err != nil {
            continue 
        }

        deaths := responseData.Timeline.Deaths[dateStr]
        recovered := responseData.Timeline.Recovered[dateStr]

        timeSeriesData = append(timeSeriesData, CovidTimeSeriesEntry{
            Date:      parsedDate,
            Cases:     cases,
            Deaths:    deaths,
            Recovered: recovered,
        })
    }

    return timeSeriesData, nil
}
