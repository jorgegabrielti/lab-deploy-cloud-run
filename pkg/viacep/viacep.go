package viacep

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

var (
	ErrInvalidZipcode = errors.New("invalid zipcode")
	ErrZipNotFound    = errors.New("can not find zipcode")
)

type ViaCEPResponse struct {
	Localidade string `json:"localidade"`
	Erro       any    `json:"erro"`
}

type Client struct {
	HTTPClient *http.Client
	BaseURL    string
}

func NewClient() *Client {
	return &Client{
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
		BaseURL:    "https://viacep.com.br/ws",
	}
}

func ValidateZipcode(zip string) (string, error) {
	clean := strings.ReplaceAll(zip, "-", "")
	clean = strings.TrimSpace(clean)

	matched, _ := regexp.MatchString(`^\d{8}$`, clean)
	if !matched {
		return "", ErrInvalidZipcode
	}
	return clean, nil
}

func (c *Client) FetchLocation(ctx context.Context, zip string) (string, error) {
	cleanZip, err := ValidateZipcode(zip)
	if err != nil {
		return "", err
	}

	reqURL := fmt.Sprintf("%s/%s/json/", c.BaseURL, cleanZip)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return "", fmt.Errorf("erro ao criar requisicao: %w", err)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("erro na chamada ViaCEP: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", ErrZipNotFound
	}

	var res ViaCEPResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", fmt.Errorf("erro ao decodificar resposta: %w", err)
	}

	if isErroTrue(res.Erro) || strings.TrimSpace(res.Localidade) == "" {
		return "", ErrZipNotFound
	}

	return res.Localidade, nil
}

func isErroTrue(val any) bool {
	if val == nil {
		return false
	}
	switch v := val.(type) {
	case bool:
		return v
	case string:
		return v == "true" || v == "1"
	}
	return false
}

func SanitizeCityName(city string) string {
	return url.QueryEscape(city)
}
