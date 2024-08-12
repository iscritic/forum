package delivery

import (
	"net/http"

	"forum/internal/helpers/template"
	"forum/internal/repository"
	"forum/pkg/logger"
	"forum/pkg/mw"
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
	mux.HandleFunc("/category/", app.SortedByCategoryHandler)

	mux.HandleFunc("/register", app.RegisterHandler)
	mux.HandleFunc("/login", app.LoginHandler)

	// require authentication
	protected := mw.New(app.requiredAuthentication)

	mux.Handle("/post/create", protected.ThenFunc(app.CreatePostHandler))
	mux.Handle("/post/comment", protected.ThenFunc(app.CreateComment))

	mux.Handle("/createdposts", protected.ThenFunc(app.MyPostsHandler))
	mux.Handle("/likedposts", protected.ThenFunc(app.MyLikedPostsHandler))

	// mux.Handle("/likedposts", protected.ThenFunc())

	mux.Handle("/post/like", protected.ThenFunc(app.LikePostHandler))
	mux.Handle("/post/dislike", protected.ThenFunc(app.DislikePostHandler))
	mux.Handle("/comment/like", protected.ThenFunc(app.LikeCommentHandler))
	mux.Handle("/comment/dislike", protected.ThenFunc(app.DislikeCommentHandler))

	mux.Handle("/post/like", protected.ThenFunc(app.LikePostHandler))
	mux.Handle("/post/dislike", protected.ThenFunc(app.DislikePostHandler))
	mux.Handle("/comment/like", protected.ThenFunc(app.LikeCommentHandler))
	mux.Handle("/comment/dislike", protected.ThenFunc(app.DislikeCommentHandler))

	// standard midllewares for all routes
	standard := mw.New(app.logRequest, app.recoverPanic, secureHeaders)

	return standard.Then(mux)
}

// hello world
