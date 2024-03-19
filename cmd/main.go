package main

import (
	"database/sql"
	"forum/internal"
	"forum/pkg/logger"
	"log"
	"net/http"
)

func main() {
	lg := logger.NewLogger()

	db, err := sql.Open("sqlite3", "./db/store.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	srv := &http.Server{
		Addr:    "0.0.0.0:7000",
		Handler: internal.Routes(lg),
	}

	lg.InfoLog.Println("Listening serven on http://localhost:7000...")
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatalln(err)
	}
}
