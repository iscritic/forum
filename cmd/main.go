package main

import (
	"forum/internal/config"
	"forum/internal/delivery"
	"forum/internal/helpers/template"
	"forum/internal/repository"
	"forum/pkg/logger"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func main() {

	lg := logger.NewLogger()

	cfg := config.MustLoad()
	lg.InfoLog.Printf("config is set: %s\n", cfg)

	db, err := repository.New(cfg.StoragePath)
	if err != nil {
		panic(err)
	}
	templateCache := template.NewTemplateCache()
	if err := template.LoadTemplates(templateCache, "./web/html"); err != nil {
		log.Fatalf("failed to load templates: %v", err)
	}

	srv := &http.Server{
		Addr:        cfg.HTTPServer.Address,
		Handler:     delivery.Routes(lg, db, templateCache),
		ReadTimeout: config.ParseTime(cfg.HTTPServer.ReadTimeout),
		IdleTimeout: config.ParseTime(cfg.HTTPServer.IdleTimeout),
	}

	lg.InfoLog.Printf("Listening serven on http://%s...\n", cfg.HTTPServer.Address)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatalln(err)
	}
}
