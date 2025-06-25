package weatherAPI

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"
)

// mockRoundTripper implements http.RoundTripper for mocking HTTP responses
type mockRoundTripper struct {
	resp *http.Response
	err  error
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.resp, m.err
}

func newMockClient(body string, status int) *http.Client {
	return &http.Client{
		Transport: &mockRoundTripper{
			resp: &http.Response{
				StatusCode: status,
				Body:       io.NopCloser(bytes.NewBufferString(body)),
				Header:     make(http.Header),
			},
		},
	}
}

func TestWithUnits(t *testing.T) {
	opts := &GetForecastOptions{}
	WithUnits(SI)(opts)
	if opts.Units != SI {
		t.Errorf("expected SI units, got %v", opts.Units)
	}
}

func TestBuildForecastURL(t *testing.T) {
	lat, lon := 40.0, -75.0
	expected := "https://api.weather.gov/points/40.000000,-75.000000"
	got := buildForecastURL(lat, lon)
	if got != expected {
		t.Errorf("expected %s, got %s", expected, got)
	}
}

func TestSetCoordinates(t *testing.T) {
	api := &weatherGovAPI{}
	lat, lon := 10.0, 20.0
	api.SetCoordinates(&lat, &lon)
	if api.coordinates.latitude == nil || api.coordinates.longitude == nil {
		t.Fatal("coordinates not set")
	}
	if *api.coordinates.latitude != lat || *api.coordinates.longitude != lon {
		t.Errorf("expected %v,%v got %v,%v", lat, lon, *api.coordinates.latitude, *api.coordinates.longitude)
	}
}

func TestInitForecastAPI_Success(t *testing.T) {
	lat, lon := 40.0, -75.0
	mockBody := `{"properties":{"@id":"id","@type":"type","cwa":"cwa","forecastOffice":"office","gridId":"grid","gridX":1,"gridY":2,"forecast":"https://api.weather.gov/forecast","forecastHourly":"https://api.weather.gov/hourly","forecastGridData":"gridData","observationStations":"stations"}}`
	client := newMockClient(mockBody, http.StatusOK)
	api := NewWeatherGovAPI(client, &lat, &lon).(*weatherGovAPI)
	ctx := context.Background()
	err := api.InitForecastAPI(ctx, &lat, &lon)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if api.forecastProperties.Forecast != "https://api.weather.gov/forecast" {
		t.Errorf("unexpected forecast property: %v", api.forecastProperties.Forecast)
	}
}

func TestInitForecastAPI_ErrorStatus(t *testing.T) {
	lat, lon := 40.0, -75.0
	client := newMockClient("{}", http.StatusBadRequest)
	api := NewWeatherGovAPI(client, &lat, &lon).(*weatherGovAPI)
	ctx := context.Background()
	err := api.InitForecastAPI(ctx, &lat, &lon)
	if err == nil {
		t.Fatal("expected error for bad status code")
	}
}

func TestGetDailyForecast_Success(t *testing.T) {
	lat, lon := 40.0, -75.0
	forecastURL := "https://api.weather.gov/forecast"
	mockBody := `{"properties":{"units":"us","forecastGenerator":"gen","generatedAt":"2024-01-01T00:00:00Z","updateTime":"2024-01-01T01:00:00Z","validTimes":"2024-01-01T00:00:00Z/2024-01-02T00:00:00Z","elevation":{"unitCode":"unit","value":10},"periods":[]}}`
	client := newMockClient(mockBody, http.StatusOK)
	api := NewWeatherGovAPI(client, &lat, &lon).(*weatherGovAPI)
	api.forecastProperties.Forecast = forecastURL
	ctx := context.Background()
	forecast, err := api.GetDailyForecast(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if forecast.Units != "us" {
		t.Errorf("expected units 'us', got %v", forecast.Units)
	}
}

func TestGetDailyForecast_ErrorStatus(t *testing.T) {
	lat, lon := 40.0, -75.0
	client := newMockClient("{}", http.StatusBadRequest)
	api := NewWeatherGovAPI(client, &lat, &lon).(*weatherGovAPI)
	api.forecastProperties.Forecast = "https://api.weather.gov/forecast"
	ctx := context.Background()
	_, err := api.GetDailyForecast(ctx)
	if err == nil {
		t.Fatal("expected error for bad status code")
	}
}

func TestGetHourlyForecast_Success(t *testing.T) {
	lat, lon := 40.0, -75.0
	hourlyURL := "https://api.weather.gov/hourly"
	mockBody := `{"properties":{"units":"us","forecastGenerator":"gen","generatedAt":"2024-01-01T00:00:00Z","updateTime":"2024-01-01T01:00:00Z","validTimes":"2024-01-01T00:00:00Z/2024-01-02T00:00:00Z","elevation":{"unitCode":"unit","value":10},"periods":[]}}`
	client := newMockClient(mockBody, http.StatusOK)
	api := NewWeatherGovAPI(client, &lat, &lon).(*weatherGovAPI)
	api.forecastProperties.ForecastHourly = hourlyURL
	ctx := context.Background()
	forecast, err := api.GetHourlyForecast(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if forecast.Units != "us" {
		t.Errorf("expected units 'us', got %v", forecast.Units)
	}
}

func TestGetHourlyForecast_ErrorStatus(t *testing.T) {
	lat, lon := 40.0, -75.0
	client := newMockClient("{}", http.StatusBadRequest)
	api := NewWeatherGovAPI(client, &lat, &lon).(*weatherGovAPI)
	api.forecastProperties.ForecastHourly = "https://api.weather.gov/hourly"
	ctx := context.Background()
	_, err := api.GetHourlyForecast(ctx)
	if err == nil {
		t.Fatal("expected error for bad status code")
	}
}

func TestUnits_String(t *testing.T) {
	tests := []struct {
		u        Units
		expected string
	}{
		{US, "us"},
		{SI, "si"},
		{Units(99), "us"},
	}
	for _, tt := range tests {
		if got := tt.u.String(); got != tt.expected {
			t.Errorf("Units(%d).String() = %s, want %s", tt.u, got, tt.expected)
		}
	}
}

func TestReadBody_InvalidJSON(t *testing.T) {
	var v struct{}
	err := readBody(bytes.NewBufferString("{invalid json"), &v)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestHttpStatusError(t *testing.T) {
	err := httpStatusError(404)
	if err == nil || err.Error() != "received status code 404" {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestPeriodJSONUnmarshal(t *testing.T) {
	// Test that Period struct can unmarshal a minimal valid JSON
	jsonStr := `{
		"number": 1,
		"name": "Today",
		"startTime": "2024-01-01T00:00:00Z",
		"endTime": "2024-01-01T12:00:00Z",
		"isDaytime": true,
		"windDirection": "N",
		"icon": "icon",
		"shortForecast": "Sunny",
		"detailedForecast": "Clear"
	}`
	var p Period
	err := readBody(bytes.NewBufferString(jsonStr), &p)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if p.Name != "Today" || !p.IsDaytime {
		t.Errorf("unexpected period values: %+v", p)
	}
}

func TestNewWeatherGovAPI(t *testing.T) {
	client := &http.Client{}
	lat, lon := 1.0, 2.0
	api := NewWeatherGovAPI(client, &lat, &lon)
	if api == nil {
		t.Fatal("expected non-nil WeatherAPI")
	}
}
