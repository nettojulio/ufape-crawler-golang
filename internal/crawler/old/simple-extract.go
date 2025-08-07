package old

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/emirpasic/gods/queues/linkedlistqueue"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/nettojulio/ufape-crawler-golang/internal/graph"
	"github.com/nettojulio/ufape-crawler-golang/utils"
	"github.com/nettojulio/ufape-crawler-golang/utils/configs"
	"golang.org/x/net/html"
)

type QueueItem struct {
	URL   string
	Depth uint
}

var (
	GlobalQueue = linkedlistqueue.New()
	GlobalSeen  = hashset.New()
)

func InitialSearch() {
	startApplication := time.Now()

	g := graph.NewGraph()
	crawl(g, configs.InitialURL, configs.MaxDepth)

	fmt.Println("\n=== Subdomain Clusters ===")
	g.SubdomainClusters()
	for host, c := range g.Clusters {
		fmt.Printf("[%s] %d páginas\n", host, c.Total)
	}

	fmt.Println("\n=== Nós do Grafo ===")
	fmt.Printf("Total: %d\n", len(g.Nodes))

	g.SaveJSON()

	fmt.Printf("\nCrawling completed in %v\n", time.Since(startApplication))
}

func extractLinks(base *url.URL, body *html.Node) []string {
	var links []string
	uniqueLinks := hashset.New()
	curr := utils.NormalizeURL(base.String())
	curateNodes(body, base, curr, uniqueLinks, &links)
	return links
}

func crawl(g *graph.Graph, startURL string, maxDepth uint) {
	startURL = utils.NormalizeURL(startURL)
	GlobalQueue.Enqueue(QueueItem{URL: startURL, Depth: 1})
	GlobalSeen.Add(startURL)

	for !GlobalQueue.Empty() {
		itemRaw, _ := GlobalQueue.Dequeue()
		item := itemRaw.(QueueItem)

		if item.Depth > maxDepth {
			continue
		}
		curateUrl(item, g)
	}
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

func makeRequest(urlStr string) (*http.Response, error) {
	client := &http.Client{
		Timeout: configs.RequestTimeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 10 {
				return fmt.Errorf("stopped after 10 redirects")
			}
			return nil
		},
	}
	req, err := http.NewRequest("GET", urlStr, nil)

	if err != nil {
		fmt.Printf("error creating request for %s: %v\n", urlStr, err)
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:142.0) Gecko/20100101 Firefox/142.0")
	req.Header.Set("Accept", "*/*")
	req.Header.Add("Connection", "keep-alive")

	return client.Do(req)
}

func curateUrl(item QueueItem, g *graph.Graph) {
	start := time.Now()
	resp, err := makeRequest(item.URL)
	if err != nil {
		errorNode(&start, resp, err, item, g)
		return
	}

	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			fmt.Printf("Error closing response body for %s: %v\n", item.URL, err)
			return
		}
	}(resp.Body)

	successNode(resp, item, g, &start)
}

func errorNode(start *time.Time, resp *http.Response, err error, item QueueItem, g *graph.Graph) {
	elapsed := time.Since(*start)
	var contentType string
	if resp != nil {
		contentType = resp.Header.Get("Content-Type")
	}
	fmt.Printf("Failed [Depth %d] %s | Status: %d | Content-Type: %s | Time: %v: %s | Error: %v\n", item.Depth, item.URL, 0, contentType, elapsed, item.URL, err)
	g.AddNode(item.URL, item.Depth, 0, time.Since(*start), contentType, "[Error fetching page]")
}

func successNode(resp *http.Response, item QueueItem, g *graph.Graph, start *time.Time) {
	elapsed := time.Since(*start)
	contentType := resp.Header.Get("Content-Type")
	fmt.Printf("[Depth %d] %s | Status: %d | Time: %v\n", item.Depth, item.URL, resp.StatusCode, elapsed)

	if resp.StatusCode != http.StatusOK {
		return
	}

	doc, err := html.Parse(resp.Body)

	if err != nil {
		return
	}

	title := getTitle(doc)

	g.AddNode(item.URL, item.Depth, resp.StatusCode, elapsed, contentType, title)
	base, err := url.Parse(item.URL)
	if err != nil {
		fmt.Printf("Error parsing URL %s: %v\n", item.URL, err)
		return
	}
	links := extractLinks(base, doc)
	curateFoundLinks(links, g, item)
}

func getTitle(doc *html.Node) (title string) {
	title = "[Empty title]"

	titleNode := getTitleNode(doc)

	if titleNode != nil && titleNode.FirstChild != nil {
		title = strings.TrimSpace(titleNode.FirstChild.Data)
	}
	return
}

func curateFoundLinks(links []string, g *graph.Graph, item QueueItem) {
	for _, link := range links {
		if !configs.FliterLinks {
			g.AddEdge(item.URL, link)
		}
		if !GlobalSeen.Contains(link) {
			g.AddEdge(item.URL, link)
			GlobalSeen.Add(link)
			GlobalQueue.Enqueue(QueueItem{URL: link, Depth: item.Depth + 1})
		}
	}
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
	norm := utils.NormalizeURL(rawURL)
	if norm != curr && !uniqueLinks.Contains(norm) {
		uniqueLinks.Add(norm)
		*links = append(*links, norm)
	}
}
