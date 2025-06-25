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

	fmt.Println(comfortMessage(forecast.Periods[0]))
}

// Initialize Logger
func initLogs() {
	// create log output
	logFile, err := os.OpenFile("roofmail.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open log file:", err)
	}

	var secondWriter *os.File = nil
	if os.Getenv("APP_ENV") == "development" {
		secondWriter = os.Stdout
	}

	multiInfo := io.MultiWriter(logFile, secondWriter)
	multiDebug := io.MultiWriter(logFile, secondWriter)

	infoLogger = log.New(multiInfo, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	debugLogger = log.New(multiDebug, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)

	if os.Getenv("APP_ENV") != "development" {
		debugLogger.SetOutput(io.Discard)
	}
}

// Load configuration from environment variables or defaults
func loadConfig() Config {
	return Config{
		Version: "0.0.0",
	}
}

// Get the current temperature in Fahreneheit
func getTempF(period wapi.Period) float64 {
	tempF := period.Temperature.Value

	if period.Temperature.UnitCode == "wmoUnit:degC" {
		tempF = ctof(period.Temperature.Value)
	}

	return tempF
}

// Get the current Beaufort value
func getBeaufort(period wapi.Period) int {

	return beaufortScale(*period.WindSpeed)
}

// Convert Celsius to Fahrenheit
func ctof(c float64) float64 {
	return (9 * c / 5) + 32
}

// Convert KPH to MPH
func kphToMph(kph float64) float64 {
	return kph / 1.609344
}

// Determine arbitrary "comfort" value
// Returns an `int`, where 0 is most comfortable and 10 is least comfortable.
// Note: This will probably have to be tuned
func isComfortable(period wapi.Period) bool {
	beaufort := getBeaufort(period)
	temperature := getTempF(period)

	var minTempF = 70.0

	return (temperature >= minTempF && beaufort < 4)
}

func comfortMessage(period wapi.Period) string {
	isComfy := isComfortable(period)
	temp := getTempF(period)
	beau := getBeaufort(period)

	notStr := " not "

	if isComfy {
		notStr = " "
	}

	return fmt.Sprintf("It looks like the weather will%sbe comfortable. The temperature is %.0f\u00B0F with a wind-level of %d.", notStr, temp, beau)
}

// Determine Beaufort value
func beaufortScale(windSpeed wapi.WindSpeed) int {
	var wind float64
	if windSpeed.Value != nil {
		wind = *windSpeed.Value
	} else if windSpeed.MaxValue != nil {
		wind = *windSpeed.MaxValue
	} else {
		return 0 // assuming no value means no wind
	}

	// convert to mph if kph
	if windSpeed.UnitCode == "wmoUnits:km_h-1" {
		wind = kphToMph(wind)
	}

	switch {
	case wind < 1.0:
		return 0
	case wind < 4.0:
		return 1
	case wind < 8.0:
		return 2
	case wind < 13:
		return 3
	case wind < 19:
		return 4
	default:
		return 5
	}
}
