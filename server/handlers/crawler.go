package handlers

import (
	"net/http"
	"net/url"
	"strings"

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

	response, err := crawler.CrawlerService(payload, *originalUrlDetails, modifiedUrlDetails)
	if err != nil {
		/*
			Retorna 200 OK com o erro no titulo da resposta e status code 518.
			Isso é feito para manter a compatibilidade com o frontend que espera um status 200
			e trata o erro no corpo da resposta.
		*/
		return c.JSON(http.StatusOK, response)
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

	for i, domain := range *payload.AllowedDomains {
		(*payload.AllowedDomains)[i] = strings.TrimPrefix(domain, "www.")
	}

	if *payload.RemoveFragment {
		u.Fragment = ""
	}
}
