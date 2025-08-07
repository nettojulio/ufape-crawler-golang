package graph

import (
	"net/url"
	"time"
)

type PageNode struct {
	URL        string        `json:"url"`
	Depth      int           `json:"depth"`
	StatusCode int           `json:"status_code"`
	LoadTime   time.Duration `json:"load_time"`
	Links      []string      `json:"links"`
	// Content-Type
}

type Graph struct {
	Nodes map[string]*PageNode
}

func NewGraph() *Graph {
	return &Graph{Nodes: make(map[string]*PageNode)}
}

func (g *Graph) AddNode(url string, depth int, status int, loadTime time.Duration) {
	if _, exists := g.Nodes[url]; !exists {
		g.Nodes[url] = &PageNode{
			URL:        url,
			Depth:      depth,
			StatusCode: status,
			LoadTime:   loadTime,
			Links:      []string{},
		}
	}
}

func (g *Graph) AddEdge(from, to string) {
	if node, exists := g.Nodes[from]; exists {
		for _, l := range node.Links {
			if l == to {
				return
			}
		}
		node.Links = append(node.Links, to)
	}
}

func (g *Graph) SubdomainClusters() map[string][]string {
	clusters := make(map[string][]string)
	for urlStr := range g.Nodes {
		u, err := url.Parse(urlStr)
		if err != nil {
			continue
		}
		clusters[u.Host] = append(clusters[u.Host], urlStr)
	}
	return clusters
}
