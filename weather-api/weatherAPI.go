package weatherAPI

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// WeatherAPI defines the interface for interacting with the weather.gov API.
type WeatherAPI interface {
	GetForecast(ctx context.Context, latitude, longitude float64) (*Forecast, error)
}

// weatherGovAPI is a concrete implementation of the WeatherAPI interface.
type weatherGovAPI struct {
	client  *http.Client
	baseURL string
}

// NewWeatherGovAPI creates a new instance of weatherGovAPI.
func NewWeatherGovAPI(client *http.Client) WeatherAPI {
	return &weatherGovAPI{
		client:  client,
		baseURL: "https://api.weather.gov",
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

// GetForecast fetches the weather forecast for the given latitude and longitude.
func (api *weatherGovAPI) GetForecast(ctx context.Context, latitude, longitude float64) (*Forecast, error) {
	url := fmt.Sprintf("%s/points/%f,%f/forecast", api.baseURL, latitude, longitude)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/geo+json")
	resp, err := api.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var forecast Forecast
	if err := json.NewDecoder(resp.Body).Decode(&forecast); err != nil {
		return nil, err
	}

	return &forecast, nil
}
