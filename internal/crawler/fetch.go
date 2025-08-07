package crawler

import (
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"time"

	"github.com/nettojulio/ufape-crawler-golang/utils/configs"
	"golang.org/x/net/html"
)

var httpClient = &http.Client{
	Timeout: configs.RequestTimeout,
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		if len(via) >= 10 {
			return fmt.Errorf("stopped after 10 redirects")
		}
		return nil
	},
}

func (c *Crawler) processURL(item QueueItem) error {
	start := time.Now()

	resp, err := makeRequest(item.URL, "")
	if err != nil {
		c.graph.AddNode(item.URL, item.Depth, 0, time.Since(start), "", "[Error fetching page]")
		return err
	}
	defer resp.Body.Close()

	elapsed := time.Since(start)
	contentType := resp.Header.Get("Content-Type")
	fmt.Printf("[Depth %d] %s | Status: %d | Time: %v\n", item.Depth, item.URL, resp.StatusCode, elapsed)

	if resp.StatusCode != http.StatusOK {
		return nil
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return err
	}

	title := getTitle(doc)
	c.graph.AddNode(item.URL, item.Depth, resp.StatusCode, elapsed, contentType, title)

	base, err := url.Parse(item.URL)
	if err != nil {
		return err
	}

	links := extractLinks(base, doc)
	c.graph.AddEdge(item.URL, links...)
	for _, link := range links {
		//_, b := c.seen[link]
		//if b {
		//	fmt.Println("Saw")
		//	continue
		//}
		//r, result := c.canBeQueueed(link)
		//if result {
		//	fmt.Println(r)
		//	c.graph.Nodes[item.URL].SimLink = true
		//}
		c.enqueue(link, item.Depth+1)
	}
	if len(links) != 0 && slices.Contains(c.graph.Edges[item.URL].To, item.URL) {
		fmt.Printf("Skipping self-link: %s\n", item.URL)
		return nil
	}

	return nil
}

func makeRequest(urlStr string, forceScheme string) (*http.Response, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	if forceScheme == "http" || forceScheme == "https" {
		u.Scheme = forceScheme
	}

	if u.Host == "" {
		return nil, fmt.Errorf("URL inválida: %s", urlStr)
	}

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64)")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Connection", "keep-alive")

	return httpClient.Do(req)
}

func getTitle(doc *html.Node) string {
	titleNode := getTitleNode(doc)
	if titleNode != nil && titleNode.FirstChild != nil {
		return strings.TrimSpace(titleNode.FirstChild.Data)
	}
	return "[Empty title]"
}

func getTitleNode(n *html.Node) *html.Node {
	if n.Type == html.ElementNode && n.Data == "title" {
		return n
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if result := getTitleNode(c); result != nil {
			return result
		}
	}
	return nil
}

func (c *Crawler) canBeQueueed(originalLink string) (map[string]bool, bool) {
	var (
		s1    = strings.Clone(originalLink)
		l1, _ = url.Parse(s1)
		s2    = strings.Clone(originalLink)
		l2, _ = url.Parse(s2)
		s3    = strings.Clone(originalLink)
		l3, _ = url.Parse(s3)
		s4    = strings.Clone(originalLink)
		l4, _ = url.Parse(s4)
	)
	l1.Scheme = "http"
	l2.Scheme = "https"
	l3.Scheme = "http"
	l4.Scheme = "https"

	if !strings.HasPrefix(l1.Host, "www.") {
		l1.Host = "www." + l1.Host
	}

	if !strings.HasPrefix(l2.Host, "www.") {
		l2.Host = "www." + l2.Host
	}

	if strings.Contains(l3.Host, "www") {
		l3.Host = strings.TrimPrefix(l3.Host, "www.")
	}
	if strings.Contains(l4.Host, "www") {
		l4.Host = strings.TrimPrefix(l4.Host, "www.")
	}

	results := make(map[string]bool)
	_, okL1 := c.seen[l1.String()]
	_, okL2 := c.seen[l2.String()]
	_, okL3 := c.seen[l3.String()]
	_, okL4 := c.seen[l4.String()]

	results = map[string]bool{
		l1.String(): okL1,
		l2.String(): okL2,
		l3.String(): okL3,
		l4.String(): okL4,
	}
	return results, okL1 || okL2 || okL3 || okL4
}
