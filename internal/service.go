package internal

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"time"

	"github.com/emirpasic/gods/sets/hashset"
	"github.com/nettojulio/ufape-crawler-golang/utils/configs"
	"golang.org/x/net/html"
)

type CorrectPayload struct {
	Url               string    `json:"url" validate:"required,url"`  // Obrigatório
	Timeout           *int      `json:"timeout,omitempty"`            // Opcional
	RemoveFragment    *bool     `json:"remove_fragment,omitempty"`    // Opcional
	AllowedDomains    *[]string `json:"allowed_domains,omitempty"`    // Opcional
	CollectSubdomains *bool     `json:"collect_subdomains,omitempty"` // Opcional
}

type LinksResponse struct {
	Available   []string `json:"available"`
	Unavailable []string `json:"unavailable"`
}

type DetailsResponse struct {
	CorrectURL string  `json:"correctUrl"`
	Original   url.URL `json:"original"`
	Modified   url.URL `json:"modified"`
}
type ResponseCrawl struct {
	StatusCode  int             `json:"statusCode"`
	ContentType string          `json:"contentType"`
	ElapsedTime string          `json:"elapsedTime"`
	Links       LinksResponse   `json:"links"`
	Title       string          `json:"title"`
	Details     DetailsResponse `json:"details"`
}

func CrawlerService(payload CorrectPayload, originalUrlDetails, modifiedUrlDetails url.URL) (ResponseCrawl, error) {
	httpClient := http.Client{
		Timeout: time.Duration(*payload.Timeout) * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 10 {
				return fmt.Errorf("stopped after 10 redirects")
			}
			return nil
		},
	}

	if modifiedUrlDetails.Host == "" {
		//return nil, fmt.Errorf("URL inválida: %s", urlStr)
	}

	req, err := http.NewRequest("GET", modifiedUrlDetails.String(), nil)
	if err != nil {
		return ResponseCrawl{}, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64)")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Connection", "keep-alive")

	start := time.Now()
	resp, err := httpClient.Do(req)
	elapsed := time.Since(start)
	if err != nil {
		return ResponseCrawl{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Printf("Error closing response body: %v\n", err)
		}
	}(resp.Body)

	contentType := resp.Header.Get("Content-Type")

	if resp.StatusCode != http.StatusOK {
		errResp := ResponseCrawl{
			StatusCode:  resp.StatusCode,
			ContentType: contentType,
			ElapsedTime: fmt.Sprintf("%+v", elapsed),
			Links: LinksResponse{
				Available:   []string{},
				Unavailable: []string{},
			},
			Title: "",
			Details: DetailsResponse{
				CorrectURL: modifiedUrlDetails.String(),
				Original:   originalUrlDetails,
				Modified:   modifiedUrlDetails,
			},
		}
		return errResp, nil
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return ResponseCrawl{}, err
	}

	title := getTitle(doc)

	base, err := url.Parse(modifiedUrlDetails.String())
	if err != nil {
		return ResponseCrawl{}, err
	}

	links := extractLinks(base, doc, &payload)

	res := ResponseCrawl{
		StatusCode:  resp.StatusCode,
		ContentType: contentType,
		ElapsedTime: fmt.Sprintf("%+v", elapsed),
		Links:       links,
		Title:       title,
		Details: DetailsResponse{
			CorrectURL: modifiedUrlDetails.String(),
			Original:   originalUrlDetails,
			Modified:   modifiedUrlDetails,
		},
	}

	return res, nil

	//response := {
	//	"statusCode":  resp.StatusCode,
	//	"contentType": contentType,
	//	"elapsedTime": elapsed,
	//	"links":       links,
	//	"title":       title,
	//	"details": map[string]interface{}{
	//		"correctUrl": modifiedUrlDetails.String(),
	//		"original":   originalUrlDetails,
	//		"modified":   modifiedUrlDetails,
	//	},
	//}
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

func extractLinks(base *url.URL, body *html.Node, payload *CorrectPayload) LinksResponse {
	var links LinksResponse = LinksResponse{
		Available:   []string{},
		Unavailable: []string{},
	}
	uniqueLinks := hashset.New()
	curr := normalizeURL(base.String())
	curateNodes(body, base, curr, uniqueLinks, &links, payload)
	return links
}

func curateNodes(n *html.Node, base *url.URL, curr string, uniqueLinks *hashset.Set, links *LinksResponse, payload *CorrectPayload) {
	if isAnchorNode(n) {
		processAnchor(n, base, curr, uniqueLinks, links, payload)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		curateNodes(c, base, curr, uniqueLinks, links, payload)
	}
}

func isAnchorNode(n *html.Node) bool {
	return n.Type == html.ElementNode && n.Data == "a"
}

func processAnchor(n *html.Node, base *url.URL, curr string, uniqueLinks *hashset.Set, links *LinksResponse, payload *CorrectPayload) {
	for _, attr := range n.Attr {
		if attr.Key != "href" {
			continue
		}
		processHref(attr.Val, base, curr, uniqueLinks, links, payload)
	}
}

func processHref(href string, base *url.URL, curr string, uniqueLinks *hashset.Set, links *LinksResponse, payload *CorrectPayload) {
	if isNonNavigableLink(href) {
		return
	}
	u, err := base.Parse(href)
	if err != nil {
		return
	}
	if len(u.Host) != 0 && len(u.Scheme) != 0 {
		addIfUnique(u.String(), curr, uniqueLinks, links, payload)
	}
}

func isNonNavigableLink(href string) bool {
	return strings.HasPrefix(href, "mailto:") || strings.HasPrefix(href, "tel:")
}

//func isSameDomain(u, base *url.URL, payload *CorrectPayload) bool {
//	allowed := slices.Contains((*payload.AllowedDomains), u.Host)
//	fmt.Println(allowed)
//	return u.Host == base.Host || strings.HasSuffix(u.Host, "."+base.Host)
//}

func addIfUnique(rawURL, curr string, uniqueLinks *hashset.Set, links *LinksResponse, payload *CorrectPayload) {
	norm := normalizeURL(rawURL)
	if norm != curr && !uniqueLinks.Contains(norm) {
		uniqueLinks.Add(norm)
		info, err := url.Parse(norm)
		if err != nil {
			return
		}
		if strings.Contains(norm, "transparencia") {
			fmt.Println("")
		}
		if slices.Contains(*payload.AllowedDomains, info.Host) || (*payload.CollectSubdomains && isSubdomainHost(info.Host, *payload.AllowedDomains)) {
			links.Available = append(links.Available, norm)
		} else {
			links.Unavailable = append(links.Unavailable, norm)
		}
	}
}

func normalizeURL(raw string) string {
	return NormalizeURL(raw)
}

func isSubdomainHost(host string, domains []string) bool {
	for _, domain := range domains {
		if host == domain || strings.HasSuffix(host, "."+domain) {
			return true
		}
	}
	return false
}

func NormalizeURL(raw string) string {
	u, err := url.Parse(raw)
	if err != nil {
		return raw
	}
	if configs.RemoveFragment {
		u.Fragment = ""
	}
	pass, ok := u.User.Password()
	passEmpty := !ok || pass == ""
	if u.User != nil {
		if u.User.Username() == "" && passEmpty {
			u.User = nil
		}
	}

	if strings.Contains(u.Host, "www") {
		u.Host = strings.TrimPrefix(u.Host, "www.")
	}

	norm := u.String()
	if configs.LowerCaseURLs {
		norm = strings.ToLower(norm)
	}

	return strings.TrimRight(norm, "/")
}
