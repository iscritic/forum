package main

import (
	"forum/internal/config"
	"forum/internal/delivery"
	"forum/internal/repository"
	"forum/pkg/flog"
	"forum/pkg/tmpl"
	"log"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	cfg := config.MustLoad()
	lg := flog.NewLogger(cfg.LogLevel)

	lg.Debug("config is set: %s\n", cfg)

	db, err := repository.New(cfg.StoragePath)
	if err != nil {
		panic(err)
	}

	templateCache := tmpl.NewTemplateCache()
	if err := tmpl.LoadTemplates(templateCache, "./web/html"); err != nil {
		log.Fatalf("failed to load templates: %v", err)
	}

	httpPort := os.Getenv("PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	srv := &http.Server{
		Addr:        ":" + httpPort,
		Handler:     delivery.Routes(lg, db, templateCache),
		ReadTimeout: config.ParseTime(cfg.HTTPServer.ReadTimeout),
		IdleTimeout: config.ParseTime(cfg.HTTPServer.IdleTimeout),
	}

	lg.Info("Server is listening: http://localhost:%s\n", httpPort)
	err = srv.ListenAndServe()
	if err != nil {
		lg.Error(err.Error())
	}
}
