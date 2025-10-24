package crawler

import (
	"net/url"
	"reflect"
	"slices"
	"strings"
	"testing"

	"golang.org/x/net/html"
)

func parseHTML(t *testing.T, htmlString string) *html.Node {
	t.Helper()
	doc, err := html.Parse(strings.NewReader(htmlString))
	if err != nil {
		t.Fatalf("Failed to parse test HTML: %v", err)
	}
	return doc
}

func TestGetTitle(t *testing.T) {
	testCases := []struct {
		name          string
		htmlContent   string
		expectedTitle string
	}{
		{
			name:          "Simple Title",
			htmlContent:   `<html><head><title>Page Title</title></head></html>`,
			expectedTitle: "Page Title",
		},
		{
			name:          "Title with leading/trailing spaces",
			htmlContent:   `<html><head><title>  My Title with Spaces  </title></head></html>`,
			expectedTitle: "My Title with Spaces",
		},
		{
			name:          "No title tag",
			htmlContent:   `<html><head></head><body><h1>Hello</h1></body></html>`,
			expectedTitle: "[Empty title]",
		},
		{
			name:          "Empty title tag",
			htmlContent:   `<html><head><title></title></head></html>`,
			expectedTitle: "[Empty title]",
		},
		{
			name:          "Deeply nested title",
			htmlContent:   `<html><body><div><p><title>Deep Title</title></p></div></body></html>`,
			expectedTitle: "Deep Title",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			doc := parseHTML(t, tc.htmlContent)
			actualTitle := GetTitle(doc)
			if actualTitle != tc.expectedTitle {
				t.Errorf("expected title %q, got %q", tc.expectedTitle, actualTitle)
			}
		})
	}
}

func TestExtractLinks(t *testing.T) {
	baseURL, _ := url.Parse("https://example.com/some/path")

	testCases := []struct {
		name                string
		htmlContent         string
		opts                ParseOptions
		expectedAvailable   []string
		expectedUnavailable []string
	}{
		{
			name: "Basic extraction with relative and absolute links",
			htmlContent: `
                <a href="/page1">Page 1</a>
                <a href="https://example.com/page2">Page 2</a>
                <a href="https://otherdomain.com/page">Other Domain</a>
                <a href="/page1">Duplicate Page 1</a>
                <a href="mailto:test@example.com">Email</a>
                <a href="tel:+123456">Phone</a>
            `,
			opts: ParseOptions{
				BaseURL:           baseURL,
				AllowedDomains:    []string{"example.com"},
				CollectSubdomains: false,
				RemoveFragment:    true,
				LowerCaseURLs:     false,
			},
			expectedAvailable:   []string{"https://example.com/page1", "https://example.com/page2"},
			expectedUnavailable: []string{"https://otherdomain.com/page"},
		},
		{
			name: "Collect subdomains enabled",
			htmlContent: `
                <a href="https://blog.example.com/post1">Blog Post</a>
                <a href="https://docs.example.com/api">Docs</a>
                <a href="https://example.org">Different TLD</a>
            `,
			opts: ParseOptions{
				BaseURL:           baseURL,
				AllowedDomains:    []string{"example.com"},
				CollectSubdomains: true,
				RemoveFragment:    true,
				LowerCaseURLs:     false,
			},
			expectedAvailable:   []string{"https://blog.example.com/post1", "https://docs.example.com/api"},
			expectedUnavailable: []string{"https://example.org"},
		},
		{
			name: "Remove fragment enabled",
			htmlContent: `
                <a href="/page#section1">Section Link</a>
                <a href="/page#section2">Another Section</a>
            `,
			opts: ParseOptions{
				BaseURL:           baseURL,
				AllowedDomains:    []string{"example.com"},
				CollectSubdomains: false,
				RemoveFragment:    true,
				LowerCaseURLs:     false,
			},
			expectedAvailable:   []string{"https://example.com/page"},
			expectedUnavailable: []string{},
		},
		{
			name:        "Lowercase URLs enabled",
			htmlContent: `<a href="/Some/MixedCase/Path">Mixed Case</a>`,
			opts: ParseOptions{
				BaseURL:           baseURL,
				AllowedDomains:    []string{"example.com"},
				CollectSubdomains: false,
				RemoveFragment:    true,
				LowerCaseURLs:     true,
			},
			expectedAvailable:   []string{"https://example.com/some/mixedcase/path"},
			expectedUnavailable: []string{},
		},
		{
			name:        "No links found",
			htmlContent: `<p>This is a paragraph with no links.</p>`,
			opts: ParseOptions{
				BaseURL:        baseURL,
				AllowedDomains: []string{"example.com"},
			},
			expectedAvailable:   []string{},
			expectedUnavailable: []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			doc := parseHTML(t, tc.htmlContent)
			links := ExtractLinks(doc, tc.opts)

			slices.Sort(links.Available)
			slices.Sort(tc.expectedAvailable)
			slices.Sort(links.Unavailable)
			slices.Sort(tc.expectedUnavailable)

			if !reflect.DeepEqual(links.Available, tc.expectedAvailable) {
				t.Errorf("mismatch in available links\ngot:  %v\nwant: %v", links.Available, tc.expectedAvailable)
			}
			if !reflect.DeepEqual(links.Unavailable, tc.expectedUnavailable) {
				t.Errorf("mismatch in unavailable links\ngot:  %v\nwant: %v", links.Unavailable, tc.expectedUnavailable)
			}
		})
	}
}
