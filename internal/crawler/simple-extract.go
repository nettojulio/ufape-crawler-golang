package crawler

import (
	"fmt"
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
	Depth int
}

func extractLinks(base *url.URL, body *html.Node) []string {
	var links []string
	uniqueLinks := hashset.New()

	curr := utils.NormalizeURL(base.String())

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					if strings.HasPrefix(attr.Val, "mailto:") || strings.HasPrefix(attr.Val, "tel:") {
						continue
					}
					u, err := base.Parse(attr.Val)
					if err != nil {
						continue
					}
					if u.Host == base.Host || strings.HasSuffix(u.Host, "."+base.Host) {
						norm := utils.NormalizeURL(u.String())
						if norm != curr && !uniqueLinks.Contains(norm) {
							uniqueLinks.Add(norm)
							links = append(links, norm)
						}
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(body)
	return links
}

func crawl(startURL string, maxDepth int) *graph.Graph {
	seen := hashset.New()
	queue := linkedlistqueue.New()
	g := graph.NewGraph()

	startURL = utils.NormalizeURL(startURL)
	queue.Enqueue(QueueItem{URL: startURL, Depth: 0})
	seen.Add(startURL)

	for !queue.Empty() {
		itemRaw, _ := queue.Dequeue()
		item := itemRaw.(QueueItem)

		if item.Depth > maxDepth {
			continue
		}

		client := &http.Client{Timeout: configs.RequestTimeout, CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 10 {
				return fmt.Errorf("stopped after 10 redirects")
			}
			return nil
		}}
		start := time.Now()
		resp, err := client.Get(item.URL)
		if err != nil {
			elapsed := time.Since(start)
			fmt.Printf("❌ Failed [Depth %d] %s | Status: %d | Time: %v: %s | Error: %v\n", item.Depth, item.URL, 0, elapsed, item.URL, err)
			g.AddNode(item.URL, item.Depth, 0, time.Since(start))
			continue
		}
		func() {
			defer resp.Body.Close()
			elapsed := time.Since(start)
			fmt.Printf("[Depth %d] %s | Status: %d | Time: %v\n",
				item.Depth, item.URL, resp.StatusCode, elapsed)

			g.AddNode(item.URL, item.Depth, resp.StatusCode, elapsed)

			if resp.StatusCode != http.StatusOK {
				return
			}

			doc, err := html.Parse(resp.Body)
			if err != nil {
				return
			}

			base, _ := url.Parse(item.URL)
			links := extractLinks(base, doc)

			for _, link := range links {
				if !configs.FliterLinks {
					g.AddEdge(item.URL, link)
				}
				if !seen.Contains(link) {
					g.AddEdge(item.URL, link)
					seen.Add(link)
					queue.Enqueue(QueueItem{URL: link, Depth: item.Depth + 1})
				}
			}
		}()
	}

	return g
}

func InitialSearch() {
	start := time.Now()
	gr := crawl(configs.InitialURL, configs.MaxDepth)

	fmt.Println("\n=== Subdomain Clusters ===")
	for host, urls := range gr.SubdomainClusters() {
		fmt.Printf("[%s] %d páginas\n", host, len(urls))
	}

	fmt.Println("\n=== Nós do Grafo ===")
	fmt.Printf("Total: %d\n", len(gr.Nodes))

	fmt.Printf("\nCrawling completed in %v\n", time.Since(start))
}
