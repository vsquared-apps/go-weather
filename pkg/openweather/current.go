package openweather

import (
	"context"
	"fmt"
	"net/http"
)

// CurrentService communicates with the current weather data related methods of the OpenWeatherMap API.
//
// API docs: https://openweathermap.org/current
type CurrentService struct {
	client *Client
}

// Location holds information for a location.
type Location struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`

	// Seconds east of UTC.
	UTCOffset   int          `json:"timezone"`
	Coordinates *Coordinates `json:"coord,omitempty"`

	Weather    *Weather `json:"main,omitempty"`
	Visibility int      `json:"visibility"`
	Wind       *Wind    `json:"wind,omitempty"`
	Clouds     *Clouds  `json:"clouds,omitempty"`
	Rain       *Rain    `json:"rain,omitempty"`
	Snow       *Snow    `json:"snow,omitempty"`
}

// Coordinates holds geographic coordinates for a location.
type Coordinates struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lon"`
}

// Weather holds weather information for a location.
type Weather struct {
	Temperature         float64 `json:"temp"`
	FeelsLike           float64 `json:"feels_like"`
	Min                 float64 `json:"temp_min"`
	Max                 float64 `json:"temp_max"`
	Pressure            float64 `json:"pressure"`
	Humidity            float64 `json:"humidity"`
	SeaLevelPressure    float64 `json:"sea_level"`
	GroundLevelPressure float64 `json:"grnd_level"`
}

// Wind holds wind information for a location.
type Wind struct {
	Speed float64 `json:"speed"`
	// Degrees.
	Direction int     `json:"deg"`
	Gust      float64 `json:"gust"`
}

// Clouds holds cloud information for a location.
type Clouds struct {
	// Percentage.
	Cloudiness int `json:"all"`
}

// Rain holds rain information (in millimeters) for a location.
type Rain struct {
	PastHour   float64 `json:"1h"`
	Past3Hours float64 `json:"3h"`
}

// Snow holds snow information (in millimeters) for a location.
type Snow struct {
	PastHour   float64 `json:"1h"`
	Past3Hours float64 `json:"3h"`
}

// ByCity gets the current weather data by a city's name.
func (s *CurrentService) ByCity(ctx context.Context, name string) (*Location, *http.Response, error) {
	path := fmt.Sprintf("data/2.5/weather?q=%s", name)

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(Location)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, nil
}
