package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/jorgegabrielti/desafio-clima-cep/pkg/viacep"
	"github.com/jorgegabrielti/desafio-clima-cep/pkg/weather"
)

type LocationFetcher interface {
	FetchLocation(ctx context.Context, zip string) (string, error)
}

type WeatherFetcher interface {
	FetchTemperature(ctx context.Context, city string) (float64, error)
}

type WeatherHandler struct {
	LocationClient LocationFetcher
	WeatherClient  WeatherFetcher
}

func NewWeatherHandler(loc LocationFetcher, wea WeatherFetcher) *WeatherHandler {
	return &WeatherHandler{
		LocationClient: loc,
		WeatherClient:  wea,
	}
}

func (h *WeatherHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	zipParam := extractZipcode(r)

	if zipParam == "" {
		sendPlainError(w, http.StatusUnprocessableEntity, "invalid zipcode")
		return
	}

	cleanZip, err := viacep.ValidateZipcode(zipParam)
	if err != nil {
		sendPlainError(w, http.StatusUnprocessableEntity, "invalid zipcode")
		return
	}

	city, err := h.LocationClient.FetchLocation(r.Context(), cleanZip)
	if err != nil {
		if errors.Is(err, viacep.ErrInvalidZipcode) {
			sendPlainError(w, http.StatusUnprocessableEntity, "invalid zipcode")
			return
		}
		sendPlainError(w, http.StatusNotFound, "can not find zipcode")
		return
	}

	tempC, err := h.WeatherClient.FetchTemperature(r.Context(), city)
	if err != nil {
		sendPlainError(w, http.StatusNotFound, "can not find zipcode")
		return
	}

	response := weather.CalculateTemperatures(tempC)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

func extractZipcode(r *http.Request) string {
	if zip := r.URL.Query().Get("zipcode"); zip != "" {
		return zip
	}
	if zip := r.URL.Query().Get("cep"); zip != "" {
		return zip
	}

	path := strings.Trim(r.URL.Path, "/")
	parts := strings.Split(path, "/")
	if len(parts) > 0 && parts[len(parts)-1] != "" && parts[len(parts)-1] != "weather" {
		return parts[len(parts)-1]
	}

	return ""
}

func sendPlainError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(statusCode)
	_, _ = w.Write([]byte(message))
}
