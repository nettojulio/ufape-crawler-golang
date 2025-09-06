package crawler

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"
)

type mockRoundTripper struct {
	Response *http.Response
	Err      error
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.Response, m.Err
}

func TestCrawlerService(t *testing.T) {
	originalURL, _ := url.Parse("http://example.com")
	modifiedURL, _ := url.Parse("http://example.com")
	timeout := 10
	payload := CorrectPayload{Timeout: &timeout, LowerCaseURLs: new(bool), RemoveFragment: new(bool)}

	*payload.LowerCaseURLs = false
	*payload.RemoveFragment = false

	t.Run("Cenário de Sucesso 200 OK", func(t *testing.T) {
		htmlBody := `<html><head><title>Título de Teste</title></head><body></body></html>`
		mockResponse := &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"text/html"}},
			Body:       io.NopCloser(bytes.NewReader([]byte(htmlBody))),
		}

		mockClient := &http.Client{
			Transport: &mockRoundTripper{Response: mockResponse},
		}

		response, err := CrawlerService(mockClient, payload, *originalURL, *modifiedURL)

		if err != nil {
			t.Fatalf("CrawlerService retornou um erro inesperado: %v", err)
		}
		if response.StatusCode != http.StatusOK {
			t.Errorf("StatusCode incorreto: got %d, want %d", response.StatusCode, http.StatusOK)
		}
		if response.Title != "Título de Teste" {
			t.Errorf("Título incorreto: got %q, want %q", response.Title, "Título de Teste")
		}
	})

	t.Run("Cenário de Falha com Status 404", func(t *testing.T) {
		mockResponse := &http.Response{
			StatusCode: http.StatusNotFound,
			Header:     http.Header{"Content-Type": []string{"text/html"}},
			Body:       io.NopCloser(bytes.NewReader([]byte(""))),
		}

		mockClient := &http.Client{
			Transport: &mockRoundTripper{Response: mockResponse},
		}

		response, err := CrawlerService(mockClient, payload, *originalURL, *modifiedURL)

		if err != nil {
			t.Fatalf("CrawlerService retornou um erro inesperado: %v", err)
		}
		if response.StatusCode != http.StatusNotFound {
			t.Errorf("StatusCode incorreto: got %d, want %d", response.StatusCode, http.StatusNotFound)
		}
		if response.Title != "" {
			t.Errorf("Título deveria ser vazio para respostas de erro, mas foi %q", response.Title)
		}
	})

	t.Run("Cenário de Falha de Rede", func(t *testing.T) {
		expectedErr := errors.New("falha de conexão")
		mockClient := &http.Client{
			Transport: &mockRoundTripper{Err: expectedErr},
		}

		response, err := CrawlerService(mockClient, payload, *originalURL, *modifiedURL)

		if err == nil {
			t.Fatal("CrawlerService deveria retornar um erro, mas não retornou.")
		}
		if !errors.Is(err, expectedErr) {
			t.Errorf("Erro incorreto: got %v, want %v", err, expectedErr)
		}
		if response.StatusCode != 518 {
			t.Errorf("StatusCode incorreto para erro de rede: got %d, want %d", response.StatusCode, 518)
		}
		if !strings.Contains(response.Title, expectedErr.Error()) {
			t.Errorf("Título deveria conter a mensagem de erro: got %q, want %q", response.Title, expectedErr.Error())
		}
	})
}

func TestBuildErrorResponse(t *testing.T) {
	originalURL, _ := url.Parse("http://original.com")
	modifiedURL, _ := url.Parse("http://modified.com")
	elapsed := time.Second

	response := buildErrorResponse(404, "text/plain", elapsed, *originalURL, *modifiedURL)

	if response.StatusCode != 404 {
		t.Errorf("StatusCode incorreto: got %d, want %d", response.StatusCode, 404)
	}
	if response.ElapsedTime != elapsed.Nanoseconds() {
		t.Errorf("ElapsedTime incorreto: got %d, want %d", response.ElapsedTime, elapsed.Nanoseconds())
	}
	if len(response.Links.Available) != 0 || len(response.Links.Unavailable) != 0 {
		t.Error("Links não deveriam estar vazios em uma resposta de erro")
	}
}
