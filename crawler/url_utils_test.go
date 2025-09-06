package crawler

import (
	"testing"
)

func TestNormalizeURL(t *testing.T) {
	payloadRemoveFragment := &CorrectPayload{RemoveFragment: new(bool), LowerCaseURLs: new(bool)}
	*payloadRemoveFragment.RemoveFragment = true

	payloadKeepFragment := &CorrectPayload{RemoveFragment: new(bool), LowerCaseURLs: new(bool)}
	*payloadKeepFragment.RemoveFragment = false

	payloadLowerCase := &CorrectPayload{RemoveFragment: new(bool), LowerCaseURLs: new(bool)}
	*payloadLowerCase.RemoveFragment = true
	*payloadLowerCase.LowerCaseURLs = true

	testCases := []struct {
		name        string
		rawURL      string
		payload     *CorrectPayload
		expectedURL string
	}{
		{
			name:        "Remove Fragmento",
			rawURL:      "http://example.com/page#section1",
			payload:     payloadRemoveFragment,
			expectedURL: "http://example.com/page",
		},
		{
			name:        "Mantém Fragmento",
			rawURL:      "http://example.com/page#section1",
			payload:     payloadKeepFragment,
			expectedURL: "http://example.com/page#section1",
		},
		{
			name:        "Remove www.",
			rawURL:      "https://www.example.com/path",
			payload:     payloadRemoveFragment,
			expectedURL: "https://example.com/path",
		},
		{
			name:        "Remove Barra Final",
			rawURL:      "https://example.com/path/",
			payload:     payloadRemoveFragment,
			expectedURL: "https://example.com/path",
		},
		{
			name:        "Converte para LowerCase",
			rawURL:      "HTTP://EXAMPLE.COM/Path",
			payload:     payloadLowerCase,
			expectedURL: "http://example.com/path",
		},
		{
			name:        "Remove Informação de Usuário Vazia",
			rawURL:      "https://@example.com",
			payload:     payloadRemoveFragment,
			expectedURL: "https://example.com",
		},
		{
			name:        "Mantém Informação de Usuário Válida",
			rawURL:      "https://user:pass@example.com",
			payload:     payloadRemoveFragment,
			expectedURL: "https://user:pass@example.com",
		},
		{
			name:        "Combinação de Regras",
			rawURL:      "HTTPS://WWW.Example.com/Path/?query=1#Fragment/",
			payload:     payloadLowerCase,
			expectedURL: "https://example.com/path/?query=1",
		},
		{
			name:        "URL Inválida",
			rawURL:      "://inválido",
			payload:     payloadRemoveFragment,
			expectedURL: "://inválido",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			normalizedURL := NormalizeURL(tc.rawURL, tc.payload)

			if normalizedURL != tc.expectedURL {
				t.Errorf("URL normalizada incorreta: got %q, want %q", normalizedURL, tc.expectedURL)
			}
		})
	}
}

func TestIsSubdomainHost(t *testing.T) {
	domains := []string{"example.com", "another.org"}

	testCases := []struct {
		name     string
		host     string
		domains  []string
		expected bool
	}{
		{
			name:     "Host Exato",
			host:     "example.com",
			domains:  domains,
			expected: true,
		},
		{
			name:     "Subdomínio Direto",
			host:     "blog.example.com",
			domains:  domains,
			expected: true,
		},
		{
			name:     "Subdomínio Profundo",
			host:     "api.staging.example.com",
			domains:  domains,
			expected: true,
		},
		{
			name:     "Outro Domínio Válido",
			host:     "another.org",
			domains:  domains,
			expected: true,
		},
		{
			name:     "Host não é subdomínio",
			host:     "notrelated.net",
			domains:  domains,
			expected: false,
		},
		{
			name:     "Similar mas incorreto (sem o ponto)",
			host:     "badexample.com",
			domains:  domains,
			expected: false,
		},
		{
			name:     "Host Vazio",
			host:     "",
			domains:  domains,
			expected: false,
		},
		{
			name:     "Lista de Domínios Vazia",
			host:     "example.com",
			domains:  []string{},
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			isSubdomain := IsSubdomainHost(tc.host, tc.domains)

			if isSubdomain != tc.expected {
				t.Errorf("Resultado incorreto para host %q: got %v, want %v", tc.host, isSubdomain, tc.expected)
			}
		})
	}
}
