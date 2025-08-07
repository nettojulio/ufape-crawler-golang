package crawler

import (
	"fmt"
	"time"

	"github.com/nettojulio/ufape-crawler-golang/internal/graph"
	"github.com/nettojulio/ufape-crawler-golang/utils/configs"
)

type Crawler struct {
	queue *Queue
	seen  map[string]struct{}
	graph *graph.Graph
}

func NewCrawler() *Crawler {
	return &Crawler{
		queue: NewQueue(),
		seen:  make(map[string]struct{}),
		graph: graph.NewGraph(),
	}
}

func (c *Crawler) Run(startURL string, maxDepth uint) {
	start := time.Now()
	startURL = normalizeURL(startURL)

	c.enqueue(startURL, 1)

	for !c.queue.Empty() {
		item := c.queue.Dequeue()

		if item.Depth > maxDepth {
			continue
		}

		if err := c.processURL(item); err != nil {
			fmt.Printf("Error processing %s: %v\n", item.URL, err)
		}
	}

	fmt.Println("\n=== Subdomain Clusters ===")
	c.graph.SubdomainClusters()
	for host, cluster := range c.graph.Clusters {
		fmt.Printf("[%s] %d páginas\n", host, cluster.Total)
	}

	fmt.Println("\n=== Nós do Grafo ===")
	fmt.Printf("Total: %d\n", len(c.graph.Nodes))

	c.graph.SaveJSON()
	fmt.Printf("\nCrawling completed in %v\n", time.Since(start))
}

func (c *Crawler) enqueue(url string, depth uint) {
	if _, exists := c.seen[url]; exists {
		return
	}
	c.seen[url] = struct{}{}
	c.queue.Enqueue(QueueItem{URL: url, Depth: depth})
}

func Execute() {
	c := NewCrawler()
	c.Run(configs.InitialURL, configs.MaxDepth)
}
