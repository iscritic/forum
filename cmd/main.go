package main

import (
	"forum/pkg/tmpl"
	"log"
	"net/http"

	"forum/internal/config"
	"forum/internal/delivery"
	"forum/internal/repository"
	"forum/pkg/flog"

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

	srv := &http.Server{
		Addr:        cfg.HTTPServer.Address,
		Handler:     delivery.Routes(lg, db, templateCache),
		ReadTimeout: config.ParseTime(cfg.HTTPServer.ReadTimeout),
		IdleTimeout: config.ParseTime(cfg.HTTPServer.IdleTimeout),
	}

	lg.Info("Server is listening: http://%s\n", srv.Addr)
	err = srv.ListenAndServe()
	if err != nil {
		lg.Error(err.Error())
	}
}
