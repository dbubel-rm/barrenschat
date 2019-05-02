package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dbubel/bchat/internal/platform/web"
	"github.com/stretchr/testify/assert"
)

func TestHealthEndpoint(t *testing.T) {
	resetTestDB()
	a := API(l, d).(*web.App)
	r := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	a.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code, "Response code should be 200")

}
