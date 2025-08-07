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
	norm := u.String()
	if configs.LowerCaseURLs {
		norm = strings.ToLower(norm)
	}
	return strings.TrimRight(norm, "/")
}
