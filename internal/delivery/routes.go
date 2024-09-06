package delivery

import (
	"net/http"

	"forum/pkg/tmpl"

	"forum/internal/repository"
	"forum/pkg/flog"
	"forum/pkg/mw"
)

type application struct {
	log       *flog.Logger
	storage   *repository.Storage
	tmplcache *tmpl.TemplateCache
}

func Routes(l *flog.Logger, db *repository.Storage, tc *tmpl.TemplateCache) http.Handler {
	mux := http.NewServeMux()

	app := &application{
		log:       l,
		storage:   db,
		tmplcache: tc,
	}

	fileServer := http.FileServer(http.Dir("./web/static"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", app.HomeHandler)
	mux.HandleFunc("/post/", app.ViewPostHandler)
	mux.HandleFunc("/category/", app.SortedByCategoryHandler)

	mux.HandleFunc("/register", app.RegisterHandler)
	mux.HandleFunc("/login", app.LoginHandler)

	// require authentication
	protected := mw.New(app.requireRole("user", "moderator", "admin"))

	mux.Handle("/post/create", protected.ThenFunc(app.CreatePostHandler))
	mux.Handle("/post/comment", protected.ThenFunc(app.CreateComment))

	mux.Handle("/post/edit/", protected.ThenFunc(app.EditPostHandler))
	mux.Handle("/post/delete/", protected.ThenFunc(app.DeletePostHandler))

	mux.Handle("/comment/edit/", protected.ThenFunc(app.EditCommentHandler))
	mux.Handle("/comment/delete/", protected.ThenFunc(app.DeleteCommentHandler))

	mux.Handle("/createdposts", protected.ThenFunc(app.MyPostsHandler))
	mux.Handle("/icommented", protected.ThenFunc(app.IcommentedPostsHandler))
	mux.Handle("/likedposts", protected.ThenFunc(app.MyLikedPostsHandler))
	mux.Handle("/dislikedposts", protected.ThenFunc(app.MyDislikedPostsHandler))

	mux.Handle("/post/like", protected.ThenFunc(app.LikePostHandler))
	mux.Handle("/post/dislike", protected.ThenFunc(app.DislikePostHandler))
	mux.Handle("/comment/like", protected.ThenFunc(app.LikeCommentHandler))
	mux.Handle("/comment/dislike", protected.ThenFunc(app.DislikeCommentHandler))

	mux.Handle("/logout", protected.ThenFunc(app.LogoutHandler))

	// standard midllewares for all routes
	standard := mw.New(app.logRequest, app.recoverPanic, secureHeaders, app.sessionManager)

	return standard.Then(mux)
}
