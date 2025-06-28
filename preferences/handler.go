package preferences

type PrefHandler struct {
	repo *PrefRepository
}

// NewPrefHandler creates a new PrefHandler instance
// with the provided PrefRepository
// It initializes the handler with the repository
// to manage user preferences.
func NewPrefHandler(repo *PrefRepository) *PrefHandler {
	return &PrefHandler{repo}
}

// NewPreference creates a new Preference instance with default values
func (h *PrefHandler) NewPreference(userID int64) *Preference {
	return &Preference{
		UserID:           userID,
		Units:            "us",
		TemperatureMin:   23.9,
		TemperatureMax:   38.0,
		WindMin:          0,
		WindMax:          4,
		PrecipitationMin: 0.0,
		PrecipitationMax: 10.0,
		Locale:           "en-US",
	}
}
