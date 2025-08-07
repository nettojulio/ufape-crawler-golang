package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/nettojulio/ufape-crawler-golang/config"
	"github.com/nettojulio/ufape-crawler-golang/server"
)

var Version = "development"

func main() {
	flag.Parse()
	fmt.Printf("Version: %s\n", Version)

	cfg, err := config.Load(Version)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	e := server.New(cfg)
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", cfg.Port)))
}
