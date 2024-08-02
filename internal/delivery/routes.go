package delivery

import (
	"forum/internal/helpers/template"
	"forum/internal/repository"
	"forum/pkg/logger"
	"forum/pkg/mw"
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
	mux.HandleFunc("/category/", app.SortedByCategory)

	mux.HandleFunc("/register", app.RegisterHandler)
	mux.HandleFunc("/login", app.LoginHandler)

	// require authentification
	protected := mw.New(app.requiredAuthentication)

	mux.Handle("/post/create", protected.ThenFunc(app.CreatePostHandler))
	mux.Handle("/post/comment", protected.ThenFunc(app.CreateComment))
	mux.Handle("/post/like", protected.ThenFunc(app.LikeHandler))

	// standard midllewares for all routes
	standard := mw.New(app.logRequest, app.recoverPanic, secureHeaders)

	return standard.Then(mux)
}
