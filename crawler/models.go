package crawler

import "net/url"

// CorrectPayload define a estrutura do corpo da requisição para o endpoint de crawling.
// Esta struct contém todas as configurações para a operação de crawling.
type CorrectPayload struct {
	Url               string    `json:"url" validate:"required,url" example:"http://ufape.edu.br"`
	Timeout           *int      `json:"timeout,omitempty" example:"60"`
	RemoveFragment    *bool     `json:"remove_fragment,omitempty" example:"false"`
	AllowedDomains    *[]string `json:"allowed_domains,omitempty" example:"ufape.edu.br"`
	CollectSubdomains *bool     `json:"collect_subdomains,omitempty" example:"true"`
	LowerCaseURLs     *bool     `json:"lower_case_urls,omitempty" example:"false"`
	CanRetry          *bool     `json:"can_retry,omitempty" example:"false"`
	MaxAttempts       *int      `json:"max_attempts,omitempty" example:"1"`
}

// LinksResponse agrupa os links encontrados na página entre disponíveis e indisponíveis.
type LinksResponse struct {
	Available   []string `json:"available" example:"http://ufape.edu.br/link-valido"`
	Unavailable []string `json:"unavailable" example:"http://ufape.edu.br/link-quebrado"`
}

// DetailsResponseDTO é a struct de resposta para a API, usando tipos simples.
type DetailsResponseDTO struct {
	CorrectURL string     `json:"correctUrl" example:"http://ufape.edu.br"`
	Original   URLDetails `json:"original"`
	Modified   URLDetails `json:"modified"`
}

// ResponseCrawlDTO é a resposta principal da API, usando o DTO de detalhes.
type ResponseCrawlDTO struct {
	StatusCode  int                `json:"statusCode" example:"200"`
	ContentType string             `json:"contentType" example:"text/html; charset=utf-8"`
	ElapsedTime int64              `json:"elapsedTime" example:"150"`
	Links       LinksResponse      `json:"links"`
	Title       string             `json:"title" example:"Universidade Federal do Agreste de Pernambuco"`
	Details     DetailsResponseDTO `json:"details"`
}

// DetailsResponse fornece detalhes sobre as URLs original e modificada usadas no crawling.
type DetailsResponse struct {
	CorrectURL string  `json:"correctUrl" example:"http://ufape.edu.br"`
	Original   url.URL `json:"original"`
	Modified   url.URL `json:"modified"`
}

// ResponseCrawl define a estrutura da resposta completa do endpoint de crawling.
type ResponseCrawl struct {
	StatusCode  int             `json:"statusCode" example:"200"`
	ContentType string          `json:"contentType" example:"text/html; charset=utf-8"`
	ElapsedTime int64           `json:"elapsedTime" example:"150"`
	Links       LinksResponse   `json:"links"`
	Title       string          `json:"title" example:"Universidade Federal do Agreste de Pernambuco"`
	Details     DetailsResponse `json:"details"`
}

// APIHealth define a estrutura da resposta do endpoint de verificação de saúde da API.
type APIHealth struct {
	Status  string `json:"status" example:"OK"`
	Version string `json:"version" example:"1.0.0"`
}

// URLDetails fornece uma representação detalhada de uma URL, incluindo seus componentes.
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
