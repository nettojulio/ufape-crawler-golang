package crawler

import "net/url"

type CorrectPayload struct {
	Url               string    `json:"url" validate:"required,url"`
	Timeout           *int      `json:"timeout,omitempty"`
	RemoveFragment    *bool     `json:"remove_fragment,omitempty"`
	AllowedDomains    *[]string `json:"allowed_domains,omitempty"`
	CollectSubdomains *bool     `json:"collect_subdomains,omitempty"`
	LowerCaseURLs     *bool     `json:"lower_case_urls,omitempty"`
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
