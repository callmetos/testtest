package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"navmate-backend/config"
	"navmate-backend/internal/routes"

	"github.com/gin-gonic/gin"
)

func TestHealthSmoke(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// Create a mock config for testing
	cfg := &config.Config{}
	cfg.Google.GoogleMapsAPIKey = "test-api-key"

	routes.SetupRouter(r, nil, cfg)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}
