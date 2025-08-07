package utils

import (
	"net/url"
	"strings"

	"github.com/nettojulio/ufape-crawler-golang/utils/configs"
)

func NormalizeURL(raw string) string {
	u, err := url.Parse(raw)
	if err != nil {
		return raw
	}
	if configs.RemoveFragment {
		u.Fragment = ""
	}
	pass, ok := u.User.Password()
	passEmpty := !ok || pass == ""
	if u.User != nil {
		if u.User.Username() == "" && passEmpty {
			u.User = nil
		}
	}

	if strings.Contains(u.Host, "www") {
		u.Host = strings.TrimPrefix(u.Host, "www.")
	}

	norm := u.String()
	if configs.LowerCaseURLs {
		norm = strings.ToLower(norm)
	}

	return strings.TrimRight(norm, "/")
}
