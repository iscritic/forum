package main

import (
	"fmt"
	"forum/internal/config"
)

func main() {

	//TODO init config

	cfg := config.MustLoad()

	fmt.Println(cfg)

	//TODO init logger

	//TODO init storage

	//TODO init router

	//TODO run server

	//lg := logger.NewLogger()
	//
	//db, err := sql.Open("sqlite3", "./db/store.db")
	//if err != nil {
	//	panic(err)
	//}
	//defer db.Close()
	//
	//srv := &http.Server{
	//	Addr:    "0.0.0.0:7000",
	//	Handler: internal.Routes(lg),
	//}
	//
	//lg.InfoLog.Println("Listening serven on http://localhost:7000...")
	//err = srv.ListenAndServe()
	//if err != nil {
	//	log.Fatalln(err)
	//}
}
