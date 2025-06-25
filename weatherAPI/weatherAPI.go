package weatherAPI

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// GLOBALS
var BASE_URL = "https://api.weather.gov"

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
	ValidTimes        string    `json:"validTimes"`
	Elevation         struct {
		UnitCode string  `json:"unitCode"`
		Value    float64 `json:"value"`
	} `json:"elevation"`
	Periods []Period `json:"periods"`
}

type WindSpeed struct {
	UnitCode string   `json:"unitCode"`
	Value    *float64 `json:"value,omitempty"`
	MaxValue *float64 `json:"maxValue,omitempty"`
	MinValue *float64 `json:"minValue,omitempty"`
}

// Period represents a forecast period, compatible with both daily and hourly responses
type Period struct {
	Number      int       `json:"number"`
	Name        string    `json:"name"`
	StartTime   time.Time `json:"startTime"`
	EndTime     time.Time `json:"endTime"`
	IsDaytime   bool      `json:"isDaytime"`
	Temperature *struct {
		UnitCode string  `json:"unitCode"`
		Value    float64 `json:"value"`
	} `json:"temperature,omitempty"`
	TemperatureUnit            string `json:"temperatureUnit,omitempty"` // Optional for daily
	TemperatureTrend           string `json:"temperatureTrend,omitempty"`
	ProbabilityOfPrecipitation *struct {
		UnitCode string  `json:"unitCode"`
		Value    float64 `json:"value,omitempty"`
	} `json:"probabilityOfPrecipitation,omitempty"`
	Dewpoint *struct {
		UnitCode string  `json:"unitCode"`
		Value    float64 `json:"value"`
	} `json:"dewpoint,omitempty"`
	RelativeHumidity *struct {
		UnitCode string  `json:"unitCode"`
		Value    float64 `json:"value"`
	} `json:"relativeHumidity,omitempty"`
	WindSpeed *WindSpeed `json:"windSpeed,omitempty"`
	WindGust  *struct {
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

// WeatherAPI defines the interface for interacting with the weather.gov API.
type WeatherAPI interface {
	GetDailyForecast(ctx context.Context, opts ...GetForcastOption) (DailyForecast, error)
	InitForecastAPI(ctx context.Context, latitude, longitude *float64) error
	SetCoordinates(latitude, longitude *float64)
}

// weatherGovAPI is a concrete implementation of the WeatherAPI interface.
type weatherGovAPI struct {
	client      *http.Client
	baseURL     string
	coordinates struct {
		latitude  *float64
		longitude *float64
	}
	forecastProperties ForecastAPIProps
}

// NewWeatherGovAPI creates a new instance of weatherGovAPI.
func NewWeatherGovAPI(client *http.Client, latitude, longitude *float64) WeatherAPI {
	return &weatherGovAPI{
		client:  client,
		baseURL: BASE_URL,
		coordinates: struct {
			latitude  *float64
			longitude *float64
		}{
			latitude:  latitude,
			longitude: longitude,
		},
	}
}

// Forecast represents the weather forecast data.
type Forecast struct {
	Properties struct {
		Periods []struct {
			Name          string `json:"name"`
			Temperature   int    `json:"temperature"`
			WindSpeed     string `json:"windSpeed"`
			ShortForecast string `json:"shortForecast"`
		} `json:"periods"`
	} `json:"properties"`
}

type Units int

const (
	US Units = iota
	SI
)

func (u Units) String() string {
	switch u {
	case US:
		return "us"
	case SI:
		return "si"
	default:
		return "us" // default units
	}
}

// defines type for functional options
type GetForcastOption func(*GetForecastOptions)

// hold options for Forecast options
type GetForecastOptions struct {
	Units          Units
	UseQuantValues bool
}

// Sets the units option
func WithUnits(units Units) GetForcastOption {
	return func(opts *GetForecastOptions) {
		opts.Units = units
	}
}

// Build the URL used for forecasts
func buildForecastURL(latitude, longitude float64) string {
	forecastResource := fmt.Sprintf("%s/points/%f,%f", BASE_URL, latitude, longitude)
	return forecastResource
}

// Initialize the API to fetch weather forecasts
func (api *weatherGovAPI) InitForecastAPI(ctx context.Context, latitude, longitude *float64) error {
	if api.coordinates.latitude == nil || api.coordinates.longitude == nil {
		// no available coordinate data, error
		if latitude == nil || longitude == nil {
			return fmt.Errorf("no available latitude and longitude")
		}

		// if the coordinates aren't already set, set them with the given values
		api.SetCoordinates(latitude, longitude)
	}

	activeLat := *api.coordinates.latitude
	activeLong := *api.coordinates.longitude

	if latitude != nil && longitude != nil {
		activeLat = *latitude
		activeLong = *longitude
	}

	url := buildForecastURL(activeLat, activeLong)

	// create a request with a context
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	response, err := api.client.Do(request)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	// make sure the response is good
	if response.StatusCode != http.StatusOK {
		return httpStatusError(response.StatusCode)
	}

	var apiResponse ForcastAPIResponse
	err = readBody(response.Body, &apiResponse)
	if err != nil {
		return err
	}

	api.forecastProperties = apiResponse.Properties

	// cool, no errors!
	return nil
}

// Set the geographic coordinates for the API to use
func (api *weatherGovAPI) SetCoordinates(latitude, longitude *float64) {
	api.coordinates.latitude = latitude
	api.coordinates.longitude = longitude
}

// Get the daily forecast for the configured latitude and longitude.
//
// The forecast is for a seven day period with weather results for that day and the "night" of that
// day.
func (api *weatherGovAPI) GetDailyForecast(ctx context.Context, opts ...GetForcastOption) (DailyForecast, error) {
	// set default values for options
	options := &GetForecastOptions{
		Units:          US,
		UseQuantValues: true,
	}

	// Apply provided options
	for _, opt := range opts {
		opt(options)
	}

	// parse base url
	u, err := url.Parse(api.forecastProperties.Forecast)
	if err != nil {
		return DailyForecast{}, err
	}

	// create query string
	params := url.Values{}
	params.Add("units", options.Units.String())

	// add query to url
	u.RawQuery = params.Encode()

	var url = u.String()
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return DailyForecast{}, err
	}

	// set quantitative values header
	if options.UseQuantValues {
		request.Header.Set("Feature-Flags", "forecast_temperature_qv,forecast_wind_speed_qv")
	}

	response, err := api.client.Do(request)
	if err != nil {
		return DailyForecast{}, err
	}

	defer response.Body.Close()

	// make sure the response is good
	if response.StatusCode != http.StatusOK {
		return DailyForecast{}, httpStatusError(response.StatusCode)
	}

	var dailyForecastResponse DailyForecastResponse
	err = readBody(response.Body, &dailyForecastResponse)
	if err != nil {
		return DailyForecast{}, err
	}

	return dailyForecastResponse.Properties, nil
}

// create an HTTP status error
func httpStatusError(statusCode int) error {
	return fmt.Errorf("received status code %d", statusCode)
}

// Read a response body
func readBody(body io.Reader, v any) error {
	return json.NewDecoder(body).Decode(v)
}
