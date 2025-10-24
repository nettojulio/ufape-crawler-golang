package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nettojulio/ufape-crawler-golang/docs"
	"github.com/nettojulio/ufape-crawler-golang/internal/config"
	"github.com/nettojulio/ufape-crawler-golang/internal/crawler"
	server "github.com/nettojulio/ufape-crawler-golang/internal/http"
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
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	cfg, err := config.Load(Version)
	if err != nil {
		logger.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	logger.Info("starting application", "version", Version)

	if err != nil {
		logger.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	docs.SwaggerInfo.Host = cfg.Host

	httpClient := crawler.NewHTTPClient(60 * time.Second)
	crawlerService := crawler.NewService(httpClient)

	e := server.NewServer(cfg, crawlerService)

	go func() {
		addr := fmt.Sprintf(":%d", cfg.Port)
		logger.Info("server starting", "address", addr)
		if err := e.Start(addr); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	logger.Info("server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}

	logger.Info("server shutdown complete")
}
