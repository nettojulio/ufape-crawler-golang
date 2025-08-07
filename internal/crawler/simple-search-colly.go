package crawler

import (
	"fmt"
	"log"
	"strings"

	"github.com/gocolly/colly"
)

var (
	DefaultUrl     = "https://ufape.edu.br/"
	AllowedDomains = []string{"ufape.edu.br"}
)

var mapLinks = make(map[string]bool)

func SimpleSearchWithColly() {
	// Instantiate default collector
	c := colly.NewCollector(
		// Visit only domains and subdomains: ufape.edu.br
		colly.AllowedDomains(AllowedDomains...),
	)

	// On every a element which has href attribute call callback
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		// Print link

		fmt.Printf("Link found: %q -> %s\n", e.Text, link)
		// Visit link found on page
		// Only those links are visited which are in AllowedDomains

		// Uncomment the next line to actually visit the link

		precious := e.Request.AbsoluteURL(link)
		if precious == "" {
			log.Println("Internal Link:", link)
			return
		} else {
			for _, domain := range AllowedDomains {
				if strings.Contains(precious, domain) && mapLinks[precious] == false {
					fmt.Println("Putting in next request:", precious)
					mapLinks[precious] = true
				} else {
					if strings.Contains(precious, "javascript") {
						log.Println("Skipping javascript link:", precious)
						return
					}
					log.Println("Skipping link:", precious)
				}
			}

		}

		//c.Visit(e.Request.AbsoluteURL(link))
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("SC is: ", r.StatusCode)
	})

	// Start scraping on https://hackerspaces.org
	mapLinks[DefaultUrl] = true
	c.Visit(DefaultUrl)
	fmt.Println("RESUME")
}
