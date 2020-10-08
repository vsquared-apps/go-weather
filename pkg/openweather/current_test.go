package openweather

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

var expectedLocation = &Location{
	ID:   6077243,
	Name: "Montreal",

	UTCOffset: -14400,
	Coordinates: &Coordinates{
		Latitude:  45.51,
		Longitude: -73.59,
	},

	Weather: &Weather{
		Temperature: 278.89,
		FeelsLike:   273.59,
		Min:         278.15,
		Max:         279.26,
		Pressure:    1018,
		Humidity:    52,
	},
	Visibility: 10000,
	Wind: &Wind{
		Speed:     4.1,
		Direction: 310,
	},
	Clouds: &Clouds{
		Cloudiness: 20,
	},
}

func TestCurrentService_ByCity(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	blob, err := readFileContents("testdata/current/location.json")
	require.NoError(t, err)

	mux.HandleFunc("/data/2.5/weather", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "Montreal", r.URL.Query().Get("q"))
		fmt.Fprint(w, blob)
	})

	location, _, err := client.Current.ByCity(ctx, "Montreal")
	require.NoError(t, err)
	require.Equal(t, expectedLocation, location)
}
