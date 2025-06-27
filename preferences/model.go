package preferences

import (
	"roofmail/users"

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

	// Define the foreign key relationship
	User *users.User `bun:"rel:belongs-to,join:user_id=id"`
}
