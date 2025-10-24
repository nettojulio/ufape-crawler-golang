package crawler

import "testing"

func TestNormalizeURL(t *testing.T) {
	testCases := []struct {
		name           string
		rawURL         string
		removeFragment bool
		lowerCaseURLs  bool
		expectedURL    string
	}{
		{
			name:           "Remove fragment",
			rawURL:         "http://example.com/page#section1",
			removeFragment: true,
			lowerCaseURLs:  false,
			expectedURL:    "http://example.com/page",
		},
		{
			name:           "Keep fragment",
			rawURL:         "http://example.com/page#section1",
			removeFragment: false,
			lowerCaseURLs:  false,
			expectedURL:    "http://example.com/page#section1",
		},
		{
			name:           "Remove www prefix (lowercase)",
			rawURL:         "https://www.example.com/path",
			removeFragment: true,
			lowerCaseURLs:  false,
			expectedURL:    "https://example.com/path",
		},
		{
			name:           "Remove WWW prefix (uppercase)",
			rawURL:         "https://WWW.example.com/path",
			removeFragment: true,
			lowerCaseURLs:  false,
			expectedURL:    "https://example.com/path",
		},
		{
			name:           "Remove trailing slash",
			rawURL:         "https://example.com/path/",
			removeFragment: true,
			lowerCaseURLs:  false,
			expectedURL:    "https://example.com/path",
		},
		{
			name:           "Keep root trailing slash",
			rawURL:         "https://example.com/",
			removeFragment: true,
			lowerCaseURLs:  false,
			expectedURL:    "https://example.com/",
		},
		{
			name:           "Convert to lowercase",
			rawURL:         "HTTP://EXAMPLE.COM/Path",
			removeFragment: true,
			lowerCaseURLs:  true,
			expectedURL:    "http://example.com/path",
		},
		{
			name:           "Combination of rules",
			rawURL:         "HTTPS://WWW.Example.com/Path/?query=1#Fragment/",
			removeFragment: true,
			lowerCaseURLs:  true,
			expectedURL:    "https://example.com/path/?query=1",
		},
		{
			name:           "Invalid URL returns raw",
			rawURL:         "://invalid",
			removeFragment: true,
			lowerCaseURLs:  true,
			expectedURL:    "://invalid",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			normalizedURL := NormalizeURL(tc.rawURL, tc.removeFragment, tc.lowerCaseURLs)
			if normalizedURL != tc.expectedURL {
				t.Errorf("expected URL %q, got %q", tc.expectedURL, normalizedURL)
			}
		})
	}
}

func TestIsSubdomainHost(t *testing.T) {
	domains := []string{"example.com", "another.org"}

	testCases := []struct {
		name     string
		host     string
		expected bool
	}{
		{name: "Exact match", host: "example.com", expected: true},
		{name: "Direct subdomain", host: "blog.example.com", expected: true},
		{name: "Deep subdomain", host: "api.staging.example.com", expected: true},
		{name: "Exact match on other domain", host: "another.org", expected: true},
		{name: "Handles www prefix correctly", host: "www.example.com", expected: true},
		{name: "Handles WWW prefix correctly (uppercase)", host: "WWW.blog.example.com", expected: true},
		{name: "Not a subdomain (unrelated)", host: "notrelated.net", expected: false},
		{name: "Similar but incorrect (no dot)", host: "badexample.com", expected: false},
		{name: "Empty host", host: "", expected: false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			isSubdomain := IsSubdomainHost(tc.host, domains)
			if isSubdomain != tc.expected {
				t.Errorf("for host %q, expected %v, got %v", tc.host, tc.expected, isSubdomain)
			}
		})
	}

	t.Run("Empty domains list", func(t *testing.T) {
		if IsSubdomainHost("example.com", []string{}) {
			t.Error("expected false when domains list is empty, got true")
		}
	})
}
