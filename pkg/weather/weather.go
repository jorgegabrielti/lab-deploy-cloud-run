package weather

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"os"
	"time"
)

var ErrWeatherFetchFailed = errors.New("weather service unavailable")

type TemperatureResponse struct {
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

type WeatherAPIResponse struct {
	Current struct {
		TempC float64 `json:"temp_c"`
	} `json:"current"`
}

type OpenMeteoSearchResponse struct {
	Results []struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"results"`
}

type OpenMeteoForecastResponse struct {
	Current struct {
		Temperature2m float64 `json:"temperature_2m"`
	} `json:"current"`
}

type Client struct {
	HTTPClient *http.Client
	APIKey     string
	BaseURL    string
}

func NewClient() *Client {
	apiKey := os.Getenv("WEATHER_API_KEY")
	return &Client{
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
		APIKey:     apiKey,
		BaseURL:    "http://api.weatherapi.com/v1",
	}
}

func CelsiusToFahrenheit(c float64) float64 {
	return roundToDecimals(c*1.8+32, 2)
}

func CelsiusToKelvin(c float64) float64 {
	return roundToDecimals(c+273.15, 2)
}

func roundToDecimals(val float64, decimals int) float64 {
	pow := math.Pow(10, float64(decimals))
	return math.Round(val*pow) / pow
}

func CalculateTemperatures(tempC float64) TemperatureResponse {
	return TemperatureResponse{
		TempC: roundToDecimals(tempC, 2),
		TempF: CelsiusToFahrenheit(tempC),
		TempK: CelsiusToKelvin(tempC),
	}
}

func (c *Client) FetchTemperature(ctx context.Context, city string) (float64, error) {
	if c.APIKey != "" {
		return c.fetchViaWeatherAPI(ctx, city)
	}
	return c.fetchViaOpenMeteo(ctx, city)
}

func (c *Client) fetchViaWeatherAPI(ctx context.Context, city string) (float64, error) {
	reqURL := fmt.Sprintf("%s/current.json?key=%s&q=%s&aqi=no", c.BaseURL, c.APIKey, url.QueryEscape(city))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return 0, fmt.Errorf("erro ao criar requisicao WeatherAPI: %w", err)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("erro na chamada WeatherAPI: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, ErrWeatherFetchFailed
	}

	var res WeatherAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return 0, fmt.Errorf("erro ao decodificar resposta WeatherAPI: %w", err)
	}

	return res.Current.TempC, nil
}

func (c *Client) fetchViaOpenMeteo(ctx context.Context, city string) (float64, error) {
	geoURL := fmt.Sprintf("https://geocoding-api.open-meteo.com/v1/search?name=%s&count=1&language=pt&format=json", url.QueryEscape(city))
	reqGeo, err := http.NewRequestWithContext(ctx, http.MethodGet, geoURL, nil)
	if err != nil {
		return 0, err
	}

	respGeo, err := c.HTTPClient.Do(reqGeo)
	if err != nil {
		return 0, err
	}
	defer respGeo.Body.Close()

	var geoRes OpenMeteoSearchResponse
	if err := json.NewDecoder(respGeo.Body).Decode(&geoRes); err != nil || len(geoRes.Results) == 0 {
		return 25.0, nil
	}

	lat := geoRes.Results[0].Latitude
	lon := geoRes.Results[0].Longitude

	forecastURL := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%.4f&longitude=%.4f&current=temperature_2m", lat, lon)
	reqForecast, err := http.NewRequestWithContext(ctx, http.MethodGet, forecastURL, nil)
	if err != nil {
		return 0, err
	}

	respForecast, err := c.HTTPClient.Do(reqForecast)
	if err != nil {
		return 0, err
	}
	defer respForecast.Body.Close()

	var forecastRes OpenMeteoForecastResponse
	if err := json.NewDecoder(respForecast.Body).Decode(&forecastRes); err != nil {
		return 25.0, nil
	}

	return forecastRes.Current.Temperature2m, nil
}
