package crawler

import (
	"net/url"
	"strings"
)

func NormalizeURL(raw string, payload *CorrectPayload) string {
	u, err := url.Parse(raw)
	if err != nil {
		return raw
	}
	if *payload.RemoveFragment {
		u.Fragment = ""
	}
	if u.User != nil {
		if user := u.User.Username(); user == "" {
			if pass, ok := u.User.Password(); !ok || pass == "" {
				u.User = nil
			}
		}
	}
	if strings.HasPrefix(u.Host, "www.") {
		u.Host = strings.TrimPrefix(u.Host, "www.")
	}

	if strings.HasPrefix(u.Host, "WWW.") {
		u.Host = strings.TrimPrefix(u.Host, "WWW.")
	}

	norm := u.String()
	if *payload.LowerCaseURLs {
		norm = strings.ToLower(norm)
	}
	return strings.TrimRight(norm, "/")
}

func IsSubdomainHost(host string, domains []string) bool {
	for _, d := range domains {
		if host == d || strings.HasSuffix(host, "."+d) {
			return true
		}
	}
	return false
}
