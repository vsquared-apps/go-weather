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

// ByCity gets the current weather data by a city's name.
// todo: return a pre-defined struct, not an empty interface
func (s *CurrentService) ByCity(ctx context.Context, name string) (interface{}, *http.Response, error) {
	path := fmt.Sprintf("data/2.5/weather?q=%s", name)

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var root interface{}
	resp, err := s.client.Do(ctx, req, &root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, nil
}
