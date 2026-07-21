package handler_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jorgegabrielti/desafio-clima-cep/pkg/handler"
	"github.com/jorgegabrielti/desafio-clima-cep/pkg/viacep"
	"github.com/jorgegabrielti/desafio-clima-cep/pkg/weather"
)

type mockLocationFetcher struct {
	city string
	err  error
}

func (m *mockLocationFetcher) FetchLocation(ctx context.Context, zip string) (string, error) {
	if m.err != nil {
		return "", m.err
	}
	return m.city, nil
}

type mockWeatherFetcher struct {
	tempC float64
	err   error
}

func (m *mockWeatherFetcher) FetchTemperature(ctx context.Context, city string) (float64, error) {
	if m.err != nil {
		return 0, m.err
	}
	return m.tempC, nil
}

func TestWeatherHandler_Success(t *testing.T) {
	h := handler.NewWeatherHandler(
		&mockLocationFetcher{city: "São Paulo"},
		&mockWeatherFetcher{tempC: 28.5},
	)

	req := httptest.NewRequest(http.MethodGet, "/?zipcode=01001000", nil)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status code = %d, esperado %d", rec.Code, http.StatusOK)
	}

	var resp weather.TemperatureResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("erro ao unmarshal resposta: %v", err)
	}

	if resp.TempC != 28.5 {
		t.Errorf("temp_C = %.2f, esperado 28.5", resp.TempC)
	}
	if resp.TempF != 83.3 {
		t.Errorf("temp_F = %.2f, esperado 83.3", resp.TempF)
	}
	if resp.TempK != 301.65 {
		t.Errorf("temp_K = %.2f, esperado 301.65", resp.TempK)
	}
}

func TestWeatherHandler_InvalidZipcode(t *testing.T) {
	h := handler.NewWeatherHandler(
		&mockLocationFetcher{},
		&mockWeatherFetcher{},
	)

	invalidZips := []string{"123", "0100100a", "123456789"}

	for _, zip := range invalidZips {
		req := httptest.NewRequest(http.MethodGet, "/?zipcode="+zip, nil)
		rec := httptest.NewRecorder()

		h.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnprocessableEntity {
			t.Errorf("zipcode %q: status code = %d, esperado 422", zip, rec.Code)
		}
		if body := rec.Body.String(); body != "invalid zipcode" {
			t.Errorf("zipcode %q: body = %q, esperado 'invalid zipcode'", zip, body)
		}
	}
}

func TestWeatherHandler_NotFoundZipcode(t *testing.T) {
	h := handler.NewWeatherHandler(
		&mockLocationFetcher{err: viacep.ErrZipNotFound},
		&mockWeatherFetcher{},
	)

	req := httptest.NewRequest(http.MethodGet, "/?zipcode=99999999", nil)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status code = %d, esperado 404", rec.Code)
	}
	if body := rec.Body.String(); body != "can not find zipcode" {
		t.Errorf("body = %q, esperado 'can not find zipcode'", body)
	}
}
