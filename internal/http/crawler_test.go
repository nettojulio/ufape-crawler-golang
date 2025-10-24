package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/nettojulio/ufape-crawler-golang/internal/crawler"
	"github.com/stretchr/testify/assert"
)

type mockCrawlerService struct {
	resultToReturn *crawler.CrawlResult
	errorToReturn  error
}

func (m *mockCrawlerService) Crawl(ctx context.Context, payload crawler.Payload, originalURL, modifiedURL *url.URL) (*crawler.CrawlResult, error) {
	return m.resultToReturn, m.errorToReturn
}

func TestCrawlerHandler_HandleCrawl(t *testing.T) {
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}

	t.Run("Cenário de Sucesso", func(t *testing.T) {
		reqBody := `{"url": "http://example.com"}`
		mockResult := &crawler.CrawlResult{
			StatusCode:  http.StatusOK,
			Title:       "Página de Teste",
			ElapsedTime: 123 * time.Millisecond,
			FinalURL:    func() *url.URL { u, _ := url.Parse("http://example.com"); return u }(),
		}
		mockService := &mockCrawlerService{resultToReturn: mockResult}
		handler := NewCrawlerHandler(mockService)
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		err := handler.HandleCrawl(c)
		assert.NoError(t, err, "HandleCrawl não deveria retornar erro")
		assert.Equal(t, http.StatusOK, rec.Code, "O código de status HTTP deveria ser 200")
		var responseDTO crawler.ResponseDTO
		err = json.Unmarshal(rec.Body.Bytes(), &responseDTO)
		assert.NoError(t, err, "O corpo da resposta deveria ser um JSON válido")
		assert.Equal(t, "Página de Teste", responseDTO.Title, "O título na resposta está incorreto")
	})

	t.Run("Cenário de Falha - URL Ausente na Requisição", func(t *testing.T) {
		reqBody := `{"timeout": 10}`
		mockService := &mockCrawlerService{}
		handler := NewCrawlerHandler(mockService)
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.HandleCrawl(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code, "O código de status HTTP deveria ser 400 para payload inválido")

		bodyString := rec.Body.String()
		assert.Contains(t, bodyString, "'Url'", "A resposta de erro deveria mencionar o campo 'Url'")
		assert.Contains(t, bodyString, "'required'", "A resposta de erro deveria mencionar a falha na regra 'required'")
	})

	t.Run("Cenário de Falha - Serviço Retorna Erro", func(t *testing.T) {
		reqBody := `{"url": "http://failing-site.com"}`
		expectedErr := errors.New("falha de conexão simulada")
		mockService := &mockCrawlerService{errorToReturn: expectedErr}
		handler := NewCrawlerHandler(mockService)
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		err := handler.HandleCrawl(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusServiceUnavailable, rec.Code, "O código de status HTTP deveria ser 503 quando o serviço falha")
		assert.Contains(t, rec.Body.String(), expectedErr.Error(), "A resposta de erro deveria conter a mensagem do serviço")
	})
}
