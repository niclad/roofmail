package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds the configuration for the application
type Config struct {
	Version string
}

type ForcastAPIResponse struct {
	Properties ForecastAPIProps `json:"properties"`
}

type ForecastAPIProps struct {
	ID                  string `json:"@id"`
	Type                string `json:"@type"`
	CWA                 string `json:"cwa"`
	ForecastOffice      string `json:"forecastOffice"`
	GridID              string `json:"gridId"`
	GridX               int    `json:"gridX"`
	GridY               int    `json:"gridY"`
	Forecast            string `json:"forecast"`
	ForecastHourly      string `json:"forecastHourly"`
	ForecastGridData    string `json:"forecastGridData"`
	ObservationStations string `json:"observationStations"`
}

type DailyForecastResponse struct {
	Properties DailyForecast `json:"properties"`
}

type DailyForecast struct {
	Units             string    `json:"units"`
	ForecastGenerator string    `json:"forecastGenerator"`
	GeneratedAt       time.Time `json:"generatedAt"`
	UpdateTime        time.Time `json:"updateTime"`
	ValidTimes        time.Time `json:"validTimes"`
	Elevation         struct {
		UnitCode string  `json:"unitCode"`
		Value    float64 `json:"value"`
	} `json:"elevation"`
	Periods []Period `json:"periods"`
}

// Period represents a forecast period, compatible with both daily and hourly responses
type Period struct {
	Number                     int       `json:"number"`
	Name                       string    `json:"name"`
	StartTime                  time.Time `json:"startTime"`
	EndTime                    time.Time `json:"endTime"`
	IsDaytime                  bool      `json:"isDaytime"`
	Temperature                float64   `json:"temperature"`
	TemperatureUnit            string    `json:"temperatureUnit,omitempty"` // Optional for daily
	TemperatureTrend           string    `json:"temperatureTrend,omitempty"`
	ProbabilityOfPrecipitation *struct {
		UnitCode string   `json:"unitCode"`
		Value    *float64 `json:"value,omitempty"`
	} `json:"probabilityOfPrecipitation,omitempty"`
	Dewpoint *struct {
		UnitCode string  `json:"unitCode"`
		Value    float64 `json:"value"`
	} `json:"dewpoint,omitempty"`
	RelativeHumidity *struct {
		UnitCode string  `json:"unitCode"`
		Value    float64 `json:"value"`
	} `json:"relativeHumidity,omitempty"`
	WindSpeed *struct {
		UnitCode string  `json:"unitCode"`
		Value    float64 `json:"value"`
	} `json:"windSpeed,omitempty"`
	WindGust *struct {
		UnitCode string  `json:"unitCode"`
		Value    float64 `json:"value"`
	} `json:"windGust,omitempty"`
	WindDirection    string `json:"windDirection"`
	Icon             string `json:"icon"`
	ShortForecast    string `json:"shortForecast"`
	DetailedForecast string `json:"detailedForecast"`
}

type HourlyForecast struct {
}

// Log instances
var (
	infoLogger  *log.Logger
	debugLogger *log.Logger
)

// Constants for the application
var (
	LATITUDE  float64
	LONGITUDE float64
)

var BASE_URL = "https://api.weather.gov"
var ForecastAPI ForecastAPIProps

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}

	// init app dev QoL
	initLogs()
	config := loadConfig()

	// info
	infoLogger.Printf("Starting Roofmail v%s", config.Version)
	debugLogger.Println("Enabled")

	// Set the app location values
	LATITUDE, err = strconv.ParseFloat(os.Getenv("LATITUDE"), 64)
	if err != nil {
		infoLogger.Println("Error parsing latitude string to float:", err)
		return
	}

	LONGITUDE, err = strconv.ParseFloat(os.Getenv("LONGITUDE"), 64)
	if err != nil {
		infoLogger.Println("Error parsing longitude string to float:", err)
		return
	}

	// Get the forecast API information
	// This information is used to get
	ForecastAPI, err = getForecastAPI(LATITUDE, LONGITUDE)
	if err != nil {
		infoLogger.Println("Error getting forecast API: ", err)
		return
	}
	debugLogger.Printf("Forecast URL: %s\n", forecastURL(LATITUDE, LONGITUDE))
	formattedForecastAPI, err := json.MarshalIndent(ForecastAPI, "", "  ")
	if err != nil {
		infoLogger.Println("Error formatting Forecast API: ", err)
		return
	}
	debugLogger.Printf("Forecast API:\n%s\n", string(formattedForecastAPI))
}

// Initialize Logger
func initLogs() {
	// create log output
	logFile, err := os.OpenFile("roofmail.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open log file:", err)
	}

	infoLogger = log.New(logFile, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	debugLogger = log.New(logFile, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)

	if os.Getenv("APP_ENV") != "debug" {
		debugLogger.SetOutput(io.Discard)
	}
}

// Load configuration from environment variables or defaults
func loadConfig() Config {
	return Config{
		Version: "0.0.0",
	}
}

// Build the weather API URL
func forecastURL(latitude, longitude float64) string {
	baseURL := fmt.Sprintf("%s/points/%f,%f", BASE_URL, latitude, longitude)

	return baseURL
}

// Get the forecast API information for the location specified
func getForecastAPI(latitude, longitude float64) (ForecastAPIProps, error) {
	response, err := http.Get(forecastURL(latitude, longitude))
	if err != nil {
		log.Println("Error fetching forecast API: ", err)
		return ForecastAPIProps{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		log.Printf("Error: received status code %d", response.StatusCode)
		return ForecastAPIProps{}, fmt.Errorf("received status code %d", response.StatusCode)
	}

	var apiResponse ForcastAPIResponse
	err = json.NewDecoder(response.Body).Decode(&apiResponse)
	if err != nil {
		log.Println("Error decoding forecast API response: ", err)
		return ForecastAPIProps{}, err
	}

	ForecastAPI = apiResponse.Properties
	return ForecastAPI, nil
}

// Get a daily forecast
func getDailyForecast() (DailyForecast, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", ForecastAPI.Forecast, nil)
	if err != nil {
		infoLogger.Println("Error creating request: ", err)
		return DailyForecast{}, err
	}

	// Add headers
	req.Header.Set("User-Agent", "Roofmail/0.0.0")
	req.Header.Set("Accept", "application/ld+json")

	response, err := client.Do(req)
	if err != nil {
		infoLogger.Println("Error fetching daily forecast: ", err)
		return DailyForecast{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		infoLogger.Printf("Error: received status code %d", response.StatusCode)
		return DailyForecast{}, fmt.Errorf("received status code %d", response.StatusCode)
	}

	var dailyForecastResponse DailyForecastResponse
	err = json.NewDecoder(response.Body).Decode(&dailyForecastResponse)
	if err != nil {
		infoLogger.Println("Error decoding daily forecast response: ", err)
		return DailyForecast{}, err
	}

	return dailyForecastResponse.Properties, nil
}
