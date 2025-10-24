package crawler

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/html"
)

type HTTPGetter interface {
	Get(ctx context.Context, url string) (*http.Response, error)
}

type Service struct {
	httpClient HTTPGetter
}

func NewService(httpClient HTTPGetter) *Service {
	return &Service{
		httpClient: httpClient,
	}
}

func (s *Service) Crawl(ctx context.Context, payload Payload, originalURL, modifiedURL *url.URL) (*CrawlResult, error) {
	start := time.Now()
	resp, err := s.httpClient.Get(ctx, modifiedURL.String())
	elapsed := time.Since(start)

	if err != nil {
		return &CrawlResult{
			StatusCode:  http.StatusServiceUnavailable,
			ElapsedTime: elapsed,
			Title:       err.Error(),
			FinalURL:    modifiedURL,
			Links: LinksResponse{
				Available:   []string{},
				Unavailable: []string{},
			},
		}, nil
	}

	result := &CrawlResult{
		StatusCode:  resp.StatusCode,
		ContentType: resp.Header.Get("Content-Type"),
		ElapsedTime: elapsed,
		Body:        resp.Body,
		FinalURL:    resp.Request.URL,
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return result, nil
	}

	doc, err := html.Parse(resp.Body)
	resp.Body.Close()
	if err != nil {
		if strings.Contains(resp.Header.Get("Content-Type"), "text/html") {
			return nil, fmt.Errorf("failed to parse html: %w", err)
		}
		result.Title = "[Empty Title]"
		result.Links.Available = []string{}
		result.Links.Unavailable = []string{}
		return result, nil

	}

	result.Title = GetTitle(doc)
	result.Links = ExtractLinks(doc, ParseOptions{
		BaseURL:           result.FinalURL,
		AllowedDomains:    *payload.AllowedDomains,
		CollectSubdomains: *payload.CollectSubdomains,
		RemoveFragment:    *payload.RemoveFragment,
		LowerCaseURLs:     *payload.LowerCaseURLs,
	})

	return result, nil
}
