package internal

import "net/http"

func routes() *http.ServeMux {

	mux := http.NewServeMux()

	mux.HandleFunc("/", Home)

	return mux

}
