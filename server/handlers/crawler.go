package handlers

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/labstack/echo/v4"
	crawler2 "github.com/nettojulio/ufape-crawler-golang/crawler"
)

func CrawlerHandler(c echo.Context) error {
	var payload crawler2.CorrectPayload
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request body"})
	}
	if err := c.Validate(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	originalUrlDetails, err := url.Parse(payload.Url)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid url"})
	}

	modifiedUrlDetails := *originalUrlDetails
	normalizeURL(&modifiedUrlDetails)

	applyDefaults(&payload, originalUrlDetails.Host, &modifiedUrlDetails)

	response, err := crawler2.CrawlerService(payload, *originalUrlDetails, modifiedUrlDetails)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error":    "internal server error",
			"details":  err.Error(),
			"response": response,
		})
	}

	return c.JSON(http.StatusOK, response)
}

func normalizeURL(u *url.URL) {
	if u.User != nil {
		pass, ok := u.User.Password()
		if u.User.Username() == "" && (!ok || pass == "") {
			u.User = nil
		}
	}
}

func applyDefaults(payload *crawler2.CorrectPayload, host string, u *url.URL) {
	if payload.Timeout == nil {
		def := 1
		payload.Timeout = &def
	}
	if payload.RemoveFragment == nil {
		def := true
		payload.RemoveFragment = &def
	}
	if payload.AllowedDomains == nil {
		def := []string{host}
		payload.AllowedDomains = &def
	}
	if payload.CollectSubdomains == nil {
		def := true
		payload.CollectSubdomains = &def
	}
	if payload.LowerCaseURLs == nil {
		def := false
		payload.LowerCaseURLs = &def
	}

	for i, domain := range *payload.AllowedDomains {
		(*payload.AllowedDomains)[i] = strings.TrimPrefix(domain, "www.")
	}

	if *payload.RemoveFragment {
		u.Fragment = ""
	}
}
