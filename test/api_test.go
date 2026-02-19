package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"cadastral-service/internal/api"
	"cadastral-service/internal/config"
)

func TestPingEndpoint(t *testing.T) {
	router := gin.New()
	cfg := &config.Config{
		Environment: "test",
		Port:        "8080",
	}

	handler := &api.Handler{
		Config: cfg,
	}

	router.GET("/api/v1/ping", handler.Ping)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, "pong", response["message"])
	assert.NotNil(t, response["time"])
}
