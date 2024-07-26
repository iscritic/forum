package transport

import (
	"forum/internal/helpers/template"
	"forum/internal/sqlite"
	"forum/pkg/logger"
	"net/http"
)

type application struct {
	logger        *logger.Logger
	storage       *sqlite.Storage
	templateCache *template.TemplateCache
}

func Routes(l *logger.Logger, db *sqlite.Storage, tc *template.TemplateCache) http.Handler {
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
	//mux.HandleFunc("/post/update", app.updatePost)
	//
	//mux.HandleFunc("/user")
	//mux.HandleFunc("/user/create", app.createUser)
	return app.SessionMiddleware(mux)
}

func (app *application) SessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// middleware logic
		next.ServeHTTP(w, r)
	})
}
