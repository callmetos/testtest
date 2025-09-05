package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"navmate-backend/internal/routes"

	"github.com/gin-gonic/gin"
)

func TestHealthSmoke(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	routes.SetupRouter(r, nil, nil)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}
