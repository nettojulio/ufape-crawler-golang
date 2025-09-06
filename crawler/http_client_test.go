package crawler

import (
	"net/http"
	"testing"
	"time"
)

func TestNewHTTPClient(t *testing.T) {
	testTimeout := 15

	client := NewHTTPClient(testTimeout)

	expectedTimeout := time.Duration(testTimeout) * time.Second
	if client.Timeout != expectedTimeout {
		t.Errorf("Timeout incorreto: got %v want %v", client.Timeout, expectedTimeout)
	}

	if client.CheckRedirect == nil {
		t.Error("CheckRedirect não deveria ser nulo, mas é.")
	}
}

func TestNewRequest(t *testing.T) {
	t.Run("Cenário de Sucesso com URL Válida", func(t *testing.T) {
		testURL := "http://example.com"

		req, err := NewRequest(testURL)

		if err != nil {
			t.Fatalf("NewRequest retornou um erro inesperado: %v", err)
		}

		if req == nil {
			t.Fatal("NewRequest retornou uma requisição nula, mas não deveria.")
		}

		if req.Method != http.MethodGet {
			t.Errorf("Método HTTP incorreto: got %q want %q", req.Method, http.MethodGet)
		}
		if req.URL.String() != testURL {
			t.Errorf("URL incorreta: got %q want %q", req.URL.String(), testURL)
		}

		expectedHeaders := map[string]string{
			"Accept":     "*/*",
			"Connection": "keep-alive",
			"User-Agent": "Mozilla/5.0 (X11; Linux x86_64; rv:143.0) Gecko/20100101 Firefox/143.0",
		}

		for key, expectedValue := range expectedHeaders {
			actualValue := req.Header.Get(key)
			if actualValue != expectedValue {
				t.Errorf("Header %q incorreto: got %q want %q", key, actualValue, expectedValue)
			}
		}
	})

	t.Run("Cenário de Falha com URL Inválida", func(t *testing.T) {
		invalidURL := "http://example.com\x7f"

		req, err := NewRequest(invalidURL)

		if err == nil {
			t.Fatal("NewRequest deveria retornar um erro para URL inválida, mas não retornou.")
		}

		if req != nil {
			t.Error("NewRequest deveria retornar uma requisição nula em caso de erro, mas não retornou.")
		}
	})
}
