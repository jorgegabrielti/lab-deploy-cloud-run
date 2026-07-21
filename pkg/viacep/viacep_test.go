package viacep_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jorgegabrielti/desafio-clima-cep/pkg/viacep"
)

func TestValidateZipcode_Valid(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{"01001000", "01001000"},
		{"01001-000", "01001000"},
		{" 70000000 ", "70000000"},
	}

	for _, tc := range cases {
		got, err := viacep.ValidateZipcode(tc.input)
		if err != nil {
			t.Errorf("ValidateZipcode(%q) retornou erro inesperado: %v", tc.input, err)
		}
		if got != tc.want {
			t.Errorf("ValidateZipcode(%q) = %q, esperado %q", tc.input, got, tc.want)
		}
	}
}

func TestValidateZipcode_Invalid(t *testing.T) {
	invalidCases := []string{
		"123",
		"1234567",
		"123456789",
		"0100100a",
		"abcde-fgh",
		"",
	}

	for _, input := range invalidCases {
		_, err := viacep.ValidateZipcode(input)
		if err != viacep.ErrInvalidZipcode {
			t.Errorf("ValidateZipcode(%q) esperava ErrInvalidZipcode, obteve %v", input, err)
		}
	}
}

func TestFetchLocation_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"localidade": "São Paulo"}`))
	}))
	defer srv.Close()

	client := &viacep.Client{
		HTTPClient: srv.Client(),
		BaseURL:    srv.URL,
	}

	city, err := client.FetchLocation(context.Background(), "01001000")
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if city != "São Paulo" {
		t.Errorf("city = %q, esperado %q", city, "São Paulo")
	}
}

func TestFetchLocation_NotFound(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"erro": true}`))
	}))
	defer srv.Close()

	client := &viacep.Client{
		HTTPClient: srv.Client(),
		BaseURL:    srv.URL,
	}

	_, err := client.FetchLocation(context.Background(), "99999999")
	if err != viacep.ErrZipNotFound {
		t.Errorf("esperado ErrZipNotFound, obteve %v", err)
	}
}
