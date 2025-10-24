package crawler

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

type errorReader struct{}

func (er *errorReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("simulated read error")
}

type mockHTTPClient struct {
	Response *http.Response
	Err      error
}

func (m *mockHTTPClient) Get(ctx context.Context, url string) (*http.Response, error) {
	if m.Response != nil && m.Response.Request == nil {
		req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
		m.Response.Request = req
	}
	return m.Response, m.Err
}

func TestService_Crawl(t *testing.T) {
	ctx := context.Background()
	originalURL, _ := url.Parse("http://example.com")
	modifiedURL, _ := url.Parse("http://example.com")

	defaultPayload := Payload{
		AllowedDomains:    &[]string{"example.com"},
		CollectSubdomains: boolPtr(false),
		RemoveFragment:    boolPtr(true),
		LowerCaseURLs:     boolPtr(false),
	}

	t.Run("successful crawl with valid html", func(t *testing.T) {
		htmlBody := `<html><head><title>Test Title</title></head><body><a href="/page1">Link</a></body></html>`
		mockResp := &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"text/html"}},
			Body:       io.NopCloser(strings.NewReader(htmlBody)),
		}

		mockClient := &mockHTTPClient{Response: mockResp}
		service := NewService(mockClient)

		result, err := service.Crawl(ctx, defaultPayload, originalURL, modifiedURL)

		if err != nil {
			t.Fatalf("Crawl() returned an unexpected error: %v", err)
		}
		if result.StatusCode != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, result.StatusCode)
		}
		if result.Title != "Test Title" {
			t.Errorf("expected title %q, got %q", "Test Title", result.Title)
		}
		if len(result.Links.Available) != 1 || result.Links.Available[0] != "http://example.com/page1" {
			t.Errorf("unexpected available links: got %v", result.Links.Available)
		}
	})

	t.Run("server response is not 200 OK", func(t *testing.T) {
		mockResp := &http.Response{
			StatusCode: http.StatusNotFound,
			Body:       io.NopCloser(bytes.NewReader(nil)),
		}

		mockClient := &mockHTTPClient{Response: mockResp}
		service := NewService(mockClient)

		result, err := service.Crawl(ctx, defaultPayload, originalURL, modifiedURL)

		if err != nil {
			t.Fatalf("Crawl() returned an unexpected error for a valid HTTP response: %v", err)
		}
		if result.StatusCode != http.StatusNotFound {
			t.Errorf("expected status code %d, got %d", http.StatusNotFound, result.StatusCode)
		}
		if result.Title != "" {
			t.Errorf("expected empty title for a non-200 response, got %q", result.Title)
		}
	})

	t.Run("http client returns a network error", func(t *testing.T) {
		expectedErr := errors.New("connection failed")
		mockClient := &mockHTTPClient{Err: expectedErr}
		service := NewService(mockClient)

		result, _ := service.Crawl(ctx, defaultPayload, originalURL, modifiedURL)

		if result.StatusCode != http.StatusServiceUnavailable {
			t.Errorf("expected status code %d for network error, got %d", http.StatusServiceUnavailable, result.StatusCode)
		}
		if !strings.Contains(result.Title, expectedErr.Error()) {
			t.Errorf("expected result title to contain error message %q, but got %q", expectedErr.Error(), result.Title)
		}
	})
}

func boolPtr(b bool) *bool {
	return &b
}
