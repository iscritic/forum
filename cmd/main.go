package main

import (
	"log"
	"net/http"

	"forum/internal/config"
	"forum/internal/delivery"
	"forum/internal/repository"
	"forum/pkg/flog"
	"forum/pkg/tmpl"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	cfg := config.MustLoad()

	lg := flog.NewLogger(cfg.LogLevel)

	lg.Info("Config and logger initialized")
	lg.Debug("config is set: %v\n", cfg)

	db, err := repository.New(cfg.StoragePath)
	if err != nil {
		panic(err)
	}

	lg.Info("Database is connected")

	templateCache := tmpl.NewTemplateCache()
	if err := tmpl.LoadTemplates(templateCache, "./web/html"); err != nil {
		log.Fatalf("failed to load templates: %v", err)
	}

	lg.Info("Template Cashe map is initialized")

	srv := &http.Server{
		Addr:        cfg.HTTPServer.Address,
		Handler:     delivery.Routes(lg, db, templateCache),
		ReadTimeout: config.ParseTime(cfg.HTTPServer.ReadTimeout),
		IdleTimeout: config.ParseTime(cfg.HTTPServer.IdleTimeout),
	}

	lg.Info("Server is listening: http://localhost%s\n", cfg.HTTPServer.Address)
	err = srv.ListenAndServe()
	if err != nil {
		lg.Error(err.Error())
	}
}
