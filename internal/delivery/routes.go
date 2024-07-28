package delivery

import (
	"forum/internal/helpers/template"
	"forum/internal/repository"
	"forum/pkg/logger"
	"net/http"
)

type application struct {
	logger        *logger.Logger
	storage       *repository.Storage
	templateCache *template.TemplateCache
}

func Routes(l *logger.Logger, db *repository.Storage, tc *template.TemplateCache) http.Handler {
	mux := http.NewServeMux()

	app := &application{
		logger:        l,
		storage:       db,
		templateCache: tc,
	}

	fileServer := http.FileServer(http.Dir("./web/static"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", app.HomeHandler)

	mux.HandleFunc("/post/", app.ViewPostHandler)
	mux.HandleFunc("/post/create", app.CreatePostHandler)
	mux.HandleFunc("/post/comment", app.CreateComment)

	mux.HandleFunc("/register", app.RegisterHandler)
	mux.HandleFunc("/login", app.LoginHandler)

	return app.SessionMiddleware(mux)
}
