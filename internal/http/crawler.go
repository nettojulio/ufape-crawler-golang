package http

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/nettojulio/ufape-crawler-golang/internal/crawler"
)

type CrawlerServicer interface {
	Crawl(ctx context.Context, payload crawler.Payload, originalURL, modifiedURL *url.URL) (*crawler.CrawlResult, error)
}

type CrawlerHandler struct {
	crawlerService CrawlerServicer
}

func NewCrawlerHandler(cs CrawlerServicer) *CrawlerHandler {
	return &CrawlerHandler{
		crawlerService: cs,
	}
}

// HandleCrawl godoc
// @Summary      Inicia o processo de crawling
// @Description  Recebe uma URL e configurações opcionais para iniciar o crawling de uma página web.
// @Tags         Crawler
// @Accept       json
// @Produce      json
// @Param        payload body crawler.Payload true "Configurações do Crawler - URL é obrigatório"
// @Success      200  {object}  crawler.ResponseDTO
// @Router       / [post]
func (h *CrawlerHandler) HandleCrawl(c echo.Context) error {
	var payload crawler.Payload
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request body"})
	}
	if err := c.Validate(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	originalURL, err := url.Parse(payload.Url)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid url"})
	}

	modifiedURL := prepareURLAndDefaults(&payload, originalURL)

	var result *crawler.CrawlResult
	for attempt := 1; attempt <= *payload.MaxAttempts; attempt++ {
		result, err = h.crawlerService.Crawl(c.Request().Context(), payload, originalURL, modifiedURL)

		if err != nil || (result.StatusCode != http.StatusOK && result.StatusCode != http.StatusNotFound) {
			if attempt < *payload.MaxAttempts {
				time.Sleep(1 * time.Second)
				continue
			}
		}
		break
	}

	if err != nil {
		return c.JSON(http.StatusServiceUnavailable, echo.Map{"error": err.Error()})
	}

	responseDTO := mapResultToDTO(result, originalURL)
	return c.JSON(http.StatusOK, responseDTO)
}

func prepareURLAndDefaults(payload *crawler.Payload, originalURL *url.URL) *url.URL {
	if payload.Timeout == nil {
		def := 60
		payload.Timeout = &def
	}
	if payload.RemoveFragment == nil {
		def := true
		payload.RemoveFragment = &def
	}
	if payload.AllowedDomains == nil {
		def := []string{originalURL.Host}
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
		def := 1
		if *payload.CanRetry {
			def = 10
		}
		payload.MaxAttempts = &def
	}

	for i, domain := range *payload.AllowedDomains {
		(*payload.AllowedDomains)[i] = strings.TrimPrefix(domain, "www.")
	}

	normalizedStr := crawler.NormalizeURL(originalURL.String(), *payload.RemoveFragment, *payload.LowerCaseURLs)
	modifiedURL, _ := url.Parse(normalizedStr)
	return modifiedURL
}

// mapResultToDTO converte o resultado interno do crawler para o DTO da API.
func mapResultToDTO(result *crawler.CrawlResult, originalURL *url.URL) crawler.ResponseDTO {
	return crawler.ResponseDTO{
		StatusCode:  result.StatusCode,
		ContentType: result.ContentType,
		ElapsedTime: result.ElapsedTime.Nanoseconds(),
		Links:       result.Links,
		Title:       result.Title,
		Details: crawler.DetailsResponseDTO{
			CorrectURL: result.FinalURL.String(),
			Original:   mapURLToDetails(originalURL),
			Modified:   mapURLToDetails(result.FinalURL),
		},
	}
}

// mapURLToDetails converte uma url.URL para a struct de detalhes do DTO.
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
