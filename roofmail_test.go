package main

import (
	"bytes"
	"context"
	"log"
	"os"
	"strconv"
	"testing"

	wapi "roofmail/weatherAPI"
)

// --- Mock WeatherAPI ---

type mockWeatherAPI struct {
	initErr       error
	forecastErr   error
	dailyForecast wapi.DailyForecast
}

func (m *mockWeatherAPI) InitForecastAPI(ctx context.Context, a, b *float64) error {
	return m.initErr
}
func (m *mockWeatherAPI) GetDailyForecast(ctx context.Context, opts ...wapi.GetForcastOption) (wapi.DailyForecast, error) {
	return m.dailyForecast, m.forecastErr
}
func (m *mockWeatherAPI) SetCoordinates(lat, lon *float64) {}

// --- Helper functions ---

func setEnv(key, value string) func() {
	old := os.Getenv(key)
	os.Setenv(key, value)
	return func() { os.Setenv(key, old) }
}

func mockLogs() (restore func()) {
	var buf bytes.Buffer
	infoLogger = log.New(&buf, "INFO: ", 0)
	debugLogger = log.New(&buf, "DEBUG: ", 0)
	return func() {
		infoLogger = nil
		debugLogger = nil
	}
}

func captureOutput(f func()) string {
	r, w, _ := os.Pipe()
	stdout := os.Stdout
	os.Stdout = w
	defer func() { os.Stdout = stdout }()

	outC := make(chan string)
	go func() {
		var buf bytes.Buffer
		_, _ = buf.ReadFrom(r)
		outC <- buf.String()
	}()

	f()
	w.Close()
	out := <-outC
	return out
}

// --- Tests ---

func TestCtof(t *testing.T) {
	setEnv("APP_ENV", "test")()
	restore := mockLogs()
	defer restore()

	if got := ctof(0); got != 32 {
		t.Errorf("ctof(0) = %v, want 32", got)
	}
	if got := ctof(100); got != 212 {
		t.Errorf("ctof(100) = %v, want 212", got)
	}
}

func TestKphToMph(t *testing.T) {
	setEnv("APP_ENV", "test")()
	restore := mockLogs()
	defer restore()

	if got := kphToMph(16.09344); got < 9.99 || got > 10.01 {
		t.Errorf("kphToMph(16.09344) = %v, want ~10", got)
	}
}

func TestBeaufortScale(t *testing.T) {
	setEnv("APP_ENV", "test")()
	restore := mockLogs()
	defer restore()

	tests := []struct {
		wind wapi.WindSpeed
		want int
	}{
		{wapi.WindSpeed{Value: floatPtr(0.2), UnitCode: "wmoUnit:m_s-1"}, 0},
		{wapi.WindSpeed{Value: floatPtr(1), UnitCode: "wmoUnit:m_s-1"}, 1},
		{wapi.WindSpeed{Value: floatPtr(2), UnitCode: "wmoUnit:m_s-1"}, 2},
		{wapi.WindSpeed{Value: floatPtr(4.25), UnitCode: "wmoUnit:m_s-1"}, 3},
		{wapi.WindSpeed{Value: floatPtr(6.5), UnitCode: "wmoUnit:m_s-1"}, 4},
		{wapi.WindSpeed{Value: floatPtr(25), UnitCode: "wmoUnit:m_s-1"}, 5},
		{wapi.WindSpeed{Value: nil, MaxValue: floatPtr(1.5), UnitCode: "wmoUnit:m_s-1"}, 1},
		{wapi.WindSpeed{}, 0},
	}
	for _, tt := range tests {
		if got := beaufortScale(tt.wind); got != tt.want {
			t.Errorf("beaufortScale(%v) = %v, want %v", tt.wind, got, tt.want)
		}
	}
}

func TestGetTempF(t *testing.T) {
	setEnv("APP_ENV", "test")()
	restore := mockLogs()
	defer restore()

	period := wapi.Period{Temperature: &wapi.UnitValue{Value: 20, UnitCode: "wmoUnit:degC"}}
	if got := getTempF(period); got < 67.9 || got > 68.1 {
		t.Errorf("getTempF(C) = %v, want ~68", got)
	}
	period = wapi.Period{Temperature: &wapi.UnitValue{Value: 70, UnitCode: "wmoUnit:degF"}}
	if got := getTempF(period); got != 70 {
		t.Errorf("getTempF(F) = %v, want 70", got)
	}
}

func TestIsComfortable(t *testing.T) {
	setEnv("APP_ENV", "test")()
	restore := mockLogs()
	defer restore()

	period := wapi.Period{
		Temperature:                &wapi.UnitValue{Value: 30, UnitCode: "wmoUnit:degC"}, // 70F
		WindSpeed:                  &wapi.WindSpeed{Value: floatPtr(0.5), UnitCode: "wmoUnit:m_s-1"},
		ProbabilityOfPrecipitation: &wapi.UnitValue{Value: 0.0, UnitCode: "wmoUnit:percent"},
	}

	if !isComfortable(period) {
		t.Error("Expected comfortable")
	}
	period.Temperature.Value = 10 // 50F
	if isComfortable(period) {
		t.Error("Expected not comfortable (too cold)")
	}
	period.Temperature.Value = 21.1
	period.WindSpeed.Value = floatPtr(10) // high wind
	if isComfortable(period) {
		t.Error("Expected not comfortable (windy)")
	}
}

func TestComfortMessage(t *testing.T) {
	setEnv("APP_ENV", "test")()
	restore := mockLogs()
	defer restore()

	period := wapi.Period{
		Temperature:                &wapi.UnitValue{Value: 21.1, UnitCode: "wmoUnit:degC"},
		WindSpeed:                  &wapi.WindSpeed{Value: floatPtr(2), UnitCode: "wmoUnit:m_s-1"},
		ProbabilityOfPrecipitation: &wapi.UnitValue{Value: 32.1, UnitCode: "wmoUnit:percent"},
	}
	msg := comfortMessage(period)
	if msg == "" || msg[0] != 'I' {
		t.Errorf("comfortMessage = %q, want non-empty string", msg)
	}
}

func TestInitLogs(t *testing.T) {
	setEnv("APP_ENV", "test")()
	restore := mockLogs()
	defer restore()

	initLogs()
	if infoLogger == nil || debugLogger == nil {
		t.Error("Loggers not initialized")
	}
}

func TestLoadConfig(t *testing.T) {
	setEnv("APP_ENV", "test")()
	restore := mockLogs()
	defer restore()

	cfg := loadConfig()
	if cfg.Version != "0.0.0" {
		t.Errorf("loadConfig() = %v, want version 0.0.0", cfg.Version)
	}
}

// --- Integration-like test for main logic ---

func TestMainLogic_BadEnv(t *testing.T) {
	setEnv("APP_ENV", "test")()
	restore := mockLogs()
	defer restore()

	unsetLat := setEnv("LATITUDE", "notafloat")
	defer unsetLat()
	unsetLon := setEnv("LONGITUDE", "-75.0")
	defer unsetLon()
	unsetEnv := setEnv("APP_ENV", "test")
	defer unsetEnv()

	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(os.Stderr)

	initLogs()
	config := loadConfig()
	infoLogger.Printf("Starting Roofmail v%s", config.Version)
	debugLogger.Println("Enabled")
	_, err := strconv.ParseFloat(os.Getenv("LATITUDE"), 64)
	if err == nil {
		t.Error("Expected error parsing LATITUDE")
	}
	log.SetOutput(os.Stderr)
}

// --- Utility ---

func floatPtr(f float64) *float64 {
	return &f
}
