package weather_test

import (
	"math"
	"testing"

	"github.com/jorgegabrielti/desafio-clima-cep/pkg/weather"
)

func TestCelsiusToFahrenheit(t *testing.T) {
	cases := []struct {
		celsius float64
		want    float64
	}{
		{0.0, 32.0},
		{28.5, 83.3},
		{100.0, 212.0},
		{-10.0, 14.0},
	}

	for _, tc := range cases {
		got := weather.CelsiusToFahrenheit(tc.celsius)
		if math.Abs(got-tc.want) > 0.05 {
			t.Errorf("CelsiusToFahrenheit(%.1f) = %.2f, esperado %.2f", tc.celsius, got, tc.want)
		}
	}
}

func TestCelsiusToKelvin(t *testing.T) {
	cases := []struct {
		celsius float64
		want    float64
	}{
		{0.0, 273.15},
		{28.5, 301.65},
		{100.0, 373.15},
		{-273.15, 0.0},
	}

	for _, tc := range cases {
		got := weather.CelsiusToKelvin(tc.celsius)
		if math.Abs(got-tc.want) > 0.05 {
			t.Errorf("CelsiusToKelvin(%.1f) = %.2f, esperado %.2f", tc.celsius, got, tc.want)
		}
	}
}

func TestCalculateTemperatures(t *testing.T) {
	res := weather.CalculateTemperatures(28.5)
	if res.TempC != 28.5 {
		t.Errorf("TempC = %.2f, esperado 28.5", res.TempC)
	}
	if math.Abs(res.TempF-83.3) > 0.05 {
		t.Errorf("TempF = %.2f, esperado 83.3", res.TempF)
	}
	if math.Abs(res.TempK-301.65) > 0.05 {
		t.Errorf("TempK = %.2f, esperado 301.65", res.TempK)
	}
}
