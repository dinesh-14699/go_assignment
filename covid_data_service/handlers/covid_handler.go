package handlers

import (
	"context"
	"covid_handler/services"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	// "github.com/go-echarts/go-echarts/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
    "github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

var redisClient *redis.Client

func init() {
    redisClient = redis.NewClient(&redis.Options{
        Addr: "localhost:6379", 
        DB:   0, 
    })
}

type CovidData struct {
    Country     string `json:"country"`
    Cases       int    `json:"cases"`
    Deaths      int    `json:"deaths"`
    Recovered   int    `json:"recovered"`
    Active      int    `json:"active"`
    Updated     int64  `json:"updated"`
}

func GetCovidData(w http.ResponseWriter, r *http.Request) {
    country := chi.URLParam(r, "country")

    ctx := context.Background()
    cacheKey := fmt.Sprintf("covid:%s", country)
    // cachedData, err := redisClient.Get(ctx, cacheKey).Result()
    // if err == nil {
    //     logrus.Info("Cache hit for country:", country)
    //     w.Header().Set("Content-Type", "application/json")
    //     w.Write([]byte(cachedData))
    //     return
    // }

    url := fmt.Sprintf("https://disease.sh/v3/covid-19/countries/%s", country)
    resp, err := http.Get(url)
    if err != nil {
        http.Error(w, "Failed to fetch data", http.StatusInternalServerError)
        logrus.Error("Error fetching data:", err)
        return
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        http.Error(w, "Error from API", http.StatusBadGateway)
        logrus.Error("Non-200 response from API:", resp.StatusCode)
        return
    }

    var covidData CovidData
    if err := json.NewDecoder(resp.Body).Decode(&covidData); err != nil {
        http.Error(w, "Failed to decode API response", http.StatusInternalServerError)
        logrus.Error("Error decoding response:", err)
        return
    }

    covidJSON, _ := json.Marshal(covidData)
    redisClient.Set(ctx, cacheKey, covidJSON, 10*time.Minute)


    services.SendLog("covid_data_service", "info", fmt.Sprintf("Fetched COVID data for %s from API", country), "1", "dinesh")

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(covidJSON)
}

func GenerateCovidReportTable(w http.ResponseWriter, r *http.Request) {
	region := r.URL.Query().Get("region")


	data, err := services.FetchCovidData(region)

    fmt.Println(data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching data: %v", err), http.StatusInternalServerError)
		return
	}

	report := fmt.Sprintf(`
		<h1>COVID Data Report for %s</h1>
		<table border="1">
			<tr>
				<th>Region</th><th>Cases</th><th>Deaths</th><th>Recovered</th><th>Active Cases</th><th>Critical</th><th>Last Updated</th>
			</tr>
			<tr>
				<td>%s</td>
				<td>%d</td>
				<td>%d</td>
				<td>%d</td>
				<td>%d</td>
				<td>%d</td>
				<td>%s</td>
			</tr>
		</table>
	`, region, data.Region, data.Cases, data.Deaths, data.Recovered, data.ActiveCases, data.Critical, data.LastUpdatedFormatted)

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(report))
}


func GenerateCovidReportGraph(w http.ResponseWriter, r *http.Request) {
	region := r.URL.Query().Get("region")

	data, err := services.FetchCovidData(region)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching data: %v", err), http.StatusInternalServerError)
		return
	}

	// Create a new bar chart
	bar := charts.NewBar()

    // Set the X-axis categories
    bar.SetXAxis([]string{"Cases", "Deaths", "Recovered", "Active"})

    // Add the data to the bar chart
    barData := []opts.BarData{
        {Value: data.Cases},
        {Value: data.Deaths},
        {Value: data.Recovered},
        {Value: data.ActiveCases},
    }

    bar.AddSeries("COVID Stats", barData)

	// Set global options (like title)
	bar.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: fmt.Sprintf("COVID-19 Trend Analysis for %s", region)}),
	)
	

	// Set the response content type to HTML
	w.Header().Set("Content-Type", "text/html")
	err = bar.Render(w)
	if err != nil {
		http.Error(w, "Error rendering chart", http.StatusInternalServerError)
	}
}



func GenerateCovidTrendGraph(w http.ResponseWriter, r *http.Request) {
    country := r.URL.Query().Get("country")
    if country == "" {
        http.Error(w, "Country parameter is required", http.StatusBadRequest)
        return
    }

    data, err := services.FetchCovidTimeSeriesData(country)
    if err != nil {
        http.Error(w, fmt.Sprintf("Error fetching data: %v", err), http.StatusInternalServerError)
        log.Println("Error fetching time series data:", err)
        return
    }

    lineChart := charts.NewLine()
    lineChart.SetGlobalOptions(
        charts.WithTitleOpts(opts.Title{Title: fmt.Sprintf("COVID-19 Trend Analysis for %s", country)}),
        charts.WithXAxisOpts(opts.XAxis{
            Name: "Date",
            Type: "category",
        }),
    )

    var dates []string
    var cases []opts.LineData
	var deaths []opts.LineData
	var recovered []opts.LineData


    for _, entry := range data {
        dates = append(dates, entry.Date.Format("2006-01-02")) 
        cases = append(cases, opts.LineData{Value: entry.Cases})
		deaths = append(deaths, opts.LineData{Value: entry.Deaths})
		recovered = append(recovered, opts.LineData{Value: entry.Recovered})
    }

    lineChart.SetXAxis(dates).
	AddSeries("Cases", cases).
	AddSeries("Deaths", deaths).
	AddSeries("Recoverd", recovered)
    
    page := components.NewPage()
    page.AddCharts(lineChart)
    page.Render(w)
}

func DownloadCovidData(w http.ResponseWriter, r *http.Request) {
	region := r.URL.Query().Get("region")
	
	data, err := services.FetchCovidData(region)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching data: %v", err), http.StatusInternalServerError)
		return
	}
   

    w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s-covid-data.csv", region))
    w.Header().Set("Content-Type", "text/csv")

    writer := csv.NewWriter(w)
    defer writer.Flush()

    writer.Write([]string{"Country", "Cases", "Deaths", "Recovered", "Active", "Updated"})

	writer.Write([]string{
		data.Country,
		strconv.Itoa(data.Cases),
		strconv.Itoa(data.Deaths),
		strconv.Itoa(data.Recovered),
		strconv.Itoa(data.ActiveCases),
	    data.LastUpdatedFormatted,
	})
}