package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	wapi "roofmail/weatherAPI"

	"github.com/joho/godotenv"
)

// Config holds the configuration for the application
type Config struct {
	Version string
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

// Weather API
var w wapi.WeatherAPI
var ctx context.Context

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

	// set up client
	var client = http.Client{}
	w = wapi.NewWeatherGovAPI(&client, &LATITUDE, &LONGITUDE)

	// create a context
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// initialize the API
	err = w.InitForecastAPI(ctx, nil, nil)
	if err != nil {
		infoLogger.Panicln("Error initializing Weather API:", err)
		return
	}

	// get the daily forecast
	forecast, err := w.GetDailyForecast(ctx)
	if err != nil {
		infoLogger.Panicln("Error getting daily forecast:", err)
		return
	}

	formattedForecast := fmt.Sprintf("DailyForecast: %+v", forecast)

	debugLogger.Println(formattedForecast)
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
