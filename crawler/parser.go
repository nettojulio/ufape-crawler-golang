package crawler

import (
	"net/url"
	"slices"
	"strings"

	"golang.org/x/net/html"
)

func GetTitle(doc *html.Node) string {
	if n := getTitleNode(doc); n != nil && n.FirstChild != nil {
		return strings.TrimSpace(n.FirstChild.Data)
	}
	return "[Empty title]"
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

func ExtractLinks(base *url.URL, doc *html.Node, payload *CorrectPayload) LinksResponse {
	links := LinksResponse{Available: []string{}, Unavailable: []string{}}
	unique := make(map[string]struct{})
	curr := NormalizeURL(base.String(), payload)
	curateNodes(doc, base, curr, unique, &links, payload)
	return links
}

func curateNodes(n *html.Node, base *url.URL, curr string, unique map[string]struct{},
	links *LinksResponse, payload *CorrectPayload) {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, attr := range n.Attr {
			if attr.Key == "href" {
				processHref(attr.Val, base, curr, unique, links, payload)
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		curateNodes(c, base, curr, unique, links, payload)
	}
}

func processHref(href string, base *url.URL, curr string,
	unique map[string]struct{}, links *LinksResponse, payload *CorrectPayload) {

	if strings.HasPrefix(href, "mailto:") || strings.HasPrefix(href, "tel:") {
		return
	}
	u, err := base.Parse(href)
	if err != nil || u.Host == "" || u.Scheme == "" {
		return
	}
	addIfUnique(u.String(), curr, unique, links, payload)
}

func addIfUnique(rawURL, curr string, unique map[string]struct{},
	links *LinksResponse, payload *CorrectPayload) {
	norm := NormalizeURL(rawURL, payload)
	if norm == curr {
		return
	}
	if _, exists := unique[norm]; exists {
		return
	}
	unique[norm] = struct{}{}

	info, err := url.Parse(norm)
	if err != nil {
		return
	}

	if slices.Contains(*payload.AllowedDomains, info.Host) ||
		(*payload.CollectSubdomains && IsSubdomainHost(info.Host, *payload.AllowedDomains)) {
		links.Available = append(links.Available, norm)
	} else {
		links.Unavailable = append(links.Unavailable, norm)
	}
}
