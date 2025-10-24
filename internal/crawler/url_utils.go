package crawler

import (
	"net/url"
	"strings"
)

// NormalizeURL limpa e padroniza uma URL com base nas opções fornecidas.
func NormalizeURL(raw string, removeFragment, lowerCaseURLs bool) string {
	u, err := url.Parse(raw)
	if err != nil {
		return raw
	}

	if removeFragment {
		u.Fragment = ""
	}

	if strings.HasPrefix(strings.ToLower(u.Host), "www.") {
		u.Host = u.Host[4:]
	}

	norm := u.String()
	if lowerCaseURLs {
		norm = strings.ToLower(norm)
	}

	if u.Path != "/" && strings.HasSuffix(norm, "/") {
		norm = norm[:len(norm)-1]
	}

	return norm
}

// IsSubdomainHost verifica se um host é um subdomínio de uma lista de domínios permitidos.
func IsSubdomainHost(host string, domains []string) bool {
	normalizedHost := host
	if strings.HasPrefix(strings.ToLower(host), "www.") {
		normalizedHost = host[4:]
	}

	for _, d := range domains {
		if normalizedHost == d || strings.HasSuffix(normalizedHost, "."+d) {
			return true
		}
	}
	return false
}
