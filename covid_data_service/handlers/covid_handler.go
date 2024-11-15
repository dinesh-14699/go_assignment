package handlers

import (
	"covid_handler/pubsubservice"
	"covid_handler/services"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/dinesh-14699/go_assignment/common_utils/logger"
)

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

    logger.Log.Infof("Received request to fetch COVID data for country: %s", country)

	covidData, err := services.FetchCovidData(country)
    if err != nil {
		logger.Log.Errorf("Error fetching COVID data for country %s: %v", country, err)
		http.Error(w, fmt.Sprintf("Error fetching data: %v", err), http.StatusInternalServerError)
		return
	}

    covidJSON, _ := json.Marshal(covidData)

    logger.Log.Infof("Successfully fetched COVID data for country: %s", country)

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(covidJSON)
}

func GenerateCovidReportTable(w http.ResponseWriter, r *http.Request) {
	region := r.URL.Query().Get("region")
    logger.Log.Infof("Received request to generate COVID report for region: %s", region)

	data, err := services.FetchCovidData(region)
	if err != nil {
		logger.Log.Errorf("Error fetching COVID data for region %s: %v", region, err)
		http.Error(w, fmt.Sprintf("Error fetching data: %v", err), http.StatusInternalServerError)
		return
	}

	report := fmt.Sprintf(`...
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

    logger.Log.Infof("Successfully generated COVID report for region: %s", region)

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(report))
}

func GenerateCovidReportGraph(w http.ResponseWriter, r *http.Request) {
	region := r.URL.Query().Get("region")
    logger.Log.Infof("Generating COVID report graph for region: %s", region)

	data, err := services.FetchCovidData(region)
	if err != nil {
		logger.Log.Errorf("Error fetching COVID data for region %s: %v", region, err)
		http.Error(w, fmt.Sprintf("Error fetching data: %v", err), http.StatusInternalServerError)
		return
	}

	bar := charts.NewBar()
	bar.SetXAxis([]string{"Cases", "Deaths", "Recovered", "Active"})

	barData := []opts.BarData{
        {Value: data.Cases},
        {Value: data.Deaths},
        {Value: data.Recovered},
        {Value: data.ActiveCases},
    }

    bar.AddSeries("COVID Stats", barData)

	bar.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: fmt.Sprintf("COVID-19 Trend Analysis for %s", region)}),
	)

	err = bar.Render(w)
	if err != nil {
		logger.Log.Errorf("Error rendering chart for region %s: %v", region, err)
		http.Error(w, "Error rendering chart", http.StatusInternalServerError)
		return
	}

    logger.Log.Infof("Successfully generated COVID report graph for region: %s", region)
}

func GenerateCovidTrendGraph(w http.ResponseWriter, r *http.Request) {
    country := r.URL.Query().Get("country")
    if country == "" {
        http.Error(w, "Country parameter is required", http.StatusBadRequest)
        return
    }

    logger.Log.Infof("Generating trend graph for country: %s", country)

    data, err := services.FetchCovidTimeSeriesData(country)
    if err != nil {
        logger.Log.Errorf("Error fetching time series data for country %s: %v", country, err)
        http.Error(w, fmt.Sprintf("Error fetching data: %v", err), http.StatusInternalServerError)
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
        AddSeries("Recovered", recovered)

    page := components.NewPage()
    page.AddCharts(lineChart)
    page.Render(w)

    logger.Log.Infof("Successfully generated trend graph for country: %s", country)
}

func DownloadCovidData(w http.ResponseWriter, r *http.Request) {
	region := r.URL.Query().Get("region")
    logger.Log.Infof("Downloading COVID data for region: %s", region)

	data, err := services.FetchCovidData(region)
	if err != nil {
		logger.Log.Errorf("Error fetching COVID data for region %s: %v", region, err)
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

    logger.Log.Infof("Successfully downloaded COVID data for region: %s", region)
}

func FetchCovidDataAndPublish(w http.ResponseWriter, r *http.Request) {
	type CovidDataPayload struct {
		Subject  string `json:"subject"`
		To       string `json:"to"`
		Country  string `json:"country"`
	}

	var payload CovidDataPayload
    decoder := json.NewDecoder(r.Body)
    err := decoder.Decode(&payload)
    if err != nil {
        logger.Log.Errorf("Error decoding request body: %v", err)
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    if payload.Subject == "" || payload.To == "" || payload.Country == "" {
        logger.Log.Error("Missing required fields in request body")
        http.Error(w, "Missing required fields: subject, to, or country", http.StatusBadRequest)
        return
    }

	msg, err := json.Marshal(payload)
	if err != nil {
		logger.Log.Errorf("Error marshaling COVID data: %v", err)
		http.Error(w, "Error marshaling COVID data", http.StatusInternalServerError)
		return
	}

	id, err := pubsubservice.PublishMessage("MyTopic", string(msg))
	if err != nil {
		logger.Log.Errorf("Failed to publish message to Pub/Sub: %v", err)
		http.Error(w, "Failed to publish message to Pub/Sub", http.StatusInternalServerError)
		return
	}

	logger.Log.Infof("Message published successfully. Pub/Sub message ID: %s", id)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("COVID data and subject published successfully"))
}
