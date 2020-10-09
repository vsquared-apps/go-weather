# openweather

openweather is a Go client library for accessing the OpenWeather API.

You can view OpenWeather's official API documentation [here](https://openweathermap.org/api).

## Usage

Make sure to have an OpenWeather account to get an API key. [Here](https://openweathermap.org/appid) is a quick guide on how to get one.

```go
package main

import "github.com/vsquared-apps/go-weather/pkg/openweather"

func main() {
    client, _ := openweather.New("[API KEY]")
}
```

You can pass in a number of options to `New` to further configure the client (see [openweather-options.go](openweather-options.go)). For example, to get responses in metric units:

```go
client, _ := openweather.New("[API KEY]", openweather.UseMetricUnits)
```
