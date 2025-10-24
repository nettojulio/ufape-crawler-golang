package crawler

import (
	"net/url"
	"slices"
	"strings"

	"golang.org/x/net/html"
)

// ParseOptions contém as configurações necessárias para o processo de parsing.
type ParseOptions struct {
	BaseURL           *url.URL
	AllowedDomains    []string
	CollectSubdomains bool
	RemoveFragment    bool
	LowerCaseURLs     bool
}

func GetTitle(doc *html.Node) string {
	var title string
	if n := getTitleNode(doc); n != nil && n.FirstChild != nil {
		title = strings.TrimSpace(n.FirstChild.Data)
	}
	if title == "" {
		return "[Empty title]"
	}
	return title
}

func ExtractLinks(doc *html.Node, opts ParseOptions) LinksResponse {
	links := LinksResponse{Available: []string{}, Unavailable: []string{}}
	unique := make(map[string]struct{})

	currentNormalizedURL := NormalizeURL(opts.BaseURL.String(), opts.RemoveFragment, opts.LowerCaseURLs)
	unique[currentNormalizedURL] = struct{}{}

	curateNodes(doc, opts, unique, &links)
	return links
}

func getTitleNode(n *html.Node) *html.Node {
	if n.Type == html.ElementNode && n.Data == "title" {
		return n
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if res := getTitleNode(c); res != nil {
			return res
		}
	}
	return nil
}

func curateNodes(n *html.Node, opts ParseOptions, unique map[string]struct{}, links *LinksResponse) {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, attr := range n.Attr {
			if attr.Key == "href" {
				processHref(attr.Val, opts, unique, links)
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		curateNodes(c, opts, unique, links)
	}
}

func processHref(href string, opts ParseOptions, unique map[string]struct{}, links *LinksResponse) {
	u, err := opts.BaseURL.Parse(href)
	if err != nil || u.Host == "" || u.Scheme == "" || strings.HasPrefix(href, "mailto:") || strings.HasPrefix(href, "tel:") {
		return
	}

	normalized := NormalizeURL(u.String(), opts.RemoveFragment, opts.LowerCaseURLs)

	if _, exists := unique[normalized]; exists {
		return
	}
	unique[normalized] = struct{}{}

	parsedNormalized, err := url.Parse(normalized)
	if err != nil {
		return
	}
	host := parsedNormalized.Host

	if slices.Contains(opts.AllowedDomains, host) || (opts.CollectSubdomains && IsSubdomainHost(host, opts.AllowedDomains)) {
		links.Available = append(links.Available, normalized)
	} else {
		links.Unavailable = append(links.Unavailable, normalized)
	}
}
