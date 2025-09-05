package handlers

import (
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/nettojulio/ufape-crawler-golang/crawler"
)

func CrawlerHandler(c echo.Context) error {
	var payload crawler.CorrectPayload
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

	attempt := 1
	var response crawler.ResponseCrawl
	for attempt <= *payload.MaxAttempts {
		response, err = crawler.CrawlerService(payload, *originalUrlDetails, modifiedUrlDetails)
		if response.StatusCode != 404 && response.StatusCode != 200 {
			attempt++
			time.Sleep(1 * time.Second)
			continue
		}
		break
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

func applyDefaults(payload *crawler.CorrectPayload, host string, u *url.URL) {
	if payload.Timeout == nil {
		def := 60
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
	if payload.CanRetry == nil {
		def := false
		payload.CanRetry = &def
	}
	if payload.MaxAttempts == nil {
		def := func() int {
			if *payload.CanRetry {
				return 10
			}
			return 1
		}()
		payload.MaxAttempts = &def
	}

	for i, domain := range *payload.AllowedDomains {
		(*payload.AllowedDomains)[i] = strings.TrimPrefix(domain, "www.")
	}

	if *payload.RemoveFragment {
		u.Fragment = ""
	}
}
