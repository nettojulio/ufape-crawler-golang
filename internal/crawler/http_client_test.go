package crawler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestNewHTTPClient(t *testing.T) {
	testTimeout := 15 * time.Second
	client := NewHTTPClient(testTimeout)

	if client == nil {
		t.Fatal("NewHTTPClient() returned nil")
	}

	if client.client == nil {
		t.Fatal("internal http.Client is nil")
	}

	if client.client.Timeout != testTimeout {
		t.Errorf("expected timeout %v, got %v", testTimeout, client.client.Timeout)
	}

	if client.userAgent == "" {
		t.Error("expected userAgent to be set, but it was empty")
	}
}

func TestHTTPClient_Get(t *testing.T) {
	expectedUserAgent := "Mozilla/5.0 (X11; Linux x86_64; rv:145.0) Gecko/20100101 Firefox/145.0"

	t.Run("successful get request", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if ua := r.Header.Get("User-Agent"); ua != expectedUserAgent {
				t.Errorf("handler expected User-Agent %q, got %q", expectedUserAgent, ua)
			}
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		client := NewHTTPClient(5 * time.Second)
		resp, err := client.Get(context.Background(), server.URL)

		if err != nil {
			t.Fatalf("Get() returned an unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, resp.StatusCode)
		}
	})

	t.Run("server error response", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		client := NewHTTPClient(5 * time.Second)
		resp, err := client.Get(context.Background(), server.URL)

		if err != nil {
			t.Fatalf("Get() returned an unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusInternalServerError {
			t.Errorf("expected status code %d, got %d", http.StatusInternalServerError, resp.StatusCode)
		}
	})

	t.Run("client timeout", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(100 * time.Millisecond)
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		client := NewHTTPClient(50 * time.Millisecond)
		_, err := client.Get(context.Background(), server.URL)

		if err == nil {
			t.Fatal("expected a timeout error, but got nil")
		}
		if !strings.Contains(err.Error(), "context deadline exceeded") {
			t.Errorf("expected error to be a timeout error, but got: %v", err)
		}
	})

	t.Run("context canceled", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		}))
		defer server.Close()

		client := NewHTTPClient(5 * time.Second)
		ctx, cancel := context.WithCancel(context.Background())

		cancel()

		_, err := client.Get(ctx, server.URL)

		if err == nil {
			t.Fatal("expected a context canceled error, but got nil")
		}

		if !errors.Is(err, context.Canceled) {
			t.Errorf("expected error to wrap context.Canceled, but got: %v", err)
		}
	})
}
