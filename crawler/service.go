package crawler

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/net/html"
)

func CrawlerService(payload CorrectPayload, original, modified url.URL) (ResponseCrawl, error) {
	client := NewHTTPClient(*payload.Timeout)

	req, err := NewRequest(modified.String())
	if err != nil {
		return ResponseCrawl{}, err
	}

	start := time.Now()
	resp, err := client.Do(req)
	elapsed := time.Since(start)
	if err != nil {
		return ResponseCrawl{}, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	contentType := resp.Header.Get("Content-Type")
	if resp.StatusCode != http.StatusOK {
		return buildErrorResponse(resp.StatusCode, contentType, elapsed, original, modified), nil
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return ResponseCrawl{}, err
	}

	title := GetTitle(doc)
	links := ExtractLinks(&modified, doc, &payload)

	return ResponseCrawl{
		StatusCode:  resp.StatusCode,
		ContentType: contentType,
		ElapsedTime: fmt.Sprintf("%+v", elapsed),
		Links:       links,
		Title:       title,
		Details: DetailsResponse{
			CorrectURL: modified.String(),
			Original:   original,
			Modified:   modified,
		},
	}, nil
}

func buildErrorResponse(status int, contentType string, elapsed time.Duration,
	original, modified url.URL) ResponseCrawl {
	return ResponseCrawl{
		StatusCode:  status,
		ContentType: contentType,
		ElapsedTime: fmt.Sprintf("%+v", elapsed),
		Links:       LinksResponse{Available: []string{}, Unavailable: []string{}},
		Title:       "",
		Details: DetailsResponse{
			CorrectURL: modified.String(),
			Original:   original,
			Modified:   modified,
		},
	}
}
