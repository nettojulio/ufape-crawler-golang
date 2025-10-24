package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/nettojulio/ufape-crawler-golang/internal/crawler"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheckHandler(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	testVersion := "v1.2.3-test"
	handler := HealthCheckHandler(testVersion)

	expectedResponse := crawler.APIHealth{
		Status:  "OK",
		Version: testVersion,
	}

	err := handler(c)

	assert.NoError(t, err, "HealthCheckHandler não deveria retornar um erro")
	assert.Equal(t, http.StatusOK, rec.Code, "O código de status HTTP deveria ser 200 OK")

	var actualResponse crawler.APIHealth
	err = json.Unmarshal(rec.Body.Bytes(), &actualResponse)

	assert.NoError(t, err, "O corpo da resposta deveria ser um JSON decodificável")
	assert.Equal(t, expectedResponse, actualResponse, "O corpo da resposta não corresponde ao esperado")
}
