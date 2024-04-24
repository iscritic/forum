package transport

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

	fileServer := http.FileServer(http.Dir("./web/static"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", app.HomeHandler)

	mux.HandleFunc("/post/view", app.ViewPostHandler)
	mux.HandleFunc("/post/create", app.CreatePostHandler)
	mux.HandleFunc("/post/comment", app.CreateComment)

	mux.HandleFunc("/register", app.RegisterHandler)
	mux.HandleFunc("/login", app.LoginHandler)
	//mux.HandleFunc("/post/update", app.updatePost)
	//
	//mux.HandleFunc("/user")
	//mux.HandleFunc("/user/create", app.createUser)
	return mux
}
