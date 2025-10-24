package main

import (
	"bytes"
	"container/list"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/nettojulio/ufape-crawler-golang/internal/crawler"
)

const (
	DEFAULT_TIMEOUT = 60 * time.Second
	INITIAL_URL     = "https://ufape.edu.br"
	MAX_DEPTH       = math.MaxInt
)

type FinalResponse struct {
	Nodes       []Node `json:"nodes"`
	Links       []Link `json:"links"`
	GeneratedAt int64  `json:"generatedAt"`
}

type Node struct {
	ID          string `json:"id"`
	Depth       int    `json:"depth"`
	StatusCode  int    `json:"statusCode"`
	ContentType string `json:"contentType"`
	ElapsedTime int64  `json:"elapsedTime"`
	Title       string `json:"title"`
	Domain      string `json:"domain"`
}

func NewNode(url string, depth int, response *crawler.ResponseDTO) Node {
	return Node{
		ID:          url,
		Depth:       depth,
		StatusCode:  response.StatusCode,
		ContentType: response.ContentType,
		ElapsedTime: response.ElapsedTime,
		Title:       response.Title,
		Domain:      response.Details.Original.Host,
	}
}

type Link struct {
	Source string `json:"source"`
	Target string `json:"target"`
}

func NewLink(source, target string) Link {
	return Link{
		Source: source,
		Target: target,
	}
}

func NewRequestPayload(url string) *crawler.Payload {
	return &crawler.Payload{
		AllowedDomains:    &[]string{"ufape.edu.br"},
		CanRetry:          boolPtr(false),
		CollectSubdomains: boolPtr(true),
		LowerCaseURLs:     boolPtr(false),
		MaxAttempts:       intPtr(1),
		RemoveFragment:    boolPtr(false),
		Timeout:           intPtr(60),
		Url:               url,
	}
}

type CrawlItem struct {
	URL   string
	Depth int
}

type Crawler struct {
	httpClient *http.Client
	baseURL    string

	queue   *list.List
	visited map[string]struct{}
	result  *FinalResponse

	maxDepth int
}

func NewCrawler(baseURL string, timeout time.Duration, maxDepth int) *Crawler {
	transport := &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     90 * time.Second,
	}

	return &Crawler{
		httpClient: &http.Client{
			Timeout:   timeout,
			Transport: transport,
		},
		baseURL: baseURL,
		queue:   list.New(),
		visited: make(map[string]struct{}),
		result: &FinalResponse{
			Nodes:       []Node{},
			Links:       []Link{},
			GeneratedAt: time.Now().UTC().UnixMilli(),
		},
		maxDepth: maxDepth,
	}
}

func (c *Crawler) Crawl(initialURL string) {
	normalizedInitialURL := c.normalizeLink(initialURL)
	c.enqueue(&CrawlItem{URL: normalizedInitialURL, Depth: 1})
	c.markAsVisited(normalizedInitialURL)

	for c.queue.Len() > 0 {
		item := c.dequeue()
		if item == nil {
			continue
		}

		if item.Depth > c.maxDepth {
			continue
		}

		fmt.Printf("Depth: %d | Crawling: %s\n", item.Depth, item.URL)

		response, err := c.fetchLinks(item.URL)
		if err != nil {
			log.Printf("AVISO: Falha ao buscar %s: %v. Continuando...", item.URL, err)
			continue
		}
		newAvailable := []string{}
		for _, link := range response.Links.Available {
			normalizedLink := c.normalizeLink(link)
			if c.shouldVisit(normalizedLink) {
				c.enqueue(&CrawlItem{URL: normalizedLink, Depth: item.Depth + 1})
				c.markAsVisited(normalizedLink)
				newAvailable = append(newAvailable, normalizedLink)
			}
		}
		response.Links.Available = newAvailable
		c.addResponseToGraph(response, item)
	}
	fmt.Println("Crawling finalizado.")
}

func (c *Crawler) fetchLinks(url string) (*crawler.ResponseDTO, error) {
	payload := NewRequestPayload(url)
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("falha ao serializar payload: %w", err)
	}

	req, err := http.NewRequest("POST", c.baseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("falha ao criar requisição: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("falha na requisição http: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status inesperado: %s", res.Status)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("falha ao ler corpo da resposta: %w", err)
	}

	var apiResponse crawler.ResponseDTO
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return nil, fmt.Errorf("falha ao decodificar JSON: %w", err)
	}

	return &apiResponse, nil
}

func (c *Crawler) addResponseToGraph(response *crawler.ResponseDTO, sourceItem *CrawlItem) {
	node := NewNode(sourceItem.URL, sourceItem.Depth, response)
	c.result.Nodes = append(c.result.Nodes, node)

	for _, targetLink := range response.Links.Available {
		link := NewLink(sourceItem.URL, targetLink)
		c.result.Links = append(c.result.Links, link)
	}
}

func (c *Crawler) normalizeLink(link string) string {
	parsedURL, err := url.Parse(link)
	if err != nil {
		return link
	}
	parsedURL.Fragment = ""
	return strings.TrimSuffix(parsedURL.String(), "/")
}

func (c *Crawler) shouldVisit(u string) bool {
	_, exists := c.visited[u]
	return !exists
}

func (c *Crawler) markAsVisited(u string) {
	c.visited[u] = struct{}{}

	if strings.HasPrefix(u, "https://") {
		httpVersion := "http" + u[5:]
		c.visited[httpVersion] = struct{}{}
	} else if strings.HasPrefix(u, "http://") {
		httpsVersion := "httpss" + u[4:]
		c.visited[httpsVersion] = struct{}{}
	}
}

func (c *Crawler) SaveResult(filename string) error {
	fileData, err := json.MarshalIndent(c.result, "", "  ")
	if err != nil {
		return fmt.Errorf("falha ao serializar resultado: %w", err)
	}

	if err := os.WriteFile(filename, fileData, 0644); err != nil {
		return fmt.Errorf("falha ao escrever no arquivo %s: %w", filename, err)
	}

	fmt.Printf("Resultado salvo com sucesso em %s\n", filename)
	return nil
}

func (c *Crawler) enqueue(item *CrawlItem) {
	c.queue.PushBack(item)
}

func (c *Crawler) dequeue() *CrawlItem {
	if c.queue.Len() == 0 {
		return nil
	}
	element := c.queue.Front()
	c.queue.Remove(element)
	return element.Value.(*CrawlItem)
}

func main() {
	apiURL := os.Getenv("API_URL")
	if apiURL == "" {
		apiURL = "http://localhost:8080/"
		log.Println("AVISO: Variável de ambiente API_URL não definida. Usando URL padrão:", apiURL)
	}
	if MAX_DEPTH == math.MaxInt {
		log.Println("AVISO: MAX_DEPTH está configurado como 'infinito' (math.MaxInt). O crawling pode demorar muito ou nunca terminar.")
	}

	crawler := NewCrawler(apiURL, DEFAULT_TIMEOUT, MAX_DEPTH)
	crawler.Crawl(INITIAL_URL)

	if err := crawler.SaveResult("grafo_salvo.json"); err != nil {
		log.Fatalf("Erro fatal ao salvar o arquivo: %v", err)
	}
}

func boolPtr(b bool) *bool {
	return &b
}

func intPtr(i int) *int {
	return &i
}
