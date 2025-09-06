package crawler

import (
	"net/url"
	"reflect"
	"slices"
	"strings"
	"testing"

	"golang.org/x/net/html"
)

func parseHTML(t *testing.T, htmlString string) *html.Node {
	doc, err := html.Parse(strings.NewReader(htmlString))
	if err != nil {
		t.Fatalf("Falha ao fazer parse do HTML de teste: %v", err)
	}
	return doc
}

func TestGetTitle(t *testing.T) {
	testCases := []struct {
		name          string
		htmlContent   string
		expectedTitle string
	}{
		{
			name:          "Título Simples",
			htmlContent:   `<html><head><title>Título da Página</title></head><body></body></html>`,
			expectedTitle: "Título da Página",
		},
		{
			name:          "Título com Espaços",
			htmlContent:   `<html><head><title>  Meu Título com Espaços  </title></head><body></body></html>`,
			expectedTitle: "Meu Título com Espaços",
		},
		{
			name:          "Sem Tag de Título",
			htmlContent:   `<html><head></head><body><h1>Olá</h1></body></html>`,
			expectedTitle: "[Empty title]",
		},
		{
			name:          "Tag de Título Vazia",
			htmlContent:   `<html><head><title></title></head><body></body></html>`,
			expectedTitle: "[Empty title]",
		},
		{
			name:          "Título Aninhado Profundamente",
			htmlContent:   `<html><body><div><p><title>Título Profundo</title></p></div></body></html>`,
			expectedTitle: "Título Profundo",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			doc := parseHTML(t, tc.htmlContent)

			actualTitle := GetTitle(doc)

			if actualTitle != tc.expectedTitle {
				t.Errorf("Título incorreto: got %q, want %q", actualTitle, tc.expectedTitle)
			}
		})
	}
}

func TestExtractLinks(t *testing.T) {
	baseURLStr := "https://example.com"
	baseURL, err := url.Parse(baseURLStr)
	if err != nil {
		t.Fatalf("URL base inválida: %v", err)
	}

	t.Run("Extração Básica", func(t *testing.T) {
		htmlContent := `
            <a href="/pagina1">Página 1</a>
            <a href="https://outrodominio.com/page">Outro Domínio</a>
            <a href="https://example.com/pagina2">Página 2</a>
            <a href="/pagina1">Página 1 (Duplicado)</a>
            <a href="mailto:teste@example.com">Email</a>
            <a href="tel:+123456">Telefone</a>
        `
		doc := parseHTML(t, htmlContent)

		collectSubdomains := false
		removeFragment := false
		lowerCaseUrls := false
		payload := &CorrectPayload{
			AllowedDomains:    &[]string{"example.com"},
			CollectSubdomains: &collectSubdomains,
			RemoveFragment:    &removeFragment,
			LowerCaseURLs:     &lowerCaseUrls,
		}

		links := ExtractLinks(baseURL, doc, payload)

		expectedAvailable := []string{"https://example.com/pagina1", "https://example.com/pagina2"}
		expectedUnavailable := []string{"https://outrodominio.com/page"}

		slices.Sort(links.Available)
		slices.Sort(expectedAvailable)
		slices.Sort(links.Unavailable)
		slices.Sort(expectedUnavailable)

		if !reflect.DeepEqual(links.Available, expectedAvailable) {
			t.Errorf("Links disponíveis incorretos:\ngot  %v\nwant %v", links.Available, expectedAvailable)
		}
		if !reflect.DeepEqual(links.Unavailable, expectedUnavailable) {
			t.Errorf("Links indisponíveis incorretos:\ngot  %v\nwant %v", links.Unavailable, expectedUnavailable)
		}
	})

	t.Run("Coletando Subdomínios", func(t *testing.T) {
		htmlContent := `<a href="https://blog.example.com/post1">Post do Blog</a>`
		doc := parseHTML(t, htmlContent)

		collectSubdomains := true
		removeFragment := false
		lowerCaseUrls := false
		payload := &CorrectPayload{
			AllowedDomains:    &[]string{"example.com"},
			CollectSubdomains: &collectSubdomains,
			RemoveFragment:    &removeFragment,
			LowerCaseURLs:     &lowerCaseUrls,
		}

		links := ExtractLinks(baseURL, doc, payload)

		expectedAvailable := []string{"https://blog.example.com/post1"}
		if !reflect.DeepEqual(links.Available, expectedAvailable) {
			t.Errorf("Links disponíveis incorretos (subdomínio): got %v want %v", links.Available, expectedAvailable)
		}
		if len(links.Unavailable) != 0 {
			t.Errorf("Não deveria haver links indisponíveis, mas foram encontrados: %v", links.Unavailable)
		}
	})
}
