package main

import (
	"github.com/gofiber/fiber/v2"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestHealthCheckAuthorized(t *testing.T) {
	// Ustawiamy HEALTH_TOKEN
	os.Setenv("HEALTH_TOKEN", "test_secret")

	app := fiber.New()
	app.Get("/health", HealthCheck)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	req.Header.Set("X-Health-Token", "test_secret")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Błąd testu: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Oczekiwano status 200, otrzymano %d", resp.StatusCode)
	}
}

func TestHealthCheckUnauthorized(t *testing.T) {
	os.Setenv("HEALTH_TOKEN", "test_secret")

	app := fiber.New()
	app.Get("/health", HealthCheck)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	// Nie ustawiamy nagłówka X-Health-Token

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Błąd testu: %v", err)
	}

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Oczekiwano status 401, otrzymano %d", resp.StatusCode)
	}
}
