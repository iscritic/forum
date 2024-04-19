package main

import (
	"forum/internal"
	"forum/internal/config"
	"forum/internal/sqlite"
	"forum/pkg/logger"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func main() {

	lg := logger.NewLogger()

	cfg := config.MustLoad()
	lg.InfoLog.Printf("config is set: %s\n", cfg)

	db, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		panic(err)
	}

	srv := &http.Server{
		Addr:        cfg.HTTPServer.Address,
		Handler:     internal.Routes(lg, db),
		ReadTimeout: config.ParseTime(cfg.HTTPServer.ReadTimeout),
		IdleTimeout: config.ParseTime(cfg.HTTPServer.IdleTimeout),
	}

	lg.InfoLog.Printf("Listening serven on http://%s...\n", cfg.HTTPServer.Address)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatalln(err)
	}
}
