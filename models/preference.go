package models

import (
	"fmt"

	"github.com/uptrace/bun"
)

type Preference struct {
	bun.BaseModel

	UserID           int64   `bun:"user_id,pk,notnull" json:"userId"`
	Units            string  `bun:",notnull,default:'us'" json:"units"`
	TemperatureMin   float64 `bun:",default:23.9" json:"temperatureMin"`
	TemperatureMax   float64 `bun:",default:38.0" json:"temperatureMax"`
	WindMin          float64 `bun:",default:0"  json:"windMin"`
	WindMax          float64 `bun:",default:4" json:"windMax"`
	PrecipitationMin float64 `bun:",default:0.0" json:"precipitationMin"`
	PrecipitationMax float64 `bun:",default:10.0" json:"precipitationMax"`
	Locale           string  `bun:",notnull,default:'en-US'" json:"locale"`

	// Define the foreign key relationship
	User *User `bun:"rel:belongs-to,join:user_id=id"`
}

func (p *Preference) Validate() error {
	// Add validation logic if needed
	// For example, check if temperature ranges are valid
	if p.TemperatureMin >= p.TemperatureMax {
		return fmt.Errorf("temperatureMin must be less than temperatureMax")
	}
	if p.WindMin >= p.WindMax {
		return fmt.Errorf("windMin must be less than windMax")
	}
	if p.PrecipitationMin >= p.PrecipitationMax {
		return fmt.Errorf("precipitationMin must be less than precipitationMax")
	}
	return nil
}
