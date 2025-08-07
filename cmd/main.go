package main

import (
	"flag"
	"log"
	"time"

	"github.com/nettojulio/ufape-crawler-golang/internal/crawler"
	"github.com/nettojulio/ufape-crawler-golang/utils/configs"
)

var (
	Version = "development"
)

func init() {
	flag.IntVar(&configs.MaxDepth, "depth", 1, "Profundidade máxima de crawling.")
	flag.BoolVar(&configs.RemoveFragment, "remove-fragment", true, "Remove fragmentos de URL (#...).")
	flag.BoolVar(&configs.LowerCaseURLs, "lower-case", false, "Converte todas as URLs para minúsculas.")
	flag.DurationVar(&configs.RequestTimeout, "timeout", 90*time.Second, "Timeout para as requisições HTTP.")
	flag.BoolVar(&configs.FliterLinks, "filter-links", true, "Filtra links já visitados.")
	flag.StringVar(&configs.InitialURL, "initial-url", "https://ufape.edu.br", "URL inicial para o crawler.")
}

func main() {
	flag.Parse()
	log.Printf("Version: %s", Version)
	configs.DisplayConfigs()
	crawler.InitialSearch()
}
