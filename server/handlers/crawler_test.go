package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/nettojulio/ufape-crawler-golang/crawler"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func TestCrawlerHandler(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintln(w, `<!DOCTYPE html>
<html>
<head>
    <title>Página de Teste</title>
</head>
<body>
    <h1>Olá, Mundo!</h1>
    <a href="/outra-pagina">Link Válido</a>
    <a href="http://link-externo-invalido">Link Inválido</a>
</body>
</html>`)
	}))
	defer ts.Close()

	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}

	t.Run("Cenário de Sucesso", func(t *testing.T) {
		payload := fmt.Sprintf(`{"url": "%s"}`, ts.URL)
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(payload))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := CrawlerHandler(c)

		if err != nil {
			t.Fatalf("CrawlerHandler retornou um erro inesperado: %v", err)
		}

		if rec.Code != http.StatusOK {
			t.Errorf("Código de status incorreto: got %v want %v", rec.Code, http.StatusOK)
		}

		var response crawler.ResponseCrawlDTO
		if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
			t.Fatalf("Falha ao decodificar o corpo da resposta JSON: %v", err)
		}

		if response.Title != "Página de Teste" {
			t.Errorf("Título incorreto: got %q want %q", response.Title, "Página de Teste")
		}

		if response.StatusCode != http.StatusOK {
			t.Errorf("StatusCode no corpo incorreto: got %d want %d", response.StatusCode, http.StatusOK)
		}
	})

	t.Run("Cenário de Falha - JSON Inválido", func(t *testing.T) {
		invalidPayload := `{"url": "http://example.com"`
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(invalidPayload))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		_ = CrawlerHandler(c)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("Código de status incorreto: got %v want %v", rec.Code, http.StatusBadRequest)
		}
	})

	t.Run("Cenário de Falha - URL Ausente", func(t *testing.T) {
		payloadWithoutURL := `{"timeout": 10}`
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(payloadWithoutURL))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		_ = CrawlerHandler(c)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("Código de status incorreto: got %v want %v", rec.Code, http.StatusBadRequest)
		}
	})
}
