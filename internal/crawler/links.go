package crawler

import (
	"net/url"
	"strings"

	"github.com/emirpasic/gods/sets/hashset"
	"github.com/nettojulio/ufape-crawler-golang/utils"
	"golang.org/x/net/html"
)

func extractLinks(base *url.URL, body *html.Node) []string {
	var links []string
	uniqueLinks := hashset.New()
	curr := normalizeURL(base.String())
	curateNodes(body, base, curr, uniqueLinks, &links)
	return links
}

func curateNodes(n *html.Node, base *url.URL, curr string, uniqueLinks *hashset.Set, links *[]string) {
	if isAnchorNode(n) {
		processAnchor(n, base, curr, uniqueLinks, links)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		curateNodes(c, base, curr, uniqueLinks, links)
	}
}

func isAnchorNode(n *html.Node) bool {
	return n.Type == html.ElementNode && n.Data == "a"
}

func processAnchor(n *html.Node, base *url.URL, curr string, uniqueLinks *hashset.Set, links *[]string) {
	for _, attr := range n.Attr {
		if attr.Key != "href" {
			continue
		}
		processHref(attr.Val, base, curr, uniqueLinks, links)
	}
}

func processHref(href string, base *url.URL, curr string, uniqueLinks *hashset.Set, links *[]string) {
	if isNonNavigableLink(href) {
		return
	}
	u, err := base.Parse(href)
	if err != nil {
		return
	}
	if isSameDomain(u, base) {
		addIfUnique(u.String(), curr, uniqueLinks, links)
	}
}

func isNonNavigableLink(href string) bool {
	return strings.HasPrefix(href, "mailto:") || strings.HasPrefix(href, "tel:")
}

func isSameDomain(u, base *url.URL) bool {
	return u.Host == base.Host || strings.HasSuffix(u.Host, "."+base.Host)
}

func addIfUnique(rawURL, curr string, uniqueLinks *hashset.Set, links *[]string) {
	norm := normalizeURL(rawURL)
	if norm != curr && !uniqueLinks.Contains(norm) {
		uniqueLinks.Add(norm)
		*links = append(*links, norm)
	}
}

func normalizeURL(raw string) string {
	return utils.NormalizeURL(raw)
}
