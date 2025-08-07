package graph

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"time"
)

type Node struct {
	URL         string        `json:"url"`
	Depth       uint          `json:"depth"`
	StatusCode  int           `json:"status_code"`
	LoadTime    time.Duration `json:"load_time"`
	Links       uint          `json:"links"`
	ContentType string        `json:"content_type,omitempty"`
	Title       string        `json:"title,omitempty"`
	SimLink     bool          `json:"sim_link,omitempty"`
}

type Edge struct {
	From string   `json:"from"`
	To   []string `json:"to"`
}

type Cluster struct {
	Host  string `json:"host"`
	Total int    `json:"total"`
}

type Tag struct {
	Key         string `json:"key"`
	Description string `json:"description"`
}

type Graph struct {
	Nodes    map[string]*Node    `json:"nodes"`
	Edges    map[string]*Edge    `json:"edges"`
	Clusters map[string]*Cluster `json:"clusters"`
	Tags     map[string]*Tag     `json:"tags"`
}

func NewGraph() *Graph {
	return &Graph{
		Nodes:    make(map[string]*Node),
		Edges:    make(map[string]*Edge),
		Clusters: make(map[string]*Cluster),
		Tags:     make(map[string]*Tag),
	}
}

func (g *Graph) AddNode(url string, depth uint, status int, loadTime time.Duration, contentType string, title string) {
	if _, exists := g.Nodes[url]; !exists {
		g.Nodes[url] = &Node{
			URL:         url,
			Depth:       depth,
			StatusCode:  status,
			LoadTime:    loadTime,
			Links:       0,
			ContentType: contentType,
			Title:       title,
		}
	}
}

func (g *Graph) AddCluster(hosts string) {
	if cluster, exists := g.Clusters[hosts]; exists {
		cluster.Total += 1
	} else {
		g.Clusters[hosts] = &Cluster{
			Host:  hosts,
			Total: 1,
		}
	}
}

func (g *Graph) AddEdge(from string, to ...string) {
	if edge, exists := g.Edges[from]; exists {
		edge.To = append(edge.To, to...)
	} else {
		g.Edges[from] = &Edge{
			From: from,
			To:   append(make([]string, 0), to...),
		}
	}
	if len(to) != 0 {
		g.Nodes[from].Links++
	}
}

func (g *Graph) SubdomainClusters() {
	for urlStr := range g.Nodes {
		u, err := url.Parse(urlStr)
		if err != nil {
			continue
		}
		g.AddCluster(u.Host)
	}
}

func (g *Graph) SaveJSON() {
	fileName := fmt.Sprintf("graph_%s.json", time.Now().Format("02-01-2006 15:04:05"))
	file, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err = encoder.Encode(g); err != nil {
		panic(err)
	}

	fmt.Printf("Arquivo '%s' salvo com sucesso!", fileName)
}

//func (g *Graph) ShortestPath(start, end string) ([]string, bool) {
//	start = utils.NormalizeURL(start)
//	end = utils.NormalizeURL(end)
//	if _, ok := g.Nodes[start]; !ok {
//		return nil, false
//	}
//	if _, ok := g.Nodes[end]; !ok {
//		return nil, false
//	}
//
//	queue := []string{start}
//	prev := make(map[string]string)
//	visited := hashset.New()
//	visited.Add(start)
//
//	for len(queue) > 0 {
//		curr := queue[0]
//		queue = queue[1:]
//
//		if curr == end {
//			var path []string
//			for at := end; at != ""; at = prev[at] {
//				path = append([]string{at}, path...)
//			}
//			return path, true
//		}
//
//		for _, neigh := range g.Nodes[curr].Links {
//			if !visited.Contains(neigh) {
//				visited.Add(neigh)
//				prev[neigh] = curr
//				queue = append(queue, neigh)
//			}
//		}
//	}
//	return nil, false
//}
//
//func (g *Graph) ExportDOT(filename string) error {
//	var sb strings.Builder
//	sb.WriteString("digraph G {\n")
//	sb.WriteString("\tnode [shape=ellipse, fontname=Helvetica];\n")
//
//	for _, node := range g.Nodes {
//		sb.WriteString(fmt.Sprintf("\t\"%s\";\n", node.URL))
//		for _, link := range node.Links {
//			sb.WriteString(fmt.Sprintf("\t\"%s\" -> \"%s\";\n", node.URL, link))
//		}
//	}
//
//	sb.WriteString("}\n")
//	return os.WriteFile(filename, []byte(sb.String()), 0644)
//}
//
//func (g *Graph) PageRank(iterations int, d float64) map[string]float64 {
//	n := len(g.Nodes)
//	rank := make(map[string]float64)
//
//	for url := range g.Nodes {
//		rank[url] = 1.0 / float64(n)
//	}
//
//	for i := 0; i < iterations; i++ {
//		newRank := make(map[string]float64)
//		for url := range g.Nodes {
//			newRank[url] = (1 - d) / float64(n)
//		}
//
//		for url, node := range g.Nodes {
//			if len(node.Links) == 0 {
//				continue
//			}
//			share := rank[url] / float64(len(node.Links))
//			for _, link := range node.Links {
//				newRank[link] += d * share
//			}
//		}
//		rank = newRank
//	}
//	return rank
//}
