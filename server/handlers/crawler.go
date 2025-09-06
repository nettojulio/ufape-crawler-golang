package handlers

import (
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/nettojulio/ufape-crawler-golang/crawler"
)

// CrawlerHandler godoc
// @Summary      Inicia o processo de crawling
// @Description  Recebe uma URL e configurações opcionais para iniciar o crawling de uma página web.
// @Tags         Crawler
// @Accept       json
// @Produce      json
// @Param        payload body crawler.CorrectPayload true "Configurações do Crawler - URL é obrigatório"
// @Success      200  {object}  crawler.ResponseCrawlDTO
// @Router       / [post]
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

	client := crawler.NewHTTPClient(*payload.Timeout)

	attempt := 1
	var response crawler.ResponseCrawl
	for attempt <= *payload.MaxAttempts {
		response, err = crawler.CrawlerService(&client, payload, *originalUrlDetails, modifiedUrlDetails)
		if response.StatusCode != 404 && response.StatusCode != 200 {
			attempt++
			time.Sleep(1 * time.Second)
			continue
		}
		break
	}

	responseDTO := crawler.ResponseCrawlDTO{
		StatusCode:  response.StatusCode,
		ContentType: response.ContentType,
		ElapsedTime: response.ElapsedTime,
		Links:       response.Links,
		Title:       response.Title,
		Details: crawler.DetailsResponseDTO{
			CorrectURL: response.Details.CorrectURL,
			Original:   mapURLToDetails(&response.Details.Original),
			Modified:   mapURLToDetails(&response.Details.Modified),
		},
	}

	return c.JSON(http.StatusOK, responseDTO)
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

func mapURLToDetails(u *url.URL) crawler.URLDetails {
	if u == nil {
		return crawler.URLDetails{}
	}

	details := crawler.URLDetails{
		Scheme:      u.Scheme,
		Opaque:      u.Opaque,
		Host:        u.Host,
		Path:        u.Path,
		RawPath:     u.RawPath,
		OmitHost:    u.OmitHost,
		ForceQuery:  u.ForceQuery,
		RawQuery:    u.RawQuery,
		Fragment:    u.Fragment,
		RawFragment: u.RawFragment,
	}

	if u.User != nil {
		password, isSet := u.User.Password()
		details.User = &crawler.UserURLDetails{
			Username:    u.User.Username(),
			Password:    password,
			PasswordSet: isSet,
		}
	}

	return details
}
