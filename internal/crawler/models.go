package crawler

import (
	"io"
	"net/url"
	"time"
)

// Payload define a estrutura do corpo da requisição para o endpoint de crawling.
type Payload struct {
	Url               string    `json:"url" validate:"required,url" example:"http://ufape.edu.br"`
	Timeout           *int      `json:"timeout,omitempty" example:"60"`
	RemoveFragment    *bool     `json:"remove_fragment,omitempty" example:"false"`
	AllowedDomains    *[]string `json:"allowed_domains,omitempty" example:"ufape.edu.br"`
	CollectSubdomains *bool     `json:"collect_subdomains,omitempty" example:"true"`
	LowerCaseURLs     *bool     `json:"lower_case_urls,omitempty" example:"false"`
	CanRetry          *bool     `json:"can_retry,omitempty" example:"false"`
	MaxAttempts       *int      `json:"max_attempts,omitempty" example:"1"`
}

// LinksResponse agrupa os links encontrados.
type LinksResponse struct {
	Available   []string `json:"available" example:"http://ufape.edu.br/link-valido"`
	Unavailable []string `json:"unavailable" example:"http://ufape.edu.br/link-quebrado"`
}

// URLDetails fornece uma representação detalhada de uma URL.
type URLDetails struct {
	Scheme      string          `json:"Scheme" example:"https"`
	Opaque      string          `json:"Opaque" example:""`
	User        *UserURLDetails `json:"User"`
	Host        string          `json:"Host" example:"ufape.edu.br"`
	Path        string          `json:"Path" example:"/"`
	RawPath     string          `json:"RawPath" example:""`
	OmitHost    bool            `json:"OmitHost" example:"false"`
	ForceQuery  bool            `json:"ForceQuery" example:"false"`
	RawQuery    string          `json:"RawQuery" example:""`
	Fragment    string          `json:"Fragment" example:""`
	RawFragment string          `json:"RawFragment" example:""`
}

// UserURLDetails fornece detalhes sobre as credenciais de usuário em uma URL.
type UserURLDetails struct {
	Username    string `json:"username" example:"user"`
	Password    string `json:"password" example:"password"`
	PasswordSet bool   `json:"passwordSet" example:"true"`
}

// DetailsResponseDTO é a struct de detalhes da URL para a resposta da API.
type DetailsResponseDTO struct {
	CorrectURL string     `json:"correctUrl" example:"http://ufape.edu.br"`
	Original   URLDetails `json:"original"`
	Modified   URLDetails `json:"modified"`
}

// ResponseDTO é a resposta principal da API.
type ResponseDTO struct {
	StatusCode  int                `json:"statusCode" example:"200"`
	ContentType string             `json:"contentType" example:"text/html; charset=utf-8"`
	ElapsedTime int64              `json:"elapsedTime" example:"150"`
	Links       LinksResponse      `json:"links"`
	Title       string             `json:"title" example:"Universidade Federal do Agreste de Pernambuco"`
	Details     DetailsResponseDTO `json:"details"`
}

// CrawlResult é um modelo interno para transportar o resultado do crawling.
type CrawlResult struct {
	StatusCode  int
	ContentType string
	ElapsedTime time.Duration
	Links       LinksResponse
	Title       string
	Body        io.ReadCloser
	FinalURL    *url.URL
}

// APIHealth define a estrutura da resposta do endpoint de verificação de saúde.
type APIHealth struct {
	Status  string `json:"status" example:"OK"`
	Version string `json:"version" example:"1.0.0"`
}
