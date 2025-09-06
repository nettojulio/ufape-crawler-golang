package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/nettojulio/ufape-crawler-golang/crawler"
)

func TestHealthCheckHandler(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	testVersion := "v1.2.3-native"
	handler := HealthCheckHandler(testVersion)

	if err := handler(c); err != nil {
		t.Fatalf("HealthCheckHandler retornou um erro inesperado: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("Código de status incorreto: got %v want %v", rec.Code, http.StatusOK)
	}

	var response crawler.APIHealth
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("Falha ao decodificar o corpo da resposta JSON: %v", err)
	}

	expectedResponse := crawler.APIHealth{
		Status:  "OK",
		Version: testVersion,
	}

	if !reflect.DeepEqual(response, expectedResponse) {
		t.Errorf("Corpo da resposta incorreto: got %+v want %+v", response, expectedResponse)
	}
}
