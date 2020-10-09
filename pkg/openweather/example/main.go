package main

import (
	"context"
	"fmt"
	"log"

	"github.com/vsquared-apps/go-weather/pkg/openweather"
)

var ctx = context.Background()

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	client, err := openweather.New("[API KEY]", openweather.UseMetricUnits)
	if err != nil {
		return err
	}

	loc, _, err := client.Current.ByCity(ctx, "Montreal")
	if err != nil {
		return err
	}

	fmt.Printf("It is currently %.1f degrees Celsius in %s.\n", loc.Weather.Temperature, loc.Name)
	return nil
}
