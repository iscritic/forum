package internal

import (
	"forum/internal/sqlite"
	"forum/pkg/logger"
	"net/http"
)

type application struct {
	logger  *logger.Logger
	storage *sqlite.Storage
}

func Routes(l *logger.Logger, db *sqlite.Storage) *http.ServeMux {
	mux := http.NewServeMux()

	app := &application{
		logger:  l,
		storage: db,
	}

	mux.HandleFunc("/", app.home)
	//mux.HandleFunc("/post/view", app.viewPost)
	mux.HandleFunc("/post/create", app.createPost)
	//mux.HandleFunc("/post/update", app.updatePost)
	//
	//mux.HandleFunc("/user")
	//mux.HandleFunc("/user/create", app.createUser)
	return mux
}
