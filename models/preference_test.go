package models

import (
	"testing"
)

func TestValidate(t *testing.T) {
	pref := &Preference{
		TemperatureMin:   50.0,
		TemperatureMax:   30.0,
		WindMin:          10.0,
		WindMax:          5.0,
		PrecipitationMin: 20.0,
		PrecipitationMax: 15.0,
	}

	if err := pref.Validate(); err == nil {
		t.Error("expected validation error for invalid preference, got nil")
	}
	pref.TemperatureMin = 20.0
	pref.TemperatureMax = 30.0

	if err := pref.Validate(); err == nil {
		t.Error("expected validation error for invalid preference, got:", err)
	}

	pref.WindMin = 0.0
	pref.WindMax = 10.0

	if err := pref.Validate(); err == nil {
		t.Error("expected validation error for invalid preference, got:", err)
	}

	pref.PrecipitationMin = 0.0
	pref.PrecipitationMax = 10.0

	if err := pref.Validate(); err != nil {
		t.Error("expected valid preference, got error:", err)
	}
}
