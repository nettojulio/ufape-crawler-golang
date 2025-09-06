package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/nettojulio/ufape-crawler-golang/config"
	"github.com/nettojulio/ufape-crawler-golang/server"

	"github.com/nettojulio/ufape-crawler-golang/docs"
)

var Version = "development"

// @title UFAPE Crawler API
// @version 0.0.1
// @description API para realizar crawling de websites. Recebe uma URL e retorna o conteúdo e os links encontrados.
// @termsOfService http://swagger.io/terms/

// @contact.name Júlio Netto
// @contact.url https://github.com/nettojulio/ufape-crawler-golang
// @contact.email nettojulio@hotmail.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @BasePath /
func main() {
	flag.Parse()
	fmt.Printf("Version: %s\n", Version)

	cfg, err := config.Load(Version)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	docs.SwaggerInfo.Host = cfg.Host

	e := server.New(cfg)
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", cfg.Port)))
}
