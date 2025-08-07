package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nettojulio/ufape-crawler-golang/internal"
	"github.com/nettojulio/ufape-crawler-golang/utils/configs"
)

var (
	Version = "development"
)

func init() {
	flag.UintVar(&configs.MaxDepth, "depth", 3, "Profundidade máxima de crawling.")
	flag.BoolVar(&configs.RemoveFragment, "remove-fragment", true, "Remove fragmentos de URL (#...).")
	flag.BoolVar(&configs.LowerCaseURLs, "lower-case", false, "Converte todas as URLs para minúsculas.")
	flag.DurationVar(&configs.RequestTimeout, "timeout", 90*time.Second, "Timeout para as requisições HTTP.")
	flag.BoolVar(&configs.FliterLinks, "filter-links", true, "Filtra links já visitados.")
	flag.StringVar(&configs.InitialURL, "initial-url", "https://ufape.edu.br", "URL inicial para o crawler.")
}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func main() {
	flag.Parse()
	fmt.Printf("Version: %s\n", Version)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Validator = &CustomValidator{validator: validator.New()}

	e.POST("/", crawlerHandler)
	e.GET("/", healthCheckHandler)

	e.Logger.Fatal(e.Start(":8080"))
}

func crawlerHandler(c echo.Context) error {
	var payload internal.CorrectPayload

	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "invalid request body",
		})
	}

	if err := c.Validate(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	originalUrlDetails, err := url.Parse(payload.Url)
	modifiedUrlDetails := *originalUrlDetails

	pass, ok := modifiedUrlDetails.User.Password()
	passEmpty := !ok || pass == ""
	if modifiedUrlDetails.User != nil {
		if modifiedUrlDetails.User.Username() == "" && passEmpty {
			modifiedUrlDetails.User = nil
		}
	}

	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "invalid url",
		})
	}

	if payload.Timeout == nil {
		def := 60
		payload.Timeout = &def
	}

	if payload.RemoveFragment == nil {
		def := true
		payload.RemoveFragment = &def
	}

	if payload.AllowedDomains == nil {
		def := []string{originalUrlDetails.Host}
		payload.AllowedDomains = &def
	}

	if payload.CollectSubdomains == nil {
		def := true
		payload.CollectSubdomains = &def
	}

	for i, domain := range *payload.AllowedDomains {
		(*payload.AllowedDomains)[i] = strings.TrimPrefix(domain, "www.")
	}

	if *payload.RemoveFragment {
		modifiedUrlDetails.Fragment = ""
	}

	response, err := internal.CrawlerService(payload, *originalUrlDetails, modifiedUrlDetails)

	return c.JSON(http.StatusOK, response)
}

func healthCheckHandler(c echo.Context) error {
	return c.JSON(200, map[string]string{
		"status":  "OK",
		"version": Version,
	})
}
