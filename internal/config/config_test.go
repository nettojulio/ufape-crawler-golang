package config

import (
	"testing"
)

func TestLoad(t *testing.T) {
	const testVersion = "v1.0.0-test"

	t.Run("should load default values when no env vars are set", func(t *testing.T) {

		cfg, err := Load(testVersion)
		if err != nil {
			t.Fatalf("Load() returned an unexpected error: %v", err)
		}

		if cfg.Version != testVersion {
			t.Errorf("expected Version to be %q, got %q", testVersion, cfg.Version)
		}

		expectedPort := 8080
		if cfg.Port != expectedPort {
			t.Errorf("expected Port to be %d, got %d", expectedPort, cfg.Port)
		}

		expectedHost := "localhost:8080"
		if cfg.Host != expectedHost {
			t.Errorf("expected Host to be %q, got %q", expectedHost, cfg.Host)
		}
	})

	t.Run("should override defaults with environment variables", func(t *testing.T) {
		t.Setenv("APP_PORT", "9999")
		t.Setenv("APP_HOST", "test.example.com")
		t.Setenv("APP_VERSION", "v2.0.0-env")

		cfg, err := Load(testVersion)
		if err != nil {
			t.Fatalf("Load() returned an unexpected error: %v", err)
		}

		expectedVersion := "v2.0.0-env"
		if cfg.Version != expectedVersion {
			t.Errorf("expected Version from env to be %q, got %q", expectedVersion, cfg.Version)
		}

		expectedPort := 9999
		if cfg.Port != expectedPort {
			t.Errorf("expected Port from env to be %d, got %d", expectedPort, cfg.Port)
		}

		expectedHost := "test.example.com"
		if cfg.Host != expectedHost {
			t.Errorf("expected Host from env to be %q, got %q", expectedHost, cfg.Host)
		}
	})

	t.Run("should return an error for invalid port value", func(t *testing.T) {
		t.Setenv("APP_PORT", "not-a-number")

		_, err := Load(testVersion)

		if err == nil {
			t.Fatal("expected an error for invalid port, but got nil")
		}
	})
}
